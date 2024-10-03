# Producer Stage

## Overview

This Producer stage is responsible for provisioning various Google Cloud Platform (GCP) producer services for consumer applications. It currently supports the creation and configuration of the following services:

- [Cloud SQL](https://cloud.google.com/sql?hl=en): Managed relational databases.
- [AlloyDB](https://cloud.google.com/alloydb?hl=en): Managed PostgreSQL-compatible database optimized for demanding workloads.
- [Memorystore for Redis](https://cloud.google.com/memorystore/docs/cluster/memorystore-for-redis-cluster-overview): Managed Redis instances for in-memory data storage and caching.
- [GKE](https://cloud.google.com/kubernetes-engine/docs) : Deploy, manage, and scale containerized applications on Kubernetes, powered by Google Cloud.
- Vertex AI

   - [Vector Search](https://cloud.google.com/vertex-ai/docs/vector-search/overview) : With Vector Search you can leverage the same infrastructure that provides a foundation for Google products such as Google Search, YouTube, and Play.
   - [Online Prediction Endpoint](https://cloud.google.com/vertex-ai/docs/predictions/overview) : Online predictions are synchronous requests made to a model that is deployed to an endpoint.

The stage utilizes Terraform modules to streamline the provisioning process and ensure consistency across different service types.

## Prerequisites

- **Completed Prior Stages:** Successful deployment of networking resources depends on the completion of the following stages:

    - **01-organization:** This stage handles the activation of required Google Cloud APIs.
    - **02-networking:** This stage handles the creation of networking resources such as VPCs, HA-VPNs etc.
    - **03-security:** This stage handles the creation of key security components such firewall rules. For GCE, the folder to use is 03-security/GCE.

- **Enable APIs** : Based on the producer that you plan to provision, ensure the following APIs should be enabled :

    - Cloud SQL Admin API: This is the primary API for managing Cloud SQL instances, including creation, configuration, backups, and more.
    - AlloyDB API: This is the core API for managing AlloyDB clusters and instances.
    - Cloud Memorystore for Redis API: This is the API for creating, configuring and managing Redis instances in Google Cloud
    - Kubernetes Engine API: This API is used for building and managing container-based applications, deployed in the google kubernetes engine.
    - Vertex AI API : This API helps to use Vertex AI services allowing us to train high-quality custom machine learning models with minimal machine learning expertise and effort.

- **Required Permissions** : Based on the producer that you plan to provision, ensure the following permissions are granted to the user runing the terraform modules :

    - [Cloud SQL Admin](https://cloud.google.com/sql/docs/mysql/iam-roles#:~:text=roles/cloudsql.admin) (roles/cloudsql.admin): This role provides full control over Cloud SQL instances, including creation, configuration, deletion, and management of databases, users, and backups.
    - [AlloyDB Admin](https://cloud.google.com/alloydb/docs/reference/iam-roles-permissions#:~:text=Description%0AAlloyDB%20permissions-,roles/alloydb.admin,-Cloud%20AlloyDB%20Admin) (roles/alloydb.admin): Grants full control over AlloyDB clusters, including creation, configuration, scaling, and management of instances.
    - [Redis Admin](https://cloud.google.com/memorystore/docs/redis/access-control#:~:text=including%20Redis%20resources-,roles/redis.admin,-Redis%20Admin) (roles/redis.admin): Provides full control over Memorystore for Redis instances, including creation, configuration, scaling, and deletion.
    - [Container Admin](https://cloud.google.com/kubernetes-engine/docs/how-to/iam#predefined) (roles/container.admin):Provides access to full management of clusters and their Kubernetes API objects.
    - [Vertex AI Admin](https://cloud.google.com/vertex-ai/docs/general/access-control#aiplatform.admin) (roles/aiplatform.admin): Grants full access to all resources in Vertex AI, used to administer and manage Vertex AI services.



## Configuration

### General Configuration Notes

- YAML Configuration Files: Place YAML files defining each instance's configuration within the config/ directory of the respective service's folder under producer folder in configuration/ (e.g., configuration/producer/CloudSQL/config/).

- Terraform Variables: You can customize the input variables in the .tf files according to your project's requirements.

**NOTE** : for producer specific configuration details, please find them in the producer's README document :

  - [CloudSQL](cloudnetworking-config-solution/execution/04-producer/CloudSQL/README.md)
  - [AlloyDB](cloudnetworking-config-solution/execution/04-producer/AlloyDB/README.md)
  - [MRC]((cloudnetworking-config-solution/execution/04-producer/MRC/README.md))
  - [GKE]((cloudnetworking-config-solution/execution/04-producer/GKE/README.md))
  - Vertex AI

     - [Vector Search]((cloudnetworking-config-solution/execution/04-producer/VectorSearch/README.md))
     - [Online Prediction Endpoint]((cloudnetworking-config-solution/execution/04-producer/Vertex-AI-Online-Endpoints/README.md))

## Execution Steps (For Each Service)

1. **Configure**: Create YAML files specifying the desired configurations for each instance.

2. **Terraform Stages**:

    - Initialize: Run `terraform init`.
    - Plan: Run `terraform plan` to preview the infrastructure changes.
    - Apply: If the plan looks good, run `terraform apply` to create or update the resources.


## Additional Notes

- **Instance configuration**: Carefully review and customize the instance configuration to match your organization's requirements.
