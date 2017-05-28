#!/bin/bash

# THIS SCRIPT WAS TAKEN FROM : https://github.com/getlantern/flashlight-build/blob/devel/testandcover.bash
# This script tests multiple packages and creates a consolidated cover profile
# The list of packages to test is specified in packages.txt.

function die() {
  echo $*
  exit 1
}

# Initialize profile.cov
echo "mode: count" > coverage.out

# Initialize error tracking
ERROR=""


# Test each package and append coverage profile info to profile.cov
for pkg in `cat packages.txt`
do
    go test -v -covermode=count -coverprofile=coverage_tmp.out $pkg || ERROR="Error testing $pkg"
    tail -n +2 coverage_tmp.out >> coverage.out || die "Unable to append coverage for $pkg"
done

if [ ! -z "$ERROR" ]
then
    die "Encountered error, last error was: $ERROR"
fi
