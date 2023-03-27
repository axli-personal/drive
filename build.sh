#!/bin/bash
# backend:
echo "1) start build backend image"
cd backend

echo "1.1) start build user image"
cd user
go build -o ./build/main ./main
if [ $? -ne 0 ]; then
		exit 1
fi
cd ..
echo "1.1) finish build user image"

echo "1.2) start build drive image"
cd drive
go build -o ./build/main ./main
if [ $? -ne 0 ]; then
		exit 1
fi
cd ..
echo "1.2) finish build drive image"

echo "1.3) start build storage image"
cd storage
go build -o ./build/main ./main
if [ $? -ne 0 ]; then
		exit 1
fi
cd ..
echo "1.3) finish build storage image"

cd ..
echo "1) finish build backend image"

# fronted:
echo "2) start build frontend image"
cd frontend
npm run build
if [ $? -ne 0 ]; then
		exit 1
fi
cd ..
echo "2) finish build frontend image"

echo "3) start deploy"
docker compose down
docker compose build
if [ $? -ne 0 ]; then
		exit 1
fi
docker image prune --force
docker volume prune --force
docker compose up -d
echo "3) finish deploy"
