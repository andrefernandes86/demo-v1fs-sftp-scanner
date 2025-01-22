package main

import (
	"log"
	"os"
	"path/filepath"

	tmsdk "github.com/trendmicro/tm-v1-fs-golang-sdk"
)

const (
	defaultNFSPath   = "/nfs/share/default"
	maliciousNFSPath = "/nfs/share/malicious"
	uploadsPath      = "/var/sftp/uploads"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("TM_API_KEY")
	if apiKey == "" {
		log.Fatalf("API key not found. Set TM_API_KEY environment variable.")
	}

	// Initialize the Trend Micro SDK with the API key
	scanner := tmsdk.NewScannerWithKey(apiKey)
	defer scanner.Close()

	// Watch the upload directory for new files
	err := filepath.Walk(uploadsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		log.Printf("Scanning file: %s\n", path)

		// Scan the file
		result, err := scanner.ScanFile(path)
		if err != nil {
			log.Printf("Error scanning file %s: %v\n", path, err)
			return err
		}

		// Decide where to move the file
		if result.IsMalicious() {
			log.Printf("File %s is malicious. Moving to %s\n", path, maliciousNFSPath)
			moveFile(path, maliciousNFSPath)
		} else {
			log.Printf("File %s is clean. Moving to %s\n", path, defaultNFSPath)
			moveFile(path, defaultNFSPath)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error processing files: %v\n", err)
	}
}

// moveFile moves a file to the target directory
func moveFile(srcPath, targetDir string) {
	filename := filepath.Base(srcPath)
	targetPath := filepath.Join(targetDir, filename)

	err := os.Rename(srcPath, targetPath)
	if err != nil {
		log.Printf("Error moving file %s to %s: %v\n", srcPath, targetPath, err)
		return
	}

	log.Printf("File moved: %s -> %s\n", srcPath, targetPath)
}
