package unit

import (
	"bytes"
	"educ-gpt/config/dic"
	"educ-gpt/services"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	fileSrv services.FileService

	pathTemp string
)

func TestCanInitFileService(t *testing.T) {
	fileSrv = dic.FileService()
}

func TestCanUploadFile(t *testing.T) {
	filePath := filepath.Join("resources", "test.png")

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
		return
	}

	if _, err := io.Copy(formPart, file); err != nil {
		log.Fatal(err)
		return
	}

	formWriter.Close()

	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		log.Fatal(err)
		return
	}

	path, err := fileSrv.UploadImage(multipartForm.File["file"][0])
	if err != nil {
		t.Error(err)
		return
	}

	pathTemp = strings.TrimPrefix(path, "/")
}

func TestCanRemoveExistFile(t *testing.T) {
	bfrExist, err := fileSrv.DeleteFile(pathTemp)
	if err != nil {
		t.Error(err)
		return
	}

	if !bfrExist {
		t.Error("Expected true but got false")
		return
	}
}

func TestCanRemoveNotExistFile(t *testing.T) {
	bfrExist, err := fileSrv.DeleteFile("not_exist.png")
	if err != nil {
		t.Error(err)
		return
	}

	if bfrExist {
		t.Error("Expected false but got true")
		return
	}
}
