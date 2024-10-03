## Introduction

This aims to automate the deployment of private Google Kubernetes Engine (GKE) clusters. It leverages the `terraform-google-modules/kubernetes-engine/google` module to simplify cluster creation and configuration. The configuration is designed to be flexible and scalable, allowing you to define multiple clusters with varying settings through YAML configuration files.

## Pre-Requisites

Before deploying your GKE clusters, ensure you have the following:

1. **Google Cloud Project:** A Google Cloud Platform project with billing enabled.
2. **Terraform Installed:** Terraform CLI installed on your local machine. Download the appropriate version for your OS from the official Terraform website: [https://www.terraform.io/downloads.html](https://www.terraform.io/downloads.html).
3. **Service Account & Permissions:**  A service account with permissions to manage GKE clusters. The service account should have at least the following roles:
    * `roles/container.clusterAdmin`
    * `roles/compute.networkAdmin` 
    * `roles/iam.serviceAccountUser` (if using a custom service account)
4. **Network Configuration:** A VPC network and subnetwork(s) to host your GKE clusters.
5. **IP Address Ranges:**  Reserved IP address ranges for your pods and services within the chosen subnetwork(s).
6. **gcloud CLI:**  The Google Cloud SDK (`gcloud`) command-line tool installed. You'll use this to authenticate and manage your GCP resources. Download and install from here: [https://cloud.google.com/sdk/docs/install](https://cloud.google.com/sdk/docs/install)

## Let's Get Started! 🚀

This project uses YAML files to define cluster configurations, promoting modularity and ease of management. Here’s how to deploy your GKE clusters:

### Execution Steps

1. **Configure Cluster Definitions:**

   * Navigate to the `config` directory.
   * Create YAML files (e.g., `cluster1.yaml`, `cluster2.yaml`) defining the desired configuration for each cluster. Ensure the filename does not start with an underscore (`_`). 
   * See the [Example](#example) section or the provided example file (`_cluster.yaml.example`) for configuration options. Customize these files with your project ID, network settings, node pool configurations, and other relevant parameters.

2. **Initialize Terraform:**

   * Open your terminal.
   * Navigate to the root directory of this Terraform project.
   * Run `terraform init` to initialize Terraform and download the necessary providers.

3. **Review the Execution Plan:**

   * Use `terraform plan` to generate a detailed execution plan. This step displays the infrastructure changes Terraform will apply based on your configuration files. 
   * Carefully inspect the output to ensure it aligns with your intended cluster setup. 

   ```
    terraform plan -var-file=../../../configuration/producer/GKE/gke.tfvars
    ```

4. **Deploy the Clusters:**

   * Execute `terraform apply` to deploy your GKE clusters. Terraform will provision the clusters based on your YAML definitions.

   ```
    terraform apply -var-file=../../../configuration/producer/GKE/gke.tfvars
    ```

5. **Access & Manage Your Clusters:**
    * Once deployed, you can manage your GKE clusters using:
        * **Google Cloud Console:**  Provides a web-based interface to interact with your clusters.
        * **`gcloud` CLI:** Offers command-line control over your GKE resources.  
        * **`kubectl` CLI:** The primary tool for deploying and managing applications within your Kubernetes clusters.

### Example

Here's an example of a cluster configuration file (`cluster1.yaml`) illustrating some key settings:

```yaml
project_id: your-project-id
name: cluster1
region: us-central1
zones: 
  - us-central1-a
  - us-central1-b
  - us-central1-c
network: your-vpc-network 
subnetwork: your-subnet
description: "Production GKE cluster"
kubernetes_version: 1.27
ip_range_pods: pod-range-name
ip_range_services: services-range-name
node_pools:
  - name: pool-1
    machine_type: n1-standard-2
    node_locations: "us-central1-a,us-central1-b" 
    min_count: 3
    max_count: 10
    disk_size_gb: 50 
    disk_type: "pd-standard" 
enable_private_endpoint: true
enable_private_nodes: true 
master_ipv4_cidr_block: 172.16.0.0/28
```

**NOTE** : 

1. The GKE version chosen would be resposible for the type of control plane connectivity. For reference, check GKE [release notes](https://cloud.google.com/kubernetes-engine/docs/release-notes-new-features).
2. If you're creating Subnet secondary IP address range for Pods and Services for GKE cluster as a producer please refer to the official documentation for [Pods](https://cloud.google.com/kubernetes-engine/docs/concepts/alias-ips#cluster_sizing_secondary_range_pods) and [Services](https://cloud.google.com/kubernetes-engine/docs/concepts/alias-ips#cluster_sizing_secondary_range_pods).

## Important Notes:

* **Customization:** Modify the provided example YAML file or create new ones to match your specific requirements for each GKE cluster.
* **Security:** For production deployments, carefully review and implement appropriate security measures, including network policies, IAM roles, and secrets management.
* **Documentation:** Refer to the official Google Cloud documentation for detailed information on GKE and the `terraform-google-modules/kubernetes-engine/google` module:
    * [Google Kubernetes Engine Documentation](https://cloud.google.com/kubernetes-engine/docs/)
    * [`terraform-google-modules/kubernetes-engine/google` Module](https://registry.terraform.io/modules/terraform-google-modules/kubernetes-engine/google/latest)
* **Cleanup:** To destroy the GKE clusters and associated resources created by this Terraform project, execute `terraform destroy`. 

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_gke"></a> [gke](#module\_gke) | terraform-google-modules/kubernetes-engine/google//modules/private-cluster | n/a |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_add_cluster_firewall_rules"></a> [add\_cluster\_firewall\_rules](#input\_add\_cluster\_firewall\_rules) | Create additional firewall rules | `bool` | `false` | no |
| <a name="input_add_master_webhook_firewall_rules"></a> [add\_master\_webhook\_firewall\_rules](#input\_add\_master\_webhook\_firewall\_rules) | Create master\_webhook firewall rules for ports defined in `firewall_inbound_ports` | `bool` | `false` | no |
| <a name="input_add_shadow_firewall_rules"></a> [add\_shadow\_firewall\_rules](#input\_add\_shadow\_firewall\_rules) | Create GKE shadow firewall (the same as default firewall rules with firewall logs enabled). | `bool` | `false` | no |
| <a name="input_additional_ip_range_pods"></a> [additional\_ip\_range\_pods](#input\_additional\_ip\_range\_pods) | List of _names_ of the additional secondary subnet ip ranges to use for pods | `list(string)` | `[]` | no |
| <a name="input_authenticator_security_group"></a> [authenticator\_security\_group](#input\_authenticator\_security\_group) | The name of the RBAC security group for use with Google security groups in Kubernetes RBAC. Group name must be in format gke-security-groups@yourdomain.com | `string` | `null` | no |
| <a name="input_boot_disk_kms_key"></a> [boot\_disk\_kms\_key](#input\_boot\_disk\_kms\_key) | The Customer Managed Encryption Key used to encrypt the boot disk attached to each node in the node pool, if not overridden in `node_pools`. This should be of the form projects/[KEY\_PROJECT\_ID]/locations/[LOCATION]/keyRings/[RING\_NAME]/cryptoKeys/[KEY\_NAME]. For more information about protecting resources with Cloud KMS Keys please see: https://cloud.google.com/compute/docs/disks/customer-managed-encryption | `string` | `null` | no |
| <a name="input_cluster_autoscaling"></a> [cluster\_autoscaling](#input\_cluster\_autoscaling) | Cluster autoscaling configuration. See [more details](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.locations.clusters#clusterautoscaling) | <pre>object({<br>    enabled                     = bool<br>    autoscaling_profile         = string<br>    min_cpu_cores               = number<br>    max_cpu_cores               = number<br>    min_memory_gb               = number<br>    max_memory_gb               = number<br>    gpu_resources               = list(object({ resource_type = string, minimum = number, maximum = number }))<br>    auto_repair                 = bool<br>    auto_upgrade                = bool<br>    disk_size                   = optional(number)<br>    disk_type                   = optional(string)<br>    image_type                  = optional(string)<br>    strategy                    = optional(string)<br>    max_surge                   = optional(number)<br>    max_unavailable             = optional(number)<br>    node_pool_soak_duration     = optional(string)<br>    batch_soak_duration         = optional(string)<br>    batch_percentage            = optional(number)<br>    batch_node_count            = optional(number)<br>    enable_secure_boot          = optional(bool, false)<br>    enable_integrity_monitoring = optional(bool, true)<br>  })</pre> | <pre>{<br>  "auto_repair": true,<br>  "auto_upgrade": true,<br>  "autoscaling_profile": "BALANCED",<br>  "disk_size": 100,<br>  "disk_type": "pd-standard",<br>  "enable_integrity_monitoring": true,<br>  "enable_secure_boot": false,<br>  "enabled": false,<br>  "gpu_resources": [],<br>  "image_type": "COS_CONTAINERD",<br>  "max_cpu_cores": 0,<br>  "max_memory_gb": 0,<br>  "min_cpu_cores": 0,<br>  "min_memory_gb": 0<br>}</pre> | no |
| <a name="input_cluster_dns_domain"></a> [cluster\_dns\_domain](#input\_cluster\_dns\_domain) | The suffix used for all cluster service records. | `string` | `""` | no |
| <a name="input_cluster_dns_provider"></a> [cluster\_dns\_provider](#input\_cluster\_dns\_provider) | Which in-cluster DNS provider should be used. PROVIDER\_UNSPECIFIED (default) or PLATFORM\_DEFAULT or CLOUD\_DNS. | `string` | `"PROVIDER_UNSPECIFIED"` | no |
| <a name="input_cluster_dns_scope"></a> [cluster\_dns\_scope](#input\_cluster\_dns\_scope) | The scope of access to cluster DNS records. DNS\_SCOPE\_UNSPECIFIED (default) or CLUSTER\_SCOPE or VPC\_SCOPE. | `string` | `"DNS_SCOPE_UNSPECIFIED"` | no |
| <a name="input_cluster_ipv4_cidr"></a> [cluster\_ipv4\_cidr](#input\_cluster\_ipv4\_cidr) | The IP address range of the kubernetes pods in this cluster. Default is an automatically assigned CIDR. | `string` | `null` | no |
| <a name="input_cluster_resource_labels"></a> [cluster\_resource\_labels](#input\_cluster\_resource\_labels) | The GCE resource labels (a map of key/value pairs) to be applied to the cluster | `map(string)` | `{}` | no |
| <a name="input_config_connector"></a> [config\_connector](#input\_config\_connector) | Whether ConfigConnector is enabled for this cluster. | `bool` | `false` | no |
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding GKE configuration values. | `string` | `"./config"` | no |
| <a name="input_configure_ip_masq"></a> [configure\_ip\_masq](#input\_configure\_ip\_masq) | Enables the installation of ip masquerading, which is usually no longer required when using aliasied IP addresses. IP masquerading uses a kubectl call, so when you have a private cluster, you will need access to the API server. | `bool` | `false` | no |
| <a name="input_create_service_account"></a> [create\_service\_account](#input\_create\_service\_account) | Defines if service account specified to run nodes should be created. | `bool` | `true` | no |
| <a name="input_database_encryption"></a> [database\_encryption](#input\_database\_encryption) | Application-layer Secrets Encryption settings. The object format is {state = string, key\_name = string}. Valid values of state are: "ENCRYPTED"; "DECRYPTED". key\_name is the name of a CloudKMS key. | `list(object({ state = string, key_name = string }))` | <pre>[<br>  {<br>    "key_name": "",<br>    "state": "DECRYPTED"<br>  }<br>]</pre> | no |
| <a name="input_datapath_provider"></a> [datapath\_provider](#input\_datapath\_provider) | The desired datapath provider for this cluster. By default, `DATAPATH_PROVIDER_UNSPECIFIED` enables the IPTables-based kube-proxy implementation. `ADVANCED_DATAPATH` enables Dataplane-V2 feature. | `string` | `"DATAPATH_PROVIDER_UNSPECIFIED"` | no |
| <a name="input_default_max_pods_per_node"></a> [default\_max\_pods\_per\_node](#input\_default\_max\_pods\_per\_node) | The maximum number of pods to schedule per node | `number` | `110` | no |
| <a name="input_deletion_protection"></a> [deletion\_protection](#input\_deletion\_protection) | Whether or not to allow Terraform to destroy the cluster. | `bool` | `true` | no |
| <a name="input_deploy_using_private_endpoint"></a> [deploy\_using\_private\_endpoint](#input\_deploy\_using\_private\_endpoint) | A toggle for Terraform and kubectl to connect to the master's internal IP address during deployment. | `bool` | `false` | no |
| <a name="input_description"></a> [description](#input\_description) | The description of the cluster | `string` | `"GKE Cluster CNCS"` | no |
| <a name="input_disable_default_snat"></a> [disable\_default\_snat](#input\_disable\_default\_snat) | Whether to disable the default SNAT to support the private use of public IP addresses | `bool` | `false` | no |
| <a name="input_disable_legacy_metadata_endpoints"></a> [disable\_legacy\_metadata\_endpoints](#input\_disable\_legacy\_metadata\_endpoints) | Disable the /0.1/ and /v1beta1/ metadata server endpoints on the node. Changing this value will cause all node pools to be recreated. | `bool` | `true` | no |
| <a name="input_dns_cache"></a> [dns\_cache](#input\_dns\_cache) | The status of the NodeLocal DNSCache addon. | `bool` | `false` | no |
| <a name="input_enable_binary_authorization"></a> [enable\_binary\_authorization](#input\_enable\_binary\_authorization) | Enable BinAuthZ Admission controller | `bool` | `false` | no |
| <a name="input_enable_cilium_clusterwide_network_policy"></a> [enable\_cilium\_clusterwide\_network\_policy](#input\_enable\_cilium\_clusterwide\_network\_policy) | Enable Cilium Cluster Wide Network Policies on the cluster | `bool` | `false` | no |
| <a name="input_enable_confidential_nodes"></a> [enable\_confidential\_nodes](#input\_enable\_confidential\_nodes) | An optional flag to enable confidential node config. | `bool` | `false` | no |
| <a name="input_enable_cost_allocation"></a> [enable\_cost\_allocation](#input\_enable\_cost\_allocation) | Enables Cost Allocation Feature and the cluster name and namespace of your GKE workloads appear in the labels field of the billing export to BigQuery | `bool` | `false` | no |
| <a name="input_enable_intranode_visibility"></a> [enable\_intranode\_visibility](#input\_enable\_intranode\_visibility) | Whether Intra-node visibility is enabled for this cluster. This makes same node pod to pod traffic visible for VPC network | `bool` | `false` | no |
| <a name="input_enable_kubernetes_alpha"></a> [enable\_kubernetes\_alpha](#input\_enable\_kubernetes\_alpha) | Whether to enable Kubernetes Alpha features for this cluster. Note that when this option is enabled, the cluster cannot be upgraded and will be automatically deleted after 30 days. | `bool` | `false` | no |
| <a name="input_enable_l4_ilb_subsetting"></a> [enable\_l4\_ilb\_subsetting](#input\_enable\_l4\_ilb\_subsetting) | Enable L4 ILB Subsetting on the cluster | `bool` | `false` | no |
| <a name="input_enable_mesh_certificates"></a> [enable\_mesh\_certificates](#input\_enable\_mesh\_certificates) | Controls the issuance of workload mTLS certificates. When enabled the GKE Workload Identity Certificates controller and node agent will be deployed in the cluster. Requires Workload Identity. | `bool` | `false` | no |
| <a name="input_enable_network_egress_export"></a> [enable\_network\_egress\_export](#input\_enable\_network\_egress\_export) | Whether to enable network egress metering for this cluster. If enabled, a daemonset will be created in the cluster to meter network egress traffic. | `bool` | `false` | no |
| <a name="input_enable_private_endpoint"></a> [enable\_private\_endpoint](#input\_enable\_private\_endpoint) | Whether the master's internal IP address is used as the cluster endpoint | `bool` | `false` | no |
| <a name="input_enable_private_nodes"></a> [enable\_private\_nodes](#input\_enable\_private\_nodes) | Whether nodes have internal IP addresses only | `bool` | `false` | no |
| <a name="input_enable_resource_consumption_export"></a> [enable\_resource\_consumption\_export](#input\_enable\_resource\_consumption\_export) | Whether to enable resource consumption metering on this cluster. When enabled, a table will be created in the resource export BigQuery dataset to store resource consumption data. The resulting table can be joined with the resource usage table or with BigQuery billing export. | `bool` | `true` | no |
| <a name="input_enable_shielded_nodes"></a> [enable\_shielded\_nodes](#input\_enable\_shielded\_nodes) | Enable Shielded Nodes features on all nodes in this cluster | `bool` | `true` | no |
| <a name="input_enable_tpu"></a> [enable\_tpu](#input\_enable\_tpu) | Enable Cloud TPU resources in the cluster. WARNING: changing this after cluster creation is destructive! | `bool` | `false` | no |
| <a name="input_enable_vertical_pod_autoscaling"></a> [enable\_vertical\_pod\_autoscaling](#input\_enable\_vertical\_pod\_autoscaling) | Vertical Pod Autoscaling automatically adjusts the resources of pods controlled by it | `bool` | `false` | no |
| <a name="input_filestore_csi_driver"></a> [filestore\_csi\_driver](#input\_filestore\_csi\_driver) | The status of the Filestore CSI driver addon, which allows the usage of filestore instance as volumes | `bool` | `false` | no |
| <a name="input_firewall_inbound_ports"></a> [firewall\_inbound\_ports](#input\_firewall\_inbound\_ports) | List of TCP ports for admission/webhook controllers. Either flag `add_master_webhook_firewall_rules` or `add_cluster_firewall_rules` (also adds egress rules) must be set to `true` for inbound-ports firewall rules to be applied. | `list(string)` | <pre>[<br>  "8443",<br>  "9443",<br>  "15017"<br>]</pre> | no |
| <a name="input_firewall_priority"></a> [firewall\_priority](#input\_firewall\_priority) | Priority rule for firewall rules | `number` | `1000` | no |
| <a name="input_fleet_project"></a> [fleet\_project](#input\_fleet\_project) | (Optional) Register the cluster with the fleet in this project. | `string` | `null` | no |
| <a name="input_gateway_api_channel"></a> [gateway\_api\_channel](#input\_gateway\_api\_channel) | The gateway api channel of this cluster. Accepted values are `CHANNEL_STANDARD` and `CHANNEL_DISABLED`. | `string` | `null` | no |
| <a name="input_gce_pd_csi_driver"></a> [gce\_pd\_csi\_driver](#input\_gce\_pd\_csi\_driver) | Whether this cluster should enable the Google Compute Engine Persistent Disk Container Storage Interface (CSI) Driver. | `bool` | `true` | no |
| <a name="input_gcs_fuse_csi_driver"></a> [gcs\_fuse\_csi\_driver](#input\_gcs\_fuse\_csi\_driver) | Whether GCE FUSE CSI driver is enabled for this cluster. | `bool` | `false` | no |
| <a name="input_gke_backup_agent_config"></a> [gke\_backup\_agent\_config](#input\_gke\_backup\_agent\_config) | Whether Backup for GKE agent is enabled for this cluster. | `bool` | `false` | no |
| <a name="input_grant_registry_access"></a> [grant\_registry\_access](#input\_grant\_registry\_access) | Grants created cluster-specific service account storage.objectViewer and artifactregistry.reader roles. | `bool` | `false` | no |
| <a name="input_horizontal_pod_autoscaling"></a> [horizontal\_pod\_autoscaling](#input\_horizontal\_pod\_autoscaling) | Enable horizontal pod autoscaling addon | `bool` | `true` | no |
| <a name="input_http_load_balancing"></a> [http\_load\_balancing](#input\_http\_load\_balancing) | Enable httpload balancer addon | `bool` | `false` | no |
| <a name="input_identity_namespace"></a> [identity\_namespace](#input\_identity\_namespace) | The workload pool to attach all Kubernetes service accounts to. (Default value of `enabled` automatically sets project-based pool `[project_id].svc.id.goog`) | `string` | `"enabled"` | no |
| <a name="input_initial_node_count"></a> [initial\_node\_count](#input\_initial\_node\_count) | The number of nodes to create in this cluster's default node pool. | `number` | `0` | no |
| <a name="input_ip_masq_link_local"></a> [ip\_masq\_link\_local](#input\_ip\_masq\_link\_local) | Whether to masquerade traffic to the link-local prefix (169.254.0.0/16). | `bool` | `false` | no |
| <a name="input_ip_masq_resync_interval"></a> [ip\_masq\_resync\_interval](#input\_ip\_masq\_resync\_interval) | The interval at which the agent attempts to sync its ConfigMap file from the disk. | `string` | `"60s"` | no |
| <a name="input_issue_client_certificate"></a> [issue\_client\_certificate](#input\_issue\_client\_certificate) | Issues a client certificate to authenticate to the cluster endpoint. To maximize the security of your cluster, leave this option disabled. Client certificates don't automatically rotate and aren't easily revocable. WARNING: changing this after cluster creation is destructive! | `bool` | `false` | no |
| <a name="input_kubernetes_version"></a> [kubernetes\_version](#input\_kubernetes\_version) | The Kubernetes version of the masters. If set to 'latest' it will pull latest available version in the selected region. | `string` | `"latest"` | no |
| <a name="input_logging_enabled_components"></a> [logging\_enabled\_components](#input\_logging\_enabled\_components) | List of services to monitor: SYSTEM\_COMPONENTS, WORKLOADS. Empty list is default GKE configuration. | `list(string)` | `[]` | no |
| <a name="input_logging_service"></a> [logging\_service](#input\_logging\_service) | The logging service that the cluster should write logs to. Available options include logging.googleapis.com, logging.googleapis.com/kubernetes (beta), and none | `string` | `"logging.googleapis.com/kubernetes"` | no |
| <a name="input_maintenance_end_time"></a> [maintenance\_end\_time](#input\_maintenance\_end\_time) | Time window specified for recurring maintenance operations in RFC3339 format | `string` | `""` | no |
| <a name="input_maintenance_exclusions"></a> [maintenance\_exclusions](#input\_maintenance\_exclusions) | List of maintenance exclusions. A cluster can have up to three | `list(object({ name = string, start_time = string, end_time = string, exclusion_scope = string }))` | `[]` | no |
| <a name="input_maintenance_recurrence"></a> [maintenance\_recurrence](#input\_maintenance\_recurrence) | Frequency of the recurring maintenance window in RFC5545 format. | `string` | `""` | no |
| <a name="input_maintenance_start_time"></a> [maintenance\_start\_time](#input\_maintenance\_start\_time) | Time window specified for daily or recurring maintenance operations in RFC3339 format | `string` | `"05:00"` | no |
| <a name="input_master_authorized_networks"></a> [master\_authorized\_networks](#input\_master\_authorized\_networks) | List of master authorized networks. If none are provided, disallow external access (except the cluster node IPs, which GKE automatically whitelists). | `list(object({ cidr_block = string, display_name = string }))` | `[]` | no |
| <a name="input_master_ipv4_cidr_block"></a> [master\_ipv4\_cidr\_block](#input\_master\_ipv4\_cidr\_block) | The IP range in CIDR notation to use for the hosted master network. Optional for Autopilot clusters. | `string` | `"10.0.0.0/28"` | no |
| <a name="input_monitoring_enable_managed_prometheus"></a> [monitoring\_enable\_managed\_prometheus](#input\_monitoring\_enable\_managed\_prometheus) | Configuration for Managed Service for Prometheus. Whether or not the managed collection is enabled. | `bool` | `false` | no |
| <a name="input_monitoring_enable_observability_metrics"></a> [monitoring\_enable\_observability\_metrics](#input\_monitoring\_enable\_observability\_metrics) | Whether or not the advanced datapath metrics are enabled. | `bool` | `false` | no |
| <a name="input_monitoring_enabled_components"></a> [monitoring\_enabled\_components](#input\_monitoring\_enabled\_components) | List of services to monitor: SYSTEM\_COMPONENTS, WORKLOADS. Empty list is default GKE configuration. | `list(string)` | `[]` | no |
| <a name="input_monitoring_observability_metrics_relay_mode"></a> [monitoring\_observability\_metrics\_relay\_mode](#input\_monitoring\_observability\_metrics\_relay\_mode) | Mode used to make advanced datapath metrics relay available. | `string` | `null` | no |
| <a name="input_monitoring_service"></a> [monitoring\_service](#input\_monitoring\_service) | The monitoring service that the cluster should write metrics to. Automatically send metrics from pods in the cluster to the Google Cloud Monitoring API. VM metrics will be collected by Google Compute Engine regardless of this setting Available options include monitoring.googleapis.com, monitoring.googleapis.com/kubernetes (beta) and none | `string` | `"monitoring.googleapis.com/kubernetes"` | no |
| <a name="input_network_policy"></a> [network\_policy](#input\_network\_policy) | Enable network policy addon | `bool` | `false` | no |
| <a name="input_network_policy_provider"></a> [network\_policy\_provider](#input\_network\_policy\_provider) | The network policy provider. | `string` | `"CALICO"` | no |
| <a name="input_network_project_id"></a> [network\_project\_id](#input\_network\_project\_id) | The project ID of the shared VPC's host (for shared vpc support) | `string` | `""` | no |
| <a name="input_network_tags"></a> [network\_tags](#input\_network\_tags) | (Optional) - List of network tags applied to auto-provisioned node pools. | `list(string)` | `[]` | no |
| <a name="input_node_metadata"></a> [node\_metadata](#input\_node\_metadata) | Specifies how node metadata is exposed to the workload running on the node | `string` | `"GKE_METADATA"` | no |
| <a name="input_node_pools"></a> [node\_pools](#input\_node\_pools) | List of maps containing node pools | `list(map(any))` | <pre>[<br>  {<br>    "auto_repair": true,<br>    "auto_upgrade": true,<br>    "disk_size_gb": 100,<br>    "disk_type": "pd-standard",<br>    "enable_gcfs": false,<br>    "enable_gvnic": false,<br>    "gpu_driver_version": "LATEST",<br>    "gpu_sharing_strategy": "TIME_SHARING",<br>    "image_type": "COS_CONTAINERD",<br>    "initial_node_count": 10,<br>    "local_ssd_count": 0,<br>    "logging_variant": "DEFAULT",<br>    "machine_type": "e2-medium",<br>    "max_count": 100,<br>    "max_shared_clients_per_gpu": 2,<br>    "min_count": 1,<br>    "name": "default-node-pool-again",<br>    "node_locations": "us-central1-b,us-central1-c",<br>    "preemptible": false,<br>    "spot": false<br>  }<br>]</pre> | no |
| <a name="input_node_pools_labels"></a> [node\_pools\_labels](#input\_node\_pools\_labels) | Map of maps containing node labels by node-pool name | `map(map(string))` | <pre>{<br>  "all": {},<br>  "default-node-pool": {<br>    "default-node-pool": true<br>  }<br>}</pre> | no |
| <a name="input_node_pools_linux_node_configs_sysctls"></a> [node\_pools\_linux\_node\_configs\_sysctls](#input\_node\_pools\_linux\_node\_configs\_sysctls) | Map of maps containing linux node config sysctls by node-pool name | `map(map(string))` | <pre>{<br>  "all": {},<br>  "default-node-pool": {}<br>}</pre> | no |
| <a name="input_node_pools_metadata"></a> [node\_pools\_metadata](#input\_node\_pools\_metadata) | Map of maps containing node metadata by node-pool name | `map(map(string))` | <pre>{<br>  "all": {},<br>  "default-node-pool": {}<br>}</pre> | no |
| <a name="input_node_pools_oauth_scopes"></a> [node\_pools\_oauth\_scopes](#input\_node\_pools\_oauth\_scopes) | Map of lists containing node oauth scopes by node-pool name | `map(list(string))` | <pre>{<br>  "all": [<br>    "https://www.googleapis.com/auth/logging.write",<br>    "https://www.googleapis.com/auth/monitoring"<br>  ]<br>}</pre> | no |
| <a name="input_node_pools_resource_labels"></a> [node\_pools\_resource\_labels](#input\_node\_pools\_resource\_labels) | Map of maps containing resource labels by node-pool name | `map(map(string))` | <pre>{<br>  "all": {},<br>  "default-node-pool": {}<br>}</pre> | no |
| <a name="input_node_pools_tags"></a> [node\_pools\_tags](#input\_node\_pools\_tags) | Map of lists containing node network tags by node-pool name | `map(list(string))` | <pre>{<br>  "all": [],<br>  "default-node-pool": [<br>    "default-node-pool-again"<br>  ]<br>}</pre> | no |
| <a name="input_node_pools_taints"></a> [node\_pools\_taints](#input\_node\_pools\_taints) | Map of lists containing node taints by node-pool name | `map(list(object({ key = string, value = string, effect = string })))` | <pre>{<br>  "all": [],<br>  "default-node-pool": [<br>    {<br>      "effect": "PREFER_NO_SCHEDULE",<br>      "key": "default-node-pool-again",<br>      "value": true<br>    }<br>  ]<br>}</pre> | no |
| <a name="input_non_masquerade_cidrs"></a> [non\_masquerade\_cidrs](#input\_non\_masquerade\_cidrs) | List of strings in CIDR notation that specify the IP address ranges that do not use IP masquerading. | `list(string)` | <pre>[<br>  "10.0.0.0/8",<br>  "172.16.0.0/12",<br>  "192.168.0.0/16"<br>]</pre> | no |
| <a name="input_notification_config_topic"></a> [notification\_config\_topic](#input\_notification\_config\_topic) | The desired Pub/Sub topic to which notifications will be sent by GKE. Format is projects/{project}/topics/{topic}. | `string` | `""` | no |
| <a name="input_notification_filter_event_type"></a> [notification\_filter\_event\_type](#input\_notification\_filter\_event\_type) | Choose what type of notifications you want to receive. If no filters are applied, you'll receive all notification types. Can be used to filter what notifications are sent. Accepted values are UPGRADE\_AVAILABLE\_EVENT, UPGRADE\_EVENT, and SECURITY\_BULLETIN\_EVENT. | `list(string)` | `[]` | no |
| <a name="input_region"></a> [region](#input\_region) | The region to host the cluster in (optional if zonal cluster / required if regional) | `string` | `"us-central1"` | no |
| <a name="input_regional"></a> [regional](#input\_regional) | Whether is a regional cluster (zonal cluster if set false. WARNING: changing this after cluster creation is destructive!) | `bool` | `true` | no |
| <a name="input_registry_project_ids"></a> [registry\_project\_ids](#input\_registry\_project\_ids) | Projects holding Google Container Registries. If empty, we use the cluster project. If a service account is created and the `grant_registry_access` variable is set to `true`, the `storage.objectViewer` and `artifactregsitry.reader` roles are assigned on these projects. | `list(string)` | `[]` | no |
| <a name="input_release_channel"></a> [release\_channel](#input\_release\_channel) | The release channel of this cluster. Accepted values are `UNSPECIFIED`, `RAPID`, `REGULAR` and `STABLE`. Defaults to `REGULAR`. | `string` | `"REGULAR"` | no |
| <a name="input_remove_default_node_pool"></a> [remove\_default\_node\_pool](#input\_remove\_default\_node\_pool) | Remove default node pool while setting up the cluster | `bool` | `false` | no |
| <a name="input_resource_usage_export_dataset_id"></a> [resource\_usage\_export\_dataset\_id](#input\_resource\_usage\_export\_dataset\_id) | The ID of a BigQuery Dataset for using BigQuery as the destination of resource usage export. | `string` | `""` | no |
| <a name="input_security_posture_mode"></a> [security\_posture\_mode](#input\_security\_posture\_mode) | Security posture mode.  Accepted values are `DISABLED` and `BASIC`. Defaults to `DISABLED`. | `string` | `"DISABLED"` | no |
| <a name="input_security_posture_vulnerability_mode"></a> [security\_posture\_vulnerability\_mode](#input\_security\_posture\_vulnerability\_mode) | Security posture vulnerability mode.  Accepted values are `VULNERABILITY_DISABLED`, `VULNERABILITY_BASIC`, and `VULNERABILITY_ENTERPRISE`. Defaults to `VULNERABILITY_DISABLED`. | `string` | `"VULNERABILITY_DISABLED"` | no |
| <a name="input_service_account"></a> [service\_account](#input\_service\_account) | The service account to run nodes as if not overridden in `node_pools`. The create\_service\_account variable default value (true) will cause a cluster-specific service account to be created. This service account should already exists and it will be used by the node pools. If you wish to only override the service account name, you can use service\_account\_name variable. | `string` | `""` | no |
| <a name="input_service_account_name"></a> [service\_account\_name](#input\_service\_account\_name) | The name of the service account that will be created if create\_service\_account is true. If you wish to use an existing service account, use service\_account variable. | `string` | `""` | no |
| <a name="input_service_external_ips"></a> [service\_external\_ips](#input\_service\_external\_ips) | Whether external ips specified by a service will be allowed in this cluster | `bool` | `false` | no |
| <a name="input_shadow_firewall_rules_log_config"></a> [shadow\_firewall\_rules\_log\_config](#input\_shadow\_firewall\_rules\_log\_config) | The log\_config for shadow firewall rules. You can set this variable to `null` to disable logging. | <pre>object({<br>    metadata = string<br>  })</pre> | <pre>{<br>  "metadata": "INCLUDE_ALL_METADATA"<br>}</pre> | no |
| <a name="input_shadow_firewall_rules_priority"></a> [shadow\_firewall\_rules\_priority](#input\_shadow\_firewall\_rules\_priority) | The firewall priority of GKE shadow firewall rules. The priority should be less than default firewall, which is 1000. | `number` | `999` | no |
| <a name="input_stack_type"></a> [stack\_type](#input\_stack\_type) | The stack type to use for this cluster. Either `IPV4` or `IPV4_IPV6`. Defaults to `IPV4`. | `string` | `"IPV4"` | no |
| <a name="input_stateful_ha"></a> [stateful\_ha](#input\_stateful\_ha) | Whether the Stateful HA Addon is enabled for this cluster. | `bool` | `false` | no |
| <a name="input_stub_domains"></a> [stub\_domains](#input\_stub\_domains) | Map of stub domains and their resolvers to forward DNS queries for a certain domain to an external DNS server | `map(list(string))` | `{}` | no |
| <a name="input_timeouts"></a> [timeouts](#input\_timeouts) | Timeout for cluster operations. | `map(string)` | `{}` | no |
| <a name="input_upstream_nameservers"></a> [upstream\_nameservers](#input\_upstream\_nameservers) | If specified, the values replace the nameservers taken by default from the node’s /etc/resolv.conf | `list(string)` | `[]` | no |
| <a name="input_windows_node_pools"></a> [windows\_node\_pools](#input\_windows\_node\_pools) | List of maps containing Windows node pools | `list(map(string))` | `[]` | no |
| <a name="input_zones"></a> [zones](#input\_zones) | The zones to host the cluster in (optional if regional cluster / required if zonal) | `list(string)` | <pre>[<br>  "us-central1-a",<br>  "us-central1-b",<br>  "us-central1-c"<br>]</pre> | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_gke_clusters"></a> [gke\_clusters](#output\_gke\_clusters) | Details of the provisioned GKE clusters |

<!-- END_TF_DOCS -->