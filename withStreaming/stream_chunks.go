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

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v6"
)

// TestStreamChunks_Internal demonstrates the SDK streaming feature for chunks by
// creating a Rendering Context with both Model and Error channels.
func TestStreamChunks_Internal[T chunkType](mode ...string) {
	_ = mode
	logger.Info("Entering TestStreamChunks_Internal")
	// Define options for this specific route
	opts := sdk.ChunksOptions{
		// EXISTING_CODE
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
		_, _, err = opts.ChunksAddresses()
	}
	// EXISTING_CODE
	if err != nil {
		logger.Info(err.Error())
	}
	logger.Info("Leaving TestStreamChunks_Internal")
}

type chunkType interface {
	// EXISTING_CODE
	// EXISTING_CODE
}

// TestStreamChunks calls into _Internal
func TestStreamChunks() {
	// EXISTING_CODE
	TestStreamChunks_Internal[*types.ChunkAddress]()
	// EXISTING_CODE
}

// EXISTING_CODE
// TODO: Add streaming examples
// ChunksAddresses ([]types.ChunkAddress, *types.MetaData, error) {
// ChunksAppearances ([]types.ChunkAppearance, *types.MetaData, error) {
// ChunksBlooms ([]types.ChunkBloom, *types.MetaData, error) {
// ChunksDiff ([]types.Message, *types.MetaData, error) {
// ChunksIndex ([]types.ChunkIndex, *types.MetaData, error) {
// ChunksManifest ([]types.ChunkManifest, *types.MetaData, error) {
// ChunksPins ([]types.ChunkPin, *types.MetaData, error) {
// ChunksStats ([]types.ChunkStats, *types.MetaData, error) {
// ChunksTag(val string) ([]types.Message, *types.MetaData, error) {
// ChunksTruncate(val base.Blknum) ([]types.Message, *types.MetaData, error
// EXISTING_CODE
