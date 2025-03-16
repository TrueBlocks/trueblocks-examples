package main

import (
	"fmt"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/ledger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/ledger3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/walk"
)

func main() {
	if len(os.Args) < 2 {
		logger.Fatal("usage: accounting <address>")
	}

	var updater = monitor.NewUpdater("mainnet", true, false, os.Args[1:])
	monitorArray := make([]monitor.Monitor, 0, len(os.Args[1:]))
	if canceled, err := updater.FreshenMonitors(&monitorArray); err != nil || canceled {
		logger.Fatal(err)
	} else if len(monitorArray) == 0 {
		logger.Fatal("no monitors")
	}

	mon := monitorArray[0]
	conn := rpc.NewConnection("mainnet", false, map[walk.CacheType]bool{})
	ledgerOpts := &ledger.ReconcilerOptions{
		AccountFor:   mon.Address,
		FirstBlock:   0,
		LastBlock:    base.Blknum(base.NOPOS),
		AsEther:      false,
		UseTraces:    false,
		Reversed:     false,
		AssetFilters: []base.Address{},
	}
	r := ledger3.NewReconciler(conn, ledgerOpts)
	r.InitData()

	contents := file.AsciiFileToLines("transfers.csv")
	if len(contents) > 0 {
		fmt.Println(contents[0])
	}

	modelChan := make(chan types.Modeler, 1000)
	go func() {
		defer close(modelChan)
		filter := types.NewFilter(
			false,
			false,
			[]string{},
			base.BlockRange{First: 0, Last: 20000000},
			base.RecordRange{First: 0, Last: base.NOPOS},
		)
		if apps, cnt, err := mon.ReadAndFilterAppearances(filter, false /* withCount */); err != nil {
			logger.Fatal(err)

		} else if cnt == 0 {
			logger.Warn("no blocks found for the query")

		} else {
			r.ProcessStream(apps, modelChan)
		}
	}()

	PrintHeader()
	for p := range modelChan {
		p.Model("mainnet", "text", false, map[string]any{})
	}
}

func PrintHeader() {
	fmt.Println("asset\tholder\tblockNumber\ttransactionIndex\tlogIndex\trowIndex\tcorrectionIndex\tcorrectionReason\tbegBal\tamountNet\tendBalCalc\tendBal\tcheck1\tcheck2\treconciled\tcheckpoint")
}
