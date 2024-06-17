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

package unittest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"golang.org/x/exp/slices"
)

const (
	terraformDirectoryPath = "../../../../03-security/GCE"
	network                = "projects/dummy-project/global/networks/dummy-vpc-network01"
	firewallRuleName       = "allow-ssh-custom-ranges-gce-test"
)

var (
	projectID = "dummy-project-id"
	tfVars    = map[string]any{
		"project_id":   projectID,
		"network_name": network,
		"ingress_rules": []any{
			map[string]any{
				"name":        firewallRuleName,
				"description": "Allow SSH access from specific networks",
				"priority":    1000,
				"source_ranges": []string{
					"192.168.1.0/24", // Replace with your actual source ranges
					"10.0.0.0/8",
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

func TestInitAndPlanRunWithTfVars(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	want := 2
	got := planExitCode
	if got != want {
		t.Errorf("Test Plan Exit Code = %v, want = %v", got, want)
	}
}

func TestInitAndPlanRunWithoutTfVarsExpectFailureScenario(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
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

func TestResourcesCount(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCount := terraform.GetResourceCount(t, planStruct)
	if got, want := resourceCount.Add, 1; got != want {
		t.Errorf("Test Resource Count Add = %v, want = %v", got, want)
	}
	if got, want := resourceCount.Change, 0; got != want {
		t.Errorf("Test Resource Count Change = %v, want = %v", got, want)
	}
	if got, want := resourceCount.Destroy, 0; got != want {
		t.Errorf("Test Resource Count Destroy = %v, want = %v", got, want)
	}
}

func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	expectedModulesAddress := []string{"module.ssh_firewall"} // Changed module name
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planStruct := terraform.InitAndPlanAndShow(t, terraformOptions)
	content, err := terraform.ParsePlanJSON(planStruct)
	if err != nil {
		t.Fatal(err.Error())
	}
	actualModuleAddresses := make([]string, 0)
	for _, element := range content.ResourceChangesMap {
		if element.ModuleAddress != "" && !slices.Contains(actualModuleAddresses, element.ModuleAddress) {
			actualModuleAddresses = append(actualModuleAddresses, element.ModuleAddress)
		}
	}
	want := expectedModulesAddress
	got := actualModuleAddresses
	if !cmp.Equal(got, want) {
		t.Errorf("Test Element Mismatch = %v, want = %v", got, want)
	}
}
