project_name = "medsync"
project_id = "project-69c9d34a-92bd-47c1-873"
gcp_region = "us-central1"
gcp_zone   = "us-central1-a"
machine_type      = "e2-micro"
boot_image        = "ubuntu-os-cloud/ubuntu-2204-lts"
boot_disk_size_gb = 30

# Generate this with: ssh-keygen -t rsa -b 4096 -f $env:USERPROFILE/.ssh/medsync_gcp
ssh_user            = "alikhan"
ssh_public_key_path = "C:/Users/Alikhan/.ssh/medsync_gcp.pub"

# Restrict access to trusted source CIDRs only.
allowed_ssh_cidr        = "203.0.113.10/32"
monitoring_allowed_cidr = "203.0.113.10/32"
backend_allowed_cidr    = "203.0.113.10/32"
