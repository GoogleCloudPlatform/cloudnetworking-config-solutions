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
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
	"math/rand"

	"os"
	"testing"
	"time"
)

const (
	terraformDirectoryPath = "../../../../../06-consumer/CloudRun/Service"
	region                 = "us-central1"
	configFolderPath       = "../../../test/integration/consumer/CloudRun/Service/config"
	image                  = "us-docker.pkg.dev/cloudrun/container/hello"
)

var (
	projectID   = os.Getenv("TF_VAR_project_id")
	serviceName = fmt.Sprintf("test-%d", rand.Int())
	tfVars      = map[string]any{
		"config_folder_path": configFolderPath,
	}
)

type ContainerNameStruct struct {
	Image string `yaml:"image"`
}

type ContainersStruct struct {
	ContainerName ContainerNameStruct `yaml:"container-name"`
}

type CloudRunServiceStruct struct {
	ProjectID  string           `yaml:"project_id"`
	Region     string           `yaml:"region"`
	Name       string           `yaml:"name"`
	Containers ContainersStruct `yaml:"containers"`
}

func TestCreateCloudRunService(t *testing.T) {
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

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 60 seconds to let resource acheive stable state.
	time.Sleep(60 * time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	cloudRunServiceOutputValue := terraform.OutputJson(t, terraformOptions, "cloud_run_service_details")
	if !gjson.Valid(cloudRunServiceOutputValue) {
		t.Errorf("Error parsing output, invalid json: %s", cloudRunServiceOutputValue)
	}
	result := gjson.Parse(cloudRunServiceOutputValue)
	serviceIDPath := fmt.Sprintf("%s.service.id", serviceName)
	t.Log(" ========= Terraform resource creation completed ========= ")
	t.Log(" ========= Verify Service ID ========= ")
	got := gjson.Get(result.String(), serviceIDPath).String()
	want := fmt.Sprintf("projects/%s/locations/%s/services/%s", projectID, region, serviceName)
	if got != want {
		t.Errorf("Cloud Run Service with invalid ID created = %v, want = %v", got, want)
	}
	t.Log(" ========= Verify Service Location ========= ")
	serviceLocationPath := fmt.Sprintf("%s.service.location", serviceName)
	got = gjson.Get(result.String(), serviceLocationPath).String()
	want = region
	if got != want {
		t.Errorf("Cloud Run Service with invalid Location created = %v, want = %v", got, want)
	}
	t.Log(" ========= Verify Service Name ========= ")
	serviceNamePath := fmt.Sprintf("%s.id", serviceName)
	got = gjson.Get(result.String(), serviceNamePath).String()
	want = serviceName
	if got != want {
		t.Errorf("Cloud Run Service with invalid Name created = %v, want = %v", got, want)
	}
}

/*
createConfigYAML is a helper function which creates the configuration YAML file.
*/
func createConfigYAML(t *testing.T) {
	t.Log("========= YAML File =========")

	containerNameList := ContainerNameStruct{
		Image: image,
	}

	containersStructList := ContainersStruct{
		ContainerName: containerNameList,
	}
	instance1 := CloudRunServiceStruct{
		Name:       serviceName,
		ProjectID:  projectID,
		Region:     region,
		Containers: containersStructList,
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
