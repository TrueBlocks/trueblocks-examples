// Copyright 2016, 2024 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated from template examples_withStreaming_stream_route.go.tmpl
 * Modify only inside the 'EXISTING_CODE' tags to preserve code changes.
 */

package main

import (
	"fmt"
	"time"

	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/logger"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/output"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v6"
)

// TestStreamTransactions_Internal demonstrates the SDK streaming feature for transactions by
// creating a Rendering Context with both Model and Error channels.
func TestStreamTransactions_Internal[T transactionType](mode ...string) {
	_ = mode
	logger.Info("Entering TestStreamTransactions_Internal")
	// Define options for this specific route
	opts := sdk.TransactionsOptions{
		// EXISTING_CODE
		TransactionIds: []string{"11000000.*", "21000000.*"},
		// EXISTING_CODE
		RenderCtx: output.NewStreamingContext(),
	}
	opts.Globals.Cache = true

	// Set up timeout and cancel the context when it hits.
	runFor := 1 * time.Second
	go func() {
		time.Sleep(runFor)
		opts.RenderCtx.Cancel()
	}()

	// Process the streaming data
	startTime := time.Now()
	go func() {
		for {
			select {
			case model := <-opts.RenderCtx.ModelChan:
				// Type assertion based on expected return type
				elapsedSeconds := time.Since(startTime).Seconds()
				value := model.Model("mainnet", "csv", false, nil)
				logger.Info(fmt.Sprintf("%.3fs: received item of type: %s", elapsedSeconds, value))
			case err := <-opts.RenderCtx.ErrorChan:
				logger.Info("Error returned by fetchData:", err)
			}
		}
	}()

	// Generate the stream
	var err error
	// EXISTING_CODE
	var v T
	switch any(v).(type) {
	default:
		_, _, err = opts.Transactions()
	}
	// EXISTING_CODE
	if err != nil {
		logger.Info(err.Error())
	}
	logger.Info("Leaving TestStreamTransactions_Internal")
}

type transactionType interface {
	// EXISTING_CODE
	// EXISTING_CODE
}

// TestStreamTransactions calls into _Internal
func TestStreamTransactions() {
	// EXISTING_CODE
	TestStreamTransactions_Internal[*types.Transaction]()
	// EXISTING_CODE
}

// EXISTING_CODE
// TODO: Add streaming examples
// Transactions ([]types.Transaction, *types.MetaData, error) {
// TransactionsLogs ([]types.Log, *types.MetaData, error) {
// TransactionsTraces ([]types.Trace, *types.MetaData, error) {
// TransactionsUniq ([]types.Appearance, *types.MetaData, error) {
// EXISTING_CODE
