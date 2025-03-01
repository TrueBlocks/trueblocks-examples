#!/bin/bash

FOLDER=./ accounting $1 >output.txt

# Compare to gold tests
CURRENT_DIR=$(pwd)
RELATIVE_PATH=${CURRENT_DIR#*/tests/}
DIR1="$CURRENT_DIR"
DIR2="/Users/jrush/Development.2/trueblocks-core/examples/accounting/tests/$RELATIVE_PATH"
diff -r "$DIR1" "$DIR2"
