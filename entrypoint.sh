#!/bin/bash

echo "Running migrations..."
go run ./migration/create/create_migration.go
if [ $? -ne 0 ]; then
  echo "Migration failed!"
  exit 1
fi

echo "Running seeding company profile..."
go run ./seeders/company-profile/company_profile_seed.go
if [ $? -ne 0 ]; then
  echo "Seeding company profile failed!"
  exit 1
fi

echo "Running seeding data..."
go run ./seeders/data/data.go
if [ $? -ne 0 ]; then
  echo "Seeding data failed!"
  exit 1
fi

echo "Starting Go server..."
./my-go-server
if [ $? -ne 0 ]; then
  echo "Go server failed to start!"
  exit 1
fi

echo "All processes completed successfully!"
