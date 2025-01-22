# demo-v1fs-sftp-scanner

docker run -p 2222:2222-e TM_API_KEY="your_actual_api_key" sftp-scanner


docker build -t sftp-scanner .



docker run -d \
    -e TM_API_KEY="your_actual_api_key" \
    -v /path/to/nfs/default:/nfs/share/default \
    -v /path/to/nfs/malicious:/nfs/share/malicious \
    -p 22:22 \
    --name sftp-scanner \
    sftp-scanner



    docker run -e TM_API_KEY="your_api_key" -e SERVICE_PORT=2222 -p 2222:2222 sftp-scanner
