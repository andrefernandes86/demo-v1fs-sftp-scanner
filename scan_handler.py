import os
import shutil
import logging
from pathlib import Path
from tmv1fs import FileScanner

# Configure logging
logging.basicConfig(level=logging.INFO)

# Environment variables
API_KEY = os.getenv("TM_API_KEY")
UPLOADS_DIR = "/var/sftp/uploads"
DEFAULT_DIR = "/nfs/share/default"
MALICIOUS_DIR = "/nfs/share/malicious"

def scan_file(scanner, file_path):
    """Scan a file and return whether it is malicious."""
    logging.info(f"Scanning file: {file_path}")
    try:
        result = scanner.scan_file(file_path)
        if result.get("malicious", False):
            logging.info(f"File {file_path} is malicious.")
            return True
        logging.info(f"File {file_path} is clean.")
        return False
    except Exception as e:
        logging.error(f"Error scanning file {file_path}: {e}")
        return False

def move_file(src, dest_dir):
    """Move a file to the destination directory."""
    dest_path = Path(dest_dir) / Path(src).name
    shutil.move(src, dest_path)
    logging.info(f"Moved {src} to {dest_path}")

def main():
    if not API_KEY:
        logging.error("API_KEY is not set. Exiting.")
        return

    # Initialize the Trend Micro scanner
    scanner = FileScanner(api_key=API_KEY)

    # Ensure the directories exist
    os.makedirs(DEFAULT_DIR, exist_ok=True)
    os.makedirs(MALICIOUS_DIR, exist_ok=True)

    # Process files in the uploads directory
    for file_path in Path(UPLOADS_DIR).glob("*"):
        if file_path.is_file():
            is_malicious = scan_file(scanner, str(file_path))
            if is_malicious:
                move_file(str(file_path), MALICIOUS_DIR)
            else:
                move_file(str(file_path), DEFAULT_DIR)

if __name__ == "__main__":
    main()
