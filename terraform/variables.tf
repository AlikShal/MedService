variable "project_name" {
  description = "Name prefix for Google Cloud resources"
  type        = string
  default     = "medsync"
}

variable "project_id" {
  description = "Google Cloud project ID"
  type        = string
}

variable "gcp_region" {
  description = "Google Cloud region. Free Tier e2-micro is available in us-west1, us-central1, and us-east1."
  type        = string
  default     = "us-central1"
}

variable "gcp_zone" {
  description = "Google Cloud zone for the VM"
  type        = string
  default     = "us-central1-a"
}

variable "machine_type" {
  description = "Compute Engine machine type. e2-micro is Free Tier eligible in supported US regions."
  type        = string
  default     = "e2-micro"
}

variable "boot_image" {
  description = "Boot image for the VM"
  type        = string
  default     = "ubuntu-os-cloud/ubuntu-2204-lts"
}

variable "boot_disk_size_gb" {
  description = "Boot disk size in GB. Google Cloud Free Tier includes 30 GB-months standard persistent disk."
  type        = number
  default     = 30
}

variable "subnet_cidr" {
  description = "CIDR block for the MedSync subnet"
  type        = string
  default     = "10.10.0.0/24"
}

variable "ssh_user" {
  description = "Linux username for SSH metadata"
  type        = string
  default     = "alikhan"
}

variable "ssh_public_key_path" {
  description = "Local path to the SSH public key used to access the VM"
  type        = string
}

variable "allowed_ssh_cidr" {
  description = "CIDR block allowed to SSH into the instance"
  type        = string
  default     = "203.0.113.10/32"
}

variable "monitoring_allowed_cidr" {
  description = "CIDR block allowed to access Grafana and Prometheus"
  type        = string
  default     = "203.0.113.10/32"
}

variable "backend_allowed_cidr" {
  description = "CIDR block allowed to access demo backend ports 8080-8083"
  type        = string
  default     = "203.0.113.10/32"
}
