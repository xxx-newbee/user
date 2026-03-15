package utils

import (
	"strings"
	"testing"
)

func TestGenerateReferralCode(t *testing.T) {
	t.Run("TestGenerateReferralCode", func(t *testing.T) {
		code, err := GenerateReferralCode()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(code)
	})
}

func TestGenerateCode(t *testing.T) {
	t.Run("TestGenerateCode", func(t *testing.T) {
		code, err := GenerateCode(6)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(strings.ToUpper(code))
	})
}
