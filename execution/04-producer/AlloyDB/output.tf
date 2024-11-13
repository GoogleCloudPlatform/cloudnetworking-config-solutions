/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

output "cluster_details" {
  description = "Details of each AlloyDB cluster instance."
  value = {
    for instance in local.instance_list :
    instance.cluster_display_name => {
      cluster_id                    = instance.cluster_id
      cluster_display_name          = instance.cluster_display_name
      database_version              = instance.database_version
      network_id                    = instance.network_id
      region                        = instance.region
      allocated_ip_range            = try(instance.allocated_ip_range, "N/A")
      psc_enabled                   = try(instance.psc_enabled, false)
      psc_allowed_consumer_projects = try(instance.psc_allowed_consumer_projects, ["cncs-sridharshini-23", "pm-singleproject-30"])
    }
  }
}