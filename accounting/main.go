package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

// ---------------------------------------------------------
type LedgerEntry struct {
	Holder           base.Address
	StatementId      int
	CorrectionId     int
	CorrectionReason string
	Opening          base.Wei
	Amount           base.Wei
	Verified         base.Wei
	AssetAddress     base.Address `json:"assetAddress"`
	BlockNumber      base.Blknum  `json:"blockNumber"`
	LogIndex         base.Lognum  `json:"logIndex"`
	TransactionIndex base.Txnum   `json:"transactionIndex"`
}

// ---------------------------------------------------------
func (p LedgerEntry) Calculated() base.Wei {
	return *new(base.Wei).Add(&p.Opening, &p.Amount)
}

// ---------------------------------------------------------
func (p LedgerEntry) Reconciled() (base.Wei, base.Wei, bool, bool) {
	calc := p.Calculated()
	checkVal := *new(base.Wei).Add(&p.Opening, &p.Amount)
	tentativeDiff := *new(base.Wei).Sub(&checkVal, &calc)
	checkpointDiff := *new(base.Wei).Sub(&checkVal, &p.Verified)

	checkpointEqual := checkVal.Equal(&p.Verified)
	if checkpointEqual {
		return tentativeDiff, checkpointDiff, true, true
	}

	tentativeEqual := checkVal.Equal(&calc)
	return tentativeDiff, checkpointDiff, tentativeEqual, false
}

// ---------------------------------------------------------
func PrintHeader() {
	// fmt.Println("Asset\tHolder\tBlock\tTx\tLog\tRow\tCorr\tReason\tOpening\tAmount\tCalculated\tVerified\tCheck1\tCheck2\tRec\tCp")
	fmt.Println("Asset\tHolder\tBlock\tTx\tLog\tRow\tCorr\tReason\tBegBal\tAmount\tTenBal\tChkBal\tCheck1\tCheck2\tRec\tCp")
}

// ---------------------------------------------------------
func (p *LedgerEntry) Model(chain, format string, verbose bool, extraOpts map[string]any) types.Model {
	_, _, _, _ = chain, format, verbose, extraOpts
	check1, check2, reconciles, byCheckpoint := p.Reconciled()
	calc := p.Calculated()
	fmt.Printf("%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%t\t%t\n",
		p.AssetAddress.Display(0, 1),
		p.Holder.Display(0, 1),
		p.BlockNumber,
		p.TransactionIndex,
		p.LogIndex,
		p.StatementId,
		p.CorrectionId,
		p.CorrectionReason,
		p.Opening.Text(10),
		p.Amount.Text(10),
		calc.Text(10),
		p.Verified.Text(10),
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
	accountLedger     map[assetHolderKey]base.Wei
	transfers         map[blockTxKey][]LedgerEntry
	correctionCounter int
	entryCounter      int
	hasStartBlock     bool
	ledgerAssets      map[base.Address]bool
}

// ---------------------------------------------------------
func NewReconciler(addr base.Address) *Reconciler {
	r := &Reconciler{
		account:       addr,
		accountLedger: make(map[assetHolderKey]base.Wei),
		transfers:     make(map[blockTxKey][]LedgerEntry),
		conn:          NewConnection(),
		hasStartBlock: false,
		ledgerAssets:  make(map[base.Address]bool),
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
func (r *Reconciler) getPostingChannel(app *types.Appearance) <-chan LedgerEntry {
	ch := make(chan LedgerEntry)
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
func (r *Reconciler) correctingEntry(reason string, onChain, currentBal base.Wei, p *LedgerEntry) LedgerEntry {
	correction := *p
	correction.Amount = *new(base.Wei).Sub(&onChain, &currentBal)
	correction.Opening = currentBal
	correction.Verified = onChain
	correction.CorrectionReason = reason
	return correction
}

// ---------------------------------------------------------
func (r *Reconciler) flushBlock(postings []LedgerEntry, modelChan chan<- types.Modeler) {
	blockProcessedAssets := make(map[base.Address]bool)
	assetLastSeen := make(map[base.Address]int)
	for i, p := range postings {
		key := assetHolderKey{Asset: p.AssetAddress, Holder: p.Holder}
		if !blockProcessedAssets[p.AssetAddress] {
			if r.hasStartBlock && !r.ledgerAssets[p.AssetAddress] {
				if onChain, ok := r.conn.GetBalanceAtToken(p.AssetAddress, p.Holder, p.BlockNumber-1); ok {
					r.accountLedger[key] = onChain
				}
				r.ledgerAssets[p.AssetAddress] = true
			}
			if onChain, ok := r.conn.GetBalanceAtToken(p.AssetAddress, p.Holder, p.BlockNumber-1); ok {
				currentBal := r.accountLedger[key]
				if !onChain.Equal(&currentBal) {
					correctingEntry := r.correctingEntry("mis", onChain, currentBal, &p)
					r.correctionCounter++
					correctingEntry.CorrectionId = r.correctionCounter
					r.entryCounter++
					correctingEntry.StatementId = r.entryCounter
					modelChan <- &correctingEntry
				}
			}
			blockProcessedAssets[p.AssetAddress] = true
		}

		p.Opening = r.accountLedger[key]
		r.accountLedger[key] = p.Calculated()
		r.entryCounter++
		p.StatementId = r.entryCounter
		postings[i] = p
		assetLastSeen[p.AssetAddress] = i
		modelChan <- &p
	}

	var corrections []LedgerEntry
	for _, idx := range assetLastSeen {
		p := postings[idx]
		key := assetHolderKey{Asset: p.AssetAddress, Holder: p.Holder}
		currentBal := r.accountLedger[key]
		if !p.Verified.Equal(&currentBal) {
			correctingEntry := r.correctingEntry("imb", p.Verified, currentBal, &p)
			corrections = append(corrections, correctingEntry)
		}
	}

	sort.SliceStable(corrections, func(i, j int) bool {
		if corrections[i].TransactionIndex != corrections[j].TransactionIndex {
			return corrections[i].TransactionIndex < corrections[j].TransactionIndex
		}
		if corrections[i].LogIndex != corrections[j].LogIndex {
			return corrections[i].LogIndex < corrections[j].LogIndex
		}
		return corrections[i].AssetAddress.Hex() < corrections[j].AssetAddress.Hex()
	})

	for _, correction := range corrections {
		r.correctionCounter++
		correction.CorrectionId = r.correctionCounter
		r.entryCounter++
		correction.StatementId = r.entryCounter
		modelChan <- &correction
		key := assetHolderKey{Asset: correction.AssetAddress, Holder: correction.Holder}
		r.accountLedger[key] = correction.Verified
	}
}

// ---------------------------------------------------------
func (r *Reconciler) processStream(modelChan chan<- types.Modeler) {
	postingStream := make(chan LedgerEntry, 100)
	go func() {
		defer close(postingStream)
		var prevBlock base.Blknum
		for _, app := range apps {
			bn := base.Blknum(app.BlockNumber)
			if bn != prevBlock && prevBlock != 0 {
				postingStream <- LedgerEntry{
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
			postingStream <- LedgerEntry{
				BlockNumber:  prevBlock,
				AssetAddress: EndOfBlockSentinel,
			}
		}
		postingStream <- LedgerEntry{
			AssetAddress: EndOfStreamSentinel,
		}
	}()

	var postings []LedgerEntry
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
		p := LedgerEntry{
			BlockNumber:      base.Blknum(base.MustParseUint64(record[0])),
			TransactionIndex: base.Txnum(base.MustParseUint64(record[1])),
			LogIndex:         base.Lognum(base.MustParseUint64(record[2])),
			AssetAddress:     base.HexToAddress(record[3]),
			Holder:           base.HexToAddress(record[4]),
			Amount:           *base.NewWeiStr(record[5]),
		}
		p.Verified, _ = r.conn.GetBalanceAtToken(p.AssetAddress, p.Holder, p.BlockNumber)

		key := blockTxKey{BlockNumber: p.BlockNumber, TransactionIndex: p.TransactionIndex}
		r.transfers[key] = append(r.transfers[key], p)
	}

	if firstBlock := os.Getenv("FIRST_BLOCK"); firstBlock != "" {
		r.hasStartBlock = true
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
