#!/bin/bash

# Start the SSHD server for SFTP
/usr/sbin/sshd -D &

# Run the Python file handler
python3 /app/scan_handler.py
