package main

import (
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/filter"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/ledger10"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/ledger3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/walk"
)

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "0xf")
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
	ledgerOpts := &ledger10.ReconcilerOptions{
		Connection:   conn,
		AccountFor:   mon.Address,
		FirstBlock:   0,
		LastBlock:    base.Blknum(base.NOPOS),
		AsEther:      false,
		TestMode:     false,
		UseTraces:    false,
		Reversed:     false,
		AssetFilters: []base.Address{},
	}
	r := ledger3.NewReconciler(ledgerOpts)
	r.InitData()

	modelChan := make(chan types.Modeler, 1000)
	go func() {
		defer close(modelChan)
		filter := filter.NewFilter(
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

	extraOpts := map[string]any{
		"accounting": true,
	}

	types.PrintHeader()
	for p := range modelChan {
		p.Model("mainnet", "text", false, extraOpts)
	}
}
