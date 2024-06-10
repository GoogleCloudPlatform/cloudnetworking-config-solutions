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

variable "psc_endpoints" {
  description = "List of service attachment configurations"
  type = list(object({
    endpoint_project_id          = string
    producer_instance_project_id = string
    producer_instance_name       = string
    subnetwork_name              = string
    network_name                 = string
    ip_address_literal           = optional(string, "")
    allow_psc_global_access      = optional(bool, false)     # Added optional field with default value of false
    labels                       = optional(map(string), {}) # Added optional labels field
  }))
  default = []
}