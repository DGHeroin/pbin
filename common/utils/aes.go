package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "io"
)

func AESReader(cipherKey string, r io.Reader) (io.Reader, error) {
    block, err := aes.NewCipher([]byte(cipherKey))
    if err != nil {
        return nil, err
    }
    var iv [aes.BlockSize]byte
    stream := cipher.NewOFB(block, iv[:])

    return &cipher.StreamReader{S: stream, R: r}, nil
}
func AESWriter(cipherKey string, w io.Writer) (io.Writer, error) {
    block, err := aes.NewCipher([]byte(cipherKey))
    if err != nil {
        return nil, err
    }
    var iv [aes.BlockSize]byte
    stream := cipher.NewOFB(block, iv[:])
    return &cipher.StreamWriter{S: stream, W: w}, nil
}
