package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddPackage(t *testing.T) {
	tests := []struct {
		name        string
		pkgName     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid package - validator",
			pkgName: "validator",
			wantErr: false,
		},
		{
			name:    "valid package - database",
			pkgName: "database",
			wantErr: false,
		},
		{
			name:    "valid package - jwt",
			pkgName: "jwt",
			wantErr: false,
		},
		{
			name:    "valid package - utils",
			pkgName: "utils",
			wantErr: false,
		},
		{
			name:    "valid package - cache",
			pkgName: "cache",
			wantErr: false,
		},
		{
			name:    "valid package - email",
			pkgName: "email",
			wantErr: false,
		},
		{
			name:    "valid package - upload",
			pkgName: "upload",
			wantErr: false,
		},
		{
			name:    "valid package - pagination",
			pkgName: "pagination",
			wantErr: false,
		},
		{
			name:    "valid package - testutil",
			pkgName: "testutil",
			wantErr: false,
		},
		{
			name:        "invalid package",
			pkgName:     "nonexistent",
			wantErr:     true,
			errContains: "未知的功能包",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(tmpDir)

			// Create a minimal go.mod file
			goModContent := "module testproject\n\ngo 1.24\n"
			if err := os.WriteFile("go.mod", []byte(goModContent), 0o644); err != nil {
				t.Fatalf("Failed to create go.mod: %v", err)
			}

			err := AddPackage(tt.pkgName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddPackage() expected error but got none")
				} else if tt.errContains != "" && !containsSubstring(err.Error(), tt.errContains) {
					t.Errorf("AddPackage() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("AddPackage() unexpected error = %v", err)
				return
			}

			// Verify package file was created
			pkgPath := filepath.Join("pkg", tt.pkgName, tt.pkgName+".go")
			if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
				t.Errorf("Package file %s was not created", pkgPath)
			}
		})
	}
}

func TestAddPackageWithoutGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	err := AddPackage("validator")
	if err == nil {
		t.Errorf("AddPackage() should fail without go.mod")
	}
}

func TestAddPackageIdempotency(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	goModContent := "module testproject\n\ngo 1.24\n"
	if err := os.WriteFile("go.mod", []byte(goModContent), 0o644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// First call
	err := AddPackage("validator")
	if err != nil {
		t.Fatalf("First AddPackage() call failed: %v", err)
	}

	// Get file info
	filePath := filepath.Join("pkg", "validator", "validator.go")
	info1, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	// Second call (should be idempotent)
	err = AddPackage("validator")
	if err != nil {
		t.Fatalf("Second AddPackage() call failed: %v", err)
	}

	// Verify file wasn't modified
	info2, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file after second call: %v", err)
	}

	if info1.ModTime() != info2.ModTime() {
		t.Errorf("File was modified on second call, idempotency violated")
	}
}

func TestListAvailablePackages(t *testing.T) {
	packages := ListAvailablePackages()

	expectedPackages := []string{
		"validator",
		"database",
		"jwt",
		"utils",
		"cache",
		"email",
		"upload",
		"pagination",
		"testutil",
	}

	if len(packages) != len(expectedPackages) {
		t.Errorf("ListAvailablePackages() returned %d packages, want %d", len(packages), len(expectedPackages))
	}

	packageMap := make(map[string]bool)
	for _, pkg := range packages {
		packageMap[pkg] = true
	}

	for _, expected := range expectedPackages {
		if !packageMap[expected] {
			t.Errorf("ListAvailablePackages() missing package: %s", expected)
		}
	}
}
