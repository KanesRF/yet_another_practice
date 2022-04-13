## yet_another_practice
# Main idea

This is a simple http server, that receives http POST request wiht string, that contains 
JSON objects, that separated by \n symbol. The server parses it, adds two new params and writes objects to Clickhouse database.

# How to build and run

Just use 
```sh
go build -o server ./server/server.go
./server/server
```

# Examples

There is simple esample of request, used by curl in file 
[request example](testRequest.sh)