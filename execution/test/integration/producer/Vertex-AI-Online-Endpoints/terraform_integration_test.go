// Copyright 2024 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integrationtest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/rand"
	"gopkg.in/yaml.v2"
)

var (
	projectRoot, _ = filepath.Abs("../")

	// Path to the main Terraform directory for the VertexAI module.
	terraformDirectoryPath = filepath.Join(projectRoot, "../../../04-producer/Vertex-AI-Online-Endpoints")

	// Path to the main Terraform directory for the VertexAI module.
	configFolderPath = filepath.Join(projectRoot, "Vertex-AI-Online-Endpoints/config")
	projectID        = os.Getenv("TF_VAR_project_id")
	region           = "us-central1"
	psaRangeName     = "psa-range-cncs-test"
)

type EndpointConfig struct {
	Name        string `yaml:"name"`
	Project     string `yaml:"project"`
	DisplayName string `yaml:"display_name"`
	Description string `yaml:"description"`
	Location    string `yaml:"location"`
	Region      string `yaml:"region"`
	Network     string `yaml:"network"`
}

// TestCreateEndpointWithVPC creates a VPC and then creates an Vertex AI Online Endpoint with the new VPC
func TestCreateEndpointWithVPC(t *testing.T) {

	timestamp := time.Now().Format("20060102150405")
	VPCName := fmt.Sprintf("vpc-%s-%d", timestamp, rand.Intn(100000))

	createVPC(t, projectID, VPCName)
	defer deleteVPC(t, projectID, VPCName)

	createEndpointConfigYAML(t, VPCName, "endpoint_vpc.yaml")

	var (
		tfVars = map[string]interface{}{
			"config_folder_path": configFolderPath,
		}
	)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		Vars:                 tfVars,
		TerraformDir:         terraformDirectoryPath,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})
	// Refresh the Terraform state before applying changes
	terraform.RunTerraformCommand(t, terraformOptions, "refresh")

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Read the YAML file
	yamlConfig, err := readEndpointConfigYAML("endpoint_vpc.yaml")
	if err != nil {
		t.Logf("Error reading YAML config: %v", err)
	}

	validateEndpoints(t, terraformOptions, yamlConfig.Network)
}

// Function to create a YAML config for the Online Endpoint
func createEndpointConfigYAML(t *testing.T, vpcName string, fileName string) {
	t.Log("========= YAML File =========")

	// Generate a unique endpoint name with a timestamp
	endpointName := fmt.Sprintf("vertexai-name-%s-%08d", time.Now().Format("20060102150405"), rand.Intn(100000000))

	endpointConfig := EndpointConfig{
		Name:        endpointName,
		Project:     projectID,
		DisplayName: fmt.Sprintf("vertexai-displayname-%s-%08d", time.Now().Format("20060102150405"), rand.Intn(100000000)),
		Description: "test-description",
		Location:    region,
		Region:      region,
		Network:     fmt.Sprintf("projects/%s/global/networks/%s", getProjectNumber(t, projectID), vpcName),
	}

	yamlData, err := yaml.Marshal(&endpointConfig)
	if err != nil {
		t.Errorf("Error while marshalling %v", err)
	}

	filePath := fmt.Sprintf("%s/%s", configFolderPath, fileName)
	t.Logf("Created YAML config at %s with content:\n%s", filePath, string(yamlData))

	err = os.WriteFile(filePath, []byte(yamlData), 0666)
	if err != nil {
		t.Errorf("Unable to write data into the file %v", err)
	}
}

// readEndpointConfigYAML reads the YAML file and returns the EndpointConfig struct
func readEndpointConfigYAML(fileName string) (*EndpointConfig, error) {
	filePath := fmt.Sprintf("%s/%s", configFolderPath, fileName)
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read YAML file: %v", err)
	}

	var config EndpointConfig
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal YAML data: %v", err)
	}

	return &config, nil
}

// validateEndpoints validates the endpoints created by the Terraform module
func validateEndpoints(t *testing.T, terraformOptions *terraform.Options, expectedNetwork string) {
	maxRetries := 10
	retryInterval := 10 * time.Second
	endpointsReady := false

	for i := 0; i < maxRetries; i++ {
		endpointOutputValue := terraform.OutputJson(t, terraformOptions, "endpoint_configurations")

		if !gjson.Valid(endpointOutputValue) {
			t.Errorf("Error parsing output, invalid json: %s", endpointOutputValue)
			continue
		}

		result := gjson.Parse(endpointOutputValue)
		endpointsReady = true

		result.ForEach(func(key, value gjson.Result) bool {
			expectedDisplayName := key.String()

			// Verify Endpoint Display Name
			gotDisplayName := value.Get("display_name").String()
			if gotDisplayName != expectedDisplayName {
				t.Errorf("Endpoint '%s' has invalid display_name: got %s, want %s", expectedDisplayName, gotDisplayName, expectedDisplayName)
			} else {
				t.Logf("Endpoint has valid display_name")
			}

			// Verify Endpoint Network (using the expectedNetwork argument)
			gotNetwork := value.Get("network").String()
			if gotNetwork != expectedNetwork {
				t.Errorf("Endpoint '%s' has invalid network: got %s, want %s", expectedDisplayName, gotNetwork, expectedNetwork)
			} else {
				t.Logf("Endpoint has valid network")
			}

			return true
		})

		if endpointsReady {
			break
		} else {
			t.Errorf("Endpoints not ready after %d retries", maxRetries)
			time.Sleep(retryInterval)
		}
	}
}

// getProjectNumber gets the Project Number for the Endpoint configuration
func getProjectNumber(t *testing.T, projectID string) string {
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"projects", "describe", projectID, "--format=value(projectNumber)"},
	}
	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Logf("Error getting project number for project ID %s: %s", projectID, err)
	}
	return output
}

// createVPC creates a VPC and required subnet
func createVPC(t *testing.T, projectID string, networkName string) {
	text := "compute"

	// Create VPC
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", networkName, "--project=" + projectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom", "--verbosity=none"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Logf("===Error %s Encountered while executing %s", err, text)
	}

	// Wait for VPC creation to propagate
	time.Sleep(60 * time.Second)

	// Create Subnet
	subnetName := fmt.Sprintf("%s-subnet", networkName)
	cmd = shell.Command{
		Command: "gcloud",
		Args: []string{
			text, "networks", "subnets", "create", subnetName,
			"--project=" + projectID,
			"--network=" + networkName,
			"--region=" + region,
			"--range=10.0.0.0/24",
		},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Logf("===Error %s Encountered while executing %s", err, text)
	}

	// Enable Private Service Access for the VPC
	createInternalIPRange(t, projectID, networkName)
	enablePrivateServiceAccess(t, projectID, networkName, psaRangeName)
}

// createInternalIPRange creates an IP Range for testing
func createInternalIPRange(t *testing.T, projectID, networkName string) {
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"compute", "addresses", "create", psaRangeName,
			"--global",
			"--project=" + projectID,
			"--network=" + networkName,
			"--purpose=VPC_PEERING",
			"--prefix-length=24", // Adjust prefix length if needed
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error creating internal IP range: %v", err)
	}
}

// enablePrivateServiceAccess creates peering for VPC network
func enablePrivateServiceAccess(t *testing.T, projectID, networkName, psaRangeName string) {
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"services", "vpc-peerings", "connect",
			"--service=servicenetworking.googleapis.com",
			"--project=" + projectID,
			"--network=" + networkName,
			"--ranges=" + psaRangeName,
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error creating PSA connection: %v", err)
	}
}

// deletePSARange deletes the PSA range after testing
func deletePSARange(t *testing.T, projectID string, psaRangeName string) {
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"compute", "addresses", "delete", psaRangeName,
			"--project=" + projectID,
			"--global",
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error deleting internal IP range: %v", err)
	}
}

// deletePSAConnection deletes the PSA connection after testing
func deletePSAConnection(t *testing.T, projectID string, networkName string) {
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"services", "vpc-peerings", "delete", "--service=servicenetworking.googleapis.com",
			"--network=" + networkName,
			"--project=" + projectID,
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error deleting PSA Connection: %v", err)
	}
}

// deleteVPC deletes the VPC and subnets after testing
func deleteVPC(t *testing.T, projectID string, networkName string) {
	text := "compute"
	time.Sleep(60 * time.Second)

	deletePSARange(t, projectID, psaRangeName)
	deletePSAConnection(t, projectID, networkName)

	// Delete Subnet
	subnetName := fmt.Sprintf("%s-subnet", networkName)
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			text, "networks", "subnets", "delete", subnetName,
			"--project=" + projectID,
			"--region=" + region,
			"--quiet",
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}

	// Delete VPC
	time.Sleep(120 * time.Second)
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "delete", networkName, "--project=" + projectID, "--quiet"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
}

func TestMain(m *testing.M) {
	if err := cleanupYAMLFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Error cleaning up YAML files: %v\n", err)
		// Allow tests to run, but signal a setup failure
		m.Run()
		os.Exit(1) // Exit after tests complete
	}

	os.Exit(m.Run())
}

// cleanupYAMLFiles cleans up unused old YAML files.
func cleanupYAMLFiles() error {
	files, err := filepath.Glob(filepath.Join(configFolderPath, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to list YAML files: %v", err)
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("failed to remove file %s: %v", file, err)
		}
	}

	return nil
}
