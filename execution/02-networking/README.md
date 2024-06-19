# Networking Stage

## Overview

This Terraform configuration provides a flexible and robust framework for deploying and managing essential networking components within your Google Cloud Platform (GCP) environment. It empowers you to create a secure, highly available, and customizable network infrastructure that aligns with your organization's specific requirements.

Key features of this configuration include:

- **Virtual Private Cloud (VPC) Network**: Establish a custom VPC network tailored to your needs. Designate subnets for different purposes, manage routing tables, and leverage Private Service Access (PSA) for seamless communication with Google-managed services.- **Service Connection Policies (SCP)**: Implements Private Service Connect to automate private connectivity to a managed service.
- **High Availability VPN (HA VPN)**: Create redundant VPN tunnels for secure, resilient site-to-site connectivity. Leverage BGP for dynamic routing and optimal path selection.
- **Cloud NAT**: Enable private Google Compute Engine (GCE) instances within your VPC to access the internet while maintaining the security of private IP addresses.


### Benefits 

- High Availability: This configuration promotes high availability through redundant VPN tunnels and strategically placed subnets.
- Modularity: The modular structure of this configuration allows you to easily add or remove components as needed.
- PSA and PSC configuration : This module allows you to use either or both PSA (Service Networking) and PSC (Private Service Connectivity) for your large scale deployments.

## Prerequisites

- Before creating networking resources, ensure you have the completed the following pre-requsites:

1. **Completed Prior Stages:** Successful deployment of networking resources depends on the completion of the following stages:
    * **01-organization:** This stage handles the activation of required Google Cloud APIs.

2. Enable the following APIs : 

- [Compute Engine API](https://cloud.google.com/compute/docs/reference/rest/v1): Used for creating and managing VPC networks, subnets, forwarding rules, HA VPN tunnels/gateways, Cloud NAT and firewall rules.
- [Service Networking API](https://cloud.google.com/service-infrastructure/docs/service-networking/getting-started): to manage Private Service Access (PSA) configurations.
- [Network Connectivity API](https://cloud.google.com/network-connectivity/docs/reference/networkconnectivity/rest)
    - Enables connectivity with and between Google Cloud resources.
- [Service Consumer Management API](https://cloud.google.com/service-infrastructure/docs/service-consumer-management/reference/rest) : enabled in the project that Private Service Connect endpoints are deployed in. This API lets Google Cloud create the Network Connectivity Service Account that deploys Private Service Connect endpoints.

3. Permissions required for this stage : 

    - [Compute Network Admin](https://cloud.google.com/iam/docs/understanding-roles#compute.networkAdmin) : roles/compute.networkAdmin : Grants full control over VPC networks, subnets, firewall rules, and related resources.
    - [Compute Shared VPN Admin](https://cloud.google.com/compute/docs/access/iam#compute.xpnAdmin) : roles/compute.xpnAdmin : Permissions to administer shared VPC host projects, specifically enabling the host projects and associating shared VPC service projects to the host project's network.

## Components

- `ha-vpn.tf`: 
    - Defines the HA VPN gateway, creating two redundant tunnels for high availability.
    - Configures BGP sessions for dynamic routing between your on-premises network and GCP.
    - Manages custom route advertisement to control traffic flow.

- `nat.tf`: Sets up the Cloud NAT gateway and associates it with the VPC network.
    - Sets up a Cloud NAT gateway to provide internet access for instances in your private subnets.
    - Configures NAT routing to direct outbound traffic to the internet.

- `scp.tf`: Defines SCP for provisioning Private Service Connectivity to services like Memorystore for Redis Clusters.
    - Automates private connectivity to a managed service.

- `vpc.tf`: Creates the VPC network, subnets, and routing configuration.
    - Creates a VPC network with specified subnets and IP ranges.
    - Configures routing tables, including custom routes for PSA if required.
    - Supports Private Service Access (PSA) for private connectivity to Google-managed services.

- `locals`.tf: Defines local variables for use within the configuration.

- `output.tf`: Provides outputs for easy access to information about the created resources.

- `variables`.tf: Defines input variables for customizing the network configuration.

## Configuration

To configure networking.tfvars for your environment, here's an example which can be used for your reference : 

```
project_id = "gcp-project-id"
region     = "us-central1"

## VPC input variables

network_name = "network"
subnets = [
  {
    ip_cidr_range = "10.0.1.0/16"
    name          = "subnet1"
    region        = "us-central1-a"
  },
  {
    ip_cidr_range = "192.168.0.1/16"
    name          = "subnet2"
    region        = "us-central1-b"
  }
]

# PSC/Service Connecitvity Variables

create_scp_policy      = true
subnets_for_scp_policy = ["subnet1"] 

## Cloud Nat input variables
create_nat = true

## Cloud HA VPN input variables

create_havpn = false
peer_gateways = {
  default = {
    gcp = "" # e.g. projects/<google-cloud-peer-projectid>/regions/<google-cloud-region>/vpnGateways/<peer-vpn-name>
  }
}

tunnel_1_router_bgp_session_range = "ip-cidr-range"
tunnel_1_bgp_peer_asn             = 64514
tunnel_1_bgp_peer_ip_address      = "ip-from-the-cidr-range"
tunnel_1_shared_secret            = "secret1"

tunnel_2_router_bgp_session_range = "ip-cidr-range"
tunnel_2_bgp_peer_asn             = 64514
tunnel_2_bgp_peer_ip_address      = "ip-from-the-cidr-range"
tunnel_2_shared_secret            = "secret2"
```

## Usage

**NOTE** : run the terraform commands with the `-var-file` referencing the networking.tfvars present under the /configuration folder. Example : 

```
terraform plan -var-file=../configuration/networking.tfvars
terraform apply -var-file=../configuration/networking.tfvars
```

- Initialize: Run `terraform init`.
- Plan: Run `terraform plan -var-file=../configuration/networking.tfvars` to review the planned changes.
- Apply:  If the plan looks good, run `terraform apply -var-file=../configuration/networking.tfvars` to create or update the resources.


## Notes

- Dependencies: Ensure that the required GCP services are enabled in your project.
- Resource Names: Choose unique names to avoid conflicts.
- Security: Review the default firewall rules and SCPs to ensure they align with your security requirements.


<!-- BEGIN_TF_DOCS -->

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | 5.32.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_havpn"></a> [havpn](#module\_havpn) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpn-ha | v30.0.0 |
| <a name="module_nat"></a> [nat](#module\_nat) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-cloudnat | v30.0.0 |
| <a name="module_vpc_network"></a> [vpc\_network](#module\_vpc\_network) | ../modules/net-vpc | n/a |

## Resources

| Name | Type |
|------|------|
| [google_compute_route.default](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_route) | resource |
| [google_network_connectivity_service_connection_policy.policy](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/network_connectivity_service_connection_policy) | resource |
| [google_compute_network.vpc_network](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/compute_network) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_advertise_all_subnets"></a> [advertise\_all\_subnets](#input\_advertise\_all\_subnets) | Set to true if all subnets are required to be advertised. | `bool` | `false` | no |
| <a name="input_create_havpn"></a> [create\_havpn](#input\_create\_havpn) | Set to true to create a Cloud HA VPN. | `string` | `"false"` | no |
| <a name="input_create_nat"></a> [create\_nat](#input\_create\_nat) | Set to true to create a Cloud NAT. | `string` | `"true"` | no |
| <a name="input_create_network"></a> [create\_network](#input\_create\_network) | Variable to determine if a new network should be created or not. | `bool` | `true` | no |
| <a name="input_create_scp_policy"></a> [create\_scp\_policy](#input\_create\_scp\_policy) | Boolean flat to create a service connection policy. Set to true to create a service connection policy | `bool` | `false` | no |
| <a name="input_create_subnetwork"></a> [create\_subnetwork](#input\_create\_subnetwork) | Variable to determine if a new sub network should be created or not. | `bool` | `true` | no |
| <a name="input_delete_default_routes_on_create"></a> [delete\_default\_routes\_on\_create](#input\_delete\_default\_routes\_on\_create) | Set to true to delete the default routes at creation time. | `bool` | `true` | no |
| <a name="input_deletion_policy"></a> [deletion\_policy](#input\_deletion\_policy) | The deletion policy for the service networking connection. Setting to ABANDON allows the resource to be abandoned rather than deleted. This will enable a successful terraform destroy when destroying CloudSQL instances. Use with care as it can lead to dangling resources. | `string` | `""` | no |
| <a name="input_destination_range"></a> [destination\_range](#input\_destination\_range) | The destination range of outgoing packets that this route applies to. Only IPv4 is supported. | `string` | `"0.0.0.0/0"` | no |
| <a name="input_export_custom_routes"></a> [export\_custom\_routes](#input\_export\_custom\_routes) | Whether to export the custom routes to the peer network. | `bool` | `true` | no |
| <a name="input_firewall_policy_enforcement_order"></a> [firewall\_policy\_enforcement\_order](#input\_firewall\_policy\_enforcement\_order) | Order that Firewall Rules and Firewall Policies are evaluated. Can be either 'BEFORE\_CLASSIC\_FIREWALL' or 'AFTER\_CLASSIC\_FIREWALL'. | `string` | `"AFTER_CLASSIC_FIREWALL"` | no |
| <a name="input_ha_vpn_gateway1_name"></a> [ha\_vpn\_gateway1\_name](#input\_ha\_vpn\_gateway1\_name) | VPN Gateway name and prefix used for dependent resources. | `string` | `"vpn1"` | no |
| <a name="input_import_custom_routes"></a> [import\_custom\_routes](#input\_import\_custom\_routes) | Whether to import the custom routes to the peer network. | `bool` | `true` | no |
| <a name="input_nat_name"></a> [nat\_name](#input\_nat\_name) | Name of the Cloud NAT to be created. | `string` | `"internet-gateway"` | no |
| <a name="input_network_name"></a> [network\_name](#input\_network\_name) | Name of the VPC network to be created if var.create\_network is marked as true or Name of the already existing network if var.create\_network is false. | `string` | n/a | yes |
| <a name="input_next_hop_gateway"></a> [next\_hop\_gateway](#input\_next\_hop\_gateway) | URL to a gateway that should handle matching packets. Currently, you can only specify the internet gateway, using a full or partial valid URL. | `string` | `"default-internet-gateway"` | no |
| <a name="input_peer_gateways"></a> [peer\_gateways](#input\_peer\_gateways) | Configuration of the (external or GCP) peer gateway. | <pre>map(object({<br>    external = optional(object({<br>      redundancy_type = string<br>      interfaces      = list(string)<br>      description     = optional(string, "Terraform managed external VPN gateway")<br>    }))<br>    gcp = optional(string)<br>  }))</pre> | `{}` | no |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The project ID of the Google Cloud project where the VPC will be created. | `string` | n/a | yes |
| <a name="input_psa_range"></a> [psa\_range](#input\_psa\_range) | Variable to describe the CIDR range required by the PSA/Service Networking. | `string` | `"10.0.64.0/20"` | no |
| <a name="input_psa_range_name"></a> [psa\_range\_name](#input\_psa\_range\_name) | Variable to describe the name of the CIDR range required by the PSA/Service Networking. | `string` | `"psarange"` | no |
| <a name="input_region"></a> [region](#input\_region) | Name of a Google Cloud region. | `string` | n/a | yes |
| <a name="input_router1_asn"></a> [router1\_asn](#input\_router1\_asn) | ASN number required for the router1. | `number` | `64513` | no |
| <a name="input_scp_connection_limit"></a> [scp\_connection\_limit](#input\_scp\_connection\_limit) | Limit of the total number of connections to be allowed through the policy | `string` | `5` | no |
| <a name="input_service_class"></a> [service\_class](#input\_service\_class) | Allowed service class (static) | `string` | `"gcp-memorystore-redis"` | no |
| <a name="input_shared_vpc_host"></a> [shared\_vpc\_host](#input\_shared\_vpc\_host) | Enable shared VPC for this project. | `bool` | `true` | no |
| <a name="input_shared_vpc_service_projects"></a> [shared\_vpc\_service\_projects](#input\_shared\_vpc\_service\_projects) | Shared VPC service projects to register with this host. | `list(string)` | `[]` | no |
| <a name="input_subnets"></a> [subnets](#input\_subnets) | Subnet configuration. | <pre>list(object({<br>    name                  = string<br>    ip_cidr_range         = string<br>    region                = string<br>    description           = optional(string)<br>    enable_private_access = optional(bool, true)<br>    flow_logs_config = optional(object({<br>      aggregation_interval = optional(string)<br>      filter_expression    = optional(string)<br>      flow_sampling        = optional(number)<br>      metadata             = optional(string)<br>      # only if metadata == "CUSTOM_METADATA"<br>      metadata_fields = optional(list(string))<br>    }))<br>    ipv6 = optional(object({<br>      access_type = optional(string, "INTERNAL")<br>      # this field is marked for internal use in the API documentation<br>      # enable_private_access = optional(string)<br>    }))<br>    secondary_ip_ranges = optional(map(string))<br><br>    iam = optional(map(list(string)), {})<br>    iam_bindings = optional(map(object({<br>      role    = string<br>      members = list(string)<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>    iam_bindings_additive = optional(map(object({<br>      member = string<br>      role   = string<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>  }))</pre> | `[]` | no |
| <a name="input_subnets_for_scp_policy"></a> [subnets\_for\_scp\_policy](#input\_subnets\_for\_scp\_policy) | List of subnet names to apply the SCP policy to. | `list(string)` | <pre>[<br>  "default"<br>]</pre> | no |
| <a name="input_tunnel_1_bgp_peer_asn"></a> [tunnel\_1\_bgp\_peer\_asn](#input\_tunnel\_1\_bgp\_peer\_asn) | Peer BGP Autonomous System Number (ASN). | `number` | n/a | yes |
| <a name="input_tunnel_1_bgp_peer_ip_address"></a> [tunnel\_1\_bgp\_peer\_ip\_address](#input\_tunnel\_1\_bgp\_peer\_ip\_address) | Peer IP address of the BGP interface outside Google Cloud. Only IPv4 is supported. | `string` | n/a | yes |
| <a name="input_tunnel_1_gateway_interface"></a> [tunnel\_1\_gateway\_interface](#input\_tunnel\_1\_gateway\_interface) | The interface ID of the VPN gateway with which this VPN tunnel is associated. | `number` | `0` | no |
| <a name="input_tunnel_1_router_bgp_session_range"></a> [tunnel\_1\_router\_bgp\_session\_range](#input\_tunnel\_1\_router\_bgp\_session\_range) | IP address and range of the interface. | `string` | `"169.254.1.2/30"` | no |
| <a name="input_tunnel_1_shared_secret"></a> [tunnel\_1\_shared\_secret](#input\_tunnel\_1\_shared\_secret) | Shared secret used to set the secure session between the Cloud VPN gateway and the peer VPN gateway. Note: This property is sensitive and should be preserved carefully. | `string` | n/a | yes |
| <a name="input_tunnel_2_bgp_peer_asn"></a> [tunnel\_2\_bgp\_peer\_asn](#input\_tunnel\_2\_bgp\_peer\_asn) | Peer BGP Autonomous System Number (ASN). | `number` | n/a | yes |
| <a name="input_tunnel_2_bgp_peer_ip_address"></a> [tunnel\_2\_bgp\_peer\_ip\_address](#input\_tunnel\_2\_bgp\_peer\_ip\_address) | Peer IP address of the BGP interface outside Google Cloud. Only IPv4 is supported. | `string` | n/a | yes |
| <a name="input_tunnel_2_gateway_interface"></a> [tunnel\_2\_gateway\_interface](#input\_tunnel\_2\_gateway\_interface) | The interface ID of the VPN gateway with which this VPN tunnel is associated. | `number` | `1` | no |
| <a name="input_tunnel_2_router_bgp_session_range"></a> [tunnel\_2\_router\_bgp\_session\_range](#input\_tunnel\_2\_router\_bgp\_session\_range) | IP address and range of the interface. | `string` | `"169.254.2.2/30"` | no |
| <a name="input_tunnel_2_shared_secret"></a> [tunnel\_2\_shared\_secret](#input\_tunnel\_2\_shared\_secret) | Shared secret used to set the secure session between the Cloud VPN gateway and the peer VPN gateway. Note: This property is sensitive and should be preserved carefully. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_name"></a> [name](#output\_name) | Name of the VPC network. |
| <a name="output_network_id"></a> [network\_id](#output\_network\_id) | Fully qualified network ID. |
| <a name="output_service_connection_policy_details"></a> [service\_connection\_policy\_details](#output\_service\_connection\_policy\_details) | Detailed information about each service connection policy |
| <a name="output_service_connection_policy_ids"></a> [service\_connection\_policy\_ids](#output\_service\_connection\_policy\_ids) | Map of service class to service connection policy IDs |
| <a name="output_subnet_ids"></a> [subnet\_ids](#output\_subnet\_ids) | List of fully qualified subnetwork IDs. |
| <a name="output_subnet_self_links_for_scp_policy"></a> [subnet\_self\_links\_for\_scp\_policy](#output\_subnet\_self\_links\_for\_scp\_policy) | The self-links of the subnets where the SCP policy is applied. |
| <a name="output_vpc_networks"></a> [vpc\_networks](#output\_vpc\_networks) | Complete details of the VPC network. |
<!-- END_TF_DOCS -->