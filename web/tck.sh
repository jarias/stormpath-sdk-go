#! /bin/bash

set -e

if [ ! -d stormpath-framework-tck ]; then
    echo "Checking out TCK"
    git config user.email "evangelists@stormpath.com"
    git config user.name "stormpath-sdk-java TCK"
    git clone https://github.com/stormpath/stormpath-framework-tck.git stormpath-framework-tck
    echo "TCK cloned"
fi

echo "Running TCK now!"

cd stormpath-framework-tck
git fetch
git checkout master

mvn -q --fail-at-end clean verify -Dstormpath.tck.webapp.port=$1 -Dstormpath.application.href=$2 &> mvn.out