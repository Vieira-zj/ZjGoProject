package etag

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"strings"
)

const (
	blockBits = 22 // Indicate that the blocksize is 4M
	blockSize = 1 << blockBits
)

// GetEtagForText : get etag for text
func GetEtagForText(content string) (etag string, err error) {
	return etagFromReader(strings.NewReader(content), int64(len(content)))
}

func etagFromReader(f io.Reader, size int64) (etag string, err error) {
	blockCnt := blockCount(size)
	sha1Buf := make([]byte, 0, 21)

	if blockCnt <= 1 { // file size <= 4M
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf, err = calSha1(sha1Buf, f)
		if err != nil {
			return
		}
	} else { // file size > 4M
		sha1Buf = append(sha1Buf, 0x96)
		sha1BlockBuf := make([]byte, 0, blockCnt*20)
		for i := 0; i < blockCnt; i++ {
			body := io.LimitReader(f, blockSize)
			sha1BlockBuf, err = calSha1(sha1BlockBuf, body)
			if err != nil {
				return
			}
		}
		sha1Buf, _ = calSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
	}
	etag = base64.URLEncoding.EncodeToString(sha1Buf)
	return
}

func blockCount(fsize int64) int {
	return int((fsize + (blockSize - 1)) >> blockBits)
}

func calSha1(b []byte, r io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}
