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

// TestStreamTokens_Internal demonstrates the SDK streaming feature for tokens by
// creating a Rendering Context with both Model and Error channels.
func TestStreamTokens_Internal[T tokenType](mode ...string) {
	_ = mode
	logger.Info("Entering TestStreamTokens_Internal")
	// Define options for this specific route
	opts := sdk.TokensOptions{
		// EXISTING_CODE
		Addrs:    []string{"0x1f98431c8ad98523631ae4a59f267346ea31f984", "trueblocks.eth"},
		BlockIds: []string{"11000000-11500000"},
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
		_, _, err = opts.Tokens()
	}
	// EXISTING_CODE
	if err != nil {
		logger.Info(err.Error())
	}
	logger.Info("Leaving TestStreamTokens_Internal")
}

type tokenType interface {
	// EXISTING_CODE
	// EXISTING_CODE
}

// TestStreamTokens calls into _Internal
func TestStreamTokens() {
	// EXISTING_CODE
	TestStreamTokens_Internal[*types.Token]()
	// EXISTING_CODE
}

// EXISTING_CODE
// TODO: Add streaming examples
// Tokens ([]types.Token, *types.MetaData, error) {
// EXISTING_CODE
