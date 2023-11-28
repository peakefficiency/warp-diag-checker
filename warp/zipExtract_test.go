package warp_test

import (
	"archive/zip"
	"bytes"
	"os"
	"testing"

	"github.com/peakefficiency/warp-diag-checker/warp"
)

func createTestZipFile() (string, error) {
	zipFilePath := "test_zip"
	file, err := os.Create(zipFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)

	file1, err := zipWriter.Create("file1.txt")
	if err != nil {
		return "", err
	}
	_, err = file1.Write([]byte("This is file 1"))
	if err != nil {
		return "", err
	}

	file2, err := zipWriter.Create("file2.txt")
	if err != nil {
		return "", err
	}
	_, err = file2.Write([]byte("This is file 2"))
	if err != nil {
		return "", err
	}
	zipWriter.Close()

	return zipFilePath, nil
}

func TestExtractZipToMemory(t *testing.T) {
	t.Parallel()

	zipFilePath, err := createTestZipFile()
	if err != nil {
		t.Errorf("Error creating test zip file: %v", err)
		return
	}
	defer os.Remove(zipFilePath)

	contents, err := warp.ExtractToMemory(zipFilePath)
	if err != nil {
		t.Errorf("Error extracting zip file: %v", err)
		return
	}

	expectedFile1Content := []byte("This is file 1")
	expectedFile2Content := []byte("This is file 2")

	file1Content, ok := contents["file1.txt"]
	if !ok {
		t.Error("Expected file1.txt to be extracted")
		return
	}

	file2Content, ok := contents["file2.txt"]
	if !ok {
		t.Error("Expected file2.txt to be extracted")
		return
	}

	if !bytes.Equal(expectedFile1Content, file1Content.Data) {
		t.Error("file1.txt content does not match expected")
		return
	}

	if !bytes.Equal(expectedFile2Content, file2Content.Data) {
		t.Error("file2.txt content does not match expected")
		return
	}

	_, ok = contents["file3.txt"]
	if ok {
		t.Fatal("expected file not found but received 'ok'")
		return
	}

}
