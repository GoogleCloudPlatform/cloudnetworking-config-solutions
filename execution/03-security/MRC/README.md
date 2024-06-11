# Overview

To ensure seamless communication between your Memorystore for Redis Cluster and its clients (like Compute Engine instances), it's crucial to configure proper firewall rules. These rules act as gatekeepers, controlling inbound and outbound traffic to your Redis cluster.

## Description

- Firewall Rules: To control access to your Memorystore instance, you use firewall rules. These rules determine which IP addresses or networks can connect to your Redis instance.
- Authorized Networks: When you create a Memorystore instance, you can specify authorized networks. Only clients within these networks can access your Redis instance. This is the primary way to secure your Redis instance. This can be done using the `network` variable.
- Port 6379: This is the default port Redis uses for communication. By default, Memorystore for Redis instances are not accessible from the public internet.

## Example firewall rules

Using the mrc-firewall.tf module, you can setup firewall rules to achieve inbound & outbound connections given as examples below based on your networking requirements.

- INGRESS/INBOUND : (Client to Redis Cluster)

    - Allow connections to Compute Instance (consumer) from authorised source IPs on port 22 for SSH
    - Allow connections from Compute Instance to MRC instance on port 6379

- EGRESS/OUTBOUND : (Redis Cluster to Client)

    - Allow connections to MRC instance on port 6379 
    - Allow connections to port 443 to all IP addresses (for consumers such as GCE)

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_mrc_firewall"></a> [mrc\_firewall](#module\_mrc\_firewall) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpc-firewall | v31.0.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_default_rules_config"></a> [default\_rules\_config](#input\_default\_rules\_config) | Optionally created convenience rules. Set the 'disabled' attribute to true, or individual rule attributes to empty lists to disable. | <pre>object({<br>    admin_ranges = optional(list(string))<br>    disabled     = optional(bool, true)<br>    allowed_http_cidrs = optional(list(string), [<br>      "35.191.0.0/16", "130.211.0.0/22", "209.85.152.0/22", "209.85.204.0/22"]<br>    )<br>    http_tags = optional(list(string), ["http-server"])<br>    allowed_https_cidrs = optional(list(string), [<br>      "35.191.0.0/16", "130.211.0.0/22", "209.85.152.0/22", "209.85.204.0/22"]<br>    )<br>    https_tags = optional(list(string), ["https-server"])<br>    allowed_ssh_cidrs = optional(list(string), ["35.235.240.0/20"])<br>    ssh_tags   = optional(list(string), ["ssh"])<br>  })</pre> | <pre>{<br>  "disabled": true<br>}</pre> | no |
| <a name="input_egress_rules"></a> [egress\_rules](#input\_egress\_rules) | List of egress rule definitions, default to deny action. Null destination ranges will be replaced with 0/0. | <pre>map(object({<br>    deny               = optional(bool, true)<br>    description        = optional(string)<br>    destination_ranges = optional(list(string))<br>    disabled           = optional(bool, false)<br>    enable_logging = optional(object({<br>      include_metadata = optional(bool)<br>    }))<br>    priority             = optional(number, 1000)<br>    source_ranges        = optional(list(string))<br>    targets              = optional(list(string))<br>    use_service_accounts = optional(bool, false)<br>    rules = optional(list(object({<br>      protocol = string<br>      ports    = optional(list(string))<br>    })), [{ protocol = "all" }])<br>  }))</pre> | `{}` | no |
| <a name="input_ingress_rules"></a> [ingress\_rules](#input\_ingress\_rules) | List of ingress rule definitions, default to allow action. Null source ranges will be replaced with 0/0. | <pre>map(object({<br>    deny               = optional(bool, false)<br>    description        = optional(string)<br>    destination_ranges = optional(list(string), []) # empty list is needed as default to allow deletion after initial creation with a value. See https://github.com/hashicorp/terraform-provider-google/issues/14270<br>    disabled           = optional(bool, false)<br>    enable_logging     = optional(bool, false)<br>    priority             = optional(number, 1000)<br>    source_ranges        = optional(list(string))<br>    sources              = optional(list(string))<br>    targets              = optional(list(string))<br>    use_service_accounts = optional(bool, false)<br>    rules = optional(list(object({<br>      protocol = string<br>      ports    = optional(list(string))<br>    })), [{ protocol = "all" }])<br>  }))</pre> | `{}` | no |
| <a name="input_network"></a> [network](#input\_network) | Name of the network this set of firewall rules applies to. | `string` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | Project ID of the project that holds the network to which this set of firewall rules apply to. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_mrc_firewall_rules"></a> [mrc\_firewall\_rules](#output\_mrc\_firewall\_rules) | Map of firewall rules created. |
<!-- END_TF_DOCS -->