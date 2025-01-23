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
	// Get the API key from the environment variable
	apiKey := os.Getenv("TM_API_KEY")
	if apiKey == "" {
		log.Fatalf("API key not found. Set TM_API_KEY environment variable.")
	}

	// Initialize the Trend Micro File Scanner
	scanner, err := tmsdk.NewFileScanner(apiKey) // Adjust function name based on the SDK
	if err != nil {
		log.Fatalf("Failed to initialize file scanner: %v", err)
	}
	defer scanner.Close()

	// Scan files in the uploads directory
	err = filepath.Walk(uploadsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		log.Printf("Scanning file: %s", path)

		// Scan the file
		result, err := scanner.ScanFile(path) // Adjust based on the SDK
		if err != nil {
			log.Printf("Error scanning file %s: %v", path, err)
			return err
		}

		// Check if the file is malicious
		if result.IsMalicious() { // Adjust based on the SDK
			log.Printf("File %s is malicious. Moving to %s", path, maliciousNFSPath)
			moveFile(path, maliciousNFSPath)
		} else {
			log.Printf("File %s is clean. Moving to %s", path, defaultNFSPath)
			moveFile(path, defaultNFSPath)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error processing files: %v", err)
	}
}

// moveFile moves a file to the target directory
func moveFile(srcPath, targetDir string) {
	filename := filepath.Base(srcPath)
	targetPath := filepath.Join(targetDir, filename)

	err := os.Rename(srcPath, targetPath)
	if err != nil {
		log.Printf("Error moving file %s to %s: %v", srcPath, targetPath, err)
		return
	}

	log.Printf("File moved: %s -> %s", srcPath, targetPath)
}
