## Testing

The following sections describe how the modules can be tested in a Google Cloud Platform environment. There would be unit & integration tests for each of the stages(organization, networking, security etc) in this directory.

### Tests - Environment Variables

While running the integration tests locally (or in your development machines) make sure you have declared following as the environment variables.

```
export TF_VAR_project_id=<PROJECT_ID>
```

### Unit Testing

#### Running all the unit tests

To run all the tests/functions under the unit testing directory for all terraform resources created in respective stages, please follow these steps:

1. cd REPO_NAME
2. go mod init test
3. go mod tidy
4. **go test -v -json ./... | ./test-summary**

    **Note :** [test-summary](https://pkg.go.dev/gocloud.dev/internal/testing/test-summary) is used to provide summary of the test results.

#### Running specific tests

To run specific tests/functions under the unit testing directory for the different stages, please follow these steps:

1. cd /test/unit/STAGE_NAME
2. go mod init test
3. go mod tidy
4. **go test -timeout=15m -v**

    **e.g.** Here is an example demonstrating how to execute a unit test for the networking stage:
    ```
    cd /test/unit/networking
    go mod init test
    go mod tidy
    go test -timeout 30m -v
    ```

### Integration Testing

#### Running all the integration tests

To run all the tests/functions under the integration testing directory for the example, please follow these steps:

1. cd REPO_NAME
2. go mod init test
3. go mod tidy
4. **go test -v -json -timeout 60m ./... | ./test-summary**

    **Note :** [test-summary](https://pkg.go.dev/gocloud.dev/internal/testing/test-summary) is used to provide summary of the test results.


#### Running specific tests

To run specific tests/functions under the integration testing directory for the different stages, please follow these steps:

1. /test/integration/STAGE_NAME
2. go mod init test
3. go mod tidy
4. **go test -timeout 60m -v**

    **e.g.** Here is an example demonstrating how to execute a integration test for the networking stage:

    ```
    cd /test/integration/networking
    go mod init test
    go mod tidy
    go test -timeout=60m -v
    ```
