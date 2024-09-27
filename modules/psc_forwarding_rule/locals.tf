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

# Construct the full service attachment name for each endpoint, using the project, region, and connection name.

locals {
  service_attachment_name = {
    for k, v in var.psc_endpoints :
    k => v.producer_instance_name != null ?
    "projects/${v.producer_instance_project_id}/regions/${data.google_sql_database_instance.instance[k].region}/serviceAttachments/${data.google_sql_database_instance.instance[k].connection_name}" :
    null
  }
  forwarding_rule_targets = { for k, v in var.psc_endpoints :
    k => v.producer_instance_name != null ?
    try(data.google_sql_database_instance.instance[k].psc_service_attachment_link, null) :
    v.target
  }
}