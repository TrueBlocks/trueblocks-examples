package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/ledger3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
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
		r.ProcessStream(apps, modelChan)
	}()

	extraOpts := map[string]any{
		"accounting": true,
	}

	types.PrintHeader()
	for p := range modelChan {
		p.Model("mainnet", "text", false, extraOpts)
	}
}

var apps []types.Appearance

func init() {
	folder := os.Getenv("FOLDER")
	if folder == "" {
		folder = "tests"
	}
	// blockNumber,transactionIndex
	appsFn := filepath.Join(folder, "apps.csv")
	appsFile, _ := os.Open(appsFn)
	defer appsFile.Close()
	appsReader := csv.NewReader(appsFile)
	appsReader.Comment = '#'
	if appsRecords, err := appsReader.ReadAll(); err != nil {
		fmt.Println("Problem with data file:", appsFn)
		logger.Fatal(err)
	} else if len(appsRecords) == 0 {
		logger.Fatal("no transfers")
	} else {
		for _, record := range appsRecords[1:] {
			if strings.HasPrefix(record[0], "#") {
				continue
			}
			apps = append(apps, types.Appearance{
				BlockNumber:      uint32(base.MustParseInt64(record[0])),
				TransactionIndex: uint32(base.MustParseInt64(record[1])),
			})
		}
	}
}
