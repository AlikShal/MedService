variable "appointment_service_replicas" {
  description = "Controls horizontal scaling for appointment-service, the most critical business service in the healthcare platform."
  type        = number
  default     = 2
}

variable "vm_size" {
  description = "VM machine type to use when resizing infrastructure capacity instead of hardcoding a machine type."
  type        = string
  default     = "e2-micro"
}

variable "environment_name" {
  description = "Deployment environment name included in infrastructure summary outputs."
  type        = string
  default     = "development"
}

locals {
  resource_tiers = {
    small = {
      cpu       = 1
      memory_gb = 1
      services  = ["auth-service", "doctor-service", "patient-service", "chat-service"]
    }

    medium = {
      cpu       = 2
      memory_gb = 2
      services  = ["appointment-service"]
    }
  }
}

# Horizontal scaling for the most critical business service can be demonstrated with:
# docker-compose up --scale appointment-service=2
#
# In production, Kubernetes Horizontal Pod Autoscaler (HPA) would replace this manual
# Docker Compose scaling by automatically adjusting appointment-service replicas based
# on metrics such as CPU utilization, request rate, or latency.
