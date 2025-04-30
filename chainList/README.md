# Chain List Example

This example demonstrates how to use the TrueBlocks utils package to retrieve and display information about available blockchain networks.

## Prerequisites

- Go 1.19 or higher
- TrueBlocks Core installed

## Installation

```bash
go mod tidy
```

## Usage

Run the example with:

```bash
go run .
```

This will produce a markdown table showing:
- Chain name
- Chain ID
- Native currency symbol
- An RPC endpoint for each chain

## Code Structure

The example code:
1. Imports the TrueBlocks utils package
2. Calls the `GetChainList()` function to retrieve available chains
3. Formats and displays the chain information in a markdown table
```
