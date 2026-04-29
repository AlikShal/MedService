# MedSync Platform

MedSync is a containerized medical scheduling platform built to satisfy the Assignment 4-5 brief using a healthcare domain instead of the generic product/order example from the PDF.

## What is included

- `auth-service`: registration, login, JWT-based authentication, role support, seeded admin account.
- `patient-service`: patient profile creation and lookup for authenticated users.
- `doctor-service`: doctor directory with admin-only doctor creation.
- `appointment-service`: appointment creation, status workflow, doctor validation, and patient validation through service-to-service HTTP calls.
- `frontend`: a polished Nginx-served dashboard that also works as the reverse proxy / gateway.
- `postgres`: shared database container with service-owned tables.
- `prometheus` and `grafana`: monitoring and dashboard visualization.
- `terraform`: AWS VM provisioning with ports `22`, `80`, `3000`, and `9090`.
- `docs`: incident response, Terraform explanation, and deployment notes.

## Architecture

```text
Frontend + Nginx Gateway
        |
        +--> auth-service
        +--> patient-service
        +--> doctor-service
        +--> appointment-service

appointment-service --> doctor-service
appointment-service --> patient-service

All services --> PostgreSQL
Prometheus --> Service /metrics endpoints
Grafana --> Prometheus
```

## Demo access

- Admin account:
  - Email: `admin@medsync.local`
  - Password: `admin123`
- Patient account:
  - Register any new account from the frontend.

## Quick start

1. Run the full stack:

   ```bash
   docker compose up --build
   ```

2. Open:
   - Frontend: `http://localhost`
   - Grafana: `http://localhost:3000`
   - Prometheus: `http://localhost:9090`

3. Default Grafana credentials:
   - Username: `admin`
   - Password: `admin`

## Assignment deliverables in the repo

- `docker-compose.yml`
- `docker-compose.incident.yml`
- `terraform/main.tf`
- `terraform/variables.tf`
- `terraform/outputs.tf`
- `terraform/terraform.tfvars`
- `docs/deployment-guide.md`
- `docs/assignment4-incident-report.md`
- `docs/assignment5-terraform-report.md`

## Incident simulation

The incident scenario is adapted to the medical domain by treating `appointment-service` as the transactional service that is equivalent to the assignment's order service.

Run the incident version:

```bash
docker compose down
docker compose -f docker-compose.yml -f docker-compose.incident.yml up --build
```

This intentionally breaks the appointment database hostname so Prometheus, Grafana, logs, and the frontend show the outage.

## Notes for final submission

- The assignment PDF asks for PDF reports and screenshots. The repo includes ready-to-export Markdown reports and an evidence checklist in `docs/deployment-guide.md`.
- Capture screenshots after running the stack locally:
  - running containers
  - frontend before and after incident
  - Prometheus targets
  - Grafana dashboard

