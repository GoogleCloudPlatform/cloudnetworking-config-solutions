// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unittest

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

const (
	terraformDirectoryPath = "../../../05-networking-manual/"
	planFilePath           = "terraform.tfplan"
)

// tfVars: Define input variables for your Terraform module as a map.
var tfVars = map[string]interface{}{
	// psc_endpoints:  A list of endpoint configurations. Each endpoint is a map with properties:
	"psc_endpoints": []interface{}{}, // This will be populated dynamically based on env vars
}

// initTfVars initializes the tfVars map with project and endpoint data from environment variables
func initTfVars() {
	// Directly define project IDs within the function
	projectID := "your-project-id"
	producerProjectID := "your-project-id"

	endpoints := []map[string]interface{}{
		{
			"producer_instance_name":       "sql",
			"subnetwork_name":              "default",
			"network_name":                 "default",
			"ip_address_literal":           "10.128.0.5",
			"endpoint_project_id":          projectID,
			"producer_instance_project_id": producerProjectID,
		},
		{
			"producer_instance_name":       "sql-1",
			"subnetwork_name":              "default",
			"network_name":                 "default",
			"ip_address_literal":           nil,
			"endpoint_project_id":          projectID,
			"producer_instance_project_id": producerProjectID,
		},
		{
			"producer_instance_name":       "sql2",
			"subnetwork_name":              "subnetwork",
			"network_name":                 "network",
			"ip_address_literal":           nil,
			"endpoint_project_id":          projectID,
			"producer_instance_project_id": producerProjectID,
		},
	}

	// Update tfVars (no change here, as the logic is simpler now)
	tfVars["psc_endpoints"] = make([]interface{}, len(endpoints))
	for i, endpoint := range endpoints {
		tfVars["psc_endpoints"].([]interface{})[i] = endpoint
	}
}

func TestInitAndValidate(t *testing.T) {
	initTfVars()
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		NoColor:      true,
	})

	// Initialization and Validation
	_, err := terraform.InitAndValidateE(t, terraformOptions)
	if err != nil {
		t.Errorf("Failed to initialize and validate Terraform: %v", err)
	}
}

// TestResourcesCount verifies that the Terraform plan is creating the expected number of resources.
// Specifically, it checks that there are two resources to be added (one address and one forwarding rule),
// no resources to be changed, and no resources to be destroyed.
// func TestResourcesCount(t *testing.T) {

// 	// Set up Terraform options with the test directory, variables, and other settings.
// 	// The terraform.WithDefaultRetryableErrors helper is used to handle common transient errors
// 	// that might occur during Terraform operations.
// 	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
// 		TerraformDir: terraformDirectoryPath, // The path to your Terraform configuration directory
// 		Vars:         tfVars,                 // The input variables for Terraform
// 		Reconfigure:  true,                   // Force Terraform to re-evaluate the configuration
// 		Lock:         true,                   // Enable state locking (recommended for parallel runs)
// 		PlanFilePath: "./plan",               // The file path where the plan will be saved
// 		NoColor:      true,                   // Disable color output for easier parsing
// 	})

// 	// Initialize Terraform and generate an execution plan.
// 	planStruct := terraform.InitAndPlan(t, terraformOptions)

// 	// Get the resource count from the plan. This provides details about the number
// 	// of resources that will be added, changed, or destroyed by the plan.
// 	resourceCount := terraform.GetResourceCount(t, planStruct)

// 	// Verify that the plan will add two resources (as expected for one instance with one reserved IP)
// 	if got, want := resourceCount.Add, 4; got != want {
// 		t.Errorf("TestResourcesCount failed. Expected %d resources to be added, but got %d", want, got)
// 	}

// 	// Verify that the plan will not change any existing resources.
// 	if got, want := resourceCount.Change, 0; got != want {
// 		t.Errorf("TestResourcesCount failed. Expected %d resources to be changed, but got %d", want, got)
// 	}

// 	// Verify that the plan will not destroy any existing resources.
// 	if got, want := resourceCount.Destroy, 0; got != want {
// 		t.Errorf("TestResourcesCount failed. Expected %d resources to be destroyed, but got %d", want, got)
// 	}
// }

// TestPlanFailsWithoutValidResources verifies that the Terraform initialization and planning process
// fails when required input variables are not provided.
// This is a negative test case, and it expects the plan to fail with an error (exit code 1).
func TestPlanFailsWithoutValidResources(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: planFilePath,
		NoColor:      true,
	})

	_, err := terraform.InitAndPlanE(t, terraformOptions)
	if err == nil {
		t.Errorf("Expected Terraform plan to fail due to missing variables, but it succeeded")
	}

	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)

	want := 1
	if got := planExitCode; got != want {
		t.Errorf("TestPlanFailsWithoutVars: Expected plan to fail due to missing variables, but got exit code: %v", got)
	}
}

// TestTerraformModuleResourceAddressListMatch verifies that the Terraform plan output
// includes the expected module addresses for the resources being created.
// In this case, we expect two instances of "module.psc_forwarding_rules" because the module creates
// both a compute address and a forwarding rule for each instance defined in the input variables.
// func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
// 	// Define the expected module addresses for the resources in the plan output.
// 	expectedModuleAddress := []string{"module.psc_forwarding_rules", "module.psc_forwarding_rules", "module.psc_forwarding_rules", "module.psc_forwarding_rules"}

// 	// Configure Terraform options:
// 	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
// 		TerraformDir: terraformDirectoryPath, // Path to the Terraform configuration directory
// 		Vars:         tfVars,                 // Input variables for Terraform
// 		Reconfigure:  true,                   // Force re-evaluation of the configuration
// 		Lock:         true,                   // Enable state locking
// 		PlanFilePath: planFilePath,           // Save the plan to a file
// 		NoColor:      true,                   // Disable colored output for easier parsing
// 	})

// 	// Initialize Terraform and generate an execution plan.
// 	terraform.InitAndPlan(t, terraformOptions)

// 	// Get the plan output as JSON. The ShowE function provides a detailed error if any issues occur.
// 	planStruct, err := terraform.ShowE(t, terraformOptions)
// 	if err != nil {
// 		t.Errorf(err) // If there's an error getting the plan, fail the test immediately
// 	}

// 	// Parse the JSON output into a Terraform plan structure.
// 	content, err := terraform.ParsePlanJSON(planStruct)
// 	if err != nil {
// 		t.Errorf(err) // If parsing fails, fail the test immediately
// 	}

// 	// Collect the actual module addresses from the parsed plan.
// 	actualModuleAddresses := make([]string, 0) // Initialize an empty slice for actual module addresses
// 	for _, rc := range content.ResourceChangesMap {
// 		if rc.ModuleAddress != "" { // Check if resource is part of a module
// 			actualModuleAddresses = append(actualModuleAddresses, rc.ModuleAddress)
// 		}
// 	}

// 	// Sort both slices to ensure consistent order for comparison.
// 	sort.Strings(actualModuleAddresses)
// 	sort.Strings(expectedModuleAddress)

// 	// Compare the actual and expected module addresses. If they don't match, fail the test
// 	// and print the mismatched elements.
// 	if !cmp.Equal(actualModuleAddresses, expectedModuleAddress) {
// 		t.Errorf("TestTerraformModuleResourceAddressListMatch failed.\nActual module addresses: %v\nExpected module addresses: %v", actualModuleAddresses, expectedModuleAddress)
// 	}
// }
