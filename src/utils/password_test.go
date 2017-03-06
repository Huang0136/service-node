package utils

import "testing"

func TestEncryptPassword(t *testing.T) {
	pwd := "XyZ0123%#"
	slat := "0136"
	encryptPwd := EncryptPassword(pwd, slat)

	if encryptPwd != "6deeab6623dfa076c457452f07f63c500c97448a307d39cbb81a81a604e149191aed7ec36a8f22e4352c67704d576f0c0c4ef6982ae10d33ed85ff0c7e5b38fc" {
		t.Error("加密不正确", encryptPwd)
	}
}

func BenchmarkEncryptPassword(b *testing.B) {
	pwd := "XyZ0123%#"
	slat := "0136"

	for i := 0; i < b.N; i++ {
		_ = EncryptPassword(pwd, slat)
	}
}
