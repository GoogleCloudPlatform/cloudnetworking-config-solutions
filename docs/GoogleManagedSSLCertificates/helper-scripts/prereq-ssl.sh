# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

if [ -z "$GOOGLE_CLOUD_PROJECT" ]
then
   echo "Project not set!"
   echo "What project do you want to deploy the solution to?"
   read var_project_id
   gcloud config set project $var_project_id
   export GOOGLE_CLOUD_PROJECT=$var_project_id
fi

echo "Running prerequisites on project $GOOGLE_CLOUD_PROJECT for NLB Passthrough deployment."
BUCKET_NAME_NLB=$GOOGLE_CLOUD_PROJECT-tf-state-nlb-passthrough # More specific bucket name
if gsutil ls gs://$BUCKET_NAME_NLB; then
    echo "Terraform bucket gs://$BUCKET_NAME_NLB already created!"
else
    echo "Creating Terraform state bucket gs://$BUCKET_NAME_NLB..."
    gcloud storage buckets create gs://$BUCKET_NAME_NLB --project=$GOOGLE_CLOUD_PROJECT --uniform-bucket-level-access
fi

# Create provider.tf for each stage if they don't exist or overwrite
# These paths must align with how run.sh executes terraform for each stage.

# 01-organization stage
mkdir -p execution/01-organization
cat > execution/01-organization/providers.tf << EOF
terraform {
  backend "gcs" {
    bucket                      = "$BUCKET_NAME_NLB"
    prefix                      = "nlb_organization_stage"
  }
}
EOF

# 02-networking stage
mkdir -p execution/02-networking
cat > execution/02-networking/providers.tf << EOF
terraform {
  backend "gcs" {
    bucket                      = "$BUCKET_NAME_NLB"
    prefix                      = "nlb_networking_stage"
  }
}
EOF

# 03-security/Certificates/Compute-SSL-Certs/Google-Managed stage (Note: directory is for SSL certs)
mkdir -p execution/03-security/Certificates/Compute-SSL-Certs/Google-Managed
cat > execution/03-security/Certificates/Compute-SSL-Certs/Google-Managedproviders.tf << EOF
terraform {
  backend "gcs" {
    bucket                      = "$BUCKET_NAME_NLB"
    prefix                      = "nlb_security_mig_stage"
  }
}
EOF



# 06-consumer/MIG stage
mkdir -p execution/06-consumer/MIG
cat > execution/06-consumer/MIG/providers.tf << EOF
terraform {
  backend "gcs" {
    bucket                      = "$BUCKET_NAME_NLB"
    prefix                      = "nlb_consumer_mig_stage"
  }
}
EOF

# 07-consumer-load-balancing/Network/Passthrough/External stage
mkdir -p execution/07-consumer-load-balancing/Network/Passthrough/External
cat > execution/07-consumer-load-balancing/Network/Passthrough/External/providers.tf << EOF
terraform {
  backend "gcs" {
    bucket                      = "$BUCKET_NAME_NLB"
    prefix                      = "nlb_consumer_lb_ext_passthrough_stage"
  }
}
EOF

echo "Enabling required APIs for NLB Passthrough..."
gcloud services enable cloudbuild.googleapis.com \
    cloudresourcemanager.googleapis.com \
    iam.googleapis.com \
    logging.googleapis.com \
    storage.googleapis.com \
    compute.googleapis.com \
    servicenetworking.googleapis.com \
    serviceusage.googleapis.com --project=$GOOGLE_CLOUD_PROJECT

echo "Granting Cloud Build's Service Account IAM roles to deploy NLB resources..."
PROJECT_NUMBER=$(gcloud projects describe $GOOGLE_CLOUD_PROJECT --format='value(projectNumber)')
CLOUDBUILD_SA="$PROJECT_NUMBER@cloudbuild.gserviceaccount.com"

# Grant roles necessary for Cloud Build to execute Terraform and manage resources for NLB.
# Consider refining these roles based on the principle of least privilege for production.
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member="serviceAccount:$CLOUDBUILD_SA" --role="roles/editor" --condition=None
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member="serviceAccount:$CLOUDBUILD_SA" --role="roles/iam.securityAdmin" --condition=None # For managing IAM policies if TF does so
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member="serviceAccount:$CLOUDBUILD_SA" --role="roles/compute.networkAdmin" --condition=None # For VPC, subnets, FWs, LBs
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member="serviceAccount:$CLOUDBUILD_SA" --role="roles/compute.instanceAdmin.v1" --condition=None # For MIGs and instances
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member="serviceAccount:$CLOUDBUILD_SA" --role="roles/iam.serviceAccountUser" --condition=None # To use service accounts
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member="serviceAccount:$CLOUDBUILD_SA" --role="roles/serviceusage.serviceUsageAdmin" --condition=None # To enable APIs if TF does so
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member="serviceAccount:$CLOUDBUILD_SA" --role="roles/storage.admin" --condition=None # For GCS backend state bucket

echo "NLB Prerequisites script completed successfully!"
echo "Ensure your Terraform configuration files (e.g., configuration/*.tfvars, execution/**/config/*.yaml.example) are correctly filled before running the Cloud Build job."