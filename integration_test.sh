#!/bin/bash
set -e
rm -rf /opt/path

make quick

go run FileServer.go server --address http://localhost:9000 --fspath /opt/path &

sleep 5 
cd lib/mock/client

go  test . -v  

echo "finished"

exit
