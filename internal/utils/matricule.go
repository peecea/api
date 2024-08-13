package utils

import (
	"crypto/rand"
	"strconv"
	"time"
)

func GenerateMatricule() (matricule string, err error) {
	b := make([]byte, 3)
	_, err = rand.Read(b)
	if err != nil {
		return "", err
	}
	//adding unique random integer to avoid concurent function call

	randomInt := uint64(b[0])<<56 |
		uint64(b[1])<<48
	timestamp := time.Now().UnixMilli()
	matricule = strconv.FormatUint(uint64(timestamp), 10) + strconv.FormatUint(randomInt, 10)
	return matricule, nil
}
