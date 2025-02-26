#!/bin/bash

for file in ./tests/*.txt; do
    filename=$(basename "$file")
    folder=$(dirname "$file")
    address="${filename%%_*}"
    mkdir -p "./output/${filename}"
    echo "FOLDER=./output/${filename} accounting $address >./output/${filename}/${filename}.txt"
    echo "block,tx" > "./output/${filename}/apps.csv"
    cat "$file" | jq -r '.data[] | [.blockNumber, .transactionIndex] | join(",")'  >> "./output/${filename}/apps.csv"
    echo "block,tx,log,asset,holder,amount" >./output/${filename}/transfers.csv
    cat "$file" | jq -r '.data[] | [.blockNumber, .transactionIndex, .logIndex, .assetAddress, .accountedFor, .amountNet] | join(",")' >>./output/${filename}/transfers.csv
    echo "block,asset,holder,endBal" >./output/${filename}/balances.csv
    cat "$file" | jq -r '.data[] | [.blockNumber, .assetAddress, .accountedFor, .endBal] | join(",")' >>./output/${filename}/balances.csv
done

for file in ./tests.2/*.txt; do
    filename=$(basename "$file")
    folder=$(dirname "$file")
    address="${filename%%_*}"
    mkdir -p "./output/${filename}"
    echo "FOLDER=./output/${filename} accounting $address >./output/${filename}/${filename}.txt"
    echo "block,tx" > "./output/${filename}/apps.csv"
    cat "$file" | jq -r '.data[].statements[] | [.blockNumber, .transactionIndex] | join(",")'  >> "./output/${filename}/apps.csv"
    echo "block,tx,log,asset,holder,amount" >./output/${filename}/transfers.csv
    cat "$file" | jq -r '.data[].statements[] | [.blockNumber, .transactionIndex, .logIndex, .assetAddress, .accountedFor, .amountNet] | join(",")' >>./output/${filename}/transfers.csv
    echo "block,asset,holder,endBal" >./output/${filename}/balances.csv
    cat "$file" | jq -r '.data[].statements[] | [.blockNumber, .assetAddress, .accountedFor, .endBal] | join(",")' >>./output/${filename}/balances.csv
done
