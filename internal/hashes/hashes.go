package hashes

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "hash"
)

func GetHash(kind, text string) (*string, error) {
    var (
        h      hash.Hash
        hashed string
    )
    switch kind {
    case "sha256":
        h = sha256.New()
    default:
        return nil, errors.New("unsupported hashing algorithm")
    }
    h.Write([]byte(text))
    hashed = hex.EncodeToString(h.Sum(nil))
    return &hashed, nil
}