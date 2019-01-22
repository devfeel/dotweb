#!/usr/bin/env bash
set -e
export COVERAGE_PATH=$(pwd)
rm -rf "${COVERAGE_PATH}/scripts/coverage.out"
for d in $(go list ./... | grep -v vendor); do
    cd "$GOPATH/src/$d"
    if [ $(ls | grep _test.go | wc -l) -gt 0 ]; then
        go test -cover -covermode atomic -coverprofile coverage.out  
        if [ -f coverage.out ]; then
            sed '1d;$d' coverage.out >> "${COVERAGE_PATH}/scripts/coverage.out"
            rm -f coverage.out
        fi
    fi
done
