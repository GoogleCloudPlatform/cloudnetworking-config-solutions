<<<<<<< PATCH SET (ecf159 Unit testing for SCP addition to the networking stage. Major)
## Testing Your Terraform Infrastructure

This guide provides instructions for running unit and integration tests to validate your Terraform resources in a Google Cloud environment.

### Prerequisites

- Go (Golang) Installed: Make sure you have Go installed on your system. You can download it from the official website: https://go.dev/
- Terraform Project: Have your Terraform project set up and ready for testing.
- Test Files: The unit and integration test files should be located within their respective unit and integration subdirectories for each stage.
- Permission: you should have the neccessary permissions required to run each stage to be able to test it.

### Unit Tests

Unit tests verify individual Terraform resources in isolation.

#### Running All Unit Tests

To run all the tests/functions under the unit testing directory for all terraform resources created in respective stages, please follow these steps:

1. Navigate to the desired stage directory:

```
cd STAGE_NAME (such as networking)
```

2. Initialize a Go module for testing (if not already done):

```
go mod init test_file_name 
```
3. Ensure dependencies are up-to-date:

```
go mod tidy
```

4. Execute all unit tests and generate a summary:

```
go test -v -json ./... | ./test-summary**
```

**Note:** [test-summary](https://pkg.go.dev/gocloud.dev/internal/testing/test-summary) is used to provide summary of the test results.

#### Running Specific Unit Tests

To run specific tests/functions under the unit testing directory for the different stages, please follow these steps:

1. Navigate to the desired stage directory:

```
cd unit/STAGE_NAME
```

2. Initialize a Go module for testing (if not already done):

```
go mod init test
```

3. Ensure dependencies are up-to-date:

```
go mod tidy
```

4. Execute all unit tests and generate a summary:

```
go test -timeout 15m -v
```

Example : Here is an example demonstrating how to execute a unit test for the networking stage:

```
cd unit/networking
go mod init test
go mod tidy
go test -timeout 30m -v
```

### Integration Testing

Integration tests verify the interaction between multiple Terraform resources.

#### Running All Integration Tests

To run all the tests/functions under the integration testing directory for all terraform resources created in respective stages, please follow these steps:

1. Navigate to the desired stage directory:

```
cd STAGE_NAME (such as networking)
```

2. Initialize a Go module for testing (if not already done):

```
go mod init test_file_name 
```
3. Ensure dependencies are up-to-date:

```
go mod tidy
```

4. Execute all integration tests and generate a summary:

```
go test -v -json ./... | ./test-summary**
```

**Note:** [test-summary](https://pkg.go.dev/gocloud.dev/internal/testing/test-summary) is used to provide summary of the test results.


#### Running Specific Integration Tests

To run specific tests/functions under the integration testing directory for the different stages, please follow these steps:

1. Navigate to the desired stage directory:

```
cd unit/STAGE_NAME
```

2. Initialize a Go module for testing (if not already done):

```
go mod init test
```

3. Ensure dependencies are up-to-date:

```
go mod tidy
```

4. Execute all unit tests and generate a summary:

```
go test -timeout 15m -v
```

#### Important Notes

- `test-summary`: The test-summary tool is not part of the Go standard library. Ensure you have it installed.
- Timeouts: Adjust timeout values (-timeout) based on the expected execution time of your tests.
=======
>>>>>>> BASE      (703e7f Adding SCP to the networking implementation.)
