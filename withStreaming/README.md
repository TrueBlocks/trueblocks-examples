# WithStreaming Example

This example demonstrates how to use TrueBlocks SDK's streaming capabilities to process blockchain data efficiently.

## What it does

The example shows how to stream various data types (blocks, transactions, logs, ABIs, etc.) using channels to process large amounts of blockchain data asynchronously.

## Available Streaming Examples

Several examples are implemented:
- `TestStreamBlocks`: Stream block data
- `TestStreamExport`: Stream transactions with timeout cancellation
- `TestStreamNames`: Stream name data
- `TestStreamProgress`: Track progress with a progress bar
- `TestStreamAbis`: Stream and process contract ABIs

## TestStreamExport Example

The export example demonstrates how to:

1. Stream transactions for a specific address (Uniswap V3 Factory)
2. Process the transactions as they arrive
3. Cancel the streaming operation after a timeout

```go
// Create options with a streaming context
opts := sdk.ExportOptions{
    Addrs:     []string{"0x1f98431c8ad98523631ae4a59f267346ea31f984"},
    Unripe:    true,  // to the head of the chain
    RenderCtx: output.NewStreamingContext(),
}

// Set up a cancellation after 3 seconds
go func() {
    time.Sleep(3 * time.Second)
    opts.RenderCtx.Cancel()  // This gracefully stops the streaming
}()

// Process the data as it streams in
go func() {
    for {
        select {
        case model := <-opts.RenderCtx.ModelChan:
            if tx, ok := model.(*types.Transaction); ok {
                fmt.Printf("%d\t%d\n", tx.BlockNumber, tx.TransactionIndex)
            }
        case err := <-opts.RenderCtx.ErrorChan:
            // Handle errors
        }
    }
}()

// Start the stream
opts.Export()
```

This pattern is useful when you need to process large amounts of data without loading everything into memory at once.

## Running the example

```
go build -o out main.go
./out
```

By default, `TestStreamExport` runs, but you can uncomment other examples in `main.go`.

Note: If you remove this example, you will need to also remove the reference to it from ./go.work in the root folder.
