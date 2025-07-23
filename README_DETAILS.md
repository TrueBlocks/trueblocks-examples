# TrueBlocks Examples - Detailed Documentation

This file provides comprehensive details about each example in the TrueBlocks examples directory. Each example demonstrates different aspects of the TrueBlocks SDK and ecosystem, showing how to build applications that interact with blockchain data through the TrueBlocks indexing system.

## simple ([link](./simple/README.md))

This is the foundational example designed to get developers started with the TrueBlocks SDK. It demonstrates basic blockchain data retrieval by fetching blocks at monthly intervals and displaying their numbers and timestamps.

**Key Learning Points:**

- How to set up and initialize the TrueBlocks SDK in a Go application
- Basic usage of the `sdk.Blocks()` function to retrieve blockchain data
- Working with block numbers and timestamps for time-based queries
- Simple data iteration and output formatting
- Building a minimal application that connects to the TrueBlocks indexing system

## balanceChart ([link](./balanceChart/README.md))

This example demonstrates advanced data visualization by creating balance charts for Ethereum addresses over time. It retrieves prefund addresses from the genesis block and tracks their balance changes across specified time periods.

**Key Learning Points:**

- Using the TrueBlocks SDK to query historical balance data across multiple time periods
- Working with prefund addresses and genesis block data analysis
- Integrating data visualization libraries to create grouped bar charts from blockchain data
- Handling command-line parameters for configurable data analysis (number of addresses, maximum amounts, time periods)
- Processing and aggregating large datasets for meaningful visual representation
- Combining multiple SDK calls to build comprehensive financial analysis tools

## cancelContext ([link](./cancelContext/README.md))

This example serves as a test case for proper context cancellation patterns in Go applications using the TrueBlocks SDK. It demonstrates how to handle graceful shutdowns and resource cleanup.

**Key Learning Points:**

- Implementing proper context cancellation patterns in blockchain data applications
- Testing scenarios for handling interrupted operations and graceful shutdowns
- Understanding how the TrueBlocks SDK responds to context cancellation
- Best practices for resource management in long-running blockchain queries
- Building robust applications that can handle network interruptions and user cancellations

## chainList ([link](./chainList/README.md))

This utility example showcases how to retrieve and display information about all supported blockchain networks in the TrueBlocks ecosystem. It generates a formatted table of chain information.

**Key Learning Points:**

- Using the TrueBlocks utils package to access chain configuration data
- Working with multi-chain support and chain metadata
- Formatting structured data output in markdown tables
- Understanding the relationship between chain IDs, names, and RPC endpoints
- Building informational tools that help users understand available blockchain networks

## checkNodes ([link](./checkNodes/main.go))

This diagnostic tool evaluates the performance and connectivity of RPC nodes across different blockchain networks. It provides detailed timing and availability information for network endpoints.

**Key Learning Points:**

- Implementing network connectivity testing and performance benchmarking
- Working with multiple RPC endpoints and measuring response times
- Building diagnostic tools for blockchain infrastructure monitoring
- Handling network timeouts and error conditions gracefully
- Creating detailed performance reports with tabular data formatting
- Understanding RPC node reliability patterns across different chains

## chunkDiagnostics ([link](./chunkDiagnostics/README.md))

A comprehensive analysis tool that validates the integrity and consistency of TrueBlocks chunk data across different sources including manifests, indexes, blooms, and IPFS availability.

**Key Learning Points:**

- Performing complex data integrity validation across multiple data sources
- Working with TrueBlocks internal data structures (manifests, indexes, bloom filters)
- Implementing IPFS connectivity and hash verification workflows
- Building comprehensive diagnostic tools for data consistency checking
- Creating detailed analysis reports with multiple validation categories
- Understanding the TrueBlocks indexing architecture and data organization

## comparison ([link](./comparison/README.md))

This example provides a template for comparative analysis between different data sources or methodologies. It demonstrates how to structure applications that compare blockchain data across different providers or analysis methods.

**Key Learning Points:**

- Setting up comparative analysis frameworks for blockchain data
- Structuring applications for A/B testing of different data sources
- Implementing data validation and comparison methodologies
- Building extensible analysis tools that can accommodate different comparison scenarios
- Understanding performance differences between various blockchain data access patterns

## findFirst ([link](./findFirst/README.md))

This performance-focused example demonstrates efficient search patterns for finding the first transaction in a blockchain. It includes comprehensive benchmarking data showing optimization techniques for concurrent processing.

**Key Learning Points:**

- Implementing efficient search algorithms for blockchain data discovery
- Understanding and optimizing concurrent processing patterns with multiple workers
- Performance benchmarking and measurement techniques for blockchain applications
- Balancing concurrency levels for optimal throughput without overwhelming RPC endpoints
- Building search applications that can scale across different blockchain sizes
- Understanding the relationship between worker count and processing efficiency

## four_bytes ([link](./four_bytes/main.go))

A minimal example that demonstrates basic project structure and setup patterns. Currently serves as a template for function signature analysis and smart contract interaction patterns.

**Key Learning Points:**

- Basic Go project structure and initialization patterns
- Template setup for smart contract function signature analysis
- Understanding four-byte function selectors in Ethereum smart contracts
- Building foundation code for more complex smart contract interaction tools
- Project organization patterns for TrueBlocks-based applications

## keystore ([link](./keystore/README.md))

**Warning: This is an internal-use example that is explicitly not recommended for production use.**

This example demonstrates keystore and cryptographic operations but comes with strong warnings about security and auditability. It's preserved for historical and internal development purposes only.

**Key Learning Points:**

- Understanding the security considerations when working with blockchain keystores
- Learning what NOT to do when building production cryptocurrency applications
- Recognizing the importance of security audits in blockchain development
- Understanding the legal and liability considerations in blockchain application development
- Appreciating the complexity and risks involved in cryptographic key management

## monitorService ([link](./monitorService/))

This directory appears to be reserved for future development of monitoring service examples. It will likely demonstrate how to build applications that continuously monitor blockchain addresses and transactions.

**Key Learning Points:**

- (Future) Implementing continuous monitoring patterns for blockchain addresses
- (Future) Building event-driven applications that respond to blockchain state changes
- (Future) Understanding notification and alerting patterns for blockchain data
- (Future) Creating persistent monitoring services with proper error handling
- (Future) Integrating with external notification systems and APIs

## nameManager ([link](./nameManager/README.md))

A comprehensive CLI tool that demonstrates advanced name management functionality within the TrueBlocks ecosystem. It provides full CRUD operations for address-to-name mappings with publishing capabilities.

**Key Learning Points:**

- Building full-featured CLI applications using the TrueBlocks SDK
- Implementing CRUD operations for blockchain address name management
- Working with the TrueBlocks names database and custom naming systems
- Creating command-line interfaces with complex argument parsing and validation
- Understanding the relationship between addresses, names, tags, and metadata
- Building applications that can publish and share naming data across the TrueBlocks network

## withStreaming ([link](./withStreaming/main.go))

This example demonstrates the streaming capabilities of the TrueBlocks SDK, showing how to handle real-time data flows from various TrueBlocks endpoints including blocks, transactions, logs, and traces.

**Key Learning Points:**

- Implementing streaming data patterns for real-time blockchain data processing
- Working with multiple TrueBlocks streaming endpoints (export, logs, names, receipts, traces, transactions)
- Handling continuous data flows and managing streaming connections
- Building applications that can process large volumes of blockchain data efficiently
- Understanding the performance benefits of streaming vs. batch processing
- Creating robust error handling for long-running streaming operations
