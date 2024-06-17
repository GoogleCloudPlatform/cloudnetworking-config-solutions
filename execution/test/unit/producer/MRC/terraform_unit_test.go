// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unittest

import (
	// Package for comparison operations
	compare "cmp"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gruntwork-io/terratest/modules/terraform" // Terraform testing library
	"golang.org/x/exp/slices"                             // Slice manipulation utilities
	"gopkg.in/yaml.v2"                                    // YAML parsing library
)

var (
	projectRoot, _ = filepath.Abs("../../../../")

	// Path to the main Terraform directory for the MRC module.
	terraformDirectoryPath = filepath.Join(projectRoot, "04-producer/MRC")

	// Path to the main Terraform directory for the MRC module.
	configFolderPath = filepath.Join(projectRoot, "test/unit/producer/MRC/config")
)

var (
	// Get project ID
	projectID = "dummy-project-ID"

	// Terraform variables to be passed to the module.
	tfVars = map[string]any{
		"config_folder_path": configFolderPath,
	}
)

// TestInitAndPlanRunWithTfVars tests that Terraform initialization and planning
// succeed with the provided variables. It expects changes (exit code 2) as it's not applying.

func TestInitAndPlanRunWithTfVars(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})

	// Run 'terraform init' and 'terraform plan', get the exit code.
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	want := 2 // Expect changes to be applied
	got := planExitCode

	// Check if the actual exit code matches the expected one.
	if got != want {
		t.Errorf("Test Plan Exit Code = %v, want = %v", got, want)
	}
}

// TestResourcesCount verifies the number of resources to be added by the Terraform plan.
func TestResourcesCount(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})

	// Initialize and create a plan, then parse the resource count.
	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCount := terraform.GetResourceCount(t, planStruct)

	if got, want := resourceCount.Add, 2; got != want { // Expect 2 resources to be added
		t.Errorf("Test Resource Count Add = %v, want = %v", got, want)
	}

	if got, want := resourceCount.Change, 0; got != want { // Expect 0 resource to be changed
		t.Errorf("Test Resource Count Change = %v, want = %v", got, want)
	}

	if got, want := resourceCount.Destroy, 0; got != want { // Expect 0 resources to be destroyed
		t.Errorf("Test Resource Count Destroy = %v, want = %v", got, want)
	}
}

// TestTerraformModuleResourceAddressListMatch verifies that the resources defined
// in the Terraform plan match those specified in YAML configuration files.
func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	expectedModulesAddress := []string{}

	// Read the YAML files in the configuration folder.
	yamlFiles, err := os.ReadDir(configFolderPath)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Extract Redis cluster names from the YAML files.
	for _, file := range yamlFiles {
		if !file.IsDir() {
			// Read YAML content for the cluster name
			yamlData, _ := os.ReadFile(configFolderPath + "/" + file.Name())

			var config struct {
				RedisClusterName string `yaml:"redis_cluster_name"`
			}

			err := yaml.Unmarshal(yamlData, &config)
			if err != nil {
				t.Fatal(err.Error())
			}

			expectedModulesAddress = append(expectedModulesAddress, "google_redis_cluster.cluster-ha[\""+config.RedisClusterName+"\"]") // Use the cluster name
		}
	}

	// Initialize Terraform and generate a plan.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
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
		t.Fatal(err.Error())
	}

	actualModuleAddress := make([]string, 0)
	for _, element := range content.ResourceChangesMap {
		// Check only for the resource address, not module address
		if element.Address != "" && !slices.Contains(actualModuleAddress, element.Address) {
			actualModuleAddress = append(actualModuleAddress, element.Address)
		}
	}

	// Check if any modules are expected
	if len(expectedModulesAddress) > 0 {
		want := expectedModulesAddress
		got := actualModuleAddress
		if !cmp.Equal(got, want, cmpopts.SortSlices(compare.Less[string])) {
			t.Errorf("Test Element Mismatch = %v, want = %v", got, want)
		}
	} else {
		// If no modules expected, check if any actual addresses were found (should be none)
		if len(actualModuleAddress) > 0 {
			t.Errorf("Unexpected module addresses found: %v", actualModuleAddress)
		} else {
			t.Log("No modules expected, and none found in plan.")
		}
	}
}
