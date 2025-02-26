package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

// ---------------------------------------------------------
type Posting2 struct {
	Holder            base.Address
	StatementId       int
	CorrectionIndex   int
	CorrectionReason  string
	BeginBalance      base.Wei
	EventAmount       base.Wei
	TentativeBalance  base.Wei
	CheckpointBalance base.Wei
	AssetAddress      base.Address `json:"assetAddress"`
	BlockNumber       base.Blknum  `json:"blockNumber"`
	LogIndex          base.Lognum  `json:"logIndex"`
	TransactionIndex  base.Txnum   `json:"transactionIndex"`
}

// ---------------------------------------------------------
func (p Posting2) Reconciled() (base.Wei, base.Wei, bool, bool) {
	checkVal := *new(base.Wei).Add(&p.BeginBalance, &p.EventAmount)
	tentativeDiff := *new(base.Wei).Sub(&checkVal, &p.TentativeBalance)
	checkpointDiff := *new(base.Wei).Sub(&checkVal, &p.CheckpointBalance)
	tentativeEqual := checkVal.Equal(&p.TentativeBalance)
	checkpointEqual := checkVal.Equal(&p.CheckpointBalance)
	if checkpointEqual {
		return tentativeDiff, checkpointDiff, true, true
	}
	return tentativeDiff, checkpointDiff, tentativeEqual, false
}

// ---------------------------------------------------------
func PrintHeader() {
	fmt.Println("Asset\tHolder\tBlock\tTx\tLog\tRow\tCorr\tReason\tBegBal\tAmount\tTenBal\tChkBal\tCheck1\tCheck2\tRec\tCp")
}

// ---------------------------------------------------------
func (p *Posting2) Model(chain, format string, verbose bool, extraOpts map[string]any) types.Model {
	_, _, _, _ = chain, format, verbose, extraOpts
	check1, check2, reconciles, byCheckpoint := p.Reconciled()
	fmt.Printf("%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%t\t%t\n",
		p.AssetAddress.Display(0, 1),
		p.Holder.Display(0, 1),
		p.BlockNumber,
		p.TransactionIndex,
		p.LogIndex,
		p.StatementId,
		p.CorrectionIndex,
		p.CorrectionReason,
		p.BeginBalance.Text(10),
		p.EventAmount.Text(10),
		p.TentativeBalance.Text(10),
		p.CheckpointBalance.Text(10),
		check1.Text(10),
		check2.Text(10),
		reconciles,
		byCheckpoint,
	)
	return types.Model{}
}

// ---------------------------------------------------------
type Balance2 struct {
	BlockNumber base.Blknum
	Asset       base.Address
	Holder      base.Address
	Balance     base.Wei
}

// ---------------------------------------------------------
type Reconciler struct {
	conn              *Connection
	account           base.Address
	runningBal        map[assetHolderKey]base.Wei
	transfers         map[blockTxKey][]Posting2
	correctionCounter int
	statementIndex    int
}

// ---------------------------------------------------------
func NewReconciler(addr base.Address) *Reconciler {
	r := &Reconciler{
		account:    addr,
		runningBal: make(map[assetHolderKey]base.Wei),
		transfers:  make(map[blockTxKey][]Posting2),
		conn:       NewConnection(),
	}

	r.initData()

	return r
}

// ---------------------------------------------------------
var (
	apps                []types.Appearance
	EndOfBlockSentinel  = base.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	EndOfStreamSentinel = base.HexToAddress("0xbeefdeadbeefdeadbeefdeadbeefdeadbeefdead")
)

// ---------------------------------------------------------
func (r *Reconciler) getPostingChannel(app *types.Appearance) <-chan Posting2 {
	ch := make(chan Posting2)
	go func() {
		defer close(ch)
		key := blockTxKey{BlockNumber: base.Blknum(app.BlockNumber), TransactionIndex: base.Txnum(app.TransactionIndex)}
		for _, p := range r.transfers[key] {
			if p.Holder == r.account {
				ch <- p
			}
		}
	}()
	return ch
}

// ---------------------------------------------------------
func (r *Reconciler) correctingEntry(reason string, onChain, currentBal base.Wei, p *Posting2) Posting2 {
	correction := *p
	correction.EventAmount = *new(base.Wei).Sub(&onChain, &currentBal)
	correction.BeginBalance = currentBal
	correction.TentativeBalance = onChain
	correction.CheckpointBalance = onChain

	r.correctionCounter++
	correction.CorrectionIndex = r.correctionCounter
	correction.CorrectionReason = reason

	r.statementIndex++
	correction.StatementId = r.statementIndex

	key := assetHolderKey{Asset: p.AssetAddress, Holder: p.Holder}
	r.runningBal[key] = onChain
	return correction
}

// ---------------------------------------------------------
func (r *Reconciler) flushBlock(postings []Posting2, modelChan chan<- types.Modeler) {
	assetSeen := make(map[base.Address]bool)
	assetLastSeen := make(map[base.Address]int)
	for i, p := range postings {
		key := assetHolderKey{Asset: p.AssetAddress, Holder: p.Holder}
		if !assetSeen[p.AssetAddress] {
			if onChain, ok := r.conn.GetBalanceAtToken(p.AssetAddress, p.Holder, p.BlockNumber-1); ok {
				currentBal := r.runningBal[key]
				if !onChain.Equal(&currentBal) {
					correctingEntry := r.correctingEntry("mis", onChain, currentBal, &p)
					modelChan <- &correctingEntry
				}
			}
			assetSeen[p.AssetAddress] = true
		}

		p.BeginBalance = r.runningBal[key]
		p.TentativeBalance = *new(base.Wei).Add(&p.BeginBalance, &p.EventAmount)
		r.runningBal[key] = p.TentativeBalance
		r.statementIndex++
		p.StatementId = r.statementIndex

		postings[i] = p
		assetLastSeen[p.AssetAddress] = i
		modelChan <- &p
	}

	for _, idx := range assetLastSeen {
		p := postings[idx]
		key := assetHolderKey{Asset: p.AssetAddress, Holder: p.Holder}
		currentBal := r.runningBal[key]
		if !p.CheckpointBalance.Equal(&currentBal) {
			correctingEntry := r.correctingEntry("imb", p.CheckpointBalance, currentBal, &p)
			modelChan <- &correctingEntry
		}
	}
}

// ---------------------------------------------------------
func (r *Reconciler) processStream(modelChan chan<- types.Modeler) {
	postingStream := make(chan Posting2, 100)
	go func() {
		defer close(postingStream)
		var prevBlock base.Blknum
		for _, app := range apps {
			bn := base.Blknum(app.BlockNumber)
			if bn != prevBlock && prevBlock != 0 {
				postingStream <- Posting2{
					BlockNumber:  prevBlock,
					AssetAddress: EndOfBlockSentinel,
				}
			}
			for p := range r.getPostingChannel(&app) {
				postingStream <- p
			}
			prevBlock = bn
		}
		if prevBlock != 0 {
			postingStream <- Posting2{
				BlockNumber:  prevBlock,
				AssetAddress: EndOfBlockSentinel,
			}
		}
		postingStream <- Posting2{
			AssetAddress: EndOfStreamSentinel,
		}
	}()

	var postings []Posting2
	for posting := range postingStream {
		switch posting.AssetAddress {
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

// ---------------------------------------------------------
func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "0xf")
	}

	r := NewReconciler(base.HexToAddress(os.Args[1]))

	modelChan := make(chan types.Modeler, 1000)
	go func() {
		defer close(modelChan)
		r.processStream(modelChan)
	}()

	PrintHeader()
	for p := range modelChan {
		p.Model("mainnet", "text", false, nil)
	}
}

// ---------------------------------------------------------
func (r *Reconciler) initData() {
	folder := os.Getenv("FOLDER")
	if folder == "" {
		folder = "tests"
	}
	// blockNumber,transactionIndex
	appsFn := filepath.Join(folder, "apps.csv")
	appsFile, _ := os.Open(appsFn)
	defer appsFile.Close()
	appsReader := csv.NewReader(appsFile)
	appsRecords, _ := appsReader.ReadAll()
	for _, record := range appsRecords[1:] {
		if strings.HasPrefix(record[0], "#") {
			continue
		}
		apps = append(apps, types.Appearance{
			BlockNumber:      uint32(base.MustParseInt64(record[0])),
			TransactionIndex: uint32(base.MustParseInt64(record[1])),
		})
	}

	// blockNumber,assetAddress,accountedFor,endBal
	balFn := filepath.Join(folder, "balances.csv")
	balFile, _ := os.Open(balFn)
	defer balFile.Close()
	balReader := csv.NewReader(balFile)
	balRecords, _ := balReader.ReadAll()
	for _, record := range balRecords[1:] {
		if strings.HasPrefix(record[0], "#") {
			continue
		}
		b := Balance2{
			BlockNumber: base.Blknum(base.MustParseUint64(record[0])),
			Asset:       base.HexToAddress(record[1]),
			Holder:      base.HexToAddress(record[2]),
			Balance:     *base.NewWeiStr(record[3]),
		}

		key := bnAssetHolderKey{BlockNumber: b.BlockNumber, Asset: b.Asset, Holder: b.Holder}
		r.conn.balanceMap[key] = b.Balance
	}

	// blockNumber,transactionIndex,logIndex,assetAddress,accountedFor,amountNet,endBal
	transfersFn := filepath.Join(folder, "transfers.csv")
	transfersFile, _ := os.Open(transfersFn)
	defer transfersFile.Close()
	transfersReader := csv.NewReader(transfersFile)
	transfersRecords, _ := transfersReader.ReadAll()
	for _, record := range transfersRecords[1:] {
		if strings.HasPrefix(record[0], "#") {
			continue
		}
		p := Posting2{
			BlockNumber:      base.Blknum(base.MustParseUint64(record[0])),
			TransactionIndex: base.Txnum(base.MustParseUint64(record[1])),
			LogIndex:         base.Lognum(base.MustParseUint64(record[2])),
			AssetAddress:     base.HexToAddress(record[3]),
			Holder:           base.HexToAddress(record[4]),
			EventAmount:      *base.NewWeiStr(record[5]),
		}
		p.CheckpointBalance, _ = r.conn.GetBalanceAtToken(p.AssetAddress, p.Holder, p.BlockNumber)

		key := blockTxKey{BlockNumber: p.BlockNumber, TransactionIndex: p.TransactionIndex}
		r.transfers[key] = append(r.transfers[key], p)
	}
}

// ---------------------------------------------------------
// Connection provides on-chain balance lookups
type Connection struct {
	balanceMap map[bnAssetHolderKey]base.Wei
}

// ---------------------------------------------------------
func NewConnection() *Connection {
	return &Connection{
		balanceMap: make(map[bnAssetHolderKey]base.Wei),
	}
}

// ---------------------------------------------------------
func (c *Connection) GetBalanceAtToken(asset base.Address, holder base.Address, bn base.Blknum) (base.Wei, bool) {
	key := bnAssetHolderKey{BlockNumber: bn, Asset: asset, Holder: holder}
	if bal, ok := c.balanceMap[key]; ok {
		return bal, true
	}
	return *base.ZeroWei, false
}

// ---------------------------------------------------------
type blockTxKey struct {
	BlockNumber      base.Blknum
	TransactionIndex base.Txnum
}

// ---------------------------------------------------------
type assetHolderKey struct {
	Asset  base.Address
	Holder base.Address
}

type bnAssetHolderKey struct {
	BlockNumber base.Blknum
	Asset       base.Address
	Holder      base.Address
}
