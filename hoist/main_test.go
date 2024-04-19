package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestCalculateHash(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"./testing_resources/file_structure/fileTOP_01.md", "bba5b7c2eb66764791189426881953ada819b09981608f8cfeb9b0e42b74f265"},
		{"./testing_resources/file_structure/fileTOP_02.md", "872a41e0e2b42a0fb58b837fec4445d016e53339787ac25a0f6dff5fc20a72cd"},
		{"./testing_resources/file_structure/fileTOP_03.md", "a5ea2c285d223acc6ddafb2d7da889ce2167c140fc77d941d679bea67d4c9d38"},
		{"./testing_resources/file_structure/fileTOP_04.md", "e326860075ec7222aa54df239c523ffad493447de4aeb2988fe0ac5f39a7c09a"},
		{"./testing_resources/file_structure/fileTOP_05.COPYOF04.md", "e326860075ec7222aa54df239c523ffad493447de4aeb2988fe0ac5f39a7c09a"},
	}
	for _, test := range tests {
		if got, err := calculateHash(test.input); got != test.want || err != nil {
			t.Errorf("calculateHash(%q) = %v; want %v, err %v", test.input, got, test.want, err)
		}
	}
}

func TestIsSymlink(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"./testing_resources/file_structure/fileTOP_01.md", false},
		{"./testing_resources/file_structure/fileTOP_01.LINKTO01.md", true},
	}
	for _, test := range tests {
		if got, err := isSymlink(test.input); got != test.want || err != nil {
			t.Errorf("isSymlink(%q) = %v; want %v, err %v", test.input, got, test.want, err)
		}
	}
}

func TestScanDirectoryForDupes(t *testing.T) {
	tests := []struct {
		input string
		want  map[string][]string
	}{
		{"testing_resources/static_file_structure", map[string][]string{
			"bba5b7c2eb66764791189426881953ada819b09981608f8cfeb9b0e42b74f265": {
				"testing_resources/static_file_structure/fileTOP_01.md",
			},
			"872a41e0e2b42a0fb58b837fec4445d016e53339787ac25a0f6dff5fc20a72cd": {
				"testing_resources/static_file_structure/fileTOP_02.md",
			},
			"a5ea2c285d223acc6ddafb2d7da889ce2167c140fc77d941d679bea67d4c9d38": {
				"testing_resources/static_file_structure/fileTOP_03.md",
			},
			"e326860075ec7222aa54df239c523ffad493447de4aeb2988fe0ac5f39a7c09a": {
				"testing_resources/static_file_structure/fileTOP_04.md",
				"testing_resources/static_file_structure/fileTOP_05.COPYOF04.md",
			},
		}},
	}
	for _, test := range tests {
		if got, err := scanDirectoryForDupes(test.input); !reflect.DeepEqual(got, test.want) || err != nil {
			t.Errorf("scanDirectoryForDupes(%q) = %v; want %v, err %v", test.input, got, test.want, err)
		}
	}
}

func TestHoistFiles(t *testing.T) {
	existingDir := "./testing_resources/file_structure"
	tmpDir := "./testing_resources/tmp_copy_file_structure"
	// remove the existing tmp directory
	if err := os.RemoveAll(tmpDir); err != nil {
		fmt.Printf("Error removing old tmp dir [%v] for testing", tmpDir)
		os.Exit(1)
	}
	// make an exact recursive copy, including sym links, of the testing_resources/static_file_structure named testing_resources/tmp_file_structure
	cmd := exec.Command("cp", "-r", existingDir, tmpDir)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error copying directory: %v\n", err)
		os.Exit(1)
	}
	// run the tests, hoist the files

	rootDir := tmpDir
	hoistList, err := scanDirectoryForDupes(tmpDir)
	if err != nil {
		fmt.Printf("Error scanning directory for dupes in hoist test: %v\n", err)
		os.Exit(1)
	}
	if err := hoistFiles(hoistList, rootDir); err != nil {
		fmt.Printf("Error hoisting files in hoist test: %v\n", err)
	}
	// now ensure the hoist went as expected
	fmt.Printf("[STUB] Checking hoist results\n")
}
