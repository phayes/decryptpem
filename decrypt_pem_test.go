package decryptpem

import (
	"crypto/x509"
	"testing"
)

func TestDecryptPEM(t *testing.T) {
	pem, err := DecryptFileWithPassword("./test/encrypted_rsa.pem", "foobar")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = x509.ParsePKCS1PrivateKey(pem.Bytes)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = DecryptFileWithPassword("./test/encrypted_rsa.pem", "badpass")
	if err == nil {
		t.Error("Error should return with bad password")
		return
	}

	// Decrypting with prompt should error on non tty
	pem, err = DecryptFileWithPrompt("./test/encrypted_rsa.pem")
	if err == nil {
		t.Error("Error should return with non tty")
		return
	}

	// Test decrypt on non-encrypted file
	pem, err = DecryptFileWithPassword("./test/plaintext_rsa.pem", "foobar")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = x509.ParsePKCS1PrivateKey(pem.Bytes)
	if err != nil {
		t.Error(err)
		return
	}

	// Decrypting with prompt should be OK on non-encypted file even without tty
	pem, err = DecryptFileWithPrompt("./test/plaintext_rsa.pem")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = x509.ParsePKCS1PrivateKey(pem.Bytes)
	if err != nil {
		t.Error(err)
		return
	}

}
