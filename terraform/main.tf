terraform {
  required_version = ">= 1.6.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.gcp_region
  zone    = var.gcp_zone
}

resource "google_compute_network" "medsync" {
  name                    = "${var.project_name}-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "medsync" {
  name          = "${var.project_name}-subnet"
  ip_cidr_range = var.subnet_cidr
  region        = var.gcp_region
  network       = google_compute_network.medsync.id
}

resource "google_compute_firewall" "ssh" {
  name    = "${var.project_name}-allow-ssh"
  network = google_compute_network.medsync.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.allowed_ssh_cidr]
  target_tags   = ["medsync"]
}

resource "google_compute_firewall" "frontend" {
  name    = "${var.project_name}-allow-http"
  network = google_compute_network.medsync.name

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["medsync"]
}

resource "google_compute_firewall" "monitoring" {
  name    = "${var.project_name}-allow-monitoring"
  network = google_compute_network.medsync.name

  allow {
    protocol = "tcp"
    ports    = ["3000", "9090"]
  }

  source_ranges = [var.monitoring_allowed_cidr]
  target_tags   = ["medsync"]
}

resource "google_compute_firewall" "backend_demo" {
  name    = "${var.project_name}-allow-backend-demo"
  network = google_compute_network.medsync.name

  allow {
    protocol = "tcp"
    ports    = ["8080-8083"]
  }

  source_ranges = [var.backend_allowed_cidr]
  target_tags   = ["medsync"]
}

resource "google_compute_instance" "medsync" {
  name         = "${var.project_name}-vm"
  machine_type = var.machine_type
  zone         = var.gcp_zone
  tags         = ["medsync"]

  boot_disk {
    initialize_params {
      image = var.boot_image
      size  = var.boot_disk_size_gb
      type  = "pd-standard"
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.medsync.id

    access_config {
    }
  }

  metadata = {
    ssh-keys = "${var.ssh_user}:${file(var.ssh_public_key_path)}"
  }

  metadata_startup_script = <<-EOF
    #!/bin/bash
    set -eux
    id -u ${var.ssh_user} >/dev/null 2>&1 || useradd -m -s /bin/bash ${var.ssh_user}
    apt-get update
    apt-get install -y ca-certificates curl git
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" > /etc/apt/sources.list.d/docker.list
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    systemctl enable docker
    systemctl start docker
    usermod -aG docker ${var.ssh_user}
    mkdir -p /opt/medsync
    chown ${var.ssh_user}:${var.ssh_user} /opt/medsync
  EOF

  labels = {
    project = var.project_name
  }
}
