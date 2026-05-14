# MedSync Ansible Automation

This Ansible structure provides a minimal, academic automation layer for the MedSync microservices platform.

## Important Note

This setup assumes Docker Desktop and kubectl are **already installed** on your Windows system. The playbooks provide:

- **Validation** of Docker/kubectl availability
- **Environment setup** (`.env` file creation)
- **Application deployment** using Docker Compose and kubectl
- **Monitoring deployment** resources

## Prerequisites

Ansible must be installed and available in your shell PATH before running the playbooks. On Windows, use WSL (Ubuntu or another Linux distribution).

- Linux/macOS:
  ```bash
  python3 -m pip install --user ansible-core
  export PATH="$HOME/.local/bin:$PATH"
  ```
- Windows + WSL:
  1. Open PowerShell and install WSL if not already installed:
     ```powershell
     wsl --install
     ```
  2. Restart Windows and open your Linux distribution (for example, Ubuntu).
  3. In the Linux shell, install Python and pip support:
     ```bash
     sudo apt update
     sudo apt install -y python3 python3-pip python3-venv
     ```
  4. If `python3 -m pip install --user ansible-core` fails with "externally-managed-environment", create and use a virtual environment:
     ```bash
     python3 -m venv ~/.ansible-venv
     source ~/.ansible-venv/bin/activate
     python3 -m pip install --upgrade pip
     pip install ansible-core
     ```
  5. Add the virtual environment to your shell startup:
     ```bash
     echo 'source ~/.ansible-venv/bin/activate' >> ~/.bashrc
     source ~/.bashrc
     ```
  6. Run `ansible-playbook` from the WSL shell, not from PowerShell.

## Layout

- `ansible/inventory.ini` - local inventory for localhost execution
- `ansible/ansible.cfg` - Ansible configuration for this repo
- `ansible/playbooks/` - playbooks for validation, setup, deployment, and monitoring
- `ansible/roles/` - lightweight roles for environment, app deployment, and monitoring

## Playbooks

- `install_docker.yml` - validate Docker availability
- `install_kubernetes_tools.yml` - validate kubectl availability
- `setup_environment.yml` - ensure `.env` exists and validate project workspace
- `deploy_application.yml` - deploy the application using Docker Compose
- `deploy_monitoring.yml` - deploy monitoring resources using kubectl

## Run playbooks

From the repository root in WSL:

```bash
cd /mnt/c/Users/Alikhan/OneDrive/Advanced_Programming/AP2_Assignment1_AlikhanKorazbay
./run_ansible.sh ansible/playbooks/install_docker.yml
./run_ansible.sh ansible/playbooks/install_kubernetes_tools.yml
./run_ansible.sh ansible/playbooks/setup_environment.yml
./run_ansible.sh ansible/playbooks/deploy_application.yml
./run_ansible.sh ansible/playbooks/deploy_monitoring.yml
```

Or individually:
```bash
./run_ansible.sh ansible/playbooks/setup_environment.yml
./run_ansible.sh ansible/playbooks/deploy_application.yml
./run_ansible.sh ansible/playbooks/deploy_monitoring.yml
```

## Validation commands

```bash
ansible-inventory -i ansible/inventory.ini --list
ansible-playbook --syntax-check -i ansible/inventory.ini ansible/playbooks/setup_environment.yml
ansible-playbook --syntax-check -i ansible/inventory.ini ansible/playbooks/deploy_application.yml
ansible-playbook --syntax-check -i ansible/inventory.ini ansible/playbooks/deploy_monitoring.yml
```

## Verification

After running playbooks, verify:

- Docker Compose status:
  ```bash
  docker compose ps
  ```
- Kubernetes deployment status:
  ```bash
  kubectl get pods -n medsync
  kubectl get services -n medsync
  ```
- Config validation:
  ```bash
  docker compose config
  kubectl apply --dry-run=client -f k8s/
  ```
