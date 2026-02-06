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

// BenchmarkRead measures the performance of reading ZIP file comments
func BenchmarkRead(b *testing.B) {
	// Create a test ZIP file with a comment
	testZip := filepath.Join(b.TempDir(), "test.zip")
	createTestZipWithComment(b, testZip, "test comment for benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := zipcmt.Read(testZip, false)
		if err != nil {
			b.Fatalf("Read failed: %v", err)
		}
	}
}

// BenchmarkReadRaw measures the performance of reading ZIP file comments in raw mode
func BenchmarkReadRaw(b *testing.B) {
	// Create a test ZIP file with a comment
	testZip := filepath.Join(b.TempDir(), "test.zip")
	createTestZipWithComment(b, testZip, "test comment for benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := zipcmt.Read(testZip, true)
		if err != nil {
			b.Fatalf("Read failed: %v", err)
		}
	}
}

// BenchmarkWalkDir measures the performance of walking directories and processing ZIP files
func BenchmarkWalkDir(b *testing.B) {
	// Create a temporary directory with test ZIP files
	tempDir := b.TempDir()
	
	// Create multiple test ZIP files
	for i := 0; i < 10; i++ {
		createTestZipWithComment(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "comment "+string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &zipcmt.Config{
			Print: false,
			Quiet: true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}

// BenchmarkWalkDirWithDupes measures performance when showing all duplicates
func BenchmarkWalkDirWithDupes(b *testing.B) {
	// Create a temporary directory with test ZIP files
	tempDir := b.TempDir()
	
	// Create multiple test ZIP files with duplicate comments
	for i := 0; i < 10; i++ {
		createTestZipWithComment(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "duplicate comment")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &zipcmt.Config{
			Dupes: true,
			Print: false,
			Quiet: true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}

// BenchmarkWalkDirNoWalk measures performance when not walking subdirectories
func BenchmarkWalkDirNoWalk(b *testing.B) {
	// Create a temporary directory with test ZIP files
	tempDir := b.TempDir()
	
	// Create test ZIP files
	for i := 0; i < 10; i++ {
		createTestZipWithComment(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "comment")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &zipcmt.Config{
			NoWalk: true,
			Print: false,
			Quiet: true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}

// Helper function to create a test ZIP file with a comment
func createTestZipWithComment(b *testing.B, zipPath, comment string) {
	b.Helper()
	
	// Create a ZIP file with a comment
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	
	// Add a file
	fw, err := w.Create("test.txt")
	if err != nil {
		b.Fatalf("Failed to create file in zip: %v", err)
	}
	_, err = fw.Write([]byte("test content"))
	if err != nil {
		b.Fatalf("Failed to write to file in zip: %v", err)
	}
	
	// Set the comment
	w.SetComment(comment)
	
	// Close the ZIP writer
	err = w.Close()
	if err != nil {
		b.Fatalf("Failed to close zip writer: %v", err)
	}
	
	// Write to the file
	err = os.WriteFile(zipPath, buf.Bytes(), 0644)
	if err != nil {
		b.Fatalf("Failed to write zip file: %v", err)
	}
}

// BenchmarkLargeDirectory measures performance with many ZIP files
func BenchmarkLargeDirectory(b *testing.B) {
	// Create a temporary directory with many test ZIP files
	tempDir := b.TempDir()
	
	// Create 100 test ZIP files
	for i := 0; i < 100; i++ {
		createTestZipWithComment(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), fmt.Sprintf("comment %d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &zipcmt.Config{
			Print: false,
			Quiet: true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}

// BenchmarkMixedComments measures performance with mixed comment scenarios
func BenchmarkMixedComments(b *testing.B) {
	// Create a temporary directory with various test ZIP files
	tempDir := b.TempDir()
	
	// Create files with different comment scenarios
	for i := 0; i < 50; i++ {
		if i%3 == 0 {
			// No comment
			createTestZipWithComment(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "")
		} else if i%3 == 1 {
			// Short comment
			createTestZipWithComment(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "short")
		} else {
			// Long comment
			createTestZipWithComment(b, filepath.Join(tempDir, fmt.Sprintf("test%d.zip", i)), "This is a much longer comment that contains more text to test performance with larger comment sizes")
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &zipcmt.Config{
			Print: false,
			Quiet: true,
		}
		config.SetTest()
		_ = config.WalkDir(tempDir)
	}
}
