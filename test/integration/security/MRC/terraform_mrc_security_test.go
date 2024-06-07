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
	"slices"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
)

var (
	terraformDirectoryPath = "../../../../03-security/MRC" // Update with your actual path
	projectID              = os.Getenv("TF_VAR_project_id")
	uniqueID               = rand.Int()
	networkName            = fmt.Sprintf("test-vpc-security-%d", uniqueID)
	firewallName           = "test-allow-egress"
	firewallDirection      = "EGRESS"
)

func TestCreateMemorystoreRedisFirewallRule(t *testing.T) {
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
							"ports":    []string{"6379"},
						},
					},
				},
			},
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

	// Create VPC
	createVPC(t, projectID, networkName)
	defer deleteVPC(t, projectID, networkName)

	// Clean up Terraform resources
	defer terraform.Destroy(t, terraformOptions)

	// Initialize and Apply
	terraform.InitAndApply(t, terraformOptions)
	time.Sleep(60 * time.Second)

	// Get Output and Validate
	want := firewallName
	firewallRuleOutputValue := terraform.OutputJson(t, terraformOptions, "mrc_firewall_rules")
	if !gjson.Valid(firewallRuleOutputValue) {
		t.Errorf("Error parsing output, invalid json: %s", firewallRuleOutputValue)
	}

	result := gjson.Parse(firewallRuleOutputValue)
	firewallNamePath := fmt.Sprintf("%s.name", firewallName)
	got := gjson.Get(result.String(), firewallNamePath).String()

	t.Log(" ========= Verify Firewall Name ========= ")
	if got != want {
		t.Errorf("Firewall with invalid name created = %v, want = %v", got, want)
	}

	// Validate Firewall Direction
	want = firewallDirection
	firewallDirectionPath := fmt.Sprintf("%s.direction", firewallName)
	got = gjson.Get(result.String(), firewallDirectionPath).String()
	t.Log(" ========= Verify Firewall Direction ========= ")
	if got != want {
		t.Errorf("Firewall with invalid direction created = %v, want = %v", got, want)
	}
	t.Run("Destination Ranges", func(t *testing.T) {
		got := gjson.Get(result.String(), fmt.Sprintf("%s.destination_ranges", firewallName)).Array()
		want := []string{"0.0.0.0/0"} // Adjust if needed

		// Compare the contents of the arrays, not the array types directly
		if len(got) != len(want) || !slices.Contains(want, got[0].String()) {
			t.Errorf("Destination ranges mismatch: got %v, want %v", got, want)
		}
	})

	t.Run("Allow Rules", func(t *testing.T) {
		rules := gjson.Get(result.String(), fmt.Sprintf("%s.allow", firewallName)).Array()
		if len(rules) != 1 {
			t.Errorf("Expected 1 allow rule, got %d", len(rules))
		} else {
			rule := rules[0]
			gotProtocol := rule.Get("protocol").String()
			wantProtocol := "tcp"
			if gotProtocol != wantProtocol {
				t.Errorf("Allow rule protocol mismatch: got %q, want %q", gotProtocol, wantProtocol)
			}

			gotPorts := rule.Get("ports").Array()
			wantPorts := []string{"6379"}

			// Compare the contents of the arrays, not the array types directly
			if len(gotPorts) != len(wantPorts) || !slices.Contains(wantPorts, gotPorts[0].String()) {
				t.Errorf("Allow rule ports mismatch: got %v, want %v", gotPorts, wantPorts)
			}
		}
	})

	// New Assertion: Check if disabled is false
	t.Run("Disabled", func(t *testing.T) {
		got := gjson.Get(result.String(), fmt.Sprintf("%s.disabled", firewallName)).Bool()
		want := false
		if got != want {
			t.Errorf("Firewall disabled state mismatch: got %t, want %t", got, want)
		}
	})
}

// Helper Functions
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
