# Examples

This directory contains examples demonstrating how to use the TrueBlocks core SDK for various blockchain data tasks.

## Overview

These examples showcase different features and capabilities of TrueBlocks, helping developers understand how to:

- Query and analyze blockchain data
- Work with addresses, transactions, and blocks
- Efficiently index and retrieve on-chain information
- Implement common blockchain data patterns

## Available Examples

- **balanceChart** - Generate charts showing ETH balance changes over time
- **cancelContext** - Demonstrate context cancellation for long-running operations
- **chainList** - List and display information about supported chains
- **checkNodes** - Verify and check node connectivity and status
- **comparison** - Compare data between different sources or time periods  
- **findFirst** - Find the first occurrence of specific blockchain events
- **four_bytes** - Work with four-byte function signatures and method IDs
- **keystore** - Manage and work with Ethereum keystore files
- **monitorService** - Monitor blockchain addresses and events
- **nameManager** - Work with ENS names and address resolution
- **simple** - Basic example showing fundamental SDK usage
- **withStreaming** - Demonstrate streaming data capabilities

## Running Examples

Most examples can be run with:

```bash
cd example-directory
go run .
```

## Creating New Examples

### For Development/Testing (Local Branch)

When creating examples to test new features on a local branch:

1. **Create from template:**
   ```bash
   cd examples
   chifra init --example your-example-name
   ```

2. **Add local replace directives** to the generated `go.mod`:
   ```go
   replace github.com/TrueBlocks/trueblocks-sdk/v5 => ../sdk
   replace github.com/TrueBlocks/trueblocks-core/src/apps/chifra => ../src/apps/chifra
   ```

3. **Add to .gitignore** to prevent breaking CI:
   ```bash
   echo "your-example-name/" >> .gitignore
   ```

4. **Update go.work:**
   ```bash
   cd ..
   ./scripts/go-work-sync.sh
   ```

### For Production (Ready to Commit)

When your example is ready for production:

1. **Remove replace directives** from `go.mod`
2. **Remove from .gitignore** 
3. **Run go-work-sync.sh** to update dependencies to published versions
4. **Test the example** works with published dependencies
5. **Commit and submit PR**

## Dependencies

To run these examples, you'll need:

1. TrueBlocks core installed
2. Access to an Ethereum node (local or remote)
3. Any additional dependencies specified in the example's README

## Important Notes

- **Development Examples**: Examples with `replace` directives in `go.mod` should be in `.gitignore`
- **CI/CD**: Examples committed to the repo must use published dependencies only
- **go.work**: Always run `./scripts/go-work-sync.sh` after adding new examples
- **Testing**: Test examples work with both local and published dependencies

## Documentation

For more detailed information, visit the [TrueBlocks documentation](https://trueblocks.io/docs/).

## Contributing

We welcome contributions! When adding new examples:

1. Follow the development process above
2. Ensure examples work with published dependencies
3. Include a README.md explaining the example's purpose
4. Add appropriate error handling and comments
