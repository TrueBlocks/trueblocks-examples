cd acctExport_transfer_2
echo "# chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv" >transfers.csv
chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_2_asset_filt
echo "# chifra export trueblocks.eth --accounting --transfers --asset 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv --first_block 8856476 --last_block 9193814 --fmt csv" >transfers.csv
chifra export trueblocks.eth --accounting --transfers --asset 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv --first_block 8856476 --last_block 9193814 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_2_ether
echo "# chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv" >transfers.csv
chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_3
echo "# chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv" >transfers.csv
chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --transfers --last_block 4000000 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_failed_2572_1
echo "# chifra export 0x054993ab0f2b1acc0fdc65405ee203b4271bebe6 --accounting --transfers --asset 0xf5b2c59f6db42ffcdfc1625999c81fdf17953384 --last_block 15549163 --max_records 40 --fmt csv" >transfers.csv
chifra export 0x054993ab0f2b1acc0fdc65405ee203b4271bebe6 --accounting --transfers --asset 0xf5b2c59f6db42ffcdfc1625999c81fdf17953384 --last_block 15549163 --max_records 40 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_failed_2572_2
echo "# chifra export 0x65b0d5e1dc0dee0704f53f660aa865c72e986fc7 --accounting --transfers --asset 0xc713e5e149d5d0715dcd1c156a020976e7e56b88 --first_block 11670418 --last_block 11670420 --max_records 40 --fmt csv" >transfers.csv
chifra export 0x65b0d5e1dc0dee0704f53f660aa865c72e986fc7 --accounting --transfers --asset 0xc713e5e149d5d0715dcd1c156a020976e7e56b88 --first_block 11670418 --last_block 11670420 --max_records 40 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_filter_traces
echo "# chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --traces --first_block 8860513 --last_block 8860531 --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv" >transfers.csv
chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --traces --first_block 8860513 --last_block 8860531 --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_filtered
echo "# chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8860513 --last_block 8860531 --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv" >transfers.csv
chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8860513 --last_block 8860531 --asset 0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_forward
echo "# chifra export 0x868b8fd259abfcfdf9634c343593b34ef359641d --accounting --transfers --traces --last_block 8769141 --fmt csv" >transfers.csv
chifra export 0x868b8fd259abfcfdf9634c343593b34ef359641d --accounting --transfers --traces --last_block 8769141 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_nft
echo "# chifra export trueblocks.eth --accounting --transfers --first_block 8876230 --last_block 9024186 --fmt csv" >transfers.csv
chifra export trueblocks.eth --accounting --transfers --first_block 8876230 --last_block 9024186 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_token_ibt
echo "# chifra export 0xec3ef464bf821c3b10a18adf9ac7177a628e87cc --accounting --transfers --first_block 7005600 --last_block 7005780 --fmt csv" >transfers.csv
chifra export 0xec3ef464bf821c3b10a18adf9ac7177a628e87cc --accounting --transfers --first_block 7005600 --last_block 7005780 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_token_ibt_2
echo "# chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 12704455 --last_block 12705893 --fmt csv" >transfers.csv
chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 12704455 --last_block 12705893 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_tributes
echo "# chifra export 0x868b8fd259abfcfdf9634c343593b34ef359641d --accounting --transfers --first_block 8769018 --last_block 8769053 --asset 0x0ba45a8b5d5575935b8158a88c631e9f9c95a2e5 --fmt csv" >transfers.csv
chifra export 0x868b8fd259abfcfdf9634c343593b34ef359641d --accounting --transfers --first_block 8769018 --last_block 8769053 --asset 0x0ba45a8b5d5575935b8158a88c631e9f9c95a2e5 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_unfiltered
echo "# chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8860513 --last_block 8860531 --fmt csv" >transfers.csv
chifra export 0xf503017d7baf7fbc0fff7492b751025c6a78179b --accounting --transfers --first_block 8860513 --last_block 8860531 --fmt csv --append --output transfers.csv

cd ../acctExport_transfer_wei_2_1
echo "# chifra export 0x05a56e2d52c817161883f50c441c3228cfe54d9f --accounting --transfers --first_record 0 --max_records 15 --fmt csv" >transfers.csv
chifra export 0x05a56e2d52c817161883f50c441c3228cfe54d9f --accounting --transfers --first_record 0 --max_records 15 --fmt csv --append --output transfers.csv

# cd ../acctExport_transfer_wei_2_2
# echo "# chifra export 0x05a56e2d52c817161883f50c441c3228cfe54d9f --accounting --statements --first_record 250 --max_records 15 --fmt csv" >transfers.csv
# chifra export 0x05a56e2d52c817161883f50c441c3228cfe54d9f --accounting --statements --first_record 250 --max_records 15 --fmt csv --append --output transfers.csv

# cd ../acctExport_transfer_wei_2_3
# echo "# chifra export 0x05a56e2d52c817161883f50c441c3228cfe54d9f --accounting --statements --first_block 15700073 --last_block 15700075 --fmt csv" >transfers.csv
# chifra export 0x05a56e2d52c817161883f50c441c3228cfe54d9f --accounting --statements --first_block 15700073 --last_block 15700075 --fmt csv --append --output transfers.csv

# cd ../acctExport_transfer_3_accounting
# echo "# chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --last_block 4000000 --fmt csv" >transfers.csv
# chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --accounting --last_block 4000000 --fmt csv --append --output transfers.csv

# cd ../acctExport_transfer_3_bad
# echo "# chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --transfers --last_block 4000000 --fmt csv" >transfers.csv
# chifra export 0x08166f02313feae18bb044e7877c808b55b5bf58 --transfers --last_block 4000000 --fmt csv --append --output transfers.csv

# Compare to gold tests
CURRENT_DIR=$(pwd)
RELATIVE_PATH=${CURRENT_DIR#*/tests/}
DIR1="$CURRENT_DIR"
DIR2="/Users/jrush/Development.2/trueblocks-core/examples/accounting/tests/$RELATIVE_PATH"
diff -r "$DIR1" "$DIR2"

cd ..
