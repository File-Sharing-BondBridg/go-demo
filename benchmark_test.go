package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"testing"
)

// Test function to ensure tests are detected
func TestEncryption(t *testing.T) {
	// Simple test to verify the package works
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}

func BenchmarkEncrypt1KB(b *testing.B) {
	key := []byte("0123456789abcdef0123456789abcdef")
	block, err := aes.NewCipher(key)
	if err != nil {
		b.Fatalf("Failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		b.Fatalf("Failed to create GCM: %v", err)
	}

	// Prepare test data once
	testData := make([]byte, 1024)
	rand.Read(testData)

	b.ResetTimer()
	b.ReportAllocs() // Report memory allocations

	for i := 0; i < b.N; i++ {
		nonce := make([]byte, gcm.NonceSize())
		io.ReadFull(rand.Reader, nonce)
		ciphertext := gcm.Seal(nil, nonce, testData, nil)
		_ = base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
	}
}

func BenchmarkDecrypt1KB(b *testing.B) {
	key := []byte("0123456789abcdef0123456789abcdef")
	block, err := aes.NewCipher(key)
	if err != nil {
		b.Fatalf("Failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		b.Fatalf("Failed to create GCM: %v", err)
	}

	// Prepare test data and ciphertext once
	testData := make([]byte, 1024)
	rand.Read(testData)

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gcm.Seal(nil, nonce, testData, nil)
	encoded := base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		data, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			b.Fatalf("Failed to decode base64: %v", err)
		}

		nonceSize := gcm.NonceSize()
		if len(data) < nonceSize {
			b.Fatalf("Ciphertext too short")
		}

		nonce, ciphertext := data[:nonceSize], data[nonceSize:]
		_, err = gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}

func BenchmarkEncryptDecryptCycle(b *testing.B) {
	key := []byte("0123456789abcdef0123456789abcdef")
	block, err := aes.NewCipher(key)
	if err != nil {
		b.Fatalf("Failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		b.Fatalf("Failed to create GCM: %v", err)
	}

	testData := make([]byte, 1024)
	rand.Read(testData)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Encrypt
		nonce := make([]byte, gcm.NonceSize())
		io.ReadFull(rand.Reader, nonce)
		ciphertext := gcm.Seal(nil, nonce, testData, nil)
		encoded := base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))

		// Decrypt
		data, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			b.Fatalf("Failed to decode base64: %v", err)
		}

		nonceSize := gcm.NonceSize()
		nonce, ciphertext = data[:nonceSize], data[nonceSize:]
		_, err = gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}
