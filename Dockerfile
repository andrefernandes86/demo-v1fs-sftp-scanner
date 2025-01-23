FROM python:3.11-alpine

# Install necessary tools
RUN apk add --no-cache bash openssh nfs-utils git

# Set up directories for the application
WORKDIR /app

# Clone the Trend Micro Python SDK repository
RUN git clone https://github.com/trendmicro/tm-v1-fs-python-sdk.git /app/tm-v1-fs-python-sdk

# Install Python dependencies
RUN pip install requirements.txt

# Copy the custom handler script
COPY ./scan_handler.py /app

# Install and configure OpenSSH for SFTP
RUN mkdir -p /var/sftp/uploads && \
    mkdir -p /nfs/share/default && \
    mkdir -p /nfs/share/malicious && \
    addgroup -g 1001 sftpgroup && \
    adduser -D -h /var/sftp/uploads -G sftpgroup -s /sbin/nologin sftpuser && \
    echo "sftpuser:password" | chpasswd && \
    chown sftpuser:sftpgroup /var/sftp/uploads && \
    chmod 700 /var/sftp/uploads

# Copy the SSHD configuration file
COPY sshd_config /etc/ssh/sshd_config

# Copy the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set environment variables for API key and port
ENV TM_API_KEY "" \
    SERVICE_PORT 22

# Expose the configurable SFTP port
EXPOSE ${SERVICE_PORT}

# Run the entrypoint script
ENTRYPOINT ["/entrypoint.sh"]
