package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type Reconciler struct {
	mu                sync.Mutex
	addressOfInterest base.Address
	runningBalMap     map[string]int64
	seenBlockMap      map[string]base.Blknum
	balanceMap        map[string]int64
	logsMap           map[MapKey][]types.Posting
	lastPostingsMap   map[key]int
	correctionCounter int
	counterMu         sync.Mutex
	statementIndex    int
	rowIndexMu        sync.Mutex
}

func NewReconciler(addr base.Address) *Reconciler {
	r := &Reconciler{
		addressOfInterest: addr,
		runningBalMap:     make(map[string]int64),
		seenBlockMap:      make(map[string]base.Blknum),
		lastPostingsMap:   make(map[key]int),
		balanceMap:        make(map[string]int64),
		logsMap:           make(map[MapKey][]types.Posting),
	}
	r.initData()
	return r
}

var (
	apps                []types.Appearance
	EndOfBlockSentinel  = base.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	EndOfStreamSentinel = base.HexToAddress("0xbeefdeadbeefdeadbeefdeadbeefdeadbeefdead")
)

type key struct {
	asset  base.Address
	holder base.Address
}

func (r *Reconciler) GetPostingChannel(app *types.Appearance) <-chan types.Posting {
	ch := make(chan types.Posting)
	go func() {
		defer close(ch)
		logKey := logsKey(base.Blknum(app.BlockNumber), base.Txnum(app.TransactionIndex), base.Lognum(0))
		for _, p := range r.logsMap[logKey] {
			if p.Holder != r.addressOfInterest {
				continue
			}
			ch <- p
		}
	}()
	return ch
}

func (r *Reconciler) flushBlock(postings []types.Posting, modelChan chan<- types.Posting) {
	correctingEntry := func(k key, reason string, onChain, currentBal int64, p *types.Posting) types.Posting {
		r.counterMu.Lock()
		r.correctionCounter++
		correction := types.Posting{
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
		r.statementIndex++
		correction.StatementId = r.statementIndex
		r.rowIndexMu.Unlock()
		r.runningBalMap[fmt.Sprintf("%s|%s", k.asset, k.holder)] = onChain
		return correction
	}

	r.lastPostingsMap = make(map[key]int)
	for i, p := range postings {
		k := key{p.Statement.AssetAddress, p.Holder}
		seenKey := fmt.Sprintf("%d|%s|%s", p.Statement.BlockNumber, k.asset, k.holder)

		if _, seen := r.seenBlockMap[seenKey]; !seen {
			if onChain, ok := r.GetBalanceAtToken(k.asset, k.holder, p.Statement.BlockNumber-1); ok {
				r.mu.Lock()
				currentBal := r.runningBalMap[fmt.Sprintf("%s|%s", k.asset, k.holder)]
				if onChain != currentBal {
					modelChan <- correctingEntry(k, "mis", onChain, currentBal, &p)
				}
				r.mu.Unlock()
			}
			r.seenBlockMap[seenKey] = p.Statement.BlockNumber
		}

		r.mu.Lock()
		p.BeginBalance = r.runningBalMap[fmt.Sprintf("%s|%s", k.asset, k.holder)]
		p.TentativeBalance = p.BeginBalance + p.EventAmount
		r.runningBalMap[fmt.Sprintf("%s|%s", k.asset, k.holder)] = p.TentativeBalance
		r.rowIndexMu.Lock()
		r.statementIndex++
		p.StatementId = r.statementIndex
		r.rowIndexMu.Unlock()
		r.mu.Unlock()

		postings[i] = p
		r.lastPostingsMap[k] = i
	}

	for _, p := range postings {
		modelChan <- p
	}

	for k, idx := range r.lastPostingsMap {
		p := postings[idx]
		if onChain, ok := r.GetBalanceAtToken(k.asset, k.holder, p.Statement.BlockNumber); ok {
			r.mu.Lock()
			currentBal := r.runningBalMap[fmt.Sprintf("%s|%s", k.asset, k.holder)]
			if onChain != currentBal {
				modelChan <- correctingEntry(k, "imb", onChain, currentBal, &p)
			}
			r.mu.Unlock()
		}
	}
}

func (r *Reconciler) processStream(modelChan chan<- types.Posting) {
	globalStream := make(chan types.Posting)
	go func() {
		defer close(globalStream)
		var prevBlock base.Blknum
		for _, app := range apps {
			if base.Blknum(app.BlockNumber) != prevBlock && prevBlock != 0 {
				globalStream <- types.Posting{Statement: types.Statement{
					BlockNumber:  base.Blknum(prevBlock),
					AssetAddress: EndOfBlockSentinel,
				}}
			}
			for p := range r.GetPostingChannel(&app) {
				globalStream <- p
			}
			prevBlock = base.Blknum(app.BlockNumber)
		}
		if prevBlock != 0 {
			globalStream <- types.Posting{Statement: types.Statement{
				BlockNumber:  base.Blknum(prevBlock),
				AssetAddress: EndOfBlockSentinel,
			}}
		}
		globalStream <- types.Posting{Statement: types.Statement{
			AssetAddress: EndOfStreamSentinel,
		}}
	}()

	var postings []types.Posting
	for posting := range globalStream {
		switch posting.Statement.AssetAddress {
		case EndOfBlockSentinel:
			r.flushBlock(postings, modelChan)
			postings = nil
		case EndOfStreamSentinel:
			if len(postings) > 0 {
				r.flushBlock(postings, modelChan)
			}
			return
		default:
			postings = append(postings, posting)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: chifra <address>")
		os.Exit(1)
	}

	r := NewReconciler(base.HexToAddress(os.Args[1]))

	modelChan := make(chan types.Posting, 1000)
	go func() {
		defer close(modelChan)
		r.processStream(modelChan)
	}()

	types.PrintHeader()
	for p := range modelChan {
		p.PrintStatement()
	}
}

func (r *Reconciler) initData() {
	r.balanceMap = make(map[string]int64)
	r.logsMap = make(map[MapKey][]types.Posting)

	// blockNumber,transactionIndex
	appsFile, _ := os.Open("tests/apps.csv")
	defer appsFile.Close()
	appsReader := csv.NewReader(appsFile)
	appsRecords, _ := appsReader.ReadAll()
	for _, record := range appsRecords[1:] {
		block, _ := strconv.Atoi(record[0])
		tx, _ := strconv.Atoi(record[1])
		apps = append(apps, types.Appearance{
			BlockNumber:      uint32(block),
			TransactionIndex: uint32(tx),
		})
	}

	// blockNumber,transactionIndex,logIndex,assetAddress,accountedFor,amountNet,endBal
	logsFile, _ := os.Open("tests/logs.csv")
	defer logsFile.Close()
	logsReader := csv.NewReader(logsFile)
	logsRecords, _ := logsReader.ReadAll()
	for _, record := range logsRecords[1:] {
		block, _ := strconv.Atoi(record[0])
		tx, _ := strconv.Atoi(record[1])
		log, _ := strconv.Atoi(record[2])
		p := types.Posting{}
		p.Statement.BlockNumber = base.Blknum(block)
		p.Statement.TransactionIndex = base.Txnum(tx)
		p.Statement.LogIndex = base.Lognum(log)
		p.Statement.AssetAddress = base.HexToAddress(record[3])
		p.Holder = base.HexToAddress(record[4])
		p.EventAmount, _ = strconv.ParseInt(record[5], 10, 64)
		p.CheckpointBalance, _ = strconv.ParseInt(record[6], 10, 64)
		logKey := logsKey(base.Blknum(block), base.Txnum(tx), base.Lognum(0))
		r.logsMap[logKey] = append(r.logsMap[logKey], p)
	}

	// blockNumber,assetAddress,accountedFor,endBal
	mapFile, _ := os.Open("tests/balances.csv")
	defer mapFile.Close()
	mapReader := csv.NewReader(mapFile)
	mapRecords, _ := mapReader.ReadAll()
	for _, record := range mapRecords[1:] {
		asset := base.HexToAddress(record[1])
		holder := base.HexToAddress(record[2])
		key := fmt.Sprintf("%s|%s|%s", record[0], asset.Hex(), holder.Hex())
		bal, _ := strconv.ParseInt(record[3], 10, 64)
		r.balanceMap[key] = bal
	}
	// logger.Info("Data initialized", len(apps), len(logsMap), len(balanceMap))
}

var (
	cc Connection
)

// Connection provides on-chain balance lookups
type Connection struct{}

func (r *Reconciler) GetBalanceAtToken(asset base.Address, holder base.Address, bn base.Blknum) (int64, bool) {
	key := fmt.Sprintf("%d|%s|%s", bn, asset.Hex(), holder.Hex())
	if bal, ok := r.balanceMap[key]; ok {
		return bal, true
	}
	return 0, false
}

type MapKey struct {
	BlockNumber      base.Blknum
	TransactionIndex base.Txnum
	LogIndex         base.Lognum
	Asset            base.Address
	Holder           base.Address
}

func runningBalKey(asset base.Address, holder base.Address) MapKey {
	return MapKey{Asset: asset, Holder: holder}
}

func seenBlockKey(bn base.Blknum, asset base.Address, holder base.Address) MapKey {
	return MapKey{BlockNumber: bn, Asset: asset, Holder: holder}
}

func lastPostingsKey(asset base.Address, holder base.Address) MapKey {
	return MapKey{Asset: asset, Holder: holder}
}

func balanceKey(bn base.Blknum, asset base.Address, holder base.Address) MapKey {
	return MapKey{BlockNumber: bn, Asset: asset, Holder: holder}
}

func logsKey(bn base.Blknum, tx base.Txnum, log base.Lognum) MapKey {
	return MapKey{BlockNumber: bn, TransactionIndex: tx, LogIndex: log}
}
