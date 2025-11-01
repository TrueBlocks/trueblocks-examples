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

// TestStreamStatus_Internal demonstrates the SDK streaming feature for status by
// creating a Rendering Context with both Model and Error channels.
func TestStreamStatus_Internal[T statusType](mode ...string) {
	_ = mode
	logger.Info("Entering TestStreamStatus_Internal")
	// Define options for this specific route
	opts := sdk.StatusOptions{
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
		_, _, err = opts.StatusAll()
	}
	// EXISTING_CODE
	if err != nil {
		logger.Info(err.Error())
	}
	logger.Info("Leaving TestStreamStatus_Internal")
}

type statusType interface {
	// EXISTING_CODE
	// EXISTING_CODE
}

// TestStreamStatus calls into _Internal
func TestStreamStatus() {
	// EXISTING_CODE
	TestStreamStatus_Internal[*types.Status]()
	// EXISTING_CODE
}

// EXISTING_CODE
// TODO: Add streaming examples
// Status ([]types.Status, *types.MetaData, error) {
// StatusAbis ([]types.Status, *types.MetaData, error) {
// StatusAll ([]types.Status, *types.MetaData, error) {
// StatusBlocks ([]types.Status, *types.MetaData, error) {
// StatusBlooms ([]types.Status, *types.MetaData, error) {
// StatusDiagnose ([]types.Status, *types.MetaData, error) {
// StatusHealthcheck ([]types.Status, *types.MetaData, error) {
// StatusIndex ([]types.Status, *types.MetaData, error) {
// StatusLogs ([]types.Status, *types.MetaData, error) {
// StatusMaps ([]types.Status, *types.MetaData, error) {
// StatusMonitors ([]types.Status, *types.MetaData, error) {
// StatusNames ([]types.Status, *types.MetaData, error) {
// StatusResults ([]types.Status, *types.MetaData, error) {
// StatusSlurps ([]types.Status, *types.MetaData, error) {
// StatusSome ([]types.Status, *types.MetaData, error) {
// StatusStaging ([]types.Status, *types.MetaData, error) {
// StatusState ([]types.Status, *types.MetaData, error) {
// StatusStatements ([]types.Status, *types.MetaData, error) {
// StatusTokens ([]types.Status, *types.MetaData, error) {
// StatusTraces ([]types.Status, *types.MetaData, error) {
// StatusTransactions ([]types.Status, *types.MetaData, error) {
// StatusUnripe ([]types.Status, *types.MetaData, error) {
// EXISTING_CODE
