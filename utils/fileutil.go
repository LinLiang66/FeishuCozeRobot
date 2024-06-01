package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// GetFileHash 计算文件的SHA256哈希值
func GetFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashValue := fmt.Sprintf("%x", hash.Sum(nil))
	return hashValue, nil
}
