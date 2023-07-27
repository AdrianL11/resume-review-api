package resume

import (
	"bytes"
	"code.sajari.com/docconv"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

func ConvertToPlainText(base64Input string, mimeType string) (string, error) {

	// MimeTypes
	pdfMimeType := "application/pdf"
	docMimeType := "application/msword"
	docxMimeType := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

	// Is MimeType Allowed, PDF, DOC, DOCX
	if mimeType != pdfMimeType && mimeType != docxMimeType && mimeType != docMimeType {
		return "", errors.New("mimetype not allowed")
	}

	// MimeType Accepted
	base64Data := strings.Replace(base64Input, fmt.Sprintf("data:%s;base64,", mimeType), "", -1)

	// Decode Base64
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	// Read Document
	reader := bytes.NewReader(decodedData)
	res, err := docconv.Convert(reader, mimeType, true)
	if err != nil {
		return "", err
	}

	return res.Body, nil
}

func GetMimeType(base64Input string) (string, error) {

	var mimeType = ""

	if !strings.HasPrefix(base64Input, "data:") {
		return mimeType, errors.New("not a base64 data uri")
	}

	mimeType = strings.TrimPrefix(base64Input, "data:")
	splitString := strings.Split(mimeType, ";")

	if len(splitString) < 0 {
		return "", errors.New("not a base64 data uri")
	}

	mimeType = strings.Split(mimeType, ";")[0]

	return mimeType, nil
}
