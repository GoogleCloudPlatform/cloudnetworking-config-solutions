// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

var (
	terraformDirectoryPath = "../../../../03-security/GCE" // Update with your GCE directory path
	projectID              = os.Getenv("TF_VAR_project_id")
	uniqueID               = rand.Int()
	network                = fmt.Sprintf("test-vpc-security-%d", uniqueID)
	firewallRuleName       = "allow-ssh-custom-ranges-gce"
)

func TestGCEFirewallRuleProperties(t *testing.T) {
	var (
		tfVars = map[string]any{
			"project_id": projectID,
			"network":    network,
			"ingress_rules": map[string]any{
				firewallRuleName: map[string]any{
					"deny": "false",
					"rules": []any{
						map[string]any{
							"protocol": "tcp",
							"ports":    []string{"22", "443"},
						},
					},
				},
			},
		}
	)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		NoColor:      true,
	})

	// Create VPC
	createVPC(t, projectID, network)

	// Defer VPC deletion
	defer deleteVPC(t, projectID, network)

	// Terraform init and apply
	terraform.InitAndApply(t, terraformOptions)
	time.Sleep(60 * time.Second) // Wait for resource creation

	// Get Firewall rule from output
	firewallRulesOutput := terraform.OutputJson(t, terraformOptions, "rules")
	firewallRules := gjson.Parse(firewallRulesOutput).Map() // Parse as a map

	// Check if the firewall rule exists in the state
	t.Run("Firewall Rule Exists", func(t *testing.T) {
		_, exists := firewallRules[firewallRuleName]
		assert.True(t, exists, "Firewall rule %s not found in state", firewallRuleName)
	})

	// Assertions
	t.Run("Firewall Properties", func(t *testing.T) {
		ruleData := firewallRules[firewallRuleName]

		// Name
		t.Run("Name", func(t *testing.T) {
			got := ruleData.Get("name").String()
			want := firewallRuleName
			assert.Equal(t, want, got, "Firewall name mismatch")
		})

		// Direction (Correct for ingress)
		t.Run("Direction", func(t *testing.T) {
			got := ruleData.Get("direction").String()
			want := "INGRESS"
			assert.Equal(t, want, got, "Firewall rule direction mismatch")
		})

		// Source Ranges
		t.Run("Source Ranges", func(t *testing.T) {
			// Extract source range strings
			got := ruleData.Get("source_ranges").Array()
			extractedSourceRanges := make([]interface{}, len(got))
			for i, rangeResult := range got {
				extractedSourceRanges[i] = rangeResult.String() // Directly access string value
			}

			want := []interface{}{"0.0.0.0/0"}
			assert.Equal(t, want, extractedSourceRanges, "Source ranges mismatch")
		})
		// Target Tags
		t.Run("Target Tags", func(t *testing.T) {
			got := ruleData.Get("target_tags").Array()
			want := make([]interface{}, 0) // Initialize as an empty slice

			// Convert gjson.Result array to []interface{} for comparison
			for _, tag := range got {
				want = append(want, tag.String()) // Access string value of the tag
			}

			assert.Equal(t, want, want, "Target tags mismatch")
		})

		// Priority
		t.Run("Priority", func(t *testing.T) {
			got := ruleData.Get("priority").Int()
			want := int64(1000)
			assert.Equal(t, want, got, "Firewall rule priority mismatch")
		})

		// Allowed Protocols and Ports
		t.Run("Allowed Protocols and Ports", func(t *testing.T) {
			allowRuleData := ruleData.Get("allow").Array()[0].Map() // Get the first "allow" rule as before

			gotProtocol := allowRuleData["protocol"].String()
			wantProtocol := "tcp"
			assert.Equal(t, wantProtocol, gotProtocol, "Allow rule protocol mismatch")

			gotPorts := allowRuleData["ports"].Array() // `gotPorts` is still an array of gjson.Result
			extractedPorts := make([]interface{}, len(gotPorts))
			for i, portResult := range gotPorts {
				extractedPorts[i] = portResult.String() // Correctly extract the string value
			}

			wantPorts := []interface{}{"22", "443"}
			assert.Equal(t, wantPorts, extractedPorts, "Allow rule ports mismatch")
		})
	})

	// Clean up resources with "terraform destroy"
	terraform.Destroy(t, terraformOptions)
}

// Helper Functions
func deleteVPC(t *testing.T, projectID string, network string) {
	text := "compute"
	time.Sleep(60 * time.Second)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "delete", network, "--project=" + projectID, "--quiet"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
}

func createVPC(t *testing.T, projectID string, network string) {
	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", network, "--project=" + projectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom", "--verbosity=none"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
	time.Sleep(60 * time.Second)
}
