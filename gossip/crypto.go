package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 2048)
	return privkey, &privkey.PublicKey
}

func WriteToFile(content string, filename string) (bool, error) {
	f, err := os.Create(filename)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, e := f.WriteString(content)
	return true, e
}

func ExportRSAPrivateKey(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	ok, _ := WriteToFile(string(privkey_pem), "privkey.pem")
	if !ok {
		return ""
	}
	return string(privkey_pem)
}

func ParseRSAPrivateKey(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ExportRSAPublicKey(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)
	ok, err := WriteToFile(string(pubkey_pem), "pubkey.pem")
	if !ok {
		return "", err
	}
	return string(pubkey_pem), nil
}

func ParseRSAPublicKey(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break
	}
	return nil, errors.New("Key type is not RSA")
}

func Encrypt(message string, publicKey *rsa.PublicKey) []byte {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(message),
		nil)

	if err != nil {
		panic(err)
	}
	return encryptedBytes
}

func Decrypt(encryptedBytes []byte, privateKey *rsa.PrivateKey) []byte {
	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}
	return decryptedBytes
}

func Sign(message string, privateKey rsa.PrivateKey) (string, string) {
	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	checksum := msgHash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, &privateKey, crypto.SHA256, checksum, nil)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(signature), base64.StdEncoding.EncodeToString(checksum)
}

func Verify(checksum string, signature string, publicKey *rsa.PublicKey) bool {
	decodedSignature, _ := base64.StdEncoding.DecodeString(signature)
	decodedChecksum, _ := base64.StdEncoding.DecodeString(checksum)
	err := rsa.VerifyPSS(publicKey, crypto.SHA256, decodedChecksum, decodedSignature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return false
	}
	return true
}

func test_crypto() {
	// Create the keys
	priv, pub := GenerateRSAKeyPair()

	// Export the keys to pem string
	priv_pem := ExportRSAPrivateKey(priv)
	pub_pem, _ := ExportRSAPublicKey(pub)

	// Import the keys from pem string
	priv_parsed, _ := ParseRSAPrivateKey(priv_pem)
	pub_parsed, _ := ParseRSAPublicKey(pub_pem)

	// Export the newly imported keys
	priv_parsed_pem := ExportRSAPrivateKey(priv_parsed)
	pub_parsed_pem, _ := ExportRSAPublicKey(pub_parsed)

	fmt.Println(priv_parsed_pem)
	fmt.Println(pub_parsed_pem)

	message := "something to be encrypted"
	encryptedMessage := Encrypt(message, pub_parsed)

	fmt.Println("Encrypted Message", string(encryptedMessage))

	decryptedMessage := Decrypt(encryptedMessage, priv)
	fmt.Println("Decrypted Message", string(decryptedMessage))

	messageSignature, messageChecksum := Sign(string(decryptedMessage), *priv)
	fmt.Println("Signature:", messageSignature)
	fmt.Println("Checksum:", messageChecksum)

	if Verify(messageChecksum, messageSignature, pub) {
		fmt.Println("Message is verified")
	} else {
		fmt.Println("Message signature does not match")
	}

	if priv_pem != priv_parsed_pem || pub_pem != pub_parsed_pem {
		fmt.Println("Failure: Export and Import did not result in same Keys")
	} else {
		fmt.Println("Success")
	}
}
