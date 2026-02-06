package cmnt_test

import (
	"maps"
	"strings"
	"testing"

	"github.com/bengarrett/zipcmt/internal/cmnt"
)

func TestExportName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"empty", "", ""},
		{"name", "myfile.zip", "myfile-zipcomment.txt"},
		{"windows", "C:\\Users\\retro\\myfile.zip", "C:\\Users\\retro\\myfile-zipcomment.txt"},
		{"*nix", "/home/retro/myfile.zip", "/home/retro/myfile-zipcomment.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmnt.ExportName(tt.path); got != tt.want {
				t.Errorf("ExportName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExportFind(t *testing.T) {
	files := cmnt.Export{
		"file.txt":   true,
		"file_1.txt": true,
		"file_2.txt": true,
		"file_3.txt": true,
	}
	tests := []struct {
		name  string
		e     cmnt.Export
		fname string
		want  string
	}{
		{"none", files, "", ""},
		{"unique", files, "somefile.txt", "somefile.txt"},
		{"conflict", files, "file.txt", "file_4.txt"},
		{"conflict 3", files, "file_3.txt", "file_4.txt"},
		{"conflict 4", files, "file_4.txt", "file_4.txt"},
		{"underscores", files, "file_000_1.txt", "file_000_1.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Find(tt.fname); got != tt.want {
				t.Errorf("ExportFind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelf(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"expected", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmnt.Self()
			if (err != nil) != tt.wantErr {
				t.Errorf("Self() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name  string
		fname string
		want  bool
	}{
		{"empty", "", false},
		{"dir", "/somedir/", false},
		{"file", "/somedir/somefile.txt", false},
		{"zip", "somedir/somefile.zip", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmnt.Valid(tt.fname); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:funlen
func TestExport_Unique(t *testing.T) {
	tests := []struct {
		name         string
		existing     cmnt.Export
		zipPath      string
		dest         string
		wantContains string
		wantCount    int
	}{
		{
			name:         "first unique file",
			existing:     make(cmnt.Export),
			zipPath:      "test.zip",
			dest:         "/tmp",
			wantContains: "test-zipcomment.txt",
			wantCount:    1,
		},
		{
			name:         "duplicate filename",
			existing:     cmnt.Export{"/tmp/test-zipcomment.txt": true},
			zipPath:      "test.zip",
			dest:         "/tmp",
			wantContains: "test",
			wantCount:    1,
		},
		{
			name:         "multiple duplicates",
			existing:     cmnt.Export{"/tmp/test-zipcomment.txt": true, "/tmp/test_1-zipcomment.txt": true},
			zipPath:      "test.zip",
			dest:         "/tmp",
			wantContains: "test",
			wantCount:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of the original existing map for comparison
			originalExisting := make(cmnt.Export, len(tt.existing))
			maps.Copy(originalExisting, tt.existing)

			initialCount := len(originalExisting)
			gotPath := tt.existing.Unique(tt.zipPath, tt.dest)

			// Check that the path contains expected parts
			if !strings.Contains(gotPath, tt.wantContains) {
				t.Errorf("Unique() path %s doesn't contain %s", gotPath, tt.wantContains)
			}

			// Check that the path is in the correct directory
			if !strings.HasPrefix(gotPath, tt.dest) {
				t.Errorf("Unique() path %s doesn't start with expected dest %s", gotPath, tt.dest)
			}

			// Check that the export map was updated
			finalCount := len(tt.existing)
			if finalCount != initialCount+tt.wantCount {
				t.Errorf("Unique() export map count = %d, want %d", finalCount, initialCount+tt.wantCount)
			}

			// Check that the returned path is in the export map
			if !tt.existing[gotPath] {
				t.Errorf("Unique() path %s not found in export map", gotPath)
			}

			// For duplicate cases, check that the new path is different from originally existing ones
			if len(originalExisting) > 0 {
				for existingPath := range originalExisting {
					if gotPath == existingPath {
						t.Errorf("Unique() path %s should be different from original existing path %s", gotPath, existingPath)
					}
				}
			}
		})
	}
}
