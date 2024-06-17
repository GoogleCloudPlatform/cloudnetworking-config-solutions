# Overview

This Terraform module simplifies the process of creating and managing firewall rules to enable secure SSH access to your Google Compute Engine (GCE) instances. SSH (Secure Shell) is the standard protocol for remote administration of your servers, and configuring proper firewall rules is essential to protect your instances from unauthorized access.

## Key Features

- Firewall Rule Automation: Effortlessly define and deploy firewall rules specific to SSH traffic.
- Port 22 (SSH) Focus: Automatically opens the standard SSH port (TCP 22) for inbound connections, prioritizing security.
- Flexible Configuration:
    - Control allowed source IP addresses or ranges (source_ranges).
    - Apply rules to specific instances using target tags (target_tags).
    - Optionally leverage source tags (source_tags) or service accounts for filtering.
- Integration with Existing Networks: Works seamlessly with your existing GCP networks.
- Google Provider Integration: Leverages the official Terraform Google provider for reliability.


## Description

- Firewall Rules: Define how network traffic is filtered for your GCE instances. This module focuses on creating a firewall rule that specifically allows inbound SSH connections.
- Port 22 (SSH): The standard port used for SSH communication. This module's primary function is to open this port for incoming connections while maintaining security best practices.
- Source Ranges: Specify the IP addresses or ranges permitted to initiate SSH connections to your instances.
- Target Tags: Apply the firewall rule to specific GCE instances by tagging them with a designated label (e.g., "ssh-allowed").

## Example Firewall Rules

- INGRESS (Client to GCE Instance):

    - Allow SSH connections (TCP port 22) from authorized IP addresses or ranges (e.g., your office network or specific public IPs).

### Outputs

The module provides an output `firewall_rules_ingress_egress` listing the details of the created firewall rules for easier reference and management.

## Security Best Practices

- Least Privilege: Restrict source_ranges to only the IP addresses that genuinely require SSH access.
- SSH Keys: Utilize SSH key-based authentication instead of passwords for enhanced security.
- OS-Level Firewall: Implement a firewall within your GCE instance's operating system (e.g., iptables or ufw) for layered protection.
- Regular Audits: Periodically review your firewall rules to ensure they align with your security requirements.

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_ssh_firewall"></a> [ssh\_firewall](#module\_ssh\_firewall) | terraform-google-modules/network/google//modules/firewall-rules | 9.1.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_ingress_rules"></a> [ingress\_rules](#input\_ingress\_rules) | List of ingress rules. This will be ignored if variable 'rules' is non-empty | <pre>list(object({<br>    name                    = string<br>    description             = optional(string, null)<br>    disabled                = optional(bool, null)<br>    priority                = optional(number, null)<br>    destination_ranges      = optional(list(string), [])<br>    source_ranges           = optional(list(string), [])<br>    source_tags             = optional(list(string))<br>    source_service_accounts = optional(list(string))<br>    target_tags             = optional(list(string))<br>    target_service_accounts = optional(list(string))<br><br>    allow = optional(list(object({<br>      protocol = string<br>      ports    = optional(list(string))<br>    })), [])<br>    deny = optional(list(object({<br>      protocol = string<br>      ports    = optional(list(string))<br>    })), [])<br>    log_config = optional(object({<br>      metadata = string<br>    }))<br>  }))</pre> | n/a | yes |
| <a name="input_network_name"></a> [network\_name](#input\_network\_name) | The name (or self-link) of the network to create the firewall rule in. | `string` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The ID of the Google Cloud project. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_firewall_rules_ingress_egress"></a> [firewall\_rules\_ingress\_egress](#output\_firewall\_rules\_ingress\_egress) | List of ingress/egress firewall rule(s) created. |
<!-- END_TF_DOCS -->