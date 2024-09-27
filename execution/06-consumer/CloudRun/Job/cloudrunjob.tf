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

module "cloud_run_job" {
  source                 = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/cloud-run-v2?ref=v34.1.0"
  for_each               = { for job in local.instance_list : job.name => job }
  project_id             = each.value.project_id
  region                 = each.value.region
  name                   = each.value.name
  containers             = each.value.containers
  create_job             = each.value.create_job
  custom_audiences       = each.value.custom_audiences
  encryption_key         = each.value.encryption_key
  eventarc_triggers      = each.value.eventarc_triggers
  iam                    = each.value.iam
  ingress                = each.value.ingress
  labels                 = each.value.labels
  launch_stage           = each.value.launch_stage
  prefix                 = each.value.prefix
  revision               = each.value.revision
  service_account        = each.value.service_account
  service_account_create = each.value.service_account_create
  tag_bindings           = each.value.tag_bindings
  volumes                = each.value.volumes
  vpc_connector_create   = each.value.vpc_connector_create
}
