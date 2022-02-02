// © Ben Garrett https://github.com/bengarrett/zipcmt

package zipcmt_test

import (
	"os"
	"strings"
	"testing"

	zipcmt "github.com/bengarrett/zipcmt/lib"
	"github.com/gookit/color"
)

func TestConfig_Clean(t *testing.T) {
	type fields struct {
		SaveName  string
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
		{"missing", fields{SaveName: "/no/such/directory"}, true},
		{"file", fields{SaveName: "../test/test.txt"}, true},
		{"dir", fields{SaveName: "../test"}, false},
		{"home", fields{SaveName: "~"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &zipcmt.Config{
				SaveName:  tt.fields.SaveName,
				Export:    tt.fields.Export,
				Dupes:     tt.fields.Dupes,
				Overwrite: tt.fields.Overwrite,
				Raw:       tt.fields.Raw,
				Print:     tt.fields.Print,
				Quiet:     tt.fields.Quiet,
			}
			c.Zips = tt.fields.zips
			c.Cmmts = tt.fields.cmmts
			if err := c.Clean(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Read(t *testing.T) {
	type fields struct {
		Save      string
		Export    bool
		Dupes     bool
		Overwrite bool
		Raw       bool
		Print     bool
		Quiet     bool
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
			gotCmmt, err := zipcmt.Read(tt.fname, tt.fields.Raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotCmmt) > 0 != tt.wantCmmt {
				t.Errorf("Read() = %v, want %v", gotCmmt, tt.wantCmmt)
			}
		})
	}
}

func TestConfig_Scans(t *testing.T) {
	type fields struct {
		SaveName  string
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
		{"exportdir", fields{Dupes: true, SaveName: tmp}, "../test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &zipcmt.Config{
				SaveName:  tt.fields.SaveName,
				Export:    tt.fields.Export,
				Dupes:     tt.fields.Dupes,
				Overwrite: tt.fields.Overwrite,
				Raw:       tt.fields.Raw,
				Print:     tt.fields.Print,
				Quiet:     tt.fields.Quiet,
			}
			c.Zips = tt.fields.zips
			c.Cmmts = tt.fields.cmmts
			if err := c.WalkDir(tt.root); (err != nil) != tt.wantErr {
				t.Errorf("Config.Scans() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_separator(t *testing.T) {
	type fields struct {
		SaveName  string
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
			c := zipcmt.Config{
				SaveName:  tt.fields.SaveName,
				Export:    tt.fields.Export,
				Dupes:     tt.fields.Dupes,
				Overwrite: tt.fields.Overwrite,
				Raw:       tt.fields.Raw,
				Print:     tt.fields.Print,
				Quiet:     tt.fields.Quiet,
			}
			c.Zips = tt.fields.zips
			c.Cmmts = tt.fields.cmmts
			if got := strings.TrimSpace(c.Separator(tt.fname)); got != tt.want {
				t.Errorf("Config.Separator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Status(t *testing.T) {
	color.Enable = false
	type fields struct {
		SaveName  string
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
			c := zipcmt.Config{
				SaveName:  tt.fields.SaveName,
				Export:    tt.fields.Export,
				Dupes:     tt.fields.Dupes,
				Overwrite: tt.fields.Overwrite,
				Raw:       tt.fields.Raw,
				Print:     tt.fields.Print,
				Quiet:     tt.fields.Quiet,
			}
			c.Zips = tt.fields.zips
			c.Cmmts = tt.fields.cmmts
			c.SetTest()
			if got := strings.TrimSpace(c.Status()); got != tt.want {
				t.Errorf("Config.Status() = \ngot:  %v,\nwant: %v", got, tt.want)
			}
		})
	}
}
