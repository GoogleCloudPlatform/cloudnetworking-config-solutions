// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
)

var (
	terraformDirectoryPath = "../../../../03-security/CloudSQL"
	projectID              = os.Getenv("TF_VAR_project_id")
	uniqueID               = rand.Int() //included as a suffix to the VPC and subnet names.
	networkName            = fmt.Sprintf("test-vpc-security-%d", uniqueID)
	firewallName           = "test-allow-egress-cloudsql"
	firewallDirection      = "EGRESS"
)

/*
This test creates all the resources including the vpc network, subnetwork along with a PSA range.

It then validates if
1. Firewall Rule is created
2. Firewall Rule with correct direction is created
*/
func TestCreateCloudSQLFirewallRule(t *testing.T) {
	var (
		tfVars = map[string]any{
			"project_id": projectID,
			"network":    networkName,
			"egress_rules": map[string]any{
				firewallName: map[string]any{
					"deny": "false",
					"rules": []any{
						map[string]any{
							"protocol": "tcp",
							"ports":    []string{"3306"},
						},
					},
				},
			},
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
	createVPC(t, projectID, networkName)

	// Delete VPC created outside of the terraform module.
	defer deleteVPC(t, projectID, networkName)

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 60 seconds to let resource acheive stable state.
	time.Sleep(60 * time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	want := firewallName
	firewallRuleOutputValue := terraform.OutputJson(t, terraformOptions, "cloudsql_firewall_rules")
	if !gjson.Valid(firewallRuleOutputValue) {
		t.Errorf("Error parsing output, invalid json: %s", firewallRuleOutputValue)
	}
	result := gjson.Parse(firewallRuleOutputValue)
	firewallNamePath := fmt.Sprintf("%s.name", firewallName)
	got := gjson.Get(result.String(), firewallNamePath).String()

	// Validate the firewall rule name created by terraform modules.
	t.Log(" ========= Verify Firewall Name ========= ")
	if got != want {
		t.Errorf("Firewall with invalid name created = %v, want = %v", got, want)
	}

	// Validate the firewall rule direction created by terraform modules.
	want = firewallDirection
	firewallDirectionPath := fmt.Sprintf("%s.direction", firewallName)
	got = gjson.Get(result.String(), firewallDirectionPath).String()
	t.Log(" ========= Verify Firewall Direction ========= ")
	if got != want {
		t.Errorf("Firewall with invalid direction created = %v, want = %v", got, want)
	}
}

/*
deleteVPC is a helper function which deletes the VPC after
completion of the test.
*/
func deleteVPC(t *testing.T, projectID string, networkName string) {
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
 createVPC is a helper function which creates the VPC before the
 execution of the test.
*/

func createVPC(t *testing.T, projectID string, networkName string) {
	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", networkName, "--project=" + projectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom", "--verbosity=none"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
	time.Sleep(60 * time.Second)
}
