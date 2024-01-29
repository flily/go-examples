package main

import "testing"

func TestDES(t *testing.T) {
	data := uint64(0x0f1e2d3c4b5a6978)
	key := uint64(0x0011223344556677)

	encrypted := desEncryptBlockUint(data, key)
	decrypted := desDecryptBlockUint(encrypted, key)

	if data != decrypted {
		t.Errorf("data: %016x, encrypted: %016x, decrypted: %016x", data, encrypted, decrypted)
	}
}
