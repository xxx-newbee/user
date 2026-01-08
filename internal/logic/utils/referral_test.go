package utils

import "testing"

func TestGenerateReferralCode(t *testing.T) {
	t.Run("TestGenerateReferralCode", func(t *testing.T) {
		code, err := GenerateReferralCode()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(code)
	})
}
