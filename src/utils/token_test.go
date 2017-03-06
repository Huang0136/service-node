package utils

import "testing"

func TestCreateToken(t *testing.T) {
	token := CreateToken()
	if token == "" {
		t.Error("生成令牌失败")
	}
}

func BenchmarkCreateToken(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		_ = CreateToken()
	}
}
