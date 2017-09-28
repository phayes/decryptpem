# decryptpem

Golang package that decrypts encrypted PEM files and blocks. Provides (optional) TTY prompt for input for password. 

## Example
```go
// Get private key, prompt for password and decrypt if necessary
pem, err := decryptpem.DecryptFileWithPrompt("/path/to/private_key.pem")
if err != nil {
  log.Fatal(err)
}
privateKey, err := x509.ParsePKCS1PrivateKey(pem.Bytes());
if err != nil {
  log.Fatal(err)
}
```
