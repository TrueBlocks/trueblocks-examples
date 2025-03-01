#!/bin/bash

# Find all directories starting with acctExport
for dir in acctExport*/; do
    if [ -d "$dir" ] && [ -f "$dir/run_account.sh" ]; then
        echo "Running script in: $dir"
        cd "$dir"
        ./run_account.sh
        cd ..
    fi
done
