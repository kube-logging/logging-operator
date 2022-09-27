package compression

import (
	"bytes"
	"compress/gzip"
	"log"
)

func CompressString(str string) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(str)); err != nil {
		log.Fatal(err)
	}

	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}

	return b.Bytes()
}
