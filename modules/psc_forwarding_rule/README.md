<!-- BEGIN_TF_DOCS -->
Copyright 2024 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.7.0 |
| <a name="requirement_google"></a> [google](#requirement\_google) | >= 5.11.0, < 6.0.0 |
| <a name="requirement_google-beta"></a> [google-beta](#requirement\_google-beta) | >= 5.11.0, < 6.0.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | >= 5.11.0, < 6.0.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_compute_address.psc_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address) | resource |
| [google_compute_forwarding_rule.psc_forwarding_rule](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_forwarding_rule) | resource |
| [google_sql_database_instance.instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/sql_database_instance) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_psc_endpoints"></a> [psc\_endpoints](#input\_psc\_endpoints) | List of service attachment configurations | <pre>list(object({<br>    endpoint_project_id      = string<br>    producer_instance_project_id = string <br>    producer_instance_name = string<br>    subnetwork_name        = string<br>    network_name           = string<br>    ip_address_literal     = optional(string) <br>    allow_psc_global_access      = optional(bool, false)  # Added optional field with default value of false<br>    labels                       = optional(map(string), {})  # Added optional labels field<br>  }))</pre> | `[]` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_address_self_link"></a> [address\_self\_link](#output\_address\_self\_link) | Self-links of the created addresses |
| <a name="output_forwarding_rule_self_link"></a> [forwarding\_rule\_self\_link](#output\_forwarding\_rule\_self\_link) | Self-links of the created forwarding rules |
| <a name="output_ip_address_literal"></a> [ip\_address\_literal](#output\_ip\_address\_literal) | IP addresses of the created addresses |
<!-- END_TF_DOCS -->