package main

import (
	"fmt"
	"strconv"
)

var (
	conn Connection
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

func mapKey(block, txid, logid int) int {
	return block*10000001 + txid*100001 + logid
}
