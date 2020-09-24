#!/bin/bash

PID=$(ps -eaf | grep rss-exporter | grep -v grep | awk '{print $2}')

if [[ "" != "$PID" ]]; then
    echo "Killing $PID"
    kill -9 $PID
fi
