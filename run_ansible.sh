#!/usr/bin/env bash

# Ensure Ansible reads the repository config even in world-writable mounted directories.
export ANSIBLE_CONFIG="$(pwd)/ansible.cfg"
ansible-playbook -i ansible/inventory.ini "$@"
