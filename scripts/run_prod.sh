#!/bin/bash

set -e

echo "We out here in $(pwd)"

echo "==> Killing the old binary.."
./kill.sh

echo "==> Deleting the old binary..."
rm -rf ../../service/rss-exporter

mkdir -p ../../service
mv ../rss-exporter ../../service

cd ../../service

echo "==> Starting the new Service!"
./rss-exporter > service.log 2>&1 & 