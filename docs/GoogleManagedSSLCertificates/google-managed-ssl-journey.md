# Google Managed SSL Certificates

**On this page**

  1. Objectives

  2. Architecture

  3. Request flow

  4. Architecture Components

  5. Deploy the solution

  6. Prerequisites

  7. Deploy with "single click"

  8. Deploy through “terraform-cli”

  9. Optional: Delete the deployment

  10. Submit feedback

## Introduction

This Terraform module facilitates the creation and management of Google Compute Managed SSL Certificate. Google Compute Managed SSL Certificate are SSL/TLS certificates that you can provision, deploy, and manage for your domains. These certificates are obtained and renewed automatically by Google, simplifying the lifecycle management of your SSL/TLS configurations.

Managed SSL certificates are essential for enabling HTTPS traffic for your applications served via Google Cloud External HTTP(S) Load Balancers. By using Google-managed certificates, you delegate the complexities of certificate issuance and renewal to Google, ensuring your applications remain secure and accessible over HTTPS.

## Objective

The primary objectives of using SSL are:

*   **Encryption:** Encrypting the data that is transmitted between the user's browser and the web server.
*   **Authentication:** Verifying the identity of the website you are connecting to.
*   **Data Integrity:** Ensuring that the data has not been tamed with during transmission.

This user journey will guide you through the process of creating a Google-managed SSL certificate to achieve these objectives for your domain.

## Architecture

This diagram illustrates a multi-region deployment in Google Cloud using a Global External Load Balancer (HTTP(S)) in front of backend services distributed across three different regions.  It leverages Google's global network infrastructure and anycast IP addressing for efficient traffic distribution and high availability.

* **Scenario : this architecture describes a Google Cloud setup for a highly available and externally accessible service. The system uses a global external load balancer to distribute traffic across multiple regions.**

<a href="https://github.com/GoogleCloudPlatform/cloudnetworking-config-solutions/blob/main/docs/GoogleManagedSSLCertificates/images/ssl-certs.png" target="_blank">
  <img src="./images/ssl-certs.png" alt="SSL Certificate Architecture" width="600"/>
</a>

## Request Flow

### Architecture Components

*   **External Client:** The user or system initiating a request from outside the Google Cloud environment.
*   **External Load Balancer:** A global service that acts as a proxy for all incoming traffic.
*   **External IP Address:** The public entry point for all incoming traffic.
*   **Forwarding Rule:** Directs incoming traffic from the external IP address to the Target SSL Proxy.
*   **Target SSL Proxy:** This component handles SSL/TLS termination, using a Google-Managed SSL Certificate to ensure secure communication between the client and the load balancer.
*   **Backend Service:** A logical collection of backend instances. It monitors their health and directs traffic to them.
*   **Regional Backend Instance Groups:** These groups of virtual machine instances run the actual application. They are located in different regions and zones for high availability and fault tolerance.
*   **NAT (Network Address Translation):** Allows private instances to communicate with external resources while keeping their internal IP addresses private.

### Request Flow Steps

1.  **Request Initiation:** An **External Client** sends an HTTPS request to the **External IP Address** of the load balancer.
2.  **Traffic Forwarding:** The **Forwarding Rule** receives the request and directs it to the **Target SSL Proxy**.
3.  **SSL Termination:** The **Target SSL Proxy** uses the **Google-Managed SSL Certificate** to decrypt the request. It also performs SSL/TCP Health Checks on the backend instances to ensure they are healthy.
4.  **Backend Selection:** The **Backend Service** forwards the decrypted request to a healthy instance within a **Regional Backend Instance Group**. It automatically selects a backend based on proximity and health.
5.  **Response:** The backend instance processes the request and sends a response back to the load balancer, which then re-encrypts the response before sending it back to the client. The red dashed lines in the diagram represent this data path, while the blue dashed lines show the control path for health checks.

## Deploy the Solution

This section provides instructions on deploying the load balancer solution using Terraform.

### Prerequisites

For the common prerequisites for this repository, please refer to the **[prerequisites.md](../prerequisites.md)** guide. Any additional prerequisites specific to this user journey will be listed below.

### Deploy with "single-click"

This method uses Google Cloud Shell and Cloud Build to automate the deployment of the Network Passthrough External Load Balancer with a MIG backend.

1.  **Open in Cloud Shell:** Click the button below to clone the repository and open the necessary configuration files in the Cloud Shell editor.

    **Note**: For testing, ensure the `cloudshell_git_repo` and `cloudshell_git_branch` parameters in the URL point to your fork and specific branch where these "single click" files and the updated guide exist. For the final version, this will point to the main repository.

    <a href="https://ssh.cloud.google.com/cloudshell/editor?shellonly=true&cloudshell_git_repo=https://github.com/sridharshini1/cloudnetworking-config-solutions.git&cloudshell_git_branch=main&cloudshell_workspace=.&cloudshell_open_in_editor=configuration/bootstrap.tfvars,configuration/organization.tfvars,configuration/networking.tfvars,configuration/security/Certificates/Compute-SSL-Certs/Google-Managed/google_managed_ssl.tfvars,configuration/security/mig.tfvars,configuration/consumer/MIG/config/instance.yaml.example,configuration/consumer-load-balancing/Network/Passthrough/External/config/instance-lite.yaml.example&cloudshell_tutorial=docs/LoadBalancer/external-network-passthrough-lb-mig.md#deploy-with-single-click" target="_new">
    <img alt="Open in Cloud Shell" src="https://gstatic.com/cloudssh/images/open-btn.svg">
    </a>

2.  **Run NLB with SSL Prerequisites Script:**
    This script prepares your Google Cloud project: enables APIs, creates a Terraform state bucket for NLB, and sets Cloud Build permissions. From the root of the cloned `cloudnetworking-config-solutions` directory in Cloud Shell, run:

    ```bash
    sh docs/GoogleManagedSSLCertificates/helper-scripts/prereq-ssl.sh
    ```

    When prompted, enter your Google Cloud Project ID.

3.  **Review and Update Configuration Files:**
    The Cloud Shell editor will open key configuration files. Review each file and update values (project IDs, user IDs/groups, network names, regions, etc.) as per your requirements. Follow the guidance in the "Deploy through Terraform-cli" section of this document for details on each file:

      * `configuration/bootstrap.tfvars`
      * `configuration/organization.tfvars`
      * `configuration/networking.tfvars`
      * `configuration/security/Certificates/Compute-SSL-Certs/Google-Managed/google_managed_ssl.tfvars`
      * `configuration/security/mig.tfvars`
      * `execution/06-consumer/MIG/config/instance.yaml.example` (Rename to `instance.yaml` after updating.)
      * `execution/07-consumer-load-balancing/Network/Passthrough/External/config/instance-lite.yaml.example` (Rename to `instance-lite.yaml` after updating.)

    When prompted, enter your Google Cloud Project ID.

3.  **Review and Update Configuration Files:**
    The Cloud Shell editor will open key configuration files. Review each file and update values (project IDs, user IDs/groups, network names, regions, etc.) as per your requirements. Follow the guidance in the "Deploy through Terraform-cli" section of this document for details on each file:

      * `configuration/bootstrap.tfvars`
      * `configuration/organization.tfvars`
      * `configuration/networking.tfvars`
      * `configuration/security/Certificates/Compute-SSL-Certs/Google-Managed/google_managed_ssl.tfvars`
      * `configuration/security/mig.tfvars`
      * `execution/06-consumer/MIG/config/instance.yaml.example` (Rename to `instance.yaml` after updating.)
      * `execution/07-consumer-load-balancing/Network/Passthrough/External/config/instance-lite.yaml.example` (Rename to `instance-lite.yaml` after updating.)

4.  **Submit Cloud Build Job to Deploy NLB:**
    Once configurations are updated and prerequisites are met, submit the Cloud Build job. Ensure you are in the root of the cloned repository.

    Once configurations are updated and prerequisites are met, submit the Cloud Build job. Ensure you are in the root of the cloned repository.

    ```bash
    gcloud builds submit . --config docs/GoogleManagedSSLCertificates/build/cloudbuild-ssl.yaml --ignore-file=".gcloudignore"
    ```

5.  **Verify Deployment:**

    After the Cloud Build job completes, go to the "Load Balancing" section in the Google Cloud Console. Confirm your Network Passthrough External Load Balancer is created, with SSL certificates and the MIG is attached as a backend and healthy.

6.  **[Optional] Delete the Deployment using Cloud Build:**

    To remove all resources created by this deployment, run the destroy Cloud Build job:

    ```bash
    gcloud builds submit . --config docs/GoogleManagedSSLCertificates/build/cloudbuild-ssl-destroy.yaml --ignore-file=".gcloudignore"
    ```

### **Deploy through Terraform-cli**

1. Clone the repository containing the Terraform configuration files:

    ```bash
    git clone https://github.com/GoogleCloudPlatform/cloudnetworking-config-solutions.git
    ```

2. Navigate to **cloudnetworking-config-solutions** folder and update the files containing the configuration values
   * **00-bootstrap stage**
     * Update configuration/bootstrap.tfvars **\-** update the google cloud project IDs and the user IDs/groups in the tfvars.

        ```
        bootstrap_project_id                      = "your-project-id"
        network_hostproject_id                    = "your-project-id"
        network_serviceproject_id                 = "your-project-id"
        organization_stage_administrator          = ["user:user-example@example.com"]
        networking_stage_administrator            = ["user:user-example@example.com"]
        security_stage_administrator              = ["user:user-example@example.com"]
        producer_stage_administrator              = ["user:user-example@example.com"]
        producer_connectivity_stage_administrator = ["user:user-example@example.com"]
        consumer_stage_administrator              = ["user:user-example@example.com"]
        consumer_lb_administrator                 = ["user:lb-user-example@example.com"]
        ```

   * **01-organisation stage**
     * Update configuration/organization.tfvars \- update the google cloud project ID and the list of the APIs to enable for the MIG & Load Balancer.

        ```
            activate_api_identities = {
            "project-01" = {
                project_id = "your-project-id",
                activate_apis = [
                "servicenetworking.googleapis.com",
                "iam.googleapis.com",
                "compute.googleapis.com",
                "cloudresourcemanager.googleapis.com",
                "serviceusage.googleapis.com",
                ],
            },
            }

        ```
   * **02-networking stage**
     * Update `configuration/networking.tfvars` update the Google Cloud Project ID and the parameters for additional resources such as VPC, subnet, and NAT as outlined below.

        ```
        project_id  = "your-project-id"
        region      = "us-central1"

        ## VPC input variables
        network_name = "cncs-vpc"
        subnets = [
        {
            ip_cidr_range = "10.0.0.0/24"
            name          = "cncs-vpc-subnet-1"
            region        = "us-central1"
        }
        ]

        shared_vpc_host = false

        ## PSC/Service Connectivity variable
        create_scp_policy  = false

        ## Cloud Nat input variables
        create_nat = true

        ## Cloud HA VPN input variables

        create_havpn = false
        peer_gateways = {
        default = {
            gcp = "" # e.g. projects/<google-cloud-peer-projectid>/regions/<google-cloud-region>/vpnGateways/<peer-vpn-name>
        }
        }

        tunnel_1_router_bgp_session_range = ""
        tunnel_1_bgp_peer_asn             = 64514
        tunnel_1_bgp_peer_ip_address      = ""
        tunnel_1_shared_secret            = ""

        tunnel_2_router_bgp_session_range = ""
        tunnel_2_bgp_peer_asn             = 64514
        tunnel_2_bgp_peer_ip_address      = ""
        tunnel_2_shared_secret            = ""

        ## Cloud Interconnect input variables

        create_interconnect = false # Use true or false


   * **03-security stage**
     * Update configuration/security/Certificates/Compute-SSL-Certs/Google-Managed/google_managed_ssl.tfvars \- update the Google Cloud Project ID. This will facilitate the creation of essential firewall rules, granting required MIG firewall rules.

        ```
        project_id           = ""
        ssl_certificate_name = "my-managed-ssl-cert"
        ssl_managed_domains = [
        {
            domains = ["example.com", "www.example.com"]
        }
        ]
        ```

   * **03-security stage**
     * Update configuration/security/mig.tfvars file \- update the Google Cloud Project ID. This will facilitate the creation of essential firewall rules, granting required MIG firewall rules.

        ```
        project_id = "your-project-id"
        network    = "cncs-vpc"
        ingress_rules = {
        fw-allow-health-check = {
            deny               = false
            description        = "Allow health checks"
            destination_ranges = []
            disabled           = false
            enable_logging = {
            include_metadata = true
            }
            priority = 1000
            source_ranges = [
                "130.211.0.0/22",
                "35.191.0.0/16"
            ]
            targets = ["allow-health-checks"]
            rules = [{
            protocol = "tcp"
            ports    = ["80"]
            }]
        }
        }
        ```

   * **06-consumer stage**
     * Update the execution/06-consumer/MIG/config/instance.yaml.example file and rename it to instance.yaml

        ```
        name: minimal-mig
        project_id: your-project-id
        location: us-central1
        zone : us-central1-a
        vpc_name : cncs-vpc
        subnetwork_name : cncs-vpc-subnet-1
        named_ports:
            http: 80
        ```

    * **07-consumer-load-balancing stage**
      * Update the execution/07-consumer-load-balancing/Application/External/config/instance2.yaml.example file and rename it to instance2.yaml

        ```
        name: load-balancer-cncs
        project: your-project-id
        network: cncs-vpc
        backends:
        default:
            groups:
            - group: minimal-mig
                region: us-central1
        ```

3. **Execute the terraform script**
   You can now deploy the stages individually using **run.sh** or you can deploy all the stages automatically using the run.sh file. Navigate to the execution/ directory and run this command to run the automatic deployment using **run.sh .**

      ```
      ./run.sh -s all -t init-apply-auto-approve
      or
      ./run.sh --stage all --tfcommand init-apply-auto-approve
      ```

4. **Verify the SSL Certificate**

    Once the Terraform apply is complete, you can verify the SSL certificate in the Google Cloud Console.

    1.  Go to the **Load balancing** page in the Google Cloud Console.
    2.  Click on the **Advanced** tab.
    3.  Click on the **Certificates** tab.
    4.  You should see your newly created SSL certificate in the list.

## **Optional-Delete the deployment**

1. In Cloud Shell or in your terminal, make sure that the current working directory is $HOME/cloudshell\_open/\<Folder-name\>/execution. If it isn't, go to that directory.

2. Remove the resources that were provisioned by the solution guide:

    ```
    ./run.sh -s all -t destroy-auto-approve
    ```

Terraform displays a list of the resources that will be destroyed.

3. When you're prompted to perform the actions, enter yes.

Troubleshoot Errors
---
For common troubleshooting steps and solutions, please refer to the **[troubleshooting.md](../troubleshooting.md)** guide.

## **Submit feedback**

To provide feedback, please follow the instructions in our **[submit-feedback.md](../submit-feedback.md)** guide.