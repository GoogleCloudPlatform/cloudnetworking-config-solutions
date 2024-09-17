// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integrationtest

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"                        // For deep comparison of slices
	"github.com/gruntwork-io/terratest/modules/terraform" // Terraform testing library
	"github.com/stretchr/testify/assert"                  // Assertion library
)

// Constants for the Terraform directory path and plan file path.
const (
	terraformDirectoryPath = "../../../05-networking-manual/" // Replace with the actual path to your Terraform code directory
	planFilePath           = "./plan"                         // Path where Terraform will save the execution plan
)

// Define the names of the producer SQL instances to be tested with their Service Attachments

var (
	producerInstanceNameWithoutIP = "psc"                                                               // Replace with your actual SQL instance name for testing without IP
	producerInstanceNameWithIP    = "psc-instance"                                                      // Replace with your actual SQL instance name for testing with IP
	targetLink                    = "target-link"                                                       // Define in format projects/project-tp/regions/region/serviceAttachments/unique-id"
	producerInstanceNames         = []string{producerInstanceNameWithoutIP, producerInstanceNameWithIP} // Replace with your actual instance names to be tested in total (2 to be added)
	networkName                   = "default"                                                           // Replace with an existing VPC Network
	subnetworkName                = "default"                                                           // Replace with an existing subnet
	ipAddressLiteral              = "10.128.0.30"                                                       // Replace with an available IP Address
	ipAddressLiteralWithTarget    = "10.128.0.31"                                                       // Replace with an available IP Address
	region                        = "us-central1"                                                       // Replace with your chosen region
)

// Global variable to store the Terraform variables used in tests.
var tfVars map[string]interface{}

// initTfVars initializes the tfVars map with default or environment variable values.
// This function is used to configure the Terraform variables for the tests.
func initTfVars() {
	// Fetch project IDs from environment variables or set defaults.
	endpointProjectID := os.Getenv("TF_VAR_endpoint_project_id")
	producerInstanceProjectID := os.Getenv("TF_VAR_producer_instance_project_id")
	if producerInstanceProjectID == "" {
		producerInstanceProjectID = endpointProjectID // If not set, use the same as endpointProjectID
	}

	// Create an array of endpoint configurations (pscEndpoints) based on the instance names.
	pscEndpoints := make([]interface{}, len(producerInstanceNames))
	for i, name := range producerInstanceNames {
		// Each endpoint configuration is a map with details like project IDs, instance name, subnetwork, network, and IP address (auto-allocated in this case).
		// The subnetwork and network are set to "default" for now, but can be changed to other values if needed.
		pscEndpoints[i] = map[string]interface{}{
			"endpoint_project_id":          endpointProjectID,
			"producer_instance_project_id": producerInstanceProjectID,
			"producer_instance_name":       name,
			"subnetwork_name":              subnetworkName,
			"network_name":                 networkName,
			"ip_address_literal":           "", // This signifies that the IP address will be auto-allocated
		}
	}

	// Assign the pscEndpoints array to the global tfVars map under the key "psc_endpoints".
	tfVars = map[string]interface{}{
		"psc_endpoints": pscEndpoints,
	}
}
func TestInitAndPlanRunWithTfVars(t *testing.T) {
	initTfVars() // Initialize Terraform variables using environment variables or defaults

	// Create Terraform options for initialization and planning.
	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath, // Path to the Terraform configuration directory
		Vars:         tfVars,                 // Variables to pass to Terraform
		Reconfigure:  true,                   // Force re-evaluation of the backend configuration
		Lock:         true,                   // Enable state locking during operations (recommended for parallel runs)
		PlanFilePath: planFilePath,           // File to save the execution plan
		NoColor:      true,                   // Disable colored output for easier parsing
	})

	// Initialize Terraform and generate an execution plan, capturing the exit code.
	planExitCode := terraform.InitAndPlanWithExitCode(t, tfOptions)

	// Define the expected exit code for a successful plan with changes.
	want := 2 // Exit code 2 indicates a successful plan with pending changes.

	// Assert that the actual exit code matches the expected exit code.
	assert.Equal(t, want, planExitCode, "Expected Terraform plan to succeed with changes (exit code 2)")
}

func TestResourcesCount(t *testing.T) {
	initTfVars() // Initialize Terraform variables

	// Create Terraform options for initialization and planning.
	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: planFilePath,
		NoColor:      true,
	})

	// Initialize Terraform and generate an execution plan.
	planStruct := terraform.InitAndPlan(t, tfOptions)

	// Get the resource count from the plan structure.
	resourceCount := terraform.GetResourceCount(t, planStruct)

	// Assert that the plan adds the expected number of resources.
	assert.Equal(t, 4, resourceCount.Add, "TestResourcesCount failed. Expected %d resources to be added, but got %d", 4, resourceCount.Add) // Expecting 4 resources (2 addresses and 2 forwarding rules)

	// Assert that the plan doesn't change any existing resources.
	assert.Zero(t, resourceCount.Change, "TestResourcesCount failed. Expected %d resources to be changed, but got %d", 0, resourceCount.Change)

	// Assert that the plan doesn't destroy any existing resources.
	assert.Zero(t, resourceCount.Destroy, "TestResourcesCount failed. Expected %d resources to be destroyed, but got %d", 0, resourceCount.Destroy)
} // TestPlanFailsWithoutVars tests that the Terraform plan fails when required input variables are missing.
func TestPlanFailsWithoutVars(t *testing.T) {
	// Create Terraform options with default settings, but no variables provided
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath, // Path to Terraform configuration directory
		Reconfigure:  true,                   // Force re-evaluation of backend configuration
		Lock:         true,                   // Enable state locking during operations
		PlanFilePath: planFilePath,           // File to save the execution plan
		NoColor:      true,                   // Disable colored output
	})

	// Attempt to initialize and create a plan, expecting it to fail due to missing variables
	_, err := terraform.InitAndPlanE(t, terraformOptions)

	// Assert that the initialization and planning failed (err is not nil)
	assert.Error(t, err, "Expected Terraform plan to fail due to missing variables")

	// If the planning didn't fail, exit the test
	if err == nil {
		t.Error("Expected Terraform plan to fail due to missing variables, but it succeeded")
	}

	// Get the exit code of the failed plan
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)

	// Assert that the exit code is 1, indicating a failure
	assert.Equal(t, 1, planExitCode, "TestPlanFailsWithoutVars: Expected plan to fail due to missing variables, but got exit code: %v", planExitCode)
}

// TestTerraformModuleResourceAddressListMatch verifies that the Terraform plan output
// includes the expected module addresses for the resources being created.
func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	initTfVars() // Initialize Terraform variables

	// List of expected module addresses
	expectedModuleAddress := []string{
		"module.psc_forwarding_rules",
		"module.psc_forwarding_rules",
		"module.psc_forwarding_rules",
		"module.psc_forwarding_rules",
	}

	// Create Terraform options
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath, // Path to Terraform configuration
		Vars:         tfVars,                 // Variables for Terraform
		Reconfigure:  true,                   // Force re-evaluation of backend configuration
		Lock:         true,                   // Enable state locking
		PlanFilePath: planFilePath,           // File to save the execution plan
		NoColor:      true,                   // Disable colored output
	})

	// Initialize and generate a Terraform execution plan
	terraform.InitAndPlan(t, terraformOptions)

	// Get the plan output as JSON
	planStruct, err := terraform.ShowE(t, terraformOptions)
	if err != nil {
		t.Error(err) // Exit if there's an error getting the plan
	}

	// Parse the plan JSON into a struct for analysis
	content, err := terraform.ParsePlanJSON(planStruct)
	if err != nil {
		t.Error(err) // Exit if parsing fails
	}

	// Extract the actual module addresses from the parsed plan
	actualModuleAddresses := make([]string, 0)
	for _, rc := range content.ResourceChangesMap {
		if rc.ModuleAddress != "" { // Filter for resources within the module
			actualModuleAddresses = append(actualModuleAddresses, rc.ModuleAddress)
		}
	}

	// Sort the slices for comparison
	sort.Strings(actualModuleAddresses)
	sort.Strings(expectedModuleAddress)

	// Assert that the actual module addresses match the expected ones.
	assert.True(t, cmp.Equal(actualModuleAddresses, expectedModuleAddress),
		"TestTerraformModuleResourceAddressListMatch failed.\nActual module addresses: %v\nExpected module addresses: %v", actualModuleAddresses, expectedModuleAddress)
}

// TestPSCForwardingRuleModuleWithProvidedIPAddress tests the Terraform module
// that creates a forwarding rule with a user-specified IP address.
func TestPSCForwardingRuleModuleWithProvidedIPAddress(t *testing.T) {
	// Configure Terraform options for the test, including variables
	tfOptions := configureTerraformOptions(t)

	// Ensure resources are cleaned up after the test
	defer terraform.Destroy(t, tfOptions)

	// Initialize Terraform and apply the configuration
	terraform.InitAndApply(t, tfOptions)

	// Verify the created resources and their outputs
	assertOutputs(t, tfOptions)
}

// configureTerraformOptions configures Terraform options for testing.
// It reads environment variables for project IDs and sets up the required Terraform variables.
func configureTerraformOptions(t *testing.T) *terraform.Options {
	// Retrieve the project ID from environment variables
	endpointProjectID := os.Getenv("TF_VAR_endpoint_project_id")

	// Assert that the project ID is set
	assert.NotEmpty(t, endpointProjectID, "Environment variable 'TF_VAR_endpoint_project_id' must be set")

	// Retrieve or set the producer project ID (if not specified, it defaults to the same as the project ID)
	producerProjectID := os.Getenv("TF_VAR_producer_project_id")
	if producerProjectID == "" {
		producerProjectID = endpointProjectID
	}

	// Set the producer instance name

	// Create a map of Terraform variables
	tfVars := map[string]interface{}{
		"psc_endpoints": []interface{}{
			map[string]interface{}{
				"endpoint_project_id":          endpointProjectID,          // Project ID for the endpoint
				"producer_instance_project_id": producerProjectID,          // Project ID where the SQL instance resides
				"producer_instance_name":       producerInstanceNameWithIP, // Name of the SQL instance
				"subnetwork_name":              subnetworkName,             // Modifiable Subnetwork name for the forwarding rule
				"network_name":                 networkName,                // Modifiable Network name for the forwarding rule
				"ip_address_literal":           ipAddressLiteral,           // Specify the IP address to use
			},
		},
	}

	// Return the Terraform options with default retryable errors handling
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath, // Path to the Terraform code
		Vars:         tfVars,                 // Set the Terraform variables
	})
}

// assertOutputs verifies the outputs of the Terraform module.
func assertOutputs(t *testing.T, tfOptions *terraform.Options) {
	// Get the forwarding rule self-link output
	actualForwardingRuleSelfLinkMap := terraform.OutputMap(t, tfOptions, "forwarding_rule_self_link")

	// Get the IP address output
	actualIPAddressMap := terraform.OutputMap(t, tfOptions, "ip_address_literal")

	// Extract the producer instance name from the Terraform options
	producerInstanceName := tfOptions.Vars["psc_endpoints"].([]interface{})[0].(map[string]interface{})["producer_instance_name"].(string)

	// Construct the expected forwarding rule name
	expectedForwardingRuleName := fmt.Sprintf("psc-forwarding-rule-%s", producerInstanceName)

	// Get the actual forwarding rule self link and extract the name
	actualForwardingRuleSelfLink := actualForwardingRuleSelfLinkMap["0"]
	parts := strings.Split(actualForwardingRuleSelfLink, "/")
	actualForwardingRuleName := parts[len(parts)-1]

	// Get the actual IP address
	actualIPAddress := actualIPAddressMap["0"]

	// Assert that the forwarding rule name matches the expected name
	assert.Equal(t, expectedForwardingRuleName, actualForwardingRuleName, "Forwarding rule name mismatch")

	// Assert that the IP address is not nil
	assert.NotNil(t, actualIPAddress, "IP address is nil")
}

// TestPSCForwardingRuleModuleWithAutoAllocatedIPAddress tests the Terraform module for PSC forwarding rule creation
// when the IP address is NOT explicitly provided (auto-allocated).
func TestPSCForwardingRuleModuleWithAutoAllocatedIPAddress(t *testing.T) {
	// Configure Terraform options, setting the IP address to "" (empty string) for auto-allocation.
	tfOptions := configureTerraformOptionsWithNullIPAddress(t)

	// Defer the destruction of Terraform resources to clean up after the test.
	defer terraform.Destroy(t, tfOptions)

	// Initialize Terraform and apply the configuration to create resources.
	terraform.InitAndApply(t, tfOptions)

	// Assert that the outputs match the expected values.
	assertOutputsForAutoAllocatedIPAddress(t, tfOptions)
}

// configureTerraformOptionsWithNullIPAddress sets up Terraform options with a null (empty) IP address for auto-allocation.
func configureTerraformOptionsWithNullIPAddress(t *testing.T) *terraform.Options {
	// Retrieve the project ID from environment variables, or use a default if not set.
	endpointProjectID := os.Getenv("TF_VAR_endpoint_project_id")
	assert.NotEmpty(t, endpointProjectID, "Environment variable 'TF_VAR_endpoint_project_id' must be set")

	// If the producer project ID is not set, it defaults to the same as the project ID.
	producerProjectID := os.Getenv("TF_VAR_producer_project_id")
	if producerProjectID == "" {
		producerProjectID = endpointProjectID
	}

	// Configure the Terraform variables with the project IDs, instance name, subnetwork, network, and an empty IP address.
	tfVars := map[string]interface{}{
		"psc_endpoints": []interface{}{
			map[string]interface{}{
				"endpoint_project_id":          endpointProjectID,             // Project where the forwarding rule will be created
				"producer_instance_project_id": producerProjectID,             // Project where the Cloud SQL instance resides
				"producer_instance_name":       producerInstanceNameWithoutIP, // Name of the Cloud SQL instance
				"subnetwork_name":              subnetworkName,                // Modifiable Subnetwork name for the forwarding rule and address - ensure this is same region as the cloud sql instance
				"network_name":                 networkName,                   // Modifiable Network name for the forwarding rule
				"ip_address_literal":           "",                            // Empty IP address signifies auto-allocation
			},
		},
	}

	// Return the configured Terraform options with error handling.
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath, // Path to the Terraform code directory
		Vars:         tfVars,                 // Terraform variables
	})
}

// assertOutputsForAutoAllocatedIPAddress verifies the Terraform outputs when an IP address is auto-allocated.
func assertOutputsForAutoAllocatedIPAddress(t *testing.T, tfOptions *terraform.Options) {
	// Retrieve the forwarding rule self-link and IP address outputs from Terraform.
	actualForwardingRuleSelfLinkMap := terraform.OutputMap(t, tfOptions, "forwarding_rule_self_link")
	actualIPAddressMap := terraform.OutputMap(t, tfOptions, "ip_address_literal")

	// Extract the producer instance name from the Terraform options.
	producerInstanceName := tfOptions.Vars["psc_endpoints"].([]interface{})[0].(map[string]interface{})["producer_instance_name"].(string)

	// Construct the expected forwarding rule name based on the instance name.
	expectedForwardingRuleName := fmt.Sprintf("psc-forwarding-rule-%s", producerInstanceName)

	// Get the actual forwarding rule self-link from the output map.
	actualForwardingRuleSelfLink := actualForwardingRuleSelfLinkMap["0"] // Assuming only one instance is created

	// Extract the actual forwarding rule name from the self-link by splitting the URL.
	parts := strings.Split(actualForwardingRuleSelfLink, "/")
	actualForwardingRuleName := parts[len(parts)-1]

	// Extract the actual IP address from the output map.
	actualIPAddress := actualIPAddressMap["0"] // Assuming only one instance is created

	// Assert that the extracted forwarding rule name matches the expected name.
	assert.Equal(t, expectedForwardingRuleName, actualForwardingRuleName, "Forwarding rule name mismatch")

	// Assert that the IP address is not nil (since it should have been auto-allocated).
	assert.NotNil(t, actualIPAddress, "IP address is nil")
}

// TestPSCForwardingRuleModuleWithTarget tests the Terraform module
// that creates a forwarding rule with a user-specified target (service attachment link).
func TestPSCForwardingRuleModuleWithTarget(t *testing.T) {
	// Configure Terraform options for the test, including variables
	tfOptions := configureTerraformOptionsWithTarget(t)

	// Ensure resources are cleaned up after the test
	defer terraform.Destroy(t, tfOptions)

	// Initialize Terraform and apply the configuration
	terraform.InitAndApply(t, tfOptions)

	// Verify the created resources and their outputs
	assertOutputsWithTarget(t, tfOptions)
}

// configureTerraformOptionsWithTarget configures Terraform options for testing with a target specified.
// It reads environment variables for project IDs and sets up the required Terraform variables.
func configureTerraformOptionsWithTarget(t *testing.T) *terraform.Options {
	// Retrieve the project ID from environment variables
	endpointProjectID := os.Getenv("TF_VAR_endpoint_project_id")

	// Assert that the project ID is set
	assert.NotEmpty(t, endpointProjectID, "Environment variable 'TF_VAR_endpoint_project_id' must be set")

	// Retrieve or set the producer project ID (if not specified, it defaults to the same as the project ID)
	producerProjectID := os.Getenv("TF_VAR_producer_project_id")
	if producerProjectID == "" {
		producerProjectID = endpointProjectID
	}

	// Create a map of Terraform variables
	tfVars := map[string]interface{}{
		"psc_endpoints": []interface{}{
			map[string]interface{}{
				"endpoint_project_id":          endpointProjectID,          // Project ID for the endpoint
				"producer_instance_project_id": producerProjectID,          // Project ID where the service attachment resides
				"target":                       targetLink,                 // Service attachment link
				"subnetwork_name":              subnetworkName,             // Modifiable Subnetwork name for the forwarding rule
				"network_name":                 networkName,                // Modifiable Network name for the forwarding rule
				"ip_address_literal":           ipAddressLiteralWithTarget, // Specify the IP address to use
				"region":                       region,                     // Specify the region
			},
		},
	}

	// Return the Terraform options with default retryable errors handling
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath, // Path to the Terraform code
		Vars:         tfVars,                 // Set the Terraform variables
	})
}

// assertOutputsWithTarget verifies the outputs of the Terraform module when using a target.
func assertOutputsWithTarget(t *testing.T, tfOptions *terraform.Options) {
	// Get the forwarding rule self-link output
	actualForwardingRuleSelfLinkMap := terraform.OutputMap(t, tfOptions, "forwarding_rule_self_link")

	// Get the IP address output
	actualIPAddressMap := terraform.OutputMap(t, tfOptions, "ip_address_literal")

	// Extract the target from the Terraform options
	target := tfOptions.Vars["psc_endpoints"].([]interface{})[0].(map[string]interface{})["target"].(string)

	// Construct the expected forwarding rule name (using a "custom-" prefix since no instance name is provided)
	expectedForwardingRuleName := "psc-forwarding-rule-custom-0"

	// Get the actual forwarding rule self link and extract the name
	actualForwardingRuleSelfLink := actualForwardingRuleSelfLinkMap["0"]
	parts := strings.Split(actualForwardingRuleSelfLink, "/")
	actualForwardingRuleName := parts[len(parts)-1]

	// Get the actual IP address
	actualIPAddress := actualIPAddressMap["0"]

	// Assert that the forwarding rule name matches the expected name
	assert.Equal(t, expectedForwardingRuleName, actualForwardingRuleName, "Forwarding rule name mismatch")

	// Assert that the IP address is not nil
	assert.NotNil(t, actualIPAddress, "IP address is nil")

	// Assert that the target in the forwarding rule matches the provided target
	actualTargetMap := terraform.OutputMap(t, tfOptions, "forwarding_rule_target") // Get the output as a map
	actualTarget := actualTargetMap["0"]                                           // Access the target value using the key "0"
	assert.Equal(t, target, actualTarget, "Target mismatch")
}
