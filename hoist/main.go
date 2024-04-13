package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Function to calculate the SHA256 hash of a file
func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func isSymlink(filePath string) (bool, error) {
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return false, fmt.Errorf("failed stat: %w", err)
	}
	return fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}

// Function to recursively scan a directory and identify duplicate files
func scanDirectory(rootDir string) (map[string][]string, error) {
	fileHashes := make(map[string][]string)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip all links
		if isLink, _ := isSymlink(path); !isLink {
			// skip directories
			if !info.IsDir() {
				hash, err := calculateHash(path)
				if err != nil {
					return err
				}
				fileHashes[hash] = append(fileHashes[hash], path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return fileHashes, nil
}

// Function to hoist duplicate files and create links to hoisted files
func hoistFiles(fileHashes map[string][]string, rootDir string) error {
	// create a directory in the root directory to store the hoisted files
	hoistDirName := filepath.Join(rootDir, "hoisted-resources")
	for fileHash, paths := range fileHashes {
		if len(paths) > 1 {
			for _, originalPath := range paths[1:] {
				hoistHashedPath := filepath.Join(hoistDirName, fileHash)
				fmt.Println("Moving:", originalPath, hoistHashedPath)
				// create the target directory for the hoisted file
				if err := os.MkdirAll(filepath.Dir(hoistHashedPath), 0755); err != nil {
					return err
				}
				// move the file to the hoisted location
				if err := os.Rename(originalPath, hoistHashedPath); err != nil {
					return err
				}
				// get the relative path from the original file location, to the hoistPath
				relHoistedPath, err := filepath.Rel(filepath.Dir(originalPath), hoistHashedPath)
				if err != nil {
					return err
				}
				if err := os.Symlink(relHoistedPath, originalPath); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func printHelp() {
	fmt.Println("Hoist")
	fmt.Println("\tA tool to hoist duplicate files and create links to hoisted files.")
	fmt.Println("\t- Works relative to the workind directory, moving dupes to `hoisted-resources` directory and replaceing the originals with links.")
	fmt.Println("Version:", Version)
	fmt.Println("Usage: hoist <rootDir>")
}

func main() {
	if len(os.Args) != 2 {
		printHelp()
		os.Exit(1)
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		printHelp()
		os.Exit(0)
	}
	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Println("Version:", Version)
		os.Exit(0)
	}

	rootDir := os.Args[1]

	fileHashes, err := scanDirectory(rootDir)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if err := hoistFiles(fileHashes, rootDir); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Duplicate files hoisted and links created successfully.")
	fmt.Println("Version:", Version)
}
