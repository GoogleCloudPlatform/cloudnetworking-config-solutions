// Copyright 2024 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integrationtest

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	// for sorting slices
	// for comparison operations
	"github.com/google/go-cmp/cmp"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

var (
	projectRoot, _ = filepath.Abs("../../../../")
	// Path to the Terraform module directory.
	terraformDirectoryPath = filepath.Join(projectRoot, "04-producer/GKE")
	// Path to the folder containing YAML configuration files.
	configFolderPath   = filepath.Join(projectRoot, "test/integration/producer/GKE/config")
	projectID          = os.Getenv("TF_VAR_project_id")
	region             = "us-central1"
	kubernetesVersion  = "1.27.16-gke.1287000"
	instanceName       = fmt.Sprintf("gke-%d", rand.Int())
	networkName        = fmt.Sprintf("gke-cluster-vpc-%d", rand.Int())
	subnetName         = fmt.Sprintf("gke-cluster-subnetwork-%d", rand.Int())
	subnetIPRange      = "10.0.0.0/16"
	ipRangePods        = "pods"
	ipRangeServices    = "services"
	podIPRange         = "10.1.0.0/16"
	servicesIPRange    = "10.2.0.0/16"
	deletionProtection = false
	tfVars             = map[string]any{
		"config_folder_path": configFolderPath,
	}
	invalidTFVars = map[string]any{
		"config_folder_path": configFolderPath,
		"network":            "random/google/cloud/network/",
	}
)

type GKEConfig struct {
	Name               string `yaml:"name"`
	ProjectID          string `yaml:"project_id"`
	KubernetesVersion  string `yaml:"kubernetes_version"`
	Network            string `yaml:"network"`
	Subnetwork         string `yaml:"subnetwork"`
	IPRangePods        string `yaml:"ip_range_pods"`
	IPRangeServices    string `yaml:"ip_range_services"`
	Region             string `yaml:"region"`
	DeletionProtection bool   `yaml:"deletion_protection"`
}

// TestCreateGKECluster tests the creation of a GKE cluster.
func TestCreateGKECluster(t *testing.T) {
	// Initialize a GKE config YAML file to be tested.
	createGKEConfigYAML(t)

	var (
		tfVars = map[string]any{
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

	// Create network, subnet, and IP ranges
	createNetwork(t, projectID, networkName)
	createSubnet(t, projectID, networkName, subnetName)
	createIPRanges(t, projectID, region, subnetName)
	time.Sleep(60 * time.Second)

	// Delete network, subnet, and IP ranges
	defer deleteNetwork(t, projectID, networkName)
	defer deleteSubnet(t, projectID, subnetName)

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for the GKE cluster to become available with retries
	maxRetries := 10
	retryInterval := 10 * time.Second
	var clusterID string

	for i := 0; i < maxRetries; i++ {
		clusterOutput := terraform.OutputJson(t, terraformOptions, "gke_clusters")

		if !gjson.Valid(clusterOutput) {
			t.Errorf("Error parsing output, invalid json: %s", clusterOutput)
			continue // Skip to next retry if output is invalid
		}

		result := gjson.Parse(clusterOutput)

		result.ForEach(func(key, value gjson.Result) bool {
			clusterID = key.String()
			clusterData := value

			// Verify GKE Cluster Properties
			name := clusterData.Get("name").String()
			if name != instanceName {
				t.Errorf("GKE Cluster name is invalid: got %s, want %s", name, instanceName)
			} else {
				t.Logf("GKE Cluster name is valid: %s", name) // Success message
			}

			// Verify Kubernetes Version
			kubernetesVersion := clusterData.Get("master_version").String()
			if kubernetesVersion != "1.27.16-gke.1287000" {
				t.Errorf("GKE Cluster Kubernetes version is invalid: got %s, want 1.27.16-gke.1287000", kubernetesVersion)
			} else {
				t.Logf("GKE Cluster Kubernetes version is valid: %s", kubernetesVersion) // Success message
			}

			// Verify Region
			regionInLogs := clusterData.Get("region").String()
			if regionInLogs != region {
				t.Errorf("GKE Cluster region is invalid: got %s, want %s", regionInLogs, region)
			} else {
				t.Logf("GKE Cluster region is valid: %s", regionInLogs) // Success message
			}

			return false // Break out of the iteration
		})

		if clusterID != "" {
			break
		} else {
			time.Sleep(retryInterval)
		}
	}

	if clusterID == "" {
		t.Errorf("GKE cluster ID not found after %d retries", maxRetries)
	}
}

// TestTerraformModuleResourceAddressListMatch compares and verifies the list of resources,
// modules created by the Terraform solution.
func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	// 1. Read and parse the YAML config file
	yamlFile, err := ioutil.ReadFile(filepath.Join(configFolderPath, "gke-config.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal into a GKEConfig struct
	var gkeConfig GKEConfig
	err = yaml.Unmarshal(yamlFile, &gkeConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Construct the expected module address
	expectedModulesAddress := []string{fmt.Sprintf("module.gke[\"%s\"]", gkeConfig.Name)}

	// Terraform options for planning.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})

	// Run Terraform init and plan, capturing the plan structure.
	planStruct := terraform.InitAndPlanAndShow(t, terraformOptions)

	// Parse the Terraform plan JSON output.
	content, err := terraform.ParsePlanJSON(planStruct)
	if err != nil {
		t.Fatal(err)
	}

	// Extract the actual module addresses from the plan.
	actualModuleAddress := make(map[string]bool) // Use a map as a set
	for _, element := range content.ResourceChangesMap {
		if element.ModuleAddress != "" && !actualModuleAddress[element.ModuleAddress] {
			actualModuleAddress[element.ModuleAddress] = true // Only store once
		}
	}

	// Convert the set of addresses to a sorted slice for comparison
	actualModuleSlice := make([]string, 0, len(actualModuleAddress))
	for addr := range actualModuleAddress {
		actualModuleSlice = append(actualModuleSlice, addr)
	}
	sort.Strings(actualModuleSlice) // Sort the slice

	// Compare the sorted lists of expected and actual module addresses.
	want := expectedModulesAddress
	got := actualModuleSlice
	if !cmp.Equal(got, want) {
		t.Errorf("Test Element Mismatch: got = %v, want = %v", got, want)
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

	if got, want := resourceCount.Add, 7; got != want { // Expect 4 resources to be added
		t.Errorf("Test Resource Count Add = %v, want = %v", got, want)
	}

	if got, want := resourceCount.Change, 0; got != want { // Expect 0 resource to be changed
		t.Errorf("Test Resource Count Change = %v, want = %v", got, want)
	}

	if got, want := resourceCount.Destroy, 0; got != want { // Expect 0 resources to be destroyed
		t.Errorf("Test Resource Count Destroy = %v, want = %v", got, want)
	}
}

/*
createGKEConfigYAML creates the YAML configuration file for GKE.
*/
func createGKEConfigYAML(t *testing.T) {
	t.Log("========= YAML File =========")
	gkeConfig := GKEConfig{
		Name:               instanceName,
		ProjectID:          projectID,
		KubernetesVersion:  kubernetesVersion,
		Network:            networkName,
		Subnetwork:         subnetName,
		IPRangePods:        ipRangePods,
		IPRangeServices:    ipRangeServices,
		Region:             region,
		DeletionProtection: deletionProtection,
	}

	yamlData, err := yaml.Marshal(&gkeConfig)
	if err != nil {
		t.Errorf("Error while marshalling YAML: %v", err)
	}
	filePath := filepath.Join(configFolderPath, "gke-config.yaml")
	t.Logf("Created YAML config at %s with content:\n%s", filePath, string(yamlData))

	err = os.WriteFile(filePath, []byte(yamlData), 0666)
	if err != nil {
		t.Errorf("Unable to write data into the file: %v", err)
	}
}

/*
createNetwork creates the VPC network required for the GKE cluster.
*/
func createNetwork(t *testing.T, projectID string, networkName string) {
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"compute", "networks", "create", networkName,
			"--project=" + projectID,
			"--bgp-routing-mode=global",
			"--subnet-mode=custom",
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error creating network: %s", err)
	}
}

/*
createSubnet creates the subnet required for the GKE cluster.
*/
func createSubnet(t *testing.T, projectID string, networkName string, subnetName string) {
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"compute", "networks", "subnets", "create", subnetName,
			"--project=" + projectID,
			"--network=" + networkName,
			"--region=" + region,
			"--range=" + subnetIPRange,
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error creating subnet: %s", err)
	}
}

/*
createIPRanges creates the IP ranges required for the GKE cluster.
*/
func createIPRanges(t *testing.T, projectID string, region string, subnetName string) {
	// Create IP range for pods
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"compute", "networks", "subnets", "update", subnetName,
			"--project=" + projectID,
			"--region=" + region,
			"--add-secondary-ranges=" + "pods=" + podIPRange,
			"--add-secondary-ranges=" + "services=" + servicesIPRange,
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error creating IP range for pods/services: %s", err)
	}
}

/*
deleteSubnet deletes the subnet.
*/
func deleteSubnet(t *testing.T, projectID string, subnetName string) {
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"compute", "networks", "subnets", "delete", subnetName,
			"--project=" + projectID,
			"--region=" + region,
			"--quiet",
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error deleting subnet: %s", err)
	}
}

/*
deleteNetwork deletes the VPC network.
*/
func deleteNetwork(t *testing.T, projectID string, networkName string) {
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"compute", "networks", "delete", networkName,
			"--project=" + projectID,
			"--quiet",
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error deleting network: %s", err)
	}
}
