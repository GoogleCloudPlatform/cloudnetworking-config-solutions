# Consumer Stage

## Overview

This Consumer stage is responsible for provisioning consumer service instances such as Google Compute Engine (GCE) virtual machines. It uses Terraform modules to manage the creation and configuration of consumers such as VMs based on input parameters defined in YAML files.

The stage is designed to be highly flexible. For GCE, it allows customizations such as instance type, boot disk, network configuration, and attached storage.


## Prerequisites


- **Completed Prior Stages:** Successful deployment of networking resources depends on the completion of the following stages:

    - **01-organization:** This stage handles the activation of required Google Cloud APIs.
    - **02-networking:** This stage handles the creation of networking resources such as VPCs, HA-VPNs etc.
    - **03-security:** This stage handles the creation of key security components such firewall rules. For GCE, the folder to use is 03-security/GCE.

- Enable the following APIs :

    - [Compute Engine API](https://cloud.google.com/compute/docs/reference/rest/v1): Used for creating and managing GCE VMs.

- Permissions required :

    - [Compute Admin role](https://cloud.google.com/compute/docs/access/iam#compute.admin) : Used to create and manage GCE VMs.
    - [Service Account User](https://cloud.google.com/compute/docs/access/iam#iam.serviceAccountUser) : Lets a principal attach a service account to a resource.

## Configuration

### General Configuration Notes

- YAML Configuration Files: Place YAML files defining each instance's configuration within the configs/ directory of the respective service's folder (e.g., configuration/consumer/GCE/config).

- Terraform Variables: You can customize the input variables in the .tf files according to your project's requirements.

Configurations would be different for different consumer services as listed below :

1. GCE : For configuration of the GCE VM, you can read more in the [GCE README]((cloudnetworking-config-solution/execution/06-consumer/GCE/README.md)).

For every consumer, you can define .yaml files for the consumer configuration. With every .yaml file in the configs/ folder, our terraform module would create an instance. For an example, for GCE an example yaml files to create two instances are :

- instance1.yaml :

  ```
  name: instance1
  project_id: <project-id>
  region: us-central1
  zone: us-central1-a
  image: ubuntu-os-cloud/ubuntu-2204-lts
  network: projects/<project-id>/global/networks/<network-name>
  subnetwork: projects/<project-id>/regions/us-central1/subnetworks/<subnetwork-name>
  ```

- instance2.yaml :

  ```
  name: instance2
  project_id: <project-id>
  region: us-central1
  zone: us-central1-a
  image: ubuntu-os-cloud/ubuntu-2204-lts
  network: projects/<project-id>/global/networks/<network-name>
  subnetwork: projects/<project-id>/regions/us-central1/subnetworks/<subnetwork-name>
  ```

## Execution Steps

1. **Input/Configure** the yaml files based on your requirements.

2. **Terraform Stages** :

    - Initialize: Run `terraform init`.
    - Plan: Run `terraform plan` to review the planned changes.
    - Apply:  If the plan looks good, run `terraform apply` to create or update the resources.


## Additional Notes

- **Instance configuration**: Carefully review and customize the instance configuration to match your organization's requirements.
