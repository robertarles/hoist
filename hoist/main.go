package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Settings struct {
	hoistDirname string
}

func NewSettings() *Settings {
	return &Settings{
		hoistDirname: "hoisted-resources",
	}
}

// Function to calculate the SHA256 hash of a file
func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "error hashing", err
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
func scanDirectoryForDupes(rootDir string) (map[string][]string, error) {
	fileHashes := make(map[string][]string)
	// get the settings
	settings := NewSettings()

	// Walk the directory and calculate the hash of each file
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		// make the Walk method skip the hoist directory
		if err != nil {
			return err
		}
		// skip the hoist directory `settings.hoistDirname`
		if filepath.Base(path) == settings.hoistDirname {
			return filepath.SkipDir
		}
		// skip all links
		if isLink, _ := isSymlink(path); !isLink {
			// skip directories
			if !info.IsDir() {
				hash, err := calculateHash(path)
				if err != nil {
					return err
				}
				// append the this file to the list of files for this hash
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
	// keep track of the number of times each file is hoisted, then calculate the total size saved
	hoistedFileCounts := make(map[string]int)
	// create a directory in the root directory to store the hoisted files
	hoistDirName := filepath.Join(rootDir, "hoisted-resources")
	// create the target directory for the hoisted file
	if err := os.MkdirAll(hoistDirName, 0755); err != nil {
		return fmt.Errorf("failed to create hoist directory: %w", err)
	}
	// for each hash made, create the links for the files of that hash
	for fileHash, paths := range fileHashes {
		if len(paths) > 1 {
			// DEBUG fmt.Printf("[DEBUG] filehash:%v has %v files\n", fileHash, len(paths))
			// for each file with the same hash
			for _, originalPath := range paths {
				// create the hoisted full path, in the format hoistDirname/<filename>_<hash>.<ext>
				hoistPath := filepath.Join(hoistDirName, fileHash+filepath.Ext(originalPath))
				hoistedFileCounts[hoistPath]++
				// create a tmp filename in this scope to enable deferred removal
				tmpFilename := filepath.Join(filepath.Dir(originalPath), filepath.Base(originalPath)+"_"+fileHash+".tmp")
				fmt.Println("- ", originalPath, "\n\t->:", hoistPath)
				// if the file does not already exists in the hoisted directory, hoist this one up
				if _, err := os.Stat(hoistPath); os.IsNotExist(err) {
					// fmt.Println("[DEBUG] Hoisted file not found, hoisting:", originalPath, "to:", hoistPath)
					// move the file to the hoisted location
					if err := os.Rename(originalPath, hoistPath); err != nil {
						fmt.Printf("Error moving file: %v\n", err)
						return fmt.Errorf("failed to move original file: %w", err)
					}
				} else {
					// fmt.Printf("[DEBUG] Hoisted file found, renaming original %v to tmp file: %v", originalPath, hoistPath)
					// rename the file, in place, and defer delete for after the link is created
					if err := os.Rename(originalPath, tmpFilename); err != nil {
						fmt.Printf("Error renaming original file: %v\n", err)
						return fmt.Errorf("failed to make backup tmp of original file: %w", err)
					}
				}
				// get the relative path FROM the original file location TO the hoistPath
				relHoistedPath, err := filepath.Rel(filepath.Dir(originalPath), hoistPath)
				if err != nil {
					fmt.Printf("Error calculating relative path: %v\n", err)
					return fmt.Errorf("failed to calulate relative path: %w", err)
				}
				// finally, create a symlink to replace the hoisted file
				if err := os.Symlink(relHoistedPath, originalPath); err != nil {
					// fmt.Printf("[DEBUG] Error creating symlink: %v for file: %v\n", err, originalPath)
					// and recover if symlink fails and the tmpfile exist
					if errRenaming := os.Rename(tmpFilename, originalPath); errRenaming != nil {
						fmt.Printf("Error renaming temp file back to original: %v\n", errRenaming)
						return fmt.Errorf("failed to create symlink, and failed to restore original file while recovering from hoisting error: %w", errRenaming)
					}
					return fmt.Errorf("failed to create symlink: %w", err)
				}

				// check if tmpfilename exists
				if _, err := os.Stat(tmpFilename); !os.IsNotExist(err) {
					// fmt.Printf("[DEBUG] Removing tmp file: %v\n", tmpFilename)
					if err := os.Remove(tmpFilename); err != nil {
						// debug output
						fmt.Printf("Error removing tmp file: %v\n", err)
						return fmt.Errorf("failed to remove tmp file: %w", err)
					}
				}
			}
		}
	}
	// print the hoisted file counts
	for hoistedPath, count := range hoistedFileCounts {
		fmt.Println("Hoisted:", hoistedPath, "Count:\t", count)
	}
	// print the total size saved
	var totalSizeSaved int64
	for hoistedPath, count := range hoistedFileCounts {
		fileInfo, err := os.Stat(hoistedPath)
		if err != nil {
			return err
		}
		totalSizeSaved += fileInfo.Size() * int64(count)
	}
	fmt.Println("Total size saved:", totalSizeSaved, "bytes")
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

	fileHashes, err := scanDirectoryForDupes(rootDir)
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
