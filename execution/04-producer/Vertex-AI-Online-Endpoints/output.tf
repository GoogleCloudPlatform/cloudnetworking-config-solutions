# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

output "endpoint_configurations_from_yaml" {
  value       = local.endpoint_list
  description = "Endpoint configurations read from YAML files."
}

output "endpoint_configurations" {
  value = {
    for endpoint_key, endpoint_details in module.vertex_endpoints : endpoint_key => {
      name         = endpoint_details.endpoint_configuration.name
      display_name = endpoint_details.endpoint_configuration.display_name
      description  = endpoint_details.endpoint_configuration.description
      location     = endpoint_details.endpoint_configuration.location
      region       = endpoint_details.endpoint_configuration.region
      labels       = endpoint_details.endpoint_configuration.labels
      network      = endpoint_details.endpoint_configuration.network
    }
  }
  description = "Configuration details for all created Vertex AI endpoints."
}