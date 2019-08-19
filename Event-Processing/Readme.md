# ToDo API

This Directory contains the implementations of the Event Processing Use-Case.

## Running the Tests

To run the tests, esure the Go SDK is installed. After that the tests can be executed by running the following command in this directory:

```bash
go run *.go -endpoint <API Endpoint> -delay <Delay> -count <Count>
```

The tests assume the implemeńtation you want to test is already deployed. To do so follow the instructions in the directories of the
implementations.

The parameters of the test have the following purposes:

- `endpoint`: Defines the path to the Exposed api. E.g. `https://5f4snin1r2.execute-api.us-east-1.amazonaws.com/dev` on AWS. It is important to only keep the part that is always identical between the functions, since the test will use the url to build the urls to the specific functions. This parameter is **required**.
- `delay`: Defines the delay between the insertion and the validation in seconds. Setting this is **optional** if it is unset the default of 90 (Seconds) will be used.
- `count`: The number of insertions the test should perform for each event type. This is also **optional** the default is 50.
