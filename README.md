# api-simulator

## how to use it

Add your config inside reources/{your-api-name}/

run the server with go run main.go

in another terminal run

curl http://localhost:1949/{your-defined-api-path} 

## an example:

curl http://localhost:1949/example.com/say-hello

curl http://localhost:1949/blog/api/v1/posts


## todo

- Enable using variable in json: eg: `{ "message": "Hello {{macro()}}" }` should return the result of macro()
- Add proxy support
- Add variable that simulates real environment. eg delays and timeouts.
