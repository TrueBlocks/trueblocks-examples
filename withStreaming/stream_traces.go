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

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// TestStreamTraces_Internal demonstrates the SDK streaming feature for traces by
// creating a Rendering Context with both Model and Error channels.
func TestStreamTraces_Internal[T traceType](mode ...string) {
	_ = mode
	logger.Info("Entering TestStreamTraces_Internal")
	// Define options for this specific route
	opts := sdk.TracesOptions{
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
		_, _, err = opts.Traces()
	}
	// EXISTING_CODE
	if err != nil {
		logger.Info(err.Error())
	}
	logger.Info("Leaving TestStreamTraces_Internal")
}

type traceType interface {
	// EXISTING_CODE
	// EXISTING_CODE
}

// TestStreamTraces calls into _Internal
func TestStreamTraces() {
	// EXISTING_CODE
	TestStreamTraces_Internal[*types.Trace]()
	// EXISTING_CODE
}

// EXISTING_CODE
// TODO: Add streaming examples
// Traces ([]types.Trace, *types.MetaData, error) {
// TracesCount ([]types.TraceCount, *types.MetaData, error) {
// EXISTING_CODE
