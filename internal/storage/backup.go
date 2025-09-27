// internal/storage/backup.go
package storage

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BackupManager handles database backup and restore operations
type BackupManager struct {
	dataDir string
}

// BackupInfo contains information about a backup
type BackupInfo struct {
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
	TableCount  int       `json:"table_count"`
	BackupSize  int64     `json:"backup_size"`
	Description string    `json:"description"`
}

// NewBackupManager creates a new backup manager
func NewBackupManager(dataDir string) *BackupManager {
	return &BackupManager{
		dataDir: dataDir,
	}
}

// CreateBackup creates a backup of the database
func (bm *BackupManager) CreateBackup(backupPath string, description string) error {
	// Create backup directory if it doesn't exist
	backupDir := filepath.Dir(backupPath)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup file
	backupFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer backupFile.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(backupFile)
	defer gzipWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Get all .harudb files
	entries, err := os.ReadDir(bm.dataDir)
	if err != nil {
		return fmt.Errorf("failed to read data directory: %w", err)
	}

	tableCount := 0
	totalSize := int64(0)

	// Add all .harudb files to backup
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".harudb") {
			continue
		}

		filePath := filepath.Join(bm.dataDir, entry.Name())
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		// Read file content
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		// Create tar header
		header := &tar.Header{
			Name:    entry.Name(),
			Size:    fileInfo.Size(),
			Mode:    int64(fileInfo.Mode()),
			ModTime: fileInfo.ModTime(),
		}

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header: %w", err)
		}

		// Write file content
		if _, err := tarWriter.Write(fileContent); err != nil {
			return fmt.Errorf("failed to write file content: %w", err)
		}

		tableCount++
		totalSize += fileInfo.Size()
	}

	// Create backup info
	backupInfo := BackupInfo{
		Timestamp:   time.Now(),
		Version:     "v0.0.5",
		TableCount:  tableCount,
		BackupSize:  totalSize,
		Description: description,
	}

	// Serialize backup info
	infoData, err := json.MarshalIndent(backupInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup info: %w", err)
	}

	// Add backup info to tar
	infoHeader := &tar.Header{
		Name:    "backup_info.json",
		Size:    int64(len(infoData)),
		Mode:    0644,
		ModTime: time.Now(),
	}

	if err := tarWriter.WriteHeader(infoHeader); err != nil {
		return fmt.Errorf("failed to write backup info header: %w", err)
	}

	if _, err := tarWriter.Write(infoData); err != nil {
		return fmt.Errorf("failed to write backup info: %w", err)
	}

	return nil
}

// RestoreBackup restores a database from a backup
func (bm *BackupManager) RestoreBackup(backupPath string) error {
	// Open backup file
	backupFile, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer backupFile.Close()

	// Create gzip reader
	gzipReader, err := gzip.NewReader(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzipReader)

	// Clear existing data directory (except WAL and users)
	entries, err := os.ReadDir(bm.dataDir)
	if err != nil {
		return fmt.Errorf("failed to read data directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".harudb") {
			continue
		}

		filePath := filepath.Join(bm.dataDir, entry.Name())
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to remove existing file %s: %w", entry.Name(), err)
		}
	}

	// Extract files from backup
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Skip backup info file
		if header.Name == "backup_info.json" {
			continue
		}

		// Only restore .harudb files
		if !strings.HasSuffix(header.Name, ".harudb") {
			continue
		}

		// Create file
		filePath := filepath.Join(bm.dataDir, header.Name)
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", header.Name, err)
		}

		// Copy file content
		if _, err := io.Copy(file, tarReader); err != nil {
			file.Close()
			return fmt.Errorf("failed to copy file content for %s: %w", header.Name, err)
		}

		file.Close()
	}

	return nil
}

// GetBackupInfo returns information about a backup file
func (bm *BackupManager) GetBackupInfo(backupPath string) (*BackupInfo, error) {
	// Open backup file
	backupFile, err := os.Open(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open backup file: %w", err)
	}
	defer backupFile.Close()

	// Create gzip reader
	gzipReader, err := gzip.NewReader(backupFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzipReader)

	// Find backup info file
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar header: %w", err)
		}

		if header.Name == "backup_info.json" {
			// Read backup info
			infoData := make([]byte, header.Size)
			if _, err := io.ReadFull(tarReader, infoData); err != nil {
				return nil, fmt.Errorf("failed to read backup info: %w", err)
			}

			var backupInfo BackupInfo
			if err := json.Unmarshal(infoData, &backupInfo); err != nil {
				return nil, fmt.Errorf("failed to unmarshal backup info: %w", err)
			}

			return &backupInfo, nil
		}
	}

	return nil, fmt.Errorf("backup info not found in backup file")
}

// ListBackups lists all backup files in a directory
func (bm *BackupManager) ListBackups(backupDir string) ([]string, error) {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".backup") {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}
