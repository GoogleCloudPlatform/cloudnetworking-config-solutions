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
	compare "cmp"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
)

var (
	projectID              = os.Getenv("TF_VAR_project_id")
	terraformDirectoryPath = "../../../01-organization"
	apisList               = []string{"aiplatform.googleapis.com", "alloydb.googleapis.com", "compute.googleapis.com", "container.googleapis.com", "iam.googleapis.com", "servicenetworking.googleapis.com", "sqladmin.googleapis.com"}
	tfVars                 = map[string]any{
		"activate_api_identities": map[string]any{
			projectID: map[string]any{
				"project_id":    projectID,
				"activate_apis": apisList,
			},
		},
	}
)

/*
This test validates if
1. Correct Project ID is used.
2. List of Project API's has been enabled.
*/
func TestEnableAPI(t *testing.T) {
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
	activateAPIOutputValue := terraform.OutputJson(t, terraformOptions, "activated_api_identities")
	if !gjson.Valid(activateAPIOutputValue) {
		t.Errorf("Error parsing output, invalid json: %s", activateAPIOutputValue)
	}
	result := gjson.Parse(activateAPIOutputValue)
	t.Log(" ========= Terraform resource creation completed ========= ")
	t.Log(" ========= Validate project ID ========= ")
	want := projectID
	projectIDPath := fmt.Sprintf("%s.project_id", projectID)
	got := gjson.Get(result.String(), projectIDPath).String()
	if got != want {
		t.Errorf("Project APIs being enabled in an invalid project ID = %v, want = %v", got, want)
	}
	enabledAPIPath := fmt.Sprintf("%s.enabled_apis", projectID)
	t.Log(" ========= Verify list of enabled API ========= ")
	wantList, err := json.Marshal(apisList)
	if err != nil {
		t.Errorf("Cannot encode to JSON %v", err)
	}
	wantAPIList := string(wantList)
	got = gjson.Get(result.String(), enabledAPIPath).String()
	if !cmp.Equal(got, wantAPIList, cmpopts.SortSlices(compare.Less[string])) {
		t.Errorf("Test list of enabled APIs Mismatch = %v, want = %v", got, wantAPIList)
	}
}
