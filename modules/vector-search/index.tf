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

resource "google_vertex_ai_index" "index" {
  labels       = var.index_labels
  region       = var.region
  display_name = var.index_display_name
  project      = var.project_id
  description  = var.index_description
  metadata {
    contents_delta_uri = var.contents_delta_uri
    config {
      dimensions                  = var.dimensions
      approximate_neighbors_count = var.approximate_neighbors_count
      shard_size                  = var.shard_size
      distance_measure_type       = var.distance_measure_type
      algorithm_config {
        dynamic "tree_ah_config" {
          for_each = var.tree_ah_config != null ? [var.tree_ah_config] : []
          content {
            leaf_node_embedding_count    = lookup(var.tree_ah_config, "leaf_node_embedding_count", null)
            leaf_nodes_to_search_percent = lookup(var.tree_ah_config, "leaf_nodes_to_search_percent", null)
          }
        }
        dynamic "brute_force_config" {
          for_each = var.brute_force_config != null ? [var.brute_force_config] : []
          content {
          }
        }
      }
    }
  }
  index_update_method = var.index_update_method
}
