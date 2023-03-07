#!/bin/bash

cd internal/user
echo "=== start build user image ==="
go build -o ./build/main ./main
if [ $? -ne 0 ]; then
		exit 1
fi
echo "=== finish build user image ==="

cd ../drive
echo "=== start build drive image ==="
go build -o ./build/main ./main
if [ $? -ne 0 ]; then
		exit 1
fi
echo "=== finish build drive image ==="

cd ../storage
echo "=== start build storage image ==="
go build -o ./build/main ./main
if [ $? -ne 0 ]; then
		exit 1
fi
echo "=== finish build storage image ==="

cd ../..
echo "=== start deploy ==="
docker compose down
docker compose build
docker image prune --force
docker volume prune --force
docker compose up -d
echo "=== finish deploy ==="
