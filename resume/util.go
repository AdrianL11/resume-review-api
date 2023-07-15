package resume

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/dslipak/pdf"
	"github.com/google/uuid"
	"os"
	"strings"
)

func ParseResume(res string) (string, error) {

	if !strings.Contains(res, "data:application/pdf;base64,") {
		return "", errors.New("no pdf found")
	}

	var resume = strings.Replace(res, "data:application/pdf;base64,", "", -1)
	var _uuid = uuid.New().String()

	dec, err := base64.StdEncoding.DecodeString(resume)
	if err != nil {
		return "", err
	}

	f, err := os.Create(_uuid + ".pdf")
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	if _, err := f.Write(dec); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}

	// Read PDF
	r, err := pdf.Open(_uuid + ".pdf")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)

	err = os.Remove(_uuid + ".pdf")
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
