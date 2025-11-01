package main

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v6/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v6/pkg/utils"
)

func main() {
	// Get the chain list
	chainList, err := utils.UpdateChainList("./data")
	if err != nil {
		logger.Fatalf("Error getting chain list: %v", err)
	}

	// Print markdown table header
	fmt.Println("| Chain Name | Chain ID | Native Currency | RPC Endpoint |")
	fmt.Println("|------------|----------|-----------------|--------------|")

	// Iterate through chains and print details
	for _, chain := range chainList.Chains {
		// Get the first RPC endpoint if available
		rpcEndpoint := "N/A"
		if len(chain.Rpc) > 0 {
			rpcEndpoint = chain.Rpc[0]
		}

		// Print chain details as a markdown table row
		fmt.Printf("| %s | %d | %s | %s |\n",
			chain.Name,
			chain.ChainID,
			chain.NativeCurrency.Symbol,
			rpcEndpoint)
	}
}
