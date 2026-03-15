package utils

import (
	"crypto/rand"
	"math/big"
)

const chars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

func GenerateReferralCode() (string, error) {
	result := make([]byte, 9)
	charLen := big.NewInt(int64(len(chars)))
	for i := 0; i < 9; i++ {
		idx, err := rand.Int(rand.Reader, charLen)
		if err != nil {
			return "", err
		}
		result[i] = chars[idx.Int64()]
		if i == 3 {
			i++
			result[i] = '-'
		}
	}
	return string(result), nil
}

func GenerateCode(nums int) (string, error) {
	result := make([]byte, nums)
	for i := 0; i < nums; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[idx.Int64()]
	}
	return string(result), nil
}
