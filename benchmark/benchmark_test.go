package benchmark_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	zipcmt "github.com/bengarrett/zipcmt/pkg"
)

// BenchmarkRead measures the performance of reading ZIP file comments.
func BenchmarkRead(b *testing.B) {
	testZip := filepath.Join(b.TempDir(), "test.zip")
	createZip(b, testZip, "test comment for benchmark")

	b.ResetTimer()
	for b.Loop() {
		_, err := zipcmt.Read(testZip, false)
		if err != nil {
			b.Fatalf("Read failed: %v", err)
		}
	}
}

// BenchmarkReadRaw measures the performance of reading ZIP file comments in raw mode.
func BenchmarkReadRaw(b *testing.B) {
	testZip := filepath.Join(b.TempDir(), "test.zip")
	createZip(b, testZip, "test comment for benchmark")

	b.ResetTimer()
	for b.Loop() {
		_, err := zipcmt.Read(testZip, true)
		if err != nil {
			b.Fatalf("Read failed: %v", err)
		}
	}
}

// BenchmarkWalkDir measures the performance of walking directories and processing ZIP files.
func BenchmarkWalkDir(b *testing.B) {
	tempDir := b.TempDir()
	for i := range 10 {
		createZip(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "comment "+string(rune(i)))
	}

	b.ResetTimer()
	for b.Loop() {
		config := &zipcmt.Config{
			Print: false,
			Quiet: true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}

// BenchmarkWalkDirWithDupes measures performance when showing all duplicates.
func BenchmarkWalkDirWithDupes(b *testing.B) {
	tempDir := b.TempDir()
	for i := range 10 {
		createZip(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "duplicate comment")
	}

	b.ResetTimer()
	for b.Loop() {
		config := &zipcmt.Config{
			Dupes: true,
			Print: false,
			Quiet: true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}

// BenchmarkWalkDirNoWalk measures performance when not walking subdirectories.
func BenchmarkWalkDirNoWalk(b *testing.B) {
	tempDir := b.TempDir()
	for i := range 10 {
		createZip(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "comment")
	}

	b.ResetTimer()
	for b.Loop() {
		config := &zipcmt.Config{
			NoWalk: true,
			Print:  false,
			Quiet:  true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}

// Helper function to create a test ZIP file with a comment.
func createZip(b *testing.B, zipPath, comment string) {
	b.Helper()

	// create a new zipfile
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// create a new text file to insert into the archive
	fw, err := w.Create("test.txt")
	if err != nil {
		b.Fatalf("Failed to create file in zip: %v", err)
	}
	_, err = fw.Write([]byte("test content"))
	if err != nil {
		b.Fatalf("Failed to write to file in zip: %v", err)
	}
	// set a comment
	err = w.SetComment(comment)
	if err != nil {
		b.Fatalf("Failed to set the comment in zip: %v", err)
	}
	// close the new text file
	err = w.Close()
	if err != nil {
		b.Fatalf("Failed to close zip writer: %v", err)
	}
	// save the zipfile
	err = os.WriteFile(zipPath, buf.Bytes(), 0o600)
	if err != nil {
		b.Fatalf("Failed to write zip file: %v", err)
	}
}

// BenchmarkLargeDirectory measures performance with many ZIP files.
func BenchmarkLargeDirectory(b *testing.B) {
	tempDir := b.TempDir()
	// create 100 test ZIP files
	for i := range 100 {
		createZip(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), fmt.Sprintf("comment %d", i))
	}

	b.ResetTimer()
	for b.Loop() {
		config := &zipcmt.Config{
			Print: false,
			Quiet: true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}

// BenchmarkMixedComments measures performance with mixed comment scenarios.
func BenchmarkMixedComments(b *testing.B) {
	tempDir := b.TempDir()
	for i := range 50 {
		switch i % 3 {
		case 0:
			// No comment
			createZip(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "")
		case 1:
			// Short comment
			createZip(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "short")
		default:
			// Long comment
			createZip(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)),
				"This is a much longer comment that contains more text to test performance with larger comment sizes")
		}
	}

	b.ResetTimer()
	for b.Loop() {
		config := &zipcmt.Config{
			Print: false,
			Quiet: true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}
