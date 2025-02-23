package main

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type Reconciler struct {
	mu                sync.Mutex
	addressOfInterest base.Address
	runningBalances   map[string]int64
	seenBlocks        map[string]base.Blknum
	correctionCounter int
	counterMu         sync.Mutex
	rowIndexCounter   int
	rowIndexMu        sync.Mutex
}

var (
	apps                [][2]int
	EndOfBlockSentinel  = base.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	EndOfStreamSentinel = base.HexToAddress("0xbeefdeadbeefdeadbeefdeadbeefdeadbeefdead")
)

func (r *Reconciler) GetPostingChannel(block, tx int) <-chan Posting {
	ch := make(chan Posting)
	go func() {
		defer close(ch)
		for _, p := range logsByTx[mapKey(block, tx, 0)] {
			if p.Statement.Holder != r.addressOfInterest {
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
	type key struct {
		asset  base.Address
		holder base.Address
	}
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
		correction.TentativeBalance = onChain
		correction.CheckpointBalance = onChain
		r.rowIndexMu.Lock()
		r.rowIndexCounter++
		correction.RowIndex = r.rowIndexCounter
		r.rowIndexMu.Unlock()
		r.runningBalances[fmt.Sprintf("%s|%s", k.asset, k.holder.Hex())] = onChain
		wg.Add(1)
		return correction
	}

	lastPostings := make(map[key]int)
	for i, p := range buffer {
		k := key{p.Statement.AssetAddress, p.Statement.Holder}
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
		p.TentativeBalance = p.BeginBalance + p.EventAmount
		r.runningBalances[runningKey] = p.TentativeBalance
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
				globalStream <- Posting{Statement: types.Statement{
					BlockNumber:  base.Blknum(prevBlock),
					AssetAddress: EndOfBlockSentinel,
				}}
			}
			for p := range r.GetPostingChannel(block, tx) {
				globalStream <- p
			}
			prevBlock = block
		}
		if prevBlock != 0 {
			globalStream <- Posting{Statement: types.Statement{
				BlockNumber:  base.Blknum(prevBlock),
				AssetAddress: EndOfBlockSentinel,
			}}
		}
		globalStream <- Posting{Statement: types.Statement{
			AssetAddress: EndOfStreamSentinel,
		}}
	}()

	var buffer []Posting
	for posting := range globalStream {
		switch posting.Statement.AssetAddress {
		case EndOfBlockSentinel:
			r.flushBlock(buffer, modelChan, wg)
			buffer = nil
		case EndOfStreamSentinel:
			if len(buffer) > 0 {
				r.flushBlock(buffer, modelChan, wg)
			}
			return
		default:
			buffer = append(buffer, posting)
		}
	}
}

func shortenAddress(address base.Address) string {
	addr := address.Hex()
	if len(addr) == 42 {
		return addr[:2] + addr[len(addr)-1:]
	}
	return addr
}

func main() {
	initData()
	modelChan := make(chan Posting, 1000)
	var wg sync.WaitGroup
	r := &Reconciler{
		addressOfInterest: base.HexToAddress("0xf"),
		runningBalances:   make(map[string]int64),
		seenBlocks:        make(map[string]base.Blknum),
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

	printHeader()
	for _, p := range postings {
		p.printStatement()
	}
}
