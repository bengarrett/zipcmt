// © Ben Garrett https://github.com/bengarrett/zipcmt

package zipcmt

import (
	"os"
	"strings"
	"testing"

	"github.com/gookit/color"
)

func TestConfig_Clean(t *testing.T) {
	type fields struct {
		Save      string
		Export    bool
		Dupes     bool
		Overwrite bool
		Raw       bool
		Print     bool
		Quiet     bool
		zips      int
		cmmts     int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty", fields{}, false},
		{"missing", fields{Save: "/no/such/directory"}, true},
		{"file", fields{Save: "../test/test.txt"}, true},
		{"dir", fields{Save: "../test"}, false},
		{"home", fields{Save: "~"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Save:      tt.fields.Save,
				Export:    tt.fields.Export,
				Dupes:     tt.fields.Dupes,
				Overwrite: tt.fields.Overwrite,
				Raw:       tt.fields.Raw,
				Print:     tt.fields.Print,
				Quiet:     tt.fields.Quiet,
				zips:      tt.fields.zips,
				cmmts:     tt.fields.cmmts,
			}
			if err := c.Clean(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_Read(t *testing.T) {
	type fields struct {
		Save      string
		Export    bool
		Dupes     bool
		Overwrite bool
		Raw       bool
		Print     bool
		Quiet     bool
		zips      int
		cmmts     int
	}
	tests := []struct {
		name     string
		fields   fields
		fname    string
		wantCmmt bool
		wantErr  bool
	}{
		{"empty", fields{}, "", false, false},
		{"bad file", fields{}, "../missing/no_such_files.zip", false, false},
		{"no comment file", fields{}, "../test/test-no-comment.zip", false, false},
		{"file with comment", fields{}, "../test/test-with-comment.zip", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Save:      tt.fields.Save,
				Export:    tt.fields.Export,
				Dupes:     tt.fields.Dupes,
				Overwrite: tt.fields.Overwrite,
				Raw:       tt.fields.Raw,
				Print:     tt.fields.Print,
				Quiet:     tt.fields.Quiet,
				zips:      tt.fields.zips,
				cmmts:     tt.fields.cmmts,
			}
			gotCmmt, err := c.Read(tt.fname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotCmmt) > 0 != tt.wantCmmt {
				t.Errorf("Config.Read() = %v, want %v", gotCmmt, tt.wantCmmt)
			}
		})
	}
}

func TestConfig_Scans(t *testing.T) {
	type fields struct {
		Save      string
		Export    bool
		Dupes     bool
		Overwrite bool
		Raw       bool
		Print     bool
		Quiet     bool
		zips      int
		cmmts     int
	}
	tmp, err := os.MkdirTemp(os.TempDir(), "zipcmtscanstest")
	if err != nil {
		t.Errorf("Cannot create temp directory: %v", err)
	}
	defer os.RemoveAll(tmp)
	tests := []struct {
		name    string
		fields  fields
		root    string
		wantErr bool
	}{
		{"no root", fields{}, "", true},
		{"bad root", fields{}, "../test/missing", true},
		{"test dir", fields{Dupes: true}, "../test", false},
		{"exportdir", fields{Dupes: true, Save: tmp}, "../test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Save:      tt.fields.Save,
				Export:    tt.fields.Export,
				Dupes:     tt.fields.Dupes,
				Overwrite: tt.fields.Overwrite,
				Raw:       tt.fields.Raw,
				Print:     tt.fields.Print,
				Quiet:     tt.fields.Quiet,
				zips:      tt.fields.zips,
				cmmts:     tt.fields.cmmts,
			}
			if err := c.Walk(tt.root); (err != nil) != tt.wantErr {
				t.Errorf("Config.Scans() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_separator(t *testing.T) {
	type fields struct {
		Save      string
		Export    bool
		Dupes     bool
		Overwrite bool
		Raw       bool
		Print     bool
		Quiet     bool
		zips      int
		cmmts     int
	}
	tests := []struct {
		name   string
		fields fields
		fname  string
		want   string
	}{
		{"empty", fields{}, "", ""},
		{"print", fields{Print: true}, "somefile.zip", "── somefile.zip ─────────────────────────┐"},
		{"quiet", fields{Print: true, Quiet: true}, "somefile.zip", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Save:      tt.fields.Save,
				Export:    tt.fields.Export,
				Dupes:     tt.fields.Dupes,
				Overwrite: tt.fields.Overwrite,
				Raw:       tt.fields.Raw,
				Print:     tt.fields.Print,
				Quiet:     tt.fields.Quiet,
				zips:      tt.fields.zips,
				cmmts:     tt.fields.cmmts,
			}
			if got := strings.TrimSpace(c.separator(tt.fname)); got != tt.want {
				t.Errorf("Config.separator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Status(t *testing.T) {
	color.Enable = false
	type fields struct {
		Save      string
		Export    bool
		Dupes     bool
		Overwrite bool
		Raw       bool
		Print     bool
		Quiet     bool
		zips      int
		cmmts     int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"none", fields{}, "Scanned 0 zip archives and found 0 unique comments"},
		{"one", fields{zips: 1, cmmts: 1}, "Scanned 1 zip archive and found 1 unique comment"},
		{"multi", fields{zips: 5, cmmts: 2}, "Scanned 5 zip archives and found 2 unique comments"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Save:      tt.fields.Save,
				Export:    tt.fields.Export,
				Dupes:     tt.fields.Dupes,
				Overwrite: tt.fields.Overwrite,
				Raw:       tt.fields.Raw,
				Print:     tt.fields.Print,
				Quiet:     tt.fields.Quiet,
				test:      true,
				zips:      tt.fields.zips,
				cmmts:     tt.fields.cmmts,
			}
			if got := strings.TrimSpace(c.Status()); got != tt.want {
				t.Errorf("Config.Status() = \ngot:  %v,\nwant: %v", got, tt.want)
			}
		})
	}
}

func Test_valid(t *testing.T) {
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
			if got := valid(tt.fname); got != tt.want {
				t.Errorf("valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_export_find(t *testing.T) {
	files := export{
		"file.txt":   true,
		"file_1.txt": true,
		"file_2.txt": true,
		"file_3.txt": true,
	}
	tests := []struct {
		name  string
		e     export
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
			if got := tt.e.find(tt.fname); got != tt.want {
				t.Errorf("export.unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_exportName(t *testing.T) {
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
			if got := exportName(tt.path); got != tt.want {
				t.Errorf("exportName() = %v, want %v", got, tt.want)
			}
		})
	}
}
