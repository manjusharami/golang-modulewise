package module5

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
)

// sha1sig return SHA1 signature in the format "35aabcd5a32e01d18a5ef688111624f3c547e13b"
func sha1Sig(data []byte) (string, error) {
	w := sha1.New()
	r := bytes.NewReader(data)
	if _, err := io.Copy(w, r); err != nil {
		return "", err
	}

	sig := fmt.Sprintf("%x", w.Sum(nil))
	return sig, nil
}

type File struct {
	Name      string
	Content   []byte
	Signature string
}

// ValidateSigs return slice of OK files and slice of mismatched files
func ValidateSigs(files []File) ([]string, []string, error) {
	var okFiles []string
	var badFiles []string

	for _, file := range files {
		sig, err := sha1Sig(file.Content)
		if sig != file.Signature || err != nil {
			badFiles = append(badFiles, file.Name)
		} else {
			okFiles = append(okFiles, file.Name)
		}
	}

	return okFiles, badFiles, nil
}