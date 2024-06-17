# Manual Networking Setup for PSC

**NOTE** : Please skip this step if you are not provisioning a producer service using [Private Service Connect](https://cloud.google.com/vpc/docs/private-service-connect). 

## Overview

This stage establishes a Private Service Connect (PSC) connection between your consumer and the producer service you created in the previous "05-producer" step. This is done by creating a forwarding rule that directs traffic from a reserved IP address to the PSC attachment on your producer service (e.g., Cloud SQL database). PSC enables secure and private communication within Google Cloud Platform (GCP), shielding your services from the public internet. This stage configures a Private Service Connect (PSC) connection between your consumer project and the producer services you've set up. It does this by creating:

1. **Internal IP Addresses:** Reserved within specific subnetworks in your consumer project, acting as private endpoints for the PSC connection.
2. **Forwarding Rules:** Directing traffic destined for these internal IPs to your producer services through the PSC connection.

## How It Works

This configuration is designed for flexibility. You can define multiple producer services within your `configuration/psc-manual.tfvars` file.

**Key Points:**

* **`google_compute_address`:**  This resource will create a new internal IP address if `ip_address` is specified for a `psc_endpoint` in your `psc-manual.tfvars` file. If no address is provided, it will automatically reserve an address.
* **`google_compute_forwarding_rule`:**  This resource sets up the forwarding rule, connecting the internal IP (or the automatically created one) to the `psc_service_attachment_link` of your producer service. 

## Configuration

This stage uses a modularized approach. The main.tf file in the root directory orchestrates the creation of multiple forwarding rules based on the configuration provided in the psc-manual.tfvars file.

While running this stage, please carefully note the following details : 

- The variable `psc_endpoints`, is a list of objects, where each object represents a producer (such as Cloud SQL) instance:

    - `endpoint_project_id`           = "your-consumer-project-id"  # Project where the forwarding rule is created
    - `producer_instance_project_id`  = "your-producer-project-id"  # Project hosting the service (e.g., Cloud SQL)
    - `producer_instance_name`        = "your-sql-instance-name"    # Name of the producer service instance
    - `subnetwork_name`               = "your-subnetwork-name"     # Subnet for allocating the internal IP
    - `network_name`                  = "your-network-name"        # VPC network for the forwarding rule
    - `ip_address_literal`            = ""                 # (Optional) Specific internal IP, or leave empty string for automatic allocation
    - `allow_psc_global_access`       = true/false                  # (Optional) Allow global access
    - `labels`                        = { key = "value" }           # (Optional) Labels for resources

- * **Regions:** Ensure your producer service, subnetwork, and service attachment all reside within the same GCP region.
- * **IP Addresses:** Verify the specified `ip_address` values are available and not already in use.

## Example configuration/psc-manual.tfvars

```
psc_endpoints = [
{
endpoint_project_id          = "project-for-endpoint"
producer_instance_project_id = "project-producer"
producer_instance_name       = "sql-3"
subnetwork_name              = "subnetwork-3"
network_name                 = "network-3"
ip_address_literal           = ""
allow_psc_global_access      = false
labels                       = { environment = "dev" }
},
]
```

<!-- BEGIN_TF_DOCS -->
## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_psc_forwarding_rules"></a> [psc\_forwarding\_rules](#module\_psc\_forwarding\_rules) | ../modules/psc_forwarding_rule | n/a |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_psc_endpoints"></a> [psc\_endpoints](#input\_psc\_endpoints) | List of PSC Endpoint configurations | <pre>list(object({<br>    # The Google Cloud project ID where the Cloud SQL instance is located.<br>    producer_instance_project_id = string<br><br>    # The Google Cloud project ID where the forwarding rule and address will be created.<br>    endpoint_project_id = string<br><br>    # The name of the Cloud SQL instance to connect to.<br>    producer_instance_name = string<br><br>    # The name of the subnet where the internal IP address will be allocated.<br>    subnetwork_name = string<br><br>    # The name of the network where the forwarding rule will be created.<br>    network_name = string<br><br>    # Optional: The static internal IP address to use. If not provided,<br>    # Google Cloud will automatically allocate an IP address.<br>    ip_address_literal = string<br>    allow_psc_global_access      = optional(bool, false)<br>    labels                       = optional(map(string), {})<br>  }))</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_forwarding_rule_self_link"></a> [forwarding\_rule\_self\_link](#output\_forwarding\_rule\_self\_link) | Map of forwarding rule self-links, keyed by SQL instance name |
| <a name="output_ip_address_literal"></a> [ip\_address\_literal](#output\_ip\_address\_literal) | Map of IP addresses, keyed by SQL instance name |
<!-- END_TF_DOCS -->