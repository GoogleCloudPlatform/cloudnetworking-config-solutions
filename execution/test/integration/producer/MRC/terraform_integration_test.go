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
package integrationtest

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

// Test configuration (adjust as needed)
var (
	projectRoot, _ = filepath.Abs("../../../../")
	// Path to the Terraform module directory.
	terraformDirectoryPath = filepath.Join(projectRoot, "04-producer/MRC")
	// Path to the folder containing YAML configuration files.
	configFolderPath = filepath.Join(projectRoot, "test/integration/producer/MRC/config")
)

var (
	projectID    = os.Getenv("TF_VAR_project_id")
	region       = "us-central1"
	instanceName = fmt.Sprintf("mrc-%d", rand.Int())
	networkName  = fmt.Sprintf("vpc-%s-test", instanceName)
	networkID    = fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName)
)

type MRCStruct struct {
	InstanceName string `yaml:"redis_cluster_name"`
	ProjectID    string `yaml:"project_id"`
	NetworkID    string `yaml:"network_id"`
	Region       string `yaml:"region"`
}

// GetFirstNonEmptyEnvVarOrUseDefault retrieves the first non-empty environment variable
// from the provided list, or falls back to a default value if none are set.
func TestCreateMRC(t *testing.T) {
	// Initialize a MRC config YAML file to be tested.
	createConfigYAML(t)

	var (
		tfVars = map[string]any{
			"config_folder_path": configFolderPath,
			"shard_count":        3,
			"replica_count":      1, // Example: Assuming replica_count is a variable
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

	// Create VPC, subnet, and service connection policy
	createVPC(t, projectID, networkName)
	time.Sleep(60 * time.Second)

	// Delete VPC, subnet, and service connection policy
	defer deleteVPC(t, projectID, networkName)

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for the MRC cluster to become available with retries
	maxRetries := 10
	retryInterval := 10 * time.Second
	instanceReady := false // Flag to track if instance is ready

	for i := 0; i < maxRetries; i++ {

		MRCOutputValue := terraform.OutputJson(t, terraformOptions, "redis_cluster_details")

		if !gjson.Valid(MRCOutputValue) {
			t.Errorf("Error parsing output, invalid json: %s", MRCOutputValue)
			continue // Skip to next retry if output is invalid
		}

		result := gjson.Parse(MRCOutputValue)
		instanceReady = true // Assume instance is ready before checking

		// Iterate over each MRC instance details within the redis_cluster_details output
		result.ForEach(func(key, value gjson.Result) bool {
			// Extract the instance name from the key
			instanceName := key.String()

			// 1. Verify MRC Cluster Name
			gotName := value.Get("name").String()
			if gotName != instanceName {
				t.Errorf("MRC Cluster '%s' has invalid name: got %s, want %s", instanceName, gotName, instanceName)
			}

			// 2. Check if the cluster is in the ACTIVE state
			gotStatus := value.Get("state").String()
			if gotStatus != "ACTIVE" {
				t.Logf("MRC Cluster '%s' not yet ACTIVE (current state: %s). Retrying...", instanceName, gotStatus)
				instanceReady = false
				return false // Stop iterating if an instance is not ready
			}

			// 3. Verify Network ID using gcloud command
			expectedNetworkID := value.Get("network").String()
			cmd := shell.Command{
				Command: "gcloud",
				Args:    []string{"redis", "clusters", "describe", instanceName, "--project=" + projectID, "--region=" + region, "--format=json", "--verbosity=none", "--quiet"},
			}
			output, err := shell.RunCommandAndGetOutputE(t, cmd)
			if err != nil {
				t.Errorf("Error running gcloud command: %s", err)
				instanceReady = false
				return false
			}

			actualNetwork := gjson.Get(output, "pscConnections.0.network").String()

			if actualNetwork != expectedNetworkID {
				t.Errorf("MRC Cluster '%s' has invalid network ID: got %s, want %s", instanceName, actualNetwork, expectedNetworkID)
				instanceReady = false
				return false
			}

			// 4. Verify Shard Count
			gotShardCount := value.Get("shard_count").Int()
			if gotShardCount != 3 {
				t.Errorf("MRC Cluster '%s' has invalid shard count: got %d, want 3", instanceName, gotShardCount)
			}

			// 5. Verify Replica Count
			gotReplicaCount := value.Get("replica_count").Int()
			if gotReplicaCount != 1 {
				t.Errorf("MRC Cluster '%s' has invalid replica count: got %d, want 1", instanceName, gotReplicaCount)
			}
			return true // Continue iterating to the next instance
		})

		if instanceReady {
			break
		} else {
			time.Sleep(retryInterval)
		}
	}
}

/*
createVPC creates the VPC, subnet, and service connection policy before the test execution.
*/
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

	// Create Subnet
	subnetName := fmt.Sprintf("%s-subnet", networkName)
	subnetID := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s/subnetworks/%s", projectID, region, subnetName)
	cmd = shell.Command{ // Re-use cmd variable
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

	// Create Service Connection Policy
	policyName := fmt.Sprintf("%s-policy", networkName)
	cmd = shell.Command{ // Re-use cmd variable again
		Command: "gcloud",
		Args: []string{
			"network-connectivity", "service-connection-policies", "create", policyName,
			"--network=" + networkName,
			"--project=" + projectID,
			"--region=" + region,
			"--psc-connection-limit=5",
			"--service-class=gcp-memorystore-redis",
			"--subnets=" + subnetID,
		},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Logf("===Error %s Encountered while executing %s", err, text)
	}
}

/*
deleteVPC deletes the VPC, subnet, and service connection policy after the test.
*/
func deleteVPC(t *testing.T, projectID string, networkName string) {
	text := "compute"
	time.Sleep(60 * time.Second)

	// Delete Service Connection Policy
	policyName := fmt.Sprintf("%s-policy", networkName)
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"network-connectivity", "service-connection-policies", "delete", policyName,
			"--region=" + region,
			"--project=" + projectID,
			"--quiet", // Suppress prompts for confirmation
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}

	// Delete Subnet
	subnetName := fmt.Sprintf("%s-subnet", networkName)
	cmd = shell.Command{
		Command: "gcloud",
		Args: []string{
			text, "networks", "subnets", "delete", subnetName,
			"--project=" + projectID,
			"--region=" + region,
			"--quiet",
		},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}

	// 3. Delete VPC
	time.Sleep(120 * time.Second) // Wait for firewall deletion to complete
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "delete", networkName, "--project=" + projectID, "--quiet"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
}

/*
createConfigYAML is a helper function which creates the configigration YAML file
for an MRC instance.
*/
func createConfigYAML(t *testing.T) {
	t.Log("========= YAML File =========")
	instance1 := MRCStruct{
		InstanceName: instanceName,
		ProjectID:    projectID,
		NetworkID:    networkID,
		Region:       region,
	}

	yamlData, err := yaml.Marshal(&instance1)
	if err != nil {
		t.Errorf("Error while marshallaing %v", err)
	}
	filePath := fmt.Sprintf("%s/%s", "config", "instance1.yaml")
	t.Logf("Created YAML config at %s with content:\n%s", filePath, string(yamlData))

	err = os.WriteFile(filePath, []byte(yamlData), 0666)
	if err != nil {
		t.Errorf("Unable to write data into the file %v", err)
	}
}
