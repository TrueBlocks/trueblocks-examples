# Comparison Example

## Overview

This example provides a template for comparative analysis between different data sources or methodologies. It demonstrates how to structure applications that compare blockchain data across different providers or analysis methods.

## What This Example Shows

- **Setting up comparative analysis frameworks for blockchain data** - Building structured approaches for systematically comparing different data sources, time periods, or analysis methodologies
- **Structuring applications for A/B testing of different data sources** - Creating flexible architectures that can easily switch between and compare multiple blockchain data providers
- **Implementing data validation and comparison methodologies** - Developing robust techniques for ensuring data accuracy and identifying meaningful differences between datasets
- **Building extensible analysis tools that can accommodate different comparison scenarios** - Designing modular systems that can be adapted for various types of comparative analysis
- **Understanding performance differences between various blockchain data access patterns** - Learning how different approaches to accessing blockchain data impact speed, reliability, and resource usage

## Prerequisites

Ensure you have the following installed and running:

- Go Version 1.23 or higher
- The latest version of TrueBlocks Core

## Installation

Clone the repository:

```bash
git clone https://github.com/TrueBlocks/trueblocks-core.git
```

Build TrueBlocks Core by following the ([installation instructions](https://trueblocks.io/docs/install/install-core/)), and then navigate to the example directory:

```bash
cd examples/comparison
```

## Usage

Run the code from the example folder using the following command:

```bash
go run .
```

## Key Concepts

This example is particularly useful for:

- Financial analysts examining blockchain transaction patterns
- Researchers studying blockchain behavior over time
- Developers building analytics dashboards
- Portfolio managers tracking address performance
- Compliance teams identifying unusual patterns

The implementation provides a foundation for building sophisticated analysis tools that can reveal insights hidden in blockchain data through systematic comparison methodologies.

```bash
go run .
```

There are no command line options.

### Output

The code produces results similar to the following:

```bash
...
DATA
DATA
DATA
...
```

![MAKE AN IMAGE](./IMAGE_NAME.png)

## Troubleshooting

If you encounter issues, check the following:

- Ensure at least Go Version 1.23.
- Make sure you have a valid Ethereum mainnet RPC configured.
- Ensure TrueBlocks Core is properly installed and configured.
- Verify the command-line arguments are correct and within valid ranges.

For further assistance, refer to the TrueBlocks Documentation.

## License

See the LICENSE file at the root of this repo for details.
