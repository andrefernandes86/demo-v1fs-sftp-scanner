#!/bin/bash

# Start the SSHD server for SFTP
/usr/sbin/sshd -D &

# Start the file scanning handler
/app/scan-handler
