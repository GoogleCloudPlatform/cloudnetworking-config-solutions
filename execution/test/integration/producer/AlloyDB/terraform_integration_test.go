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
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	projectID              = os.Getenv("TF_VAR_project_id")
	region                 = "us-central1"
	terraformDirectoryPath = "../../../../04-producer/AlloyDB"
	configFolderPath       = "../../test/integration/producer/AlloyDB/config"
	rangeName              = "psatestrangealloydb"
	clusterDisplayName     = fmt.Sprint(rand.Int())
	networkName            = fmt.Sprintf("vpc-%s-test", clusterDisplayName)
	alloyDBClusterId       = fmt.Sprintf("cid-%s-test", clusterDisplayName)
	instanceID             = fmt.Sprintf("id-%s-test", clusterDisplayName)
	networkID              = fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName)
)

type PrimaryInstanceStruct struct {
	InstanceID string `yaml:"instance_id"`
}

type AlloyDBStruct struct {
	ClusterID          string                `yaml:"cluster_id"`
	ClusterDisplayName string                `yaml:"cluster_display_name"`
	ProjectID          string                `yaml:"project_id"`
	Region             string                `yaml:"region"`
	NetworkID          string                `yaml:"network_id"`
	PrimaryInstance    PrimaryInstanceStruct `yaml:"primary_instance"`
	AllocatedIPRange   string                `yaml:"allocated_ip_range"`
}

/*
This test creates all the pre-requsite resources including the vpc network, subnetwork along with a PSA range.
It then validates if
1. AlloyDB instance is created.
2. AlloyDB instance is created in the correct network and correct PSA range.
3. AlloyDB instance is in ACTIVE state.
*/
func TestCreateAlloyDB(t *testing.T) {
	// Initialize a AlloyDB config YAML file to be tested.
	createConfigYAML(t)
	var (
		tfVars = map[string]any{
			"config_folder_path": configFolderPath,
		}
	)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars:                 tfVars,
		TerraformDir:         terraformDirectoryPath,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})
	// Create VPC outside of the terraform module.
	err := createVPC(t, projectID, networkName)
	if err != nil {
		t.Fatal(err)
	}
	// Create PSA in the VPC.
	createPSA(t, projectID, networkName, rangeName)

	// Delete VPC created outside of the terraform module.
	defer deleteVPC(t, projectID, networkName)

	// Remove PSA from the VPC.
	defer deletePSA(t, projectID, networkName, rangeName)

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 60 seconds to let resource acheive stable state.
	time.Sleep(60 * time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	alloyDBOutputValue := terraform.OutputJson(t, terraformOptions, "cluster_details")
	t.Log(" ========= Terraform resource creation completed ========= ")
	t.Log(" ========= Verify AlloyDB Cluster ID ========= ")
	want := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, region, alloyDBClusterId)
	if !gjson.Valid(alloyDBOutputValue) {
		t.Errorf("Error parsing output, invalid json: %s", alloyDBOutputValue)
	}
	result := gjson.Parse(alloyDBOutputValue)
	clusterIDPath := fmt.Sprintf("%s.cluster_id", clusterDisplayName)
	got := gjson.Get(result.String(), clusterIDPath).String()
	if got != want {
		t.Errorf("AlloyDB Cluster with invalid Cluster ID = %v, want = %v", got, want)
	}
	t.Log(" ========= Verify AlloyDB Cluster Status ========= ")
	want = "READY"
	clusterStatusPath := fmt.Sprintf("%s.cluster_status", clusterDisplayName)
	got = gjson.Get(result.String(), clusterStatusPath).String()
	if got != want {
		t.Errorf("AlloyDB Cluster with invalid Cluster status = %v, want = %v", got, want)
	}
	t.Log(" ========= Verify AlloyDB Cluster PSA Range Name ========= ")
	want = rangeName
	clusterPSARangeNamePath := fmt.Sprintf("%s.network_config.0.allocated_ip_range", clusterDisplayName)
	got = gjson.Get(result.String(), clusterPSARangeNamePath).String()
	if got != want {
		t.Errorf("AlloyDB Cluster with invalid PSA Range Name = %v, want = %v", got, want)
	}
}

/*
deleteVPC is a helper function which deletes the VPC after
completion of the test.
*/
func deleteVPC(t *testing.T, projectID string, networkName string) {
	time.Sleep(60 * time.Second)
	text := "compute"
	time.Sleep(60 * time.Second)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "delete", networkName, "--project=" + projectID, "--quiet"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
}

/*
deletePSA is a helper function which deletes the PSA range after the
execution of the test.
*/
func deletePSA(t *testing.T, projectID string, networkName string, rangeName string) {
	// Delete PSA IP range
	time.Sleep(60 * time.Second)
	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "addresses", "delete", rangeName, "--project=" + projectID, "--global", "--verbosity=none", "--format=json", "--quiet"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
	time.Sleep(60 * time.Second)
	// Delete PSA range
	text = "services"
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "vpc-peerings", "delete", "--service=servicenetworking.googleapis.com", "--project=" + projectID, "--network=" + networkName, "--verbosity=none", "--format=json", "--quiet"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
}

/*
createVPC is a helper function which creates the VPC before the
execution of the test.
*/
func createVPC(t *testing.T, projectID string, networkName string) error {
	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", networkName, "--project=" + projectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom", "--verbosity=none"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
	return err
}

/*
createPSA is a helper function which creates the PSA range before the
execution of the test.
*/
func createPSA(t *testing.T, projectID string, networkName string, rangeName string) {
	// Create an IP range

	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "addresses", "create", rangeName, "--purpose=VPC_PEERING", "--addresses=10.0.64.0", "--prefix-length=20", "--project=" + projectID, "--network=" + networkName, "--global", "--verbosity=none", "--format=json"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}

	// Create PSA range
	text = "services"
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "vpc-peerings", "connect", "--service=servicenetworking.googleapis.com", "--ranges=" + rangeName, "--project=" + projectID, "--network=" + networkName, "--verbosity=none", "--format=json"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
}

/*
createConfigYAML is a helper function which creates the configigration YAML file
for an alloydb instance range before the.
*/
func createConfigYAML(t *testing.T) {
	t.Log("========= YAML File =========")
	instance1 := AlloyDBStruct{
		ClusterID:          alloyDBClusterId,
		ClusterDisplayName: clusterDisplayName,
		ProjectID:          projectID,
		Region:             region,
		NetworkID:          networkID,
		PrimaryInstance: PrimaryInstanceStruct{
			InstanceID: instanceID,
		},
		AllocatedIPRange: rangeName,
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
