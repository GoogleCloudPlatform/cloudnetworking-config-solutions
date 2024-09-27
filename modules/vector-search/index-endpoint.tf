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

resource "google_vertex_ai_index_endpoint" "index_endpoint" {
  display_name            = var.index_endpoint_display_name
  description             = var.index_endpoint_description
  project                 = var.project_id
  region                  = var.region
  labels                  = var.index_endpoint_labels
  network                 = var.index_endpoint_network
  public_endpoint_enabled = var.public_endpoint_enabled
  dynamic "private_service_connect_config" {
    for_each = var.private_service_connect_config != null ? [var.private_service_connect_config] : []
    content {
      enable_private_service_connect = var.private_service_connect_config["enable_private_service_connect"]
      project_allowlist              = var.private_service_connect_config["project_allowlist"]
    }
  }
}

