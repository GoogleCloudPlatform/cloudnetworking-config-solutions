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

output "redis_cluster_details" {
  description = "Detailed information about each Redis cluster"
  value = {
    for name, cluster in google_redis_cluster.cluster-ha :
    name => {
      name           = cluster.name
      region         = cluster.region
      shard_count    = cluster.shard_count
      replica_count  = cluster.replica_count
      psc_connection = try(cluster.psc_connections[0].psc_connection_id, null)
      state          = cluster.state
      network        = cluster.psc_configs[0].network
    }
  }
}
