// Package zipcmt is a viewer and an extractor of zip archive comments
package zipcmt

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func ExampleRead() {
	c := Config{
		Raw: false,
	}
	s, err := c.Read("../test/test-with-comment.zip")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(s)
	// Output:
	//This is an example test comment for zipcmmt.
	//
}

func ExampleScan() {
	c := Config{
		Print: true,
		Quiet: true,
	}
	if err := c.Scan("../test"); err != nil {
		log.Println(err)
	}
	// Output:
	//This is an example test comment for zipcmmt.[0m
	//
}

func TestConfig_Clean(t *testing.T) {
	type fields struct {
		ExportDir  string
		ExportFile bool
		NoDupes    bool
		Overwrite  bool
		Raw        bool
		Print      bool
		Quiet      bool
		zips       int
		cmmts      int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty", fields{}, false},
		{"missing", fields{ExportDir: "/no/such/directory"}, true},
		{"file", fields{ExportDir: "../test/test.txt"}, true},
		{"dir", fields{ExportDir: "../test"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ExportDir:  tt.fields.ExportDir,
				ExportFile: tt.fields.ExportFile,
				NoDupes:    tt.fields.NoDupes,
				Overwrite:  tt.fields.Overwrite,
				Raw:        tt.fields.Raw,
				Print:      tt.fields.Print,
				Quiet:      tt.fields.Quiet,
				zips:       tt.fields.zips,
				cmmts:      tt.fields.cmmts,
			}
			if err := c.Clean(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_Read(t *testing.T) {
	type fields struct {
		ExportDir  string
		ExportFile bool
		NoDupes    bool
		Overwrite  bool
		Raw        bool
		Print      bool
		Quiet      bool
		zips       int
		cmmts      int
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
				ExportDir:  tt.fields.ExportDir,
				ExportFile: tt.fields.ExportFile,
				NoDupes:    tt.fields.NoDupes,
				Overwrite:  tt.fields.Overwrite,
				Raw:        tt.fields.Raw,
				Print:      tt.fields.Print,
				Quiet:      tt.fields.Quiet,
				zips:       tt.fields.zips,
				cmmts:      tt.fields.cmmts,
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
		ExportDir  string
		ExportFile bool
		NoDupes    bool
		Overwrite  bool
		Raw        bool
		Print      bool
		Quiet      bool
		zips       int
		cmmts      int
	}
	tests := []struct {
		name    string
		fields  fields
		root    string
		wantErr bool
	}{
		{"no root", fields{}, "", true},
		{"bad root", fields{}, "../test/missing", true},
		{"test dir", fields{NoDupes: true}, "../test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ExportDir:  tt.fields.ExportDir,
				ExportFile: tt.fields.ExportFile,
				NoDupes:    tt.fields.NoDupes,
				Overwrite:  tt.fields.Overwrite,
				Raw:        tt.fields.Raw,
				Print:      tt.fields.Print,
				Quiet:      tt.fields.Quiet,
				zips:       tt.fields.zips,
				cmmts:      tt.fields.cmmts,
			}
			if err := c.Scan(tt.root); (err != nil) != tt.wantErr {
				t.Errorf("Config.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ExportDir:  tt.fields.ExportDir,
				ExportFile: tt.fields.ExportFile,
				NoDupes:    tt.fields.NoDupes,
				Overwrite:  tt.fields.Overwrite,
				Raw:        tt.fields.Raw,
				Print:      tt.fields.Print,
				Quiet:      tt.fields.Quiet,
				zips:       tt.fields.zips,
				cmmts:      tt.fields.cmmts,
			}
			if err := c.Scans(tt.root); (err != nil) != tt.wantErr {
				t.Errorf("Config.Scans() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_Separator(t *testing.T) {
	type fields struct {
		ExportDir  string
		ExportFile bool
		NoDupes    bool
		Overwrite  bool
		Raw        bool
		Print      bool
		Quiet      bool
		zips       int
		cmmts      int
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
				ExportDir:  tt.fields.ExportDir,
				ExportFile: tt.fields.ExportFile,
				NoDupes:    tt.fields.NoDupes,
				Overwrite:  tt.fields.Overwrite,
				Raw:        tt.fields.Raw,
				Print:      tt.fields.Print,
				Quiet:      tt.fields.Quiet,
				zips:       tt.fields.zips,
				cmmts:      tt.fields.cmmts,
			}
			if got := strings.TrimSpace(c.Separator(tt.fname)); got != tt.want {
				t.Errorf("Config.Separator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Status(t *testing.T) {
	type fields struct {
		ExportDir  string
		ExportFile bool
		NoDupes    bool
		Overwrite  bool
		Raw        bool
		Print      bool
		Quiet      bool
		zips       int
		cmmts      int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"none", fields{}, "Scanned 0 zip archives and found 0 comments\n"},
		{"one", fields{zips: 1, cmmts: 1}, "Scanned 1 zip archive and found 1 comment\n"},
		{"multi", fields{zips: 5, cmmts: 2}, "Scanned 5 zip archives and found 2 comments\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				ExportDir:  tt.fields.ExportDir,
				ExportFile: tt.fields.ExportFile,
				NoDupes:    tt.fields.NoDupes,
				Overwrite:  tt.fields.Overwrite,
				Raw:        tt.fields.Raw,
				Print:      tt.fields.Print,
				Quiet:      tt.fields.Quiet,
				zips:       tt.fields.zips,
				cmmts:      tt.fields.cmmts,
			}
			if got := c.Status(); got != tt.want {
				t.Errorf("Config.Status() = %v, want %v", got, tt.want)
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
