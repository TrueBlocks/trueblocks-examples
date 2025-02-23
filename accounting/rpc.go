package main

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

var (
	cc Connection
)

// Connection provides on-chain balance lookups
type Connection struct{}

func (c *Connection) GetBalanceAtToken(asset base.Address, holder base.Address, bn base.Blknum) (int64, bool) {
	key := fmt.Sprintf("%d|%s|%s", bn, asset.Hex(), holder.Hex())
	if bal, ok := balanceMap[key]; ok {
		return bal, true
	}
	return 0, false
}

func mapKey(block, txid, logid int) int {
	return block*10000001 + txid*100001 + logid
}
