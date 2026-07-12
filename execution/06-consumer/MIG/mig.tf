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

module "mig" {
  for_each              = local.mig_map
  source                = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-mig?ref=v41.1.0"
  project_id            = each.value.project_id
  location              = each.value.location
  name                  = each.value.name
  target_size           = each.value.target_size
  instance_template     = module.mig-template[each.key].template.self_link
  auto_healing_policies = each.value.auto_healing_policies # Handle auto_healing_policies
  health_check_config   = each.value.health_check_config   # Health check config handling
  autoscaler_config = {
    max_replicas    = try(each.value.autoscaler_config.max_replicas, var.autoscaler_config.max_replicas)
    min_replicas    = try(each.value.autoscaler_config.min_replicas, var.autoscaler_config.min_replicas)
    cooldown_period = try(each.value.autoscaler_config.cooldown_period, var.autoscaler_config.cooldown_period)

    scaling_signals = {
      cpu_utilization = {
        target                = try(each.value.autoscaler_config.scaling_signals.cpu_utilization.target, var.autoscaler_config.scaling_signals.cpu_utilization.target)
        optimize_availability = try(each.value.autoscaler_config.scaling_signals.cpu_utilization.optimize_availability, var.autoscaler_config.scaling_signals.cpu_utilization.optimize_availability)
      }
    }
  }
  named_ports = try(each.value.named_ports, var.named_ports)
}
