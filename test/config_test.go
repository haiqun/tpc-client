package test

import (
	"tcp-tls-project/providers"
	"testing"
)

func TestGetConfig(t *testing.T) {
	// 调用不报错-即可
	providers.GetConfig()
}