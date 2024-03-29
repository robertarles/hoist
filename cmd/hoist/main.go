package main

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

var version string = "v1.0.0-beta"

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

// Function to recursively scan a directory and identify duplicate files
func scanDirectory(rootDir string) (map[string][]string, error) {
    fileHashes := make(map[string][]string)

    err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            hash, err := calculateHash(path)
            if err != nil {
                return err
            }
            fileHashes[hash] = append(fileHashes[hash], path)
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    return fileHashes, nil
}

// Function to hoist duplicate files and create links to hoisted files
func hoistDuplicates(fileHashes map[string][]string) error {
    for _, paths := range fileHashes {
        if len(paths) > 1 {
            hoistedFilePath := paths[0] + ".hoisted"
            if err := os.Rename(paths[0], hoistedFilePath); err != nil {
                return err
            }
            for _, linkPath := range paths[1:] {
                if err := os.Link(hoistedFilePath, linkPath); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ./program <root_directory>")
        os.Exit(1)
    }

    rootDir := os.Args[1]

    fileHashes, err := scanDirectory(rootDir)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    if err := hoistDuplicates(fileHashes); err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    fmt.Println("Duplicate files hoisted and links created successfully.")
}

