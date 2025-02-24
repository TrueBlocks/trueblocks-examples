package main

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

// Posting represents a single ledger event or correction
type Posting struct {
	types.Statement
	Holder            base.Address
	StatementId       int
	CorrectionIndex   int
	CorrectionReason  string
	BeginBalance      int64
	EventAmount       int64
	TentativeBalance  int64
	CheckpointBalance int64
}

func (p Posting) Reconciled() (bool, bool) {
	if p.BeginBalance+p.EventAmount == p.CheckpointBalance {
		return true, true
	}
	return p.BeginBalance+p.EventAmount == p.TentativeBalance, false
}

func printHeader() {
	fmt.Println("Asset\tHolder\tBlock\tTx\tLog\tRow\tCorr\tReason\tBegBal\tAmount\tTenBal\tChkBal\tCheck1\tCheck2\tRec\tCp")
}

func (p *Posting) printStatement() {
	p.Holder = base.HexToAddress("0xf")
	reconciles, byCheckpoint := p.Reconciled()
	fmt.Printf("%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%d\t%d\t%d\t%d\t%d\t%d\t%t\t%t\n",
		p.Statement.AssetAddress.Display(0, 1),
		p.Holder.Display(0, 1),
		p.Statement.BlockNumber,
		p.Statement.TransactionIndex,
		p.Statement.LogIndex,
		p.StatementId,
		p.CorrectionIndex,
		p.CorrectionReason,
		p.BeginBalance,
		p.EventAmount,
		p.TentativeBalance,
		p.CheckpointBalance,
		p.BeginBalance+p.EventAmount-p.TentativeBalance,
		p.BeginBalance+p.EventAmount-p.CheckpointBalance,
		reconciles,
		byCheckpoint,
	)
}
