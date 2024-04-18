package main

import "testing"

func TestCalculateHash(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"./testing_resources/file_structure/fileTOP_01.md", "bba5b7c2eb66764791189426881953ada819b09981608f8cfeb9b0e42b74f265"},
		{"./testing_resources/file_structure/fileTOP_02.md", "872a41e0e2b42a0fb58b837fec4445d016e53339787ac25a0f6dff5fc20a72cd"},
		{"./testing_resources/file_structure/fileTOP_03.md", "a5ea2c285d223acc6ddafb2d7da889ce2167c140fc77d941d679bea67d4c9d38"},
		{"./testing_resources/file_structure/fileTOP_04.md", "e326860075ec7222aa54df239c523ffad493447de4aeb2988fe0ac5f39a7c09a"},
		{"./testing_resources/file_structure/fileTOP_05.md", "e326860075ec7222aa54df239c523ffad493447de4aeb2988fe0ac5f39a7c09a"},
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
		{"testdata1", true},
		{"testdata2", true},
	}
	for _, test := range tests {
		if got, err := isSymlink(test.input); got != test.want || err != nil {
			t.Errorf("isSymlink(%q) = %v; want %v", test.input, got, test.want)
		}
	}
}
