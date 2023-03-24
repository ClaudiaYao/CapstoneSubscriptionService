#! usr/bin/bash

cd app
cd domain
cd auth 
go test run TestCreateJWT ./...
cd ..
cd ..
cd ..