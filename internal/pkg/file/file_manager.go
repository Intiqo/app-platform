package file

import (
	"errors"
	"io"

	"github.com/gabriel-vasile/mimetype"
)

type Options struct {
	Region      string
	Bucket      string
	Filename    string
	ContentType string
	File        io.Reader
	Directory   string
}

type Manager interface {
	// UploadFile ... Uploads file to a storage bucket and returns the corresponding url
	UploadFile(opts Options) (result string, err error)
}

func GetExtensionAndContentType(file io.Reader) (string, string, error) {
	var extension string
	var contentType string
	mimeType, err := mimetype.DetectReader(file)

	if err != nil {
		return "", "", err
	} else {
		extension = mimeType.Extension()
		contentType = mimeType.String()
	}
	return extension, contentType, nil
}

func ValidateFileType(extension string) error {
	extensionMap := getSupportedExtensions()
	_, ok := extensionMap[extension[1:]]
	if !ok {
		return errors.New("this is an unsupported file type")
	}
	return nil
}

func getSupportedExtensions() map[string]string {
	return map[string]string{"jpg": "jpg", "jpeg": "jpeg", "png": "png", "jfif": "jfif", "pjpeg": "pjpeg", "pjp": "pjp", "pdf": "pdf", "doc": "doc", "docx": "docx", "xls": "xls", "xlsx": "xlsx"}
}
