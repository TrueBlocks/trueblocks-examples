package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Posting represents a single ledger event or correction
type Posting struct {
	Statement struct {
		BlockNumber       int
		TransactionIndex  int
		LogIndex          int
		AssetAddress      string
		AccountedFor      string
		CheckpointBalance int64
		TentativeBalance  int64
	}
	CorrectionIndex  int
	CorrectionReason string
	EventAmount      int64
	BeginBalance     int64
	RowIndex         int
}

func (p Posting) Reconciled(isFinal bool) string {
	if p.CorrectionIndex != 0 {
		return "true"
	}
	if !isFinal {
		return "unknown"
	}
	return strconv.FormatBool(p.BeginBalance+p.EventAmount == p.Statement.CheckpointBalance)
}

type Reconciler struct {
	mu                sync.Mutex
	addressOfInterest string
	runningBalances   map[string]int64
	seenBlocks        map[string]int
	correctionCounter int
	counterMu         sync.Mutex
	rowIndexCounter   int
	rowIndexMu        sync.Mutex
}

var (
	apps [][2]int
)

func (r *Reconciler) GetPostingChannel(block, tx int) <-chan Posting {
	ch := make(chan Posting)
	go func() {
		defer close(ch)
		for _, p := range logsByTx[mapKey(block, tx, 0)] {
			if p.Statement.AccountedFor != r.addressOfInterest {
				continue
			}
			eaKey := fmt.Sprintf("%d|%d|%d", block, tx, p.Statement.LogIndex)
			p.EventAmount = eventAmounts[eaKey]
			ch <- p
		}
	}()
	return ch
}

func (r *Reconciler) flushBlock(buffer []Posting, modelChan chan<- Posting, wg *sync.WaitGroup) {
	type key struct{ asset, holder string }
	correctingEntry := func(k key, reason string, onChain, currentBal int64, p *Posting) Posting {
		r.counterMu.Lock()
		r.correctionCounter++
		correction := Posting{
			EventAmount:      onChain - currentBal,
			BeginBalance:     currentBal,
			CorrectionIndex:  r.correctionCounter,
			CorrectionReason: reason,
		}
		r.counterMu.Unlock()
		correction.Statement = p.Statement
		correction.Statement.TentativeBalance = onChain
		correction.Statement.CheckpointBalance = onChain
		r.rowIndexMu.Lock()
		r.rowIndexCounter++
		correction.RowIndex = r.rowIndexCounter
		r.rowIndexMu.Unlock()
		r.runningBalances[fmt.Sprintf("%s|%s", k.asset, k.holder)] = onChain
		wg.Add(1)
		return correction
	}

	lastPostings := make(map[key]int)
	for i, p := range buffer {
		k := key{p.Statement.AssetAddress, p.Statement.AccountedFor}
		seenKey := fmt.Sprintf("%d|%s|%s", p.Statement.BlockNumber, k.asset, k.holder)

		if _, seen := r.seenBlocks[seenKey]; !seen {
			prevBlock := fmt.Sprintf("0x%x", p.Statement.BlockNumber-1)
			if onChain, ok := conn.GetBalanceAtToken(k.asset, k.holder, prevBlock); ok {
				r.mu.Lock()
				currentBal := r.runningBalances[fmt.Sprintf("%s|%s", k.asset, k.holder)]
				if onChain != currentBal {
					modelChan <- correctingEntry(k, "mis", onChain, currentBal, &p)
				}
				r.mu.Unlock()
			}
			r.seenBlocks[seenKey] = p.Statement.BlockNumber
		}

		r.mu.Lock()
		runningKey := fmt.Sprintf("%s|%s", k.asset, k.holder)
		p.BeginBalance = r.runningBalances[runningKey]
		p.Statement.TentativeBalance = p.BeginBalance + p.EventAmount
		r.runningBalances[runningKey] = p.Statement.TentativeBalance
		r.rowIndexMu.Lock()
		r.rowIndexCounter++
		p.RowIndex = r.rowIndexCounter
		r.rowIndexMu.Unlock()
		r.mu.Unlock()

		buffer[i] = p
		lastPostings[k] = i
	}

	for _, p := range buffer {
		wg.Add(1)
		modelChan <- p
	}

	for k, idx := range lastPostings {
		p := buffer[idx]
		hexBlock := fmt.Sprintf("0x%x", p.Statement.BlockNumber)
		if onChain, ok := conn.GetBalanceAtToken(k.asset, k.holder, hexBlock); ok {
			r.mu.Lock()
			currentBal := r.runningBalances[fmt.Sprintf("%s|%s", k.asset, k.holder)]
			if onChain != currentBal {
				modelChan <- correctingEntry(k, "imb", onChain, currentBal, &p)
			}
			r.mu.Unlock()
		}
	}
}

func (r *Reconciler) processStream(modelChan chan<- Posting, wg *sync.WaitGroup) {
	globalStream := make(chan Posting)
	go func() {
		defer close(globalStream)
		var prevBlock int
		for _, app := range apps {
			block, tx := app[0], app[1]
			if block != prevBlock && prevBlock != 0 {
				globalStream <- Posting{Statement: struct {
					BlockNumber       int
					TransactionIndex  int
					LogIndex          int
					AssetAddress      string
					AccountedFor      string
					CheckpointBalance int64
					TentativeBalance  int64
				}{BlockNumber: prevBlock, AssetAddress: "END_OF_BLOCK"}}
			}
			for p := range r.GetPostingChannel(block, tx) {
				globalStream <- p
			}
			prevBlock = block
		}
		if prevBlock != 0 {
			globalStream <- Posting{Statement: struct {
				BlockNumber       int
				TransactionIndex  int
				LogIndex          int
				AssetAddress      string
				AccountedFor      string
				CheckpointBalance int64
				TentativeBalance  int64
			}{BlockNumber: prevBlock, AssetAddress: "END_OF_BLOCK"}}
		}
		globalStream <- Posting{Statement: struct {
			BlockNumber       int
			TransactionIndex  int
			LogIndex          int
			AssetAddress      string
			AccountedFor      string
			CheckpointBalance int64
			TentativeBalance  int64
		}{AssetAddress: "END_OF_STREAM"}}
	}()

	var buffer []Posting
	for posting := range globalStream {
		switch posting.Statement.AssetAddress {
		case "END_OF_BLOCK":
			r.flushBlock(buffer, modelChan, wg)
			buffer = nil
		case "END_OF_STREAM":
			if len(buffer) > 0 {
				r.flushBlock(buffer, modelChan, wg)
			}
			return
		default:
			buffer = append(buffer, posting)
		}
	}
}

func shortenAddress(addr string) string {
	if len(addr) > 8 {
		return addr[:8]
	}
	return addr
}

func main() {
	initData()
	modelChan := make(chan Posting, 1000)
	var wg sync.WaitGroup
	r := &Reconciler{
		addressOfInterest: "0xf",
		runningBalances:   make(map[string]int64),
		seenBlocks:        make(map[string]int),
	}

	done := make(chan struct{})

	go func() {
		defer close(modelChan)
		r.processStream(modelChan, &wg)
		close(done)
	}()

	var postings []Posting
	for p := range modelChan {
		postings = append(postings, p)
		wg.Done()
	}

	<-done
	wg.Wait()

	sort.Slice(postings, func(i, j int) bool {
		if postings[i].Statement.AssetAddress == postings[j].Statement.AssetAddress {
			return postings[i].RowIndex < postings[j].RowIndex
		}
		return postings[i].Statement.AssetAddress < postings[j].Statement.AssetAddress
	})

	fmt.Println("Block\tTx\tLog\tCorrection\tReason\tAsset\tHolder\tBeginBal\tAmount\tTentative\tCheckpoint\tCheck1\tCheck2\tReconciled")
	currentAsset := ""
	var lastTentative int64
	for _, p := range postings {
		isFinal := false
		if p.CorrectionIndex == 0 {
			lastTx := 0
			for _, app := range apps {
				if app[0] == p.Statement.BlockNumber && app[1] > lastTx {
					lastTx = app[1]
				}
			}
			if p.Statement.TransactionIndex == lastTx {
				lastLog := 0
				for _, logP := range logsByTx[mapKey(p.Statement.BlockNumber, lastTx, 0)] {
					if logP.Statement.LogIndex > lastLog {
						lastLog = logP.Statement.LogIndex
					}
				}
				if p.Statement.LogIndex == lastLog {
					isFinal = true
				}
			}
		}

		assetShort := shortenAddress(p.Statement.AssetAddress)
		holderShort := shortenAddress(p.Statement.AccountedFor)

		checkpoint := "-"
		if reconciled := p.Reconciled(isFinal); reconciled != "unknown" {
			checkpoint = fmt.Sprintf("%d", p.Statement.CheckpointBalance)
		}

		check1 := p.BeginBalance + p.EventAmount - p.Statement.TentativeBalance
		check2 := "-"
		if p.CorrectionIndex != 0 || isFinal {
			check2 = fmt.Sprintf("%d", p.BeginBalance+p.EventAmount-p.Statement.CheckpointBalance)
		}

		corrIndexStr := "-"
		if p.CorrectionIndex != 0 {
			corrIndexStr = fmt.Sprintf("%d", p.CorrectionIndex)
		}

		if currentAsset != "" && currentAsset != p.Statement.AssetAddress {
			r.rowIndexMu.Lock()
			r.rowIndexCounter++
			r.rowIndexMu.Unlock()
			fmt.Println(strings.Repeat("-", 120))
			fmt.Printf("-\t-\t-\t-\t-\t%s\t-\t0\t0\t%d\t-\t0\t-\t-\n",
				shortenAddress(currentAsset),
				lastTentative,
			)
			fmt.Println()
		}

		fmt.Printf("%d\t%d\t%d\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%s\t%d\t%s\t%s\n",
			p.Statement.BlockNumber,
			p.Statement.TransactionIndex,
			p.Statement.LogIndex,
			corrIndexStr,
			p.CorrectionReason,
			assetShort,
			holderShort,
			p.BeginBalance,
			p.EventAmount,
			p.Statement.TentativeBalance,
			checkpoint,
			check1,
			check2,
			p.Reconciled(isFinal),
		)

		currentAsset = p.Statement.AssetAddress
		lastTentative = p.Statement.TentativeBalance
	}

	if currentAsset != "" {
		r.rowIndexMu.Lock()
		r.rowIndexCounter++
		r.rowIndexMu.Unlock()
		fmt.Println(strings.Repeat("-", 120))
		fmt.Printf("-\t-\t-\t-\t-\t%s\t-\t0\t0\t%d\t-\t0\t-\t-\n",
			shortenAddress(currentAsset),
			lastTentative,
		)
		fmt.Println()
	}

	if len(postings) > 0 {
		lastPosting := postings[len(postings)-1]
		fmt.Println(strings.Repeat("=", 120))
		fmt.Printf("-\t-\t-\t-\t-\tTotal\t-\t0\t0\t%d\t-\t0\t-\t-\n",
			lastPosting.Statement.TentativeBalance,
		)
		fmt.Println()
	}
}
