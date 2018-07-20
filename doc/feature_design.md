#
Simplified Architecture Version

## Basic features.
1. User upload a file with a file content and name, also
with the network id or userid.
2. File storage server ask for the file storage server to persist
the content on the specific network id section(dir).
3. If success in step 3, then the server will respond the client
with the success response and download URL.
4. If fail, such as persist error or file name conflict error, the
server will respond error codes.

Test:
1. Simulate a multi thread client connection, the client could
generate many files upload request simultaneously, to test the
server throughput.

`cd ~/workspace/go_work/src/github.com/lzbj/FileServer/lib/mock/client`
upload a file with the specified network.

`curl -F 'uploadfile=@file.input' -F 'network=97753' http://localhost:9000/storage/upload`

download a file with the specified url:

`curl http://127.0.0.1:9000/storage/download/97753/file.input -o newfile.input`


Async upload and query status

`curl -F 'uploadfile=@file.input' -F 'network=97753' http://localhost:9000/storage/async/upload`

query the upload status.

`curl  http://localhost:9000/status/query/97753/file.input -v`

#
Complex Architecture version:

## Basic features.
1. User upload a file with a file content and name, also 
with the network id or userid.
2. file storage server ask for OLTP(cache) to check if the same
name file already existed under the network id or user ID.
3. If not existed in step 2, the file storage server will ask
the FS layer to persist the content.
4. If success in step 3, the file storage server will ask
for the cache layer to persist the file name or content to a 
cache(which could be a LRU cache). 
5. If success in step 4, will ask the OLTP layer to 
update the record in DB(cache).
6. Then the server will respond the client with the success 
response and download URL.  
