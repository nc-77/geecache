#!/bin/bash

go run ./server/main.go -port=8001 &
go run ./server/main.go -port=8002 &
go run ./server/main.go -port=8003 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait