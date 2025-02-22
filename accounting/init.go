package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

var (
	eventAmounts map[string]int64
	mapping      map[string]int64
	logsByTx     map[int][]Posting
)

func initData() {
	eventAmounts = make(map[string]int64)
	mapping = make(map[string]int64)
	logsByTx = make(map[int][]Posting)

	appsFile, _ := os.Open("tests/apps.csv")
	defer appsFile.Close()
	appsReader := csv.NewReader(appsFile)
	appsRecords, _ := appsReader.ReadAll()
	for _, record := range appsRecords[1:] {
		block, _ := strconv.Atoi(record[0])
		tx, _ := strconv.Atoi(record[1])
		apps = append(apps, [2]int{block, tx})
	}

	eaFile, _ := os.Open("tests/eventAmounts.csv")
	defer eaFile.Close()
	eaReader := csv.NewReader(eaFile)
	eaRecords, _ := eaReader.ReadAll()
	for _, record := range eaRecords[1:] {
		key := fmt.Sprintf("%s|%s|%s", record[0], record[1], record[2])
		amount, _ := strconv.ParseInt(record[3], 10, 64)
		eventAmounts[key] = amount
	}

	logsFile, _ := os.Open("tests/logsByTransaction.csv")
	defer logsFile.Close()
	logsReader := csv.NewReader(logsFile)
	logsRecords, _ := logsReader.ReadAll()
	for _, record := range logsRecords[1:] {
		block, _ := strconv.Atoi(record[0])
		tx, _ := strconv.Atoi(record[1])
		log, _ := strconv.Atoi(record[2])
		p := Posting{}
		p.Statement.BlockNumber = base.Blknum(block)
		p.Statement.TransactionIndex = base.Txnum(tx)
		p.Statement.LogIndex = base.Lognum(log)
		p.Statement.AssetAddress = strings.ToLower(record[3])
		p.Statement.Holder = base.HexToAddress(record[4])
		p.CheckpointBalance, _ = strconv.ParseInt(record[5], 10, 64)
		key := mapKey(block, tx, 0)
		logsByTx[key] = append(logsByTx[key], p)
	}

	mapFile, _ := os.Open("tests/mapping.csv")
	defer mapFile.Close()
	mapReader := csv.NewReader(mapFile)
	mapRecords, _ := mapReader.ReadAll()
	for _, record := range mapRecords[1:] {
		holder := base.HexToAddress(record[2])
		key := fmt.Sprintf("%s|%s|%s", record[0], record[1], holder.Hex())
		bal, _ := strconv.ParseInt(record[3], 10, 64)
		mapping[key] = bal
	}
}
