#!/usr/bin/env bash

TEST_DIR=/tmp/test_output

# delete old test results
rm -rf /tmp/test_output/result.json

# do nothing if no test artifacts exist
if ! find $TEST_DIR/*.result > /dev/null 2>&1; then
    echo "No test artifacts found"
    exit 0
fi

# merge all `.result` files into result.json
jq -s add $TEST_DIR/*.result > $TEST_DIR/result.json

# delete all intermediate files so we don't end up collecting old test data
# next time we execute `collect-test-results`
rm -rf /tmp/test_output/*.result