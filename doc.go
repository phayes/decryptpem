// Package decryptpem decrypts encrypted PEM files and blocks. Provides (optional) TTY prompt for input for password.
//
// Installation
//
//   go get github.com/phayes/decryptpem
//
// Example
//
//   // Get private key, prompt for password and decrypt if necessary
//   pem, err := decryptpem.DecryptFileWithPrompt("/path/to/private_key.pem")
//   if err != nil {
//     log.Fatal(err)
//   }
//   privateKey, err := x509.ParsePKCS1PrivateKey(pem.Bytes());
//   if err != nil {
//     log.Fatal(err)
//   }
//
//
//   // It will also work with unencrypted plaintext PEM files
//   pem, err := decryptpem.DecryptFileWithPrompt("/path/to/plaintext_key.pem") // Will not prompt for pasword.
//   if err != nil {
//     log.Fatal(err)
//   }
//   privateKey, err := x509.ParsePKCS1PrivateKey(pem.Bytes());
//   if err != nil {
//     log.Fatal(err)
//   }
//
//
package decryptpem
