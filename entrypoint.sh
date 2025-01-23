#!/bin/bash

# Generate SSH host keys if they do not exist
if [ ! -f /etc/ssh/ssh_host_rsa_key ]; then
    ssh-keygen -A
fi

# Start the SSHD server for SFTP
/usr/sbin/sshd -D &

# Run the Python file handler
python3 /app/scan_handler.py
