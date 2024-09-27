// Copyright 2024 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// 	Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unittest

import (
	compare "cmp"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"golang.org/x/exp/slices"
)

var (
	projectRoot, _ = filepath.Abs("../")

	// Path to the main Terraform directory for the VertexAI module.
	terraformDirectoryPath = filepath.Join(projectRoot, "../../../04-producer/Vertex-AI-Online-Endpoints")

	// Path to the main Terraform directory for the VertexAI module.
	configFolderPath = filepath.Join(projectRoot, "/Vertex-AI-Online-Endpoints/config")
)

var (
	tfVars = map[string]any{
		"config_folder_path": configFolderPath,
	}
	// used to validate an expected error code if a wrong configuration file is provided.
	invalidTFVars = map[string]any{
		"config_folder_path": configFolderPath,
		"network":            "test-network-name", // giving network name, rather than self link
	}
)

/*
TestInitAndPlanRunWithTfVars performs sanity check to ensure the terraform init &&
terraform plan is executed successfully and returns a valid Succeeded run code.
*/
func TestInitAndPlanRunWithTfVars(t *testing.T) {
	/*
	 0 = Succeeded with empty diff (no changes)
	 1 = Error
	 2 = Succeeded with non-empty diff (changes present)
	*/
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	want := 2
	got := planExitCode
	if got != want {
		t.Errorf("Test Plan Exit Code = %v, want = %v", got, want)
	}
}

/*
TestInitAndPlanRunWithInvalidTfVarsExpectFailureScenario performs test runs with invalid tfvars file
to ensure the terraform init && terraform plan is executed unsuccessfully and returns an expected error run code.
*/
func TestInitAndPlanRunWithInvalidTfVarsExpectFailureScenario(t *testing.T) {
	/*
	 0 = Succeeded with empty diff (no changes)
	 1 = Error
	 2 = Succeeded with non-empty diff (changes present)
	*/
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         invalidTFVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	want := 1
	got := planExitCode
	if !cmp.Equal(got, want) {
		t.Errorf("Test Plan Exit Code = %v, want = %v", got, want)
	}
}

/*
TestResourcesCount performs validation to verify number of  resources created, deleted and updated.
*/
func TestResourcesCount(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCount := terraform.GetResourceCount(t, planStruct)
	if got, want := resourceCount.Add, 1; got != want {
		t.Errorf("Test Resource Count Add = %v, want = %v", got, want)
	}
	if got, want := resourceCount.Change, 0; got != want {
		t.Errorf("Test Resource Count Change = %v, want = %v", got, want)
	}
	if got, want := resourceCount.Destroy, 0; got != want {
		t.Errorf("Test Resource Count Destroy = %v, want = %v", got, want)
	}
}

/*
TestTerraformModuleResourceAddressListMatch compares and verifies the
list of resources, modules created by the terraform solution.
*/
func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	expectedModulesAddress := []string{"module.vertex_endpoints[\"<endpoint-display-name>\"]"}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planStruct := terraform.InitAndPlanAndShow(t, terraformOptions)
	content, err := terraform.ParsePlanJSON(planStruct)
	if err != nil {
		t.Errorf("Error parsing plan JSON : %v", err.Error())
	}
	actualModuleAddress := make([]string, 0)
	for _, element := range content.ResourceChangesMap {
		if element.ModuleAddress != "" && !slices.Contains(actualModuleAddress, element.ModuleAddress) {
			actualModuleAddress = append(actualModuleAddress, element.ModuleAddress)
		}
	}
	want := expectedModulesAddress
	got := actualModuleAddress
	if !cmp.Equal(got, want, cmpopts.SortSlices(compare.Less[string])) {
		t.Errorf("Test Element Mismatch = %v, want = %v", got, want)
	}
}
