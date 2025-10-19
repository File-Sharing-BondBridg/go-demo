package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "io"
    "net/http"
)

var key = []byte("0123456789abcdef0123456789abcdef")

type TextPayload struct {
    Text string `json:"text"`
}

type CipherPayload struct {
    Ciphertext string `json:"ciphertext"`
}

func encryptHandler(w http.ResponseWriter, r *http.Request) {
    var req TextPayload
    json.NewDecoder(r.Body).Decode(&req)
    block, _ := aes.NewCipher(key)
    gcm, _ := cipher.NewGCM(block)
    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)
    ciphertext := gcm.Seal(nil, nonce, []byte(req.Text), nil)
    combined := append(nonce, ciphertext...)
    encoded := base64.StdEncoding.EncodeToString(combined)
    json.NewEncoder(w).Encode(map[string]string{"ciphertext": encoded})
}

func decryptHandler(w http.ResponseWriter, r *http.Request) {
    var req CipherPayload
    json.NewDecoder(r.Body).Decode(&req)
    data, _ := base64.StdEncoding.DecodeString(req.Ciphertext)
    block, _ := aes.NewCipher(key)
    gcm, _ := cipher.NewGCM(block)
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
    json.NewEncoder(w).Encode(map[string]string{"text": string(plaintext)})
}

func main() {
    http.HandleFunc("/encrypt", encryptHandler)
    http.HandleFunc("/decrypt", decryptHandler)
    http.ListenAndServe(":8080", nil)
}
