output "instance_public_ip" {
  description = "Public IP address of the MedSync Google Compute Engine VM"
  value       = google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip
}

output "ssh_command" {
  description = "SSH command template for the provisioned VM"
  value       = "ssh ${var.ssh_user}@${google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip}"
}

output "frontend_url" {
  description = "HTTP endpoint for the frontend"
  value       = "http://${google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip}"
}

output "grafana_url" {
  description = "Grafana dashboard URL"
  value       = "http://${google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip}:3000"
}

output "prometheus_url" {
  description = "Prometheus targets URL"
  value       = "http://${google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip}:9090"
}

output "doctor_service_health_url" {
  description = "Doctor service health endpoint"
  value       = "http://${google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip}:8080/health"
}

output "appointment_service_health_url" {
  description = "Appointment service health endpoint"
  value       = "http://${google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip}:8081/health"
}

output "auth_service_health_url" {
  description = "Auth service health endpoint"
  value       = "http://${google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip}:8082/health"
}

output "patient_service_health_url" {
  description = "Patient service health endpoint"
  value       = "http://${google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip}:8083/health"
}

output "infrastructure_summary" {
  description = "Summary of the provisioned healthcare microservices infrastructure"
  value = {
    vm_ip            = google_compute_instance.medsync.network_interface[0].access_config[0].nat_ip
    db_endpoint      = "postgres:5432"
    environment_name = var.environment_name
  }
}
