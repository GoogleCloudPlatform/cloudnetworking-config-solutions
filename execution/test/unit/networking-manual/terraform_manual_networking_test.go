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
			"target":                       "projects/xxx-tp/regions/xx-central1/serviceAttachments/gkedpm-xxx",
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
