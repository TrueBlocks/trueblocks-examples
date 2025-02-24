package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

var (
	balanceMap map[string]int64
	logsByTx   map[int][]types.Posting
)

func initData() {
	balanceMap = make(map[string]int64)
	logsByTx = make(map[int][]types.Posting)

	appsFile, _ := os.Open("tests/apps.csv")
	defer appsFile.Close()
	appsReader := csv.NewReader(appsFile)
	appsRecords, _ := appsReader.ReadAll()
	for _, record := range appsRecords[1:] {
		block, _ := strconv.Atoi(record[0])
		tx, _ := strconv.Atoi(record[1])
		apps = append(apps, [2]int{block, tx})
	}

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
		p.CheckpointBalance, _ = strconv.ParseInt(record[5], 10, 64)
		p.EventAmount, _ = strconv.ParseInt(record[6], 10, 64)
		key := mapKey(block, tx, 0)
		logsByTx[key] = append(logsByTx[key], p)
	}

	mapFile, _ := os.Open("tests/balances.csv")
	defer mapFile.Close()
	mapReader := csv.NewReader(mapFile)
	mapRecords, _ := mapReader.ReadAll()
	for _, record := range mapRecords[1:] {
		asset := base.HexToAddress(record[1])
		holder := base.HexToAddress(record[2])
		key := fmt.Sprintf("%s|%s|%s", record[0], asset.Hex(), holder.Hex())
		bal, _ := strconv.ParseInt(record[3], 10, 64)
		balanceMap[key] = bal
	}
}

func init() {
	initData()
}
