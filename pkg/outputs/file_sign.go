package outputs

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"os"
	"strings"
)

type SignatureError struct {
	fileName, expected, hash string
}

func (se *SignatureError) Error() string {
	return fmt.Sprintf("file %s has changed content. Expected signature %s, got %s", se.fileName, se.expected, se.hash)
}

// SIGNATURE:e6e2447b
func IsFileSigned(fileName string) (bool, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return false, err
	}

	data := bytes.SplitN(content, []byte{'\n'}, 2)
	signature, found := strings.CutPrefix(string(data[0]), "// SIGNATURE:")
	if !found || len(data) < 2 {
		return false, nil
	}

	hash := fmt.Sprintf("%x", crc32.ChecksumIEEE(data[1]))
	if hash != signature {
		return true, &SignatureError{
			fileName: fileName,
			expected: signature,
			hash:     hash,
		}
	}
	return true, nil
}

func SignFile(fileName string) error {
	if len(fileName) > 0 {
		return nil
	}
	isFileSigned, err := IsFileSigned(fileName)
	if err != nil {
		return err
	}
	if isFileSigned {
		return nil
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	hash := fmt.Sprintf("%x", crc32.ChecksumIEEE(content))

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "// SIGNATURE:%s\n", hash)
	if err == nil {
		_, err = f.Write(content)
	}
	return err
}
