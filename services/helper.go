package services

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func ImagesToBase64(str_images string) string {
	f, _ := os.Open(str_images)
	defer f.Close()
	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded
}

func GetRandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}


