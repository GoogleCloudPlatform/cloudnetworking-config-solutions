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

package integrationtest

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform" // Correct import
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

// Test configuration (adjust as needed)
var (
	projectRoot, _         = filepath.Abs("../../../../")
	terraformDirectoryPath = filepath.Join(projectRoot, "06-consumer/GCE")
	configFolderPath       = filepath.Join(projectRoot, "test/integration/consumer/GCE/config")
)

var (
	projectID    = os.Getenv("TF_VAR_project_id")
	instanceName = fmt.Sprintf("gce-%d", rand.Int())
	region       = "us-central1"
	zone         = "us-central1-a"
	networkName  = fmt.Sprintf("vpc-%s-test", instanceName)
	networkID    = fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName)
	subnetworkID = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s-subnet", projectID, region, networkName)
)

// VMInstanceConfig struct
type VMInstanceConfig struct {
	Name       string `yaml:"name"`
	ProjectID  string `yaml:"project_id"`
	Region     string `yaml:"region"`
	Zone       string `yaml:"zone"`
	Image      string `yaml:"image"`
	Network    string `yaml:"network"`
	Subnetwork string `yaml:"subnetwork"`
}

func TestCreateVMInstances(t *testing.T) {
	createConfigYAML(t) // Use the updated createConfigYAML for GCE

	// Terraform Variables (GCE-Specific)
	tfVars := map[string]any{
		"config_folder_path": configFolderPath,
	}

	// Terraform Options
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		Vars:                 tfVars,
		TerraformDir:         terraformDirectoryPath,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})

	// Create VPC and Subnet Before Applying Terraform
	createVPC(t, projectID, networkName)
	time.Sleep(60 * time.Second)

	// Apply Terraform
	terraform.InitAndApply(t, terraformOptions)

	// Get Instance Information from Terraform Output
	vmInstancesOutput := terraform.OutputJson(t, terraformOptions, "vm_instances")
	vmInstances := gjson.Parse(vmInstancesOutput).Map()

	maxRetries := 5
	retryInterval := 15 * time.Second

	// Wait for Instances to be Running & Verify Configuration
	for k, instanceDetails := range vmInstances { // Iterate over keys and values
		instanceName := instanceDetails.Get("name").String() // Extract the name from the object
		zone := strings.Split(k, "/")[2]                     // Extract the zone from the key

		for i := 0; i < maxRetries; i++ {
			gcloudOutput := shell.RunCommandAndGetOutput(t, shell.Command{
				Command: "gcloud",
				Args:    []string{"compute", "instances", "describe", instanceName, "--zone", zone, "--project", projectID, "--format", "json"},
			})
			status := gjson.Get(gcloudOutput, "status").String()

			if status == "RUNNING" {
				// Verify Instance Configuration (against YAML)
				yamlFile, err := os.ReadFile(filepath.Join(configFolderPath, "instance1.yaml"))
				if err != nil {
					t.Errorf("Error reading YAML file: %s", err)
					break
				}

				var expectedInstance VMInstanceConfig
				err = yaml.Unmarshal(yamlFile, &expectedInstance)
				if err != nil {
					t.Errorf("Error unmarshaling YAML: %s", err)
					break
				}

				// Verify instance details
				t.Log("========= Verify Instance name =========")
				actualInstanceInfo := gjson.Parse(gcloudOutput)
				if actualInstanceInfo.Get("name").String() != expectedInstance.Name {
					t.Errorf("Instance name mismatch: actual=%s, expected=%s", actualInstanceInfo.Get("name").String(), expectedInstance.Name)
				}
				t.Log("========= Verify Instance zone =========")
				zoneName := filepath.Base(actualInstanceInfo.Get("zone").String())
				if zoneName != expectedInstance.Zone {
					t.Errorf("Zone mismatch: actual=%s, expected=%s", zoneName, expectedInstance.Zone)
				}

				// Check for correct image
				t.Log("========= Verify Instance image =========")
				actualImage := gjson.Get(gcloudOutput, "disks.0.licenses.0").String()
				// Get the image name from YAML config
				expectedImage := expectedInstance.Image
				// Split image name on '/'
				expectedImageParts := strings.Split(expectedImage, "/")
				// Extract only the image name
				expectedImageName := expectedImageParts[len(expectedImageParts)-1]

				if !strings.Contains(actualImage, expectedImageName) {
					t.Errorf("Image mismatch: actual=%s, expected to contain %s", actualImage, expectedImageName)
				}
				t.Log("========= Verify Instance network =========")
				// Fix for Network mismatch
				actualNetwork := gjson.Get(gcloudOutput, "networkInterfaces.0.network").String()
				if !strings.HasSuffix(actualNetwork, expectedInstance.Network) {
					t.Errorf("Network mismatch: actual=%s, expected=%s", actualNetwork, expectedInstance.Network)
				}
				t.Log("========= Verify Instance subnetwork =========")
				// Fix for Subnetwork mismatch
				actualSubnetwork := gjson.Get(gcloudOutput, "networkInterfaces.0.subnetwork").String()
				if !strings.HasSuffix(actualSubnetwork, expectedInstance.Subnetwork) {
					t.Errorf("Subnetwork mismatch: actual=%s, expected=%s", actualSubnetwork, expectedInstance.Subnetwork)
				}
				break
			}

			t.Logf("Instance '%s' not yet RUNNING (current state: %s). Retrying...", instanceName, status)
			time.Sleep(retryInterval)
		}
	}
	// Destroy Terraform Resources **First**
	terraform.Destroy(t, terraformOptions)

	// Delete VPC and Associated Resources (after Terraform destroy)
	deleteVPC(t, projectID, networkName)
}

/*
createVPC creates the VPC and subnet before the test execution.
*/
func createVPC(t *testing.T, projectID string, networkName string) {
	text := "compute"

	// Create VPC
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", networkName, "--project=" + projectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Logf("===Error %s Encountered while executing %s", err, text)
	}

	// Create Subnet
	subnetName := fmt.Sprintf("%s-subnet", networkName)
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
}

/*
deleteVPC deletes the VPC and subnet after the test.
*/
func deleteVPC(t *testing.T, projectID string, networkName string) {
	text := "compute"
	time.Sleep(120 * time.Second) // Wait for resources to be in a deletable state

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

	// 3. Delete VPC
	time.Sleep(150 * time.Second) // Wait for firewall deletion to complete
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

	// Create a GCE-specific instance configuration
	gceInstance := VMInstanceConfig{
		Name:       instanceName,
		ProjectID:  projectID,
		Region:     region,
		Zone:       zone,                              // Add zone for GCE
		Image:      "ubuntu-os-cloud/ubuntu-2204-lts", // Replace with your desired image
		Network:    networkID,                         // Use networkID for the network
		Subnetwork: subnetworkID,                      // Use subnetworkID for the subnetwork
	}

	yamlData, err := yaml.Marshal(&gceInstance)
	if err != nil {
		t.Errorf("Error while marshaling: %v", err)
	}

	// Specify a directory for config files (adjust if needed)
	configDir := "config"
	filePath := filepath.Join(configDir, "instance1.yaml") // Construct file path

	// Create the config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Errorf("Failed to create config directory: %v", err)
	}

	t.Logf("Created YAML config at %s with content:\n%s", filePath, string(yamlData))

	err = os.WriteFile(filePath, []byte(yamlData), 0644) // Use 0644 for file permissions
	if err != nil {
		t.Errorf("Unable to write data into the file: %v", err)
	}
}
