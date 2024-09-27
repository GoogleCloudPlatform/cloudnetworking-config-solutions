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
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform" // Terraform testing library
)

var (
	projectRoot, _ = filepath.Abs("../../../../")

	// Path to the main Terraform directory for the GKE module.
	terraformDirectoryPath = filepath.Join(projectRoot, "04-producer/GKE")
)

// TestTerraformConfigValidity checks if the Terraform configuration files are valid.
func TestTerraformConfigValidity(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		NoColor:      true,
	})

	// Initialize Terraform before validating
	terraform.Init(t, terraformOptions)

	// Run 'terraform validate' command and capture output.
	validateOutput := terraform.Validate(t, terraformOptions)

	// Check if the validation output contains "Success! The configuration is valid."
	if strings.Contains(validateOutput, "Success! The configuration is valid.") {
		t.Log("Terraform configuration is valid.")
	} else {
		t.Errorf("Terraform configuration is invalid. Output = %v", validateOutput)
	}
}

// TestTerraformInit checks if the Terraform initialization runs correctly.
func TestTerraformInit(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		NoColor:      true,
	})

	// Run 'terraform init' command.
	initOutput := terraform.Init(t, terraformOptions)

	// Check if the initialization output contains the success message
	if strings.Contains(initOutput, "Terraform has been successfully initialized!") {
		t.Log("Terraform initialization succeeded.")
	} else {
		t.Errorf("Terraform initialization failed. Output = %v", initOutput)
	}
}
