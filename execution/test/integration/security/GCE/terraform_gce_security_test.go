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
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
)

var (
	terraformDirectoryPath = "../../../../03-security/GCE" // Update with your GCE directory path
	projectID              = os.Getenv("TF_VAR_project_id")
	uniqueID               = rand.Int()
	networkName            = fmt.Sprintf("test-vpc-security-%d", uniqueID)
	firewallRuleName       = "allow-ssh-custom-ranges"
)

func TestGCEFirewallRuleProperties(t *testing.T) {
	var (
		tfVars = map[string]any{
			"project_id":   projectID,
			"network_name": networkName,
			"ingress_rules": []any{
				map[string]any{
					"name":        firewallRuleName,
					"description": "Allow SSH access from specific networks",
					"priority":    1000,
					"source_ranges": []string{
						"10.0.0.0/8",
						"192.168.1.0/24",
					},
					"target_tags": []string{"ssh-allowed", "https-allowed"},
					"allow": []any{
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
	createVPC(t, projectID, networkName)

	// Defer VPC deletion
	defer deleteVPC(t, projectID, networkName)

	// Terraform init and apply
	terraform.InitAndApply(t, terraformOptions)
	time.Sleep(60 * time.Second) // Wait for resource creation

	// Get Firewall rule from output
	firewallRule := terraform.OutputJson(t, terraformOptions, "firewall_rules_ingress_egress")
	if !gjson.Valid(firewallRule) {
		t.Errorf("Error parsing output, invalid JSON: %s", firewallRule)
	}
	result := gjson.Parse(firewallRule)

	// Check if the firewall rule exists in the state
	t.Run("Firewall Rule Exists", func(t *testing.T) {
		rule := result.Get(firewallRuleName)
		if !rule.Exists() {
			t.Errorf("Firewall rule %s not found in state", firewallRuleName)
		}
	})

	// Assertions
	t.Run("Firewall Properties", func(t *testing.T) {
		// Name
		t.Run("Name", func(t *testing.T) {
			got := gjson.Get(result.String(), fmt.Sprintf("%s.name", firewallRuleName)).String()
			want := firewallRuleName
			if got != want {
				t.Errorf("Firewall name mismatch: got %q, want %q", got, want)
			}
		})

		// Direction
		t.Run("Direction", func(t *testing.T) {
			got := gjson.Get(result.String(), fmt.Sprintf("%s.direction", firewallRuleName)).String()
			want := "INGRESS"
			if got != want {
				t.Errorf("Firewall rule direction mismatch: got %q, want %q", got, want)
			}
		})

		// Source Ranges (comparing contents, ignoring order)
		t.Run("Source Ranges", func(t *testing.T) {
			got := gjson.Get(result.String(), fmt.Sprintf("%s.source_ranges", firewallRuleName)).Array()
			want := tfVars["ingress_rules"].([]any)[0].(map[string]any)["source_ranges"].([]string)

			gotStr := make([]string, len(got))
			for i, v := range got {
				gotStr[i] = v.String() // Convert gjson.Result to string
			}

			// Compare contents, ignoring order
			if len(gotStr) != len(want) || !sameStringSlice(gotStr, want) {
				t.Errorf("Source ranges mismatch: got %v, want %v", gotStr, want)
			}
		})

		t.Run("Target Tags", func(t *testing.T) {
			got := gjson.Get(result.String(), fmt.Sprintf("%s.target_tags", firewallRuleName)).Array()
			want := tfVars["ingress_rules"].([]any)[0].(map[string]any)["target_tags"].([]string)

			gotStr := make([]string, len(got))
			for i, v := range got {
				gotStr[i] = v.String() // Convert gjson.Result to string
			}

			// Compare contents, ignoring order
			if len(gotStr) != len(want) || !sameStringSlice(gotStr, want) {
				t.Errorf("Target tags mismatch: got %v, want %v", gotStr, want)
			}
		})

		t.Run("Priority", func(t *testing.T) {
			got := gjson.Get(result.String(), fmt.Sprintf("%s.priority", firewallRuleName)).Int()
			var want int64 = 1000
			if got != want {
				t.Errorf("Firewall rule priority mismatch: got %d, want %d", got, want)
			}
		})

		t.Run("Allowed Protocols and Ports", func(t *testing.T) {
			allowRules := gjson.Get(result.String(), fmt.Sprintf("%s.allow", firewallRuleName)).Array()
			if len(allowRules) != 1 {
				t.Errorf("Expected 1 allow rule, but got %d", len(allowRules))
				return
			}

			rule := allowRules[0]
			gotProtocol := rule.Get("protocol").String()
			wantProtocol := "tcp"
			if gotProtocol != wantProtocol {
				t.Errorf("Allow rule protocol mismatch: got %q, want %q", gotProtocol, wantProtocol)
			}

			gotPorts := rule.Get("ports").Array()
			wantPorts := []string{"22", "443"}

			if len(gotPorts) != len(wantPorts) {
				t.Errorf("Allow rule ports mismatch: got %d ports, want %d", len(gotPorts), len(wantPorts))
			} else {
				for i, port := range gotPorts {
					if port.String() != wantPorts[i] {
						t.Errorf("Allow rule port mismatch at index %d: got %q, want %q", i, port.String(), wantPorts[i])
					}
				}
			}
		})
	})

	// Clean up resources with "terraform destroy"
	terraform.Destroy(t, terraformOptions)
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

// Helper function to compare string slices, ignoring order
func sameStringSlice(x, y []string) bool {
	if len(x) != len(y) { // Check if slices have the same length
		return false
	}

	// Create copies of the slices to sort them without modifying the originals
	xCopy := slices.Clone(x)
	yCopy := slices.Clone(y)

	// Sort both slices
	slices.Sort(xCopy)
	slices.Sort(yCopy)

	// Compare sorted slices for equality
	return slices.Equal(xCopy, yCopy)
}
