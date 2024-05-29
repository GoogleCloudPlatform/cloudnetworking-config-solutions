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

output "cloudsql_instance_details" {
  description = "Display Cloud SQL instance attributes, including name, project ID, region, connection name, IP address (public or private), and database version."
  value = { for name, instance in module.cloudsql :
    name => {
      "name" : instance.name,
      "project_id" : instance.instances.primary.project,
      "region" : instance.instances.primary.region,
      "connection_name" : instance.connection_name,
      "database_version" : instance.instances.primary.database_version,
      "public_ip_address" : try(instance.instances.primary.public_ip_address, null)
      "private_ip_address" : try(instance.instances.primary.private_ip_address, null),
  } }
  sensitive = true
}
