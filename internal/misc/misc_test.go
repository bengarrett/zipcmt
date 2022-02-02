package misc_test

import (
	"testing"

	"github.com/bengarrett/zipcmt/internal/misc"
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
			if got := misc.ExportName(tt.path); got != tt.want {
				t.Errorf("ExportName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExportFind(t *testing.T) {
	files := misc.Export{
		"file.txt":   true,
		"file_1.txt": true,
		"file_2.txt": true,
		"file_3.txt": true,
	}
	tests := []struct {
		name  string
		e     misc.Export
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
			_, err := misc.Self()
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
			if got := misc.Valid(tt.fname); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
