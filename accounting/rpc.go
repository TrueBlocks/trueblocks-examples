package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	conn         Connection
	eventAmounts map[string]int64
	mapping      map[string]int64
	logsByTx     map[int][]Posting
)

// Connection provides on-chain balance lookups
type Connection struct{}

func (c *Connection) GetBalanceAtToken(asset, holder, hexBlockNo string) (int64, bool) {
	blockNo, _ := strconv.ParseInt(hexBlockNo[2:], 16, 64)
	key := fmt.Sprintf("%d|%s|%s", blockNo, asset, holder)
	if bal, ok := mapping[key]; ok {
		return bal, true
	}
	return 0, false
}

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
		p.Statement.BlockNumber = block
		p.Statement.TransactionIndex = tx
		p.Statement.LogIndex = log
		p.Statement.AssetAddress = strings.ToLower(record[3])
		p.Statement.AccountedFor = strings.ToLower(record[4])
		p.Statement.CheckpointBalance, _ = strconv.ParseInt(record[5], 10, 64)
		key := mapKey(block, tx, 0)
		logsByTx[key] = append(logsByTx[key], p)
	}

	mapFile, _ := os.Open("tests/mapping.csv")
	defer mapFile.Close()
	mapReader := csv.NewReader(mapFile)
	mapRecords, _ := mapReader.ReadAll()
	for _, record := range mapRecords[1:] {
		key := fmt.Sprintf("%s|%s|%s", record[0], record[1], record[2])
		bal, _ := strconv.ParseInt(record[3], 10, 64)
		mapping[key] = bal
	}
}

func mapKey(block, txid, logid int) int {
	return block*10000001 + txid*100001 + logid
}
