#!/usr/bin/env bash

echo "block,tx" >$1/apps.csv
cat $1/$1.json | jq -r '.data[] | [.blockNumber, .transactionIndex] | join(",")' >>$1/apps.csv

echo "block,tx,log,asset,holder,amount" >$1/transfers.csv
cat $1/$1.json | jq -r '.data[] | [.blockNumber, .transactionIndex, .logIndex, .assetAddress, .accountedFor, .amountNet] | join(",")' >>$1/transfers.csv

echo "block,asset,holder,endBal" >$1/balances.csv
cat $1/$1.json | jq -r '.data[] | [.blockNumber, .assetAddress, .accountedFor, .endBal] | join(",")' >>$1/balances.csv
