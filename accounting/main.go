package main

import (
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/ledger3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "0xf")
	}

	r := ledger3.NewReconciler3("mainnet", base.HexToAddress(os.Args[1]))

	modelChan := make(chan types.Modeler, 1000)
	go func() {
		defer close(modelChan)
		r.ProcessStream(modelChan)
	}()

	extraOpts := map[string]any{
		"accounting": true,
	}

	types.PrintHeader()
	for p := range modelChan {
		p.Model("mainnet", "text", false, extraOpts)
	}
}
