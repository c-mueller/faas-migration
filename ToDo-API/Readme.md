# ToDo API

This Directory contains the implementations of the ToDo API on Microsoft Azure. The other implementations are based on Go and can be found under
https://github.com/c-mueller/faas-migration-go

## Running the Tests

To run the tests, esure the Go SDK is installed. After that the tests can be executed by running the following command in this directory:
```
go run *.go -endpoint <API Endpoint>
```

The tests assume the implemeńtation you want to test is already deployed. To do so follow the instructions in the directories of the
implementations.

## API Specification

The ToDo API Consists of the following functions:

- `POST /put`: Submit an Item to the API
- `GET /get`: Get an Item
- `GET /lst`: List all Items
- `POST /done`: Mark an Item as Done
- `POST /del`: Delete an Item

All Methods produce JSON responses and accept JSON requests.

### `/put` Function

The `/put` function consumes an ItemCreation Request

### ToDo Item


### ItemCreation Request
