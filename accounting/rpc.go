package main

import (
	"fmt"
	"strconv"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

var (
	conn Connection
)

// Connection provides on-chain balance lookups
type Connection struct{}

func (c *Connection) GetBalanceAtToken(asset string, holder base.Address, hexBlockNo string) (int64, bool) {
	blockNo, _ := strconv.ParseInt(hexBlockNo[2:], 16, 64)
	key := fmt.Sprintf("%d|%s|%s", blockNo, asset, holder.Hex())
	if bal, ok := mapping[key]; ok {
		return bal, true
	}
	return 0, false
}

func mapKey(block, txid, logid int) int {
	return block*10000001 + txid*100001 + logid
}
