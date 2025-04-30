package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	mrand "math/rand"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	consumers    *safeMap[string, *consumer]
	clientSignal *safeMap[string, chan struct{}]
	storage      storage
	UnimplementedLumbayLumbayServer
}

func newServer(storageImpl storage) *server {
	return &server{
		consumers:    newSafeMap[string, *consumer](),
		clientSignal: newSafeMap[string, chan struct{}](),
		storage:      storageImpl,
	}
}

func (s *server) generateKeyPair() (privateKeyPEM, publicKeyPEM string, err error) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return "", "", status.Error(codes.Internal, "unable to generate private key")
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", status.Error(codes.Internal, "unable to generate private key bytes")
	}
	privateKeyPEMBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPEM = string(pem.EncodeToMemory(privateKeyPEMBlock))

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", status.Error(codes.Internal, "unable to marshal public key")
	}
	publicKeyPEMBlock := &pem.Block{
		Type:  "PUBLIC KEy",
		Bytes: publicKeyBytes,
	}
	publicKeyPEM = string(pem.EncodeToMemory(publicKeyPEMBlock))

	return privateKeyPEM, publicKeyPEM, nil
}

func generateClientId(publicKeyPEM, salt string) string {
	publicKeyBytes := []byte(salt + publicKeyPEM)
	hash := sha256.Sum256(publicKeyBytes)
	return hex.EncodeToString(hash[:])
}

func generateClientSalt() string {
	return uuid.New().String()
}

func verifyClient(client *Client, publicKeyPEM string) error {
	dataToHash := []byte(client.Salt + publicKeyPEM)
	expectedHash := sha256.Sum256(dataToHash)
	expectedClientId := hex.EncodeToString(expectedHash[:])
	if expectedClientId == client.Id {
		return nil
	}
	return status.Error(codes.InvalidArgument, "invalid client id")
}

func generateCode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	source := mrand.NewSource(int64(time.Now().Nanosecond()))
	r := mrand.New(source)
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}
	return string(code)
}

func generateGameCode() string {
	return generateCode(8)
}
