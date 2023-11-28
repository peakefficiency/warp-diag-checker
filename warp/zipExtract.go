package warp

import (
	"archive/zip"
	"bytes"
	"io"
)

var ZipPath string

type FileContent struct {
	Data []byte
}

type FileContentMap map[string]FileContent

func ExtractToMemory(zipPath string) (zipContent FileContentMap, err error) {

	zipContent = make(FileContentMap)

	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {

		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, rc)
		if err != nil {
			return nil, err
		}

		zipContent[file.Name] = FileContent{buf.Bytes()}

	}

	return zipContent, nil

}
