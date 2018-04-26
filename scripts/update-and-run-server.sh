#/bin/bash

rm -r goconcept-files
mkdir -p goconcept-files/admin-frontend
cp -r $GOPATH/src/github.com/fivebillionmph/goconcept/assets/admin-frontend/public/* goconcept-files/admin-frontend

cd go-src2
go build -o ../server || exit
cd ..
./server
