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

resource "google_vertex_ai_index_endpoint_deployed_index" "basic_deployed_index" {
  deployed_index_id     = var.deployed_index_id
  display_name          = var.deployed_display_name
  index_endpoint        = google_vertex_ai_index_endpoint.index_endpoint.id
  index                 = google_vertex_ai_index.index.id
  reserved_ip_ranges    = var.reserved_ip_ranges
  enable_access_logging = var.enable_access_logging
  deployment_group      = var.deployment_group

  dynamic "automatic_resources" {
    for_each = var.automatic_resources != null ? [var.automatic_resources] : []
    content {
      max_replica_count = var.automatic_resources["max_replica_count"]
      min_replica_count = var.automatic_resources["min_replica_count"]
    }
  }

  dynamic "dedicated_resources" {
    for_each = var.dedicated_resources != null ? [var.dedicated_resources] : []
    content {
      max_replica_count = var.dedicated_resources["max_replica_count"]
      min_replica_count = var.dedicated_resources["min_replica_count"]

      dynamic "machine_spec" {
        for_each = var.dedicated_resources["machine_spec"] != null ? [var.dedicated_resources["machine_spec"]] : []
        content {
          machine_type = machine_spec.value.machine_type
        }
      }
    }
  }

  dynamic "deployed_index_auth_config" {
    for_each = var.deployed_index_auth_config != null ? [var.deployed_index_auth_config] : []
    content {
      dynamic "auth_provider" {
        for_each = var.deployed_index_auth_config["auth_provider"] != null ? [var.deployed_index_auth_config["auth_provider"]] : []
        content {
          audiences       = auth_provider.value.audiences
          allowed_issuers = auth_provider.value.allowed_issuers
        }
      }
    }
  }

}
