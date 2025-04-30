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

// TestStreamAbis_Internal demonstrates the SDK streaming feature for abis by
// creating a Rendering Context with both Model and Error channels.
func TestStreamAbis_Internal[T abiType](mode ...string) {
	_ = mode
	logger.Info("Entering TestStreamAbis_Internal")
	// Define options for this specific route
	opts := sdk.AbisOptions{
		// EXISTING_CODE
		Addrs: []string{
			"0x6b175474e89094c44da98b954eedeac495271d0f", // MakerDAO
			"0x5d3a536E4D6DbD6114cc1Ead35777bAB948E3643", // Compound
			"0x514910771af9ca656af840dff83e8264ecf986ca", // Chainlink
			"0x6b3595068778dd592e39a122f4f5a5cf09c90fe2", // SushiSwap
			"0xC011A72400E58ecD99Ee497CF89E3775d4bd732F", // Synthetix
			"0xba100000625a3754423978a60c9317c58a424e3d", // Balancer
			"0xd533a949740bb3306d119cc777fa900ba034cd52", // Curve Finance
		},
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
		_, _, err = opts.Abis()
	}
	// EXISTING_CODE
	if err != nil {
		logger.Info(err.Error())
	}
	logger.Info("Leaving TestStreamAbis_Internal")
}

type abiType interface {
	// EXISTING_CODE
	// EXISTING_CODE
}

// TestStreamAbis calls into _Internal
func TestStreamAbis() {
	// EXISTING_CODE
	TestStreamAbis_Internal[*types.Abi]()
	// EXISTING_CODE
}

// EXISTING_CODE
// EXISTING_CODE
