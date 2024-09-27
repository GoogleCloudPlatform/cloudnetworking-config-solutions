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
	"io"
	"math/rand"
	"os"
	"testing"
	"time"
)
// Test configuration (adjust as needed)
var (
	projectID                 = os.Getenv("TF_VAR_project_id")
	region                    = "us-central1"
	terraformDirectoryPath    = "../../../../04-producer/VectorSearch"
	configFolderPath          = "../../test/integration/producer/VectorSearch/config"
	indexUpdateMethod         = "BATCH_UPDATE"
	indexDisplayName          = fmt.Sprintf("vectorsearch%d", rand.Int())
	rangeName                 = fmt.Sprintf("psa-%s", indexDisplayName)
	indexEndpointDisplayName  = fmt.Sprintf("indexendpoint-%s", indexDisplayName)
	deployedIndexID           = fmt.Sprintf("deployedindexid_%s", indexDisplayName)
	networkName               = fmt.Sprintf("vpc-%s-test", indexDisplayName)
	dimension                 = 2
	approximateNeighborsCount = 150
)

type VectorSearchStruct struct {
	ProjectID                 string `yaml:"project_id"`
	Region                    string `yaml:"region"`
	IndexDisplayName          string `yaml:"index_display_name"`
	IndexDescription          string `yaml:"index_description"`
	Dimension                 int    `yaml:"dimension"`
	ApproximateNeighborsCount int    `yaml:"approximate_neighbors_count"`
	IndexUpdateMethod         string `yaml:"index_update_method"`
	IndexEndpointDisplayname  string `yaml:"index_endpoint_display_name"`
	IndexEndpointNetwork      string `yaml:"index_endpoint_network"`
	BruteForceConfig          string `yaml:"brute_force_config"`
	DeployedIndexId           string `yaml:"deployed_index_id"`
}

/*
TestCreateVectorSearch creates a vector search index, index endpoint and deploys the index endpoint to this index,
performs verification on successfull creation of the vector search resources.
*/
func TestCreateVectorSearch(t *testing.T) {
	// Initialize a Vector Search config YAML file to be tested.
	createConfigYAML(t)
	sourceFile := "provider.tf"
	destinationFile := terraformDirectoryPath + "/test-provider.tf"
	defer os.Remove(destinationFile)

	source, err := os.Open(sourceFile)
	if err != nil {
		t.Error(err)
	}
	defer source.Close()

	destination, err := os.Create(destinationFile)
	if err != nil {
		t.Error(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		t.Error(err)
	}

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
	err = createVPC(t, projectID, networkName)
	if err != nil {
		t.Error(err)
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
	vectorSearchOutputValue := terraform.OutputJson(t, terraformOptions, "vector_search_instance_details")
	t.Log(" ========= Terraform resource creation completed ========= ")
	if !gjson.Valid(vectorSearchOutputValue) {
		t.Errorf("Error parsing output, invalid json: %s", vectorSearchOutputValue)
	}
	result := gjson.Parse(vectorSearchOutputValue)
	indexNamePath := fmt.Sprintf("%s.index_name", indexDisplayName)
	indexEndpointNamePath := fmt.Sprintf("%s.index_endpoint_name", indexDisplayName)
	indexID := gjson.Get(result.String(), indexNamePath).String()
	indexEndpointID := gjson.Get(result.String(), indexEndpointNamePath).String()

	t.Log(" ========= Verify Vector Search Index ID name ========= ")
	indexIDPath := fmt.Sprintf("%s.index_id", indexDisplayName)
	got := gjson.Get(result.String(), indexIDPath).String()
	want := fmt.Sprintf("projects/%s/locations/%s/indexes/%s", projectID, region, indexID)
	if got != want {
		t.Errorf("Index ID with incorrect details created = %v, want = %v", got, want)
	}
	t.Log(" ========= Verify Vector Search Index Endpoint name ========= ")
	indexEndpointIDPath := fmt.Sprintf("%s.index_endpoint_id", indexDisplayName)
	got = gjson.Get(result.String(), indexEndpointIDPath).String()
	want = fmt.Sprintf("projects/%s/locations/%s/indexEndpoints/%s", projectID, region, indexEndpointID)
	if got != want {
		t.Errorf("Index Endpoint ID with incorrect details created = %v, want = %v", got, want)
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
createConfigYAML is a helper function which creates the config YAML file which is used
for creation of test instance.
 */
func createConfigYAML(t *testing.T) {
	// Fetch Project Number
	text := "projects"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "describe", projectID, "--format=value(\"projectNumber\")"},
	}
	projectNumber, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
	t.Log("========= YAML File =========")
	indexEndpointNetwork := fmt.Sprintf("projects/%s/global/networks/%s", projectNumber, networkName)
	t.Logf("Index Endpoint Network : %s", indexEndpointNetwork)
	instance1 := VectorSearchStruct{
		ProjectID:                 projectID,
		Region:                    region,
		IndexDisplayName:          indexDisplayName,
		ApproximateNeighborsCount: approximateNeighborsCount,
		Dimension:                 dimension,
		IndexUpdateMethod:         indexUpdateMethod,
		IndexEndpointDisplayname:  indexEndpointDisplayName,
		IndexEndpointNetwork:      indexEndpointNetwork,
		BruteForceConfig:          "",
		DeployedIndexId:           deployedIndexID,
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
