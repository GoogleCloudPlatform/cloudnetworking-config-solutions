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

module "gke" {
  for_each = local.cluster_map
  source   = "terraform-google-modules/kubernetes-engine/google//modules/private-cluster"
  version  = "33.1.0"

  kubernetes_version = each.value.kubernetes_version

  project_id                = each.value.project_id
  name                      = each.value.name
  region                    = each.value.region
  zones                     = each.value.zones
  regional                  = each.value.regional
  description               = each.value.description
  network                   = each.value.network
  subnetwork                = each.value.subnetwork
  ip_range_pods             = each.value.ip_range_pods
  ip_range_services         = each.value.ip_range_services
  default_max_pods_per_node = each.value.default_max_pods_per_node

  http_load_balancing        = each.value.http_load_balancing
  network_policy             = each.value.network_policy
  horizontal_pod_autoscaling = each.value.horizontal_pod_autoscaling
  filestore_csi_driver       = each.value.filestore_csi_driver

  node_pools              = each.value.node_pools
  node_pools_oauth_scopes = each.value.node_pools_oauth_scopes
  node_pools_labels       = each.value.node_pools_labels
  node_pools_metadata     = each.value.node_pools_metadata

  node_pools_taints = each.value.node_pools_taints

  node_pools_tags            = each.value.node_pools_tags
  deletion_protection        = each.value.deletion_protection
  enable_private_endpoint    = each.value.enable_private_endpoint
  enable_private_nodes       = each.value.enable_private_nodes
  master_ipv4_cidr_block     = each.value.master_ipv4_cidr_block
  master_authorized_networks = each.value.master_authorized_networks

  network_project_id                       = each.value.network_project_id
  enable_vertical_pod_autoscaling          = each.value.enable_vertical_pod_autoscaling
  service_external_ips                     = each.value.service_external_ips
  datapath_provider                        = each.value.datapath_provider
  maintenance_start_time                   = each.value.maintenance_start_time
  maintenance_exclusions                   = each.value.maintenance_exclusions
  maintenance_end_time                     = each.value.maintenance_end_time
  maintenance_recurrence                   = each.value.maintenance_recurrence
  additional_ip_range_pods                 = each.value.additional_ip_range_pods
  stack_type                               = each.value.stack_type
  windows_node_pools                       = each.value.windows_node_pools
  node_pools_resource_labels               = each.value.node_pools_resource_labels
  node_pools_linux_node_configs_sysctls    = each.value.node_pools_linux_node_configs_sysctls
  enable_cost_allocation                   = each.value.enable_cost_allocation
  resource_usage_export_dataset_id         = each.value.resource_usage_export_dataset_id
  enable_network_egress_export             = each.value.enable_network_egress_export
  enable_resource_consumption_export       = each.value.enable_resource_consumption_export
  cluster_autoscaling                      = each.value.cluster_autoscaling
  network_tags                             = each.value.network_tags
  stub_domains                             = each.value.stub_domains
  upstream_nameservers                     = each.value.upstream_nameservers
  non_masquerade_cidrs                     = each.value.non_masquerade_cidrs
  ip_masq_resync_interval                  = each.value.ip_masq_resync_interval
  ip_masq_link_local                       = each.value.ip_masq_link_local
  configure_ip_masq                        = each.value.configure_ip_masq
  logging_service                          = each.value.logging_service
  monitoring_service                       = each.value.monitoring_service
  create_service_account                   = each.value.create_service_account
  grant_registry_access                    = each.value.grant_registry_access
  registry_project_ids                     = each.value.registry_project_ids
  service_account                          = each.value.service_account
  service_account_name                     = each.value.service_account_name
  boot_disk_kms_key                        = each.value.boot_disk_kms_key
  issue_client_certificate                 = each.value.issue_client_certificate
  cluster_ipv4_cidr                        = each.value.cluster_ipv4_cidr
  cluster_resource_labels                  = each.value.cluster_resource_labels
  dns_cache                                = each.value.dns_cache
  authenticator_security_group             = each.value.authenticator_security_group
  identity_namespace                       = each.value.identity_namespace
  enable_mesh_certificates                 = each.value.enable_mesh_certificates
  release_channel                          = each.value.release_channel
  gateway_api_channel                      = each.value.gateway_api_channel
  add_cluster_firewall_rules               = each.value.add_cluster_firewall_rules
  add_master_webhook_firewall_rules        = each.value.add_master_webhook_firewall_rules
  firewall_priority                        = each.value.firewall_priority
  firewall_inbound_ports                   = each.value.firewall_inbound_ports
  add_shadow_firewall_rules                = each.value.add_shadow_firewall_rules
  shadow_firewall_rules_priority           = each.value.shadow_firewall_rules_priority
  shadow_firewall_rules_log_config         = each.value.shadow_firewall_rules_log_config
  enable_confidential_nodes                = each.value.enable_confidential_nodes
  enable_cilium_clusterwide_network_policy = each.value.enable_cilium_clusterwide_network_policy
  security_posture_mode                    = each.value.security_posture_mode
  security_posture_vulnerability_mode      = each.value.security_posture_vulnerability_mode
  disable_default_snat                     = each.value.disable_default_snat
  notification_config_topic                = each.value.notification_config_topic
  notification_filter_event_type           = each.value.notification_filter_event_type
  enable_tpu                               = each.value.enable_tpu
  network_policy_provider                  = each.value.network_policy_provider
  initial_node_count                       = each.value.initial_node_count
  remove_default_node_pool                 = each.value.remove_default_node_pool
  disable_legacy_metadata_endpoints        = each.value.disable_legacy_metadata_endpoints
  database_encryption                      = each.value.database_encryption
  enable_shielded_nodes                    = each.value.enable_shielded_nodes
  enable_binary_authorization              = each.value.enable_binary_authorization
  node_metadata                            = each.value.node_metadata
  cluster_dns_provider                     = each.value.cluster_dns_provider
  cluster_dns_scope                        = each.value.cluster_dns_scope
  cluster_dns_domain                       = each.value.cluster_dns_domain
  gce_pd_csi_driver                        = each.value.gce_pd_csi_driver
  gke_backup_agent_config                  = each.value.gke_backup_agent_config
  gcs_fuse_csi_driver                      = each.value.gcs_fuse_csi_driver
  stateful_ha                              = each.value.stateful_ha
  timeouts                                 = each.value.timeouts
  monitoring_enable_managed_prometheus     = each.value.monitoring_enable_managed_prometheus
  monitoring_enable_observability_metrics  = each.value.monitoring_enable_observability_metrics
  monitoring_enabled_components            = each.value.monitoring_enabled_components
  logging_enabled_components               = each.value.logging_enabled_components
  enable_kubernetes_alpha                  = each.value.enable_kubernetes_alpha
  config_connector                         = each.value.config_connector
  enable_intranode_visibility              = each.value.enable_intranode_visibility
  enable_l4_ilb_subsetting                 = each.value.enable_l4_ilb_subsetting
  fleet_project                            = each.value.fleet_project
}