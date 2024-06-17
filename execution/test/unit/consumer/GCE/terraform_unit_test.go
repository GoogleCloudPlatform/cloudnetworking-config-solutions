/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package unittest

// Package for comparison operations
import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform" // Terraform testing library
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var (
	projectRoot, _         = filepath.Abs("../../../../")
	terraformDirectoryPath = filepath.Join(projectRoot, "06-consumer/GCE")
	configFolderPath       = filepath.Join(projectRoot, "test/unit/consumer/GCE/config")
)

var (

	// Terraform variables to be passed to the module.
	tfVars = map[string]any{
		"config_folder_path": configFolderPath,
	}
)

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

	if got, want := resourceCount.Add, 3; got != want { // Expect 3 resources to be added (instance1, instance2, instance3)
		t.Errorf("Test Resource Count Add = %v, want = %v", got, want)
	}
}

func TestTerraformModuleVMResourceAddressListMatch(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})

	localInstanceMap := make(map[string]map[string]string) // Initialize the map
	err := filepath.Walk(configFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			// Read the YAML file
			yamlFile, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Unmarshal YAML into a map
			var instanceData map[string]string
			err = yaml.Unmarshal(yamlFile, &instanceData)
			if err != nil {
				return err
			}

			// Extract the instance name and add to the map
			instanceName := filepath.Base(path) // Assumes filename is the instance name
			instanceName = strings.TrimSuffix(instanceName, filepath.Ext(instanceName))
			localInstanceMap[instanceName] = instanceData
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error reading YAML files: %s", err)
	}

	if len(localInstanceMap) == 0 { // Check if the map is empty
		t.Error("No instances found in YAML files. Make sure the files exist and are correctly formatted.")
	}
	expectedModuleAddresses := []string{} // Start with an empty slice
	for name := range localInstanceMap {
		expectedModuleAddresses = append(expectedModuleAddresses, fmt.Sprintf("module.vm[\"%s\"]", name))
	}

	planStruct := terraform.InitAndPlanAndShow(t, terraformOptions)
	content, err := terraform.ParsePlanJSON(planStruct)
	if err != nil {
		t.Errorf("Error parsing plan JSON: %s", err) // Fail fast if parsing errors occur
	}

	actualModuleAddresses := make([]string, 0)
	for _, element := range content.ResourceChangesMap {
		if strings.HasPrefix(element.ModuleAddress, "module.vm") &&
			!slices.Contains(actualModuleAddresses, element.ModuleAddress) {
			actualModuleAddresses = append(actualModuleAddresses, element.ModuleAddress)
		}
	}

	assert.ElementsMatch(t, expectedModuleAddresses, actualModuleAddresses)
}
