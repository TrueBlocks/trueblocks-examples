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

// TestStreamBlocks_Internal demonstrates the SDK streaming feature for blocks by
// creating a Rendering Context with both Model and Error channels.
func TestStreamBlocks_Internal[T blockType](mode ...string) {
	_ = mode
	logger.Info("Entering TestStreamBlocks_Internal")
	// Define options for this specific route
	opts := sdk.BlocksOptions{
		// EXISTING_CODE
		BlockIds: []string{"3-20000003:1000000"},
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
	case *types.Block:
		if len(mode) > 0 && mode[0] == "uncles" {
			if _, _, err = opts.BlocksUncles(); err != nil {
				fmt.Println(err.Error())
			}
		} else {
			if _, _, err = opts.Blocks(); err != nil {
				fmt.Println(err.Error())
			}
		}
	case *types.LightBlock:
		if _, _, err = opts.BlocksHashes(); err != nil {
			fmt.Println(err.Error())
		}
	case *types.Log:
		if _, _, err = opts.BlocksLogs(); err != nil {
			fmt.Println(err.Error())
		}
	case *types.Trace:
		if _, _, err = opts.BlocksTraces(); err != nil {
			fmt.Println(err.Error())
		}
	case *types.Appearance:
		if _, _, err = opts.BlocksUniq(); err != nil {
			fmt.Println(err.Error())
		}
	case *types.Withdrawal:
		if _, _, err = opts.BlocksWithdrawals(); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	// EXISTING_CODE
	if err != nil {
		logger.Info(err.Error())
	}
	logger.Info("Leaving TestStreamBlocks_Internal")
}

type blockType interface {
	// EXISTING_CODE
	*types.Block | *types.LightBlock | *types.Log | *types.Trace | *types.Appearance | *types.Withdrawal
	// EXISTING_CODE
}

// TestStreamBlocks calls into _Internal
func TestStreamBlocks() {
	// EXISTING_CODE
	TestStreamBlocks_Internal[*types.Block]()
	TestStreamBlocks_Internal[*types.LightBlock]()
	TestStreamBlocks_Internal[*types.Log]()
	TestStreamBlocks_Internal[*types.Trace]()
	TestStreamBlocks_Internal[*types.Block]("uncles")
	TestStreamBlocks_Internal[*types.Appearance]()
	TestStreamBlocks_Internal[*types.Withdrawal]()
	// EXISTING_CODE
}

// EXISTING_CODE
// EXISTING_CODE
