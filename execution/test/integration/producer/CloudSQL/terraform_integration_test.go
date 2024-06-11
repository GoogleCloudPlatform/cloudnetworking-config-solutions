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
	terraformDirectoryPath = "../../../../04-producer/CloudSQL"
	configFolderPath       = "../../test/integration/producer/CloudSQL/config"
	rangeName              = "psatestrange"
	databaseVersion        = "POSTGRES_15"
	name                   = fmt.Sprintf("cloudsql-%d", rand.Int())
	networkName            = fmt.Sprintf("vpc-%s-test", name)
	networkID              = fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName)
)

// AllocatedIPRangesStruct represents the allocated IP Ranges in the PSA Configuration(PSAConfigStruct).
type AllocatedIPRangesStruct struct {
	Primary string `yaml:"primary,omitempty"`
	Replica string `yaml:"replica,omitempty"`
}

// PSAConfigStruct represents the PSA configurations in the Connectivity Struct.
type PSAConfigStruct struct {
	PrivateNetwork    string                  `yaml:"private_network"`
	AllocatedIPRanges AllocatedIPRangesStruct `yaml:"allocated_ip_ranges,omitempty"`
}

// ConnectivityStruct represents the Connectivity in the network configuration.
type ConnectivityStruct struct {
	PublicIPV4                 bool            `yaml:"public_ipv4,omitempty"`
	PSAConfig                  PSAConfigStruct `yaml:"psa_config,omitempty"`
	PSCAllowedConsumerProjects []string        `yaml:"psc_allowed_consumer_projects,omitempty"`
}

// NetworkConfigStruct represent the Network Config in the CloudSQLStruct.
type NetworkConfigStruct struct {
	AuthorizedNetwork map[string]string  `yaml:"authorized_networks,omitempty"`
	Connectivity      ConnectivityStruct `yaml:"connectivity"`
}

// CloudSQLStruct defines the structure to parse configuration data for a YAML file
// for creating Cloud SQL instances. This structure maps directly to the expected format
// of the YAML file and uses struct tags for unmarshalling.

type CloudSQLStruct struct {
	Name                        string              `yaml:"name"`
	ProjectID                   string              `yaml:"project_id"`
	Region                      string              `yaml:"region"`
	DatabaseVersion             string              `yaml:"database_version"`
	NetworkConfig               NetworkConfigStruct `yaml:"network_config"`
	TerraformDeletionProtection bool                `yaml:"terraform_deletion_protection"`
	GCPDeletionProtection       bool                `yaml:"gcp_deletion_protection"`
}

/*
This test creates all the pre-requsite resources including the vpc network, subnetwork along with a PSA range.
It then validates if
1. CloudSQL instance is created.
2. CloudSQL instance is created in the correct network, project, region and of correct version.
3. CloudSQL instance only have a private ip and does not have a public IP.
*/
func TestCreateCloudSQL(t *testing.T) {
	// Initialize a Cloud SQL config YAML file to be tested.
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
	cloudSQLOutputValue := terraform.OutputJson(t, terraformOptions, "cloudsql_instance_details")
	t.Log(" ========= Terraform resource creation completed ========= ")
	t.Log(" ========= Verify Cloud SQL instance name ========= ")
	want := name
	if !gjson.Valid(cloudSQLOutputValue) {
		t.Errorf("Error parsing output, invalid json: %s", cloudSQLOutputValue)
	}
	result := gjson.Parse(cloudSQLOutputValue)
	cloudSQLInstanceNamePath := fmt.Sprintf("%s.name", name)
	got := gjson.Get(result.String(), cloudSQLInstanceNamePath).String()
	if got != want {
		t.Errorf("Cloud SQL instance with invalid instance name created = %v, want = %v", got, want)
	}
	t.Log(" ========= Verify Cloud SQL Instance connection name ========= ")
	want = fmt.Sprintf("%s:%s:%s", projectID, region, name)
	connectionNamePath := fmt.Sprintf("%s.connection_name", name)
	got = gjson.Get(result.String(), connectionNamePath).String()
	if got != want {
		t.Errorf("Cloud SQL instance with incorrect connection name created = %v, want = %v", got, want)
	}
	t.Log(" ========= Verify Cloud SQL Instance database version ========= ")
	want = databaseVersion
	cloudSQLDatabaseVersionPath := fmt.Sprintf("%s.database_version", name)
	got = gjson.Get(result.String(), cloudSQLDatabaseVersionPath).String()
	if got != want {
		t.Errorf("Cloud SQL Instance with invalid database version created = %v, want = %v", got, want)
	}

	t.Log(" ========= Verify Cloud SQL Instance does not have a public ip ========= ")
	cloudSQLPublicIPPath := fmt.Sprintf("%s.public_ip_address", name)
	got = gjson.Get(result.String(), cloudSQLPublicIPPath).String()
	if got != "" {
		t.Errorf("Cloud SQL Instance with public ip created(should be a private ip only) = %v", got)
	}

	t.Log(" ========= Verify Cloud SQL Instance does  have a private ip ========= ")
	cloudSQLPrivateIPPath := fmt.Sprintf("%s.private_ip_address", name)
	got = gjson.Get(result.String(), cloudSQLPrivateIPPath).String()
	if got == "" {
		t.Errorf("Cloud SQL Instance does not contain private ip = %v", got)
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
		Args:    []string{text, "addresses", "delete", rangeName, "--project=" + projectID, "--global", "--verbosity=info", "--format=json", "--quiet"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Logf("===Error %s Encountered while executing %s", err, text)
	}
	time.Sleep(60 * time.Second)
	// Delete PSA range
	text = "services"
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "vpc-peerings", "delete", "--service=servicenetworking.googleapis.com", "--project=" + projectID, "--network=" + networkName, "--verbosity=info", "--format=json", "--quiet"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Logf("===Error %s Encountered while executing %s", err, text)
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
		Args:    []string{text, "networks", "create", networkName, "--project=" + projectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom", "--verbosity=info"},
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
		Args:    []string{text, "addresses", "create", rangeName, "--purpose=VPC_PEERING", "--addresses=10.0.64.0", "--prefix-length=20", "--project=" + projectID, "--network=" + networkName, "--global", "--verbosity=info", "--format=json"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
	// Create PSA range
	text = "services"
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "vpc-peerings", "connect", "--service=servicenetworking.googleapis.com", "--ranges=" + rangeName, "--project=" + projectID, "--network=" + networkName, "--verbosity=info", "--format=json"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
	time.Sleep(60 * time.Second)
}

/*
createConfigYAML is a helper function which creates the configigration YAML file
for a cloudsql instance.
*/
func createConfigYAML(t *testing.T) {
	t.Log("========= YAML File =========")
	instance1 := CloudSQLStruct{
		Name:                        name,
		ProjectID:                   projectID,
		Region:                      region,
		DatabaseVersion:             databaseVersion,
		TerraformDeletionProtection: false,
		GCPDeletionProtection:       false,
		NetworkConfig: NetworkConfigStruct{
			Connectivity: ConnectivityStruct{
				PSAConfig: PSAConfigStruct{
					PrivateNetwork: networkID,
					AllocatedIPRanges: AllocatedIPRangesStruct{
						Primary: rangeName,
					},
				},
			},
		},
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
