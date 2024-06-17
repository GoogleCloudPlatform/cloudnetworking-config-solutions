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

variable "activate_api_identities" {
  description = "Map of objects containing information required to enable API's for the Google Cloud project."
  type = map(object({
    project_id                  = string,
    activate_apis               = list(string),
    disable_dependent_services  = optional(bool, false)
    disable_services_on_destroy = optional(bool, false)
  }))
}
