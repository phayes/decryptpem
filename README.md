# Decrypt PEM

[![Build Status](https://travis-ci.org/phayes/decryptpem.svg?branch=master)](https://travis-ci.org/phayes/decryptpem)
[![Build Status](https://scrutinizer-ci.com/g/phayes/decryptpem/badges/build.png?b=master)](https://scrutinizer-ci.com/g/phayes/decryptpem/build-status/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/phayes/decryptpem)](https://goreportcard.com/report/github.com/phayes/decryptpem)
[![Scrutinizer Issues](https://img.shields.io/badge/scrutinizer-issues-blue.svg)](https://scrutinizer-ci.com/g/phayes/decryptpem/issues)
[![GoDoc](https://godoc.org/github.com/phayes/decryptpem?status.svg)](https://godoc.org/github.com/phayes/decryptpem)

Golang package that decrypts encrypted PEM files and blocks. Provides (optional) TTY prompt for input for password. 

## Installation

```
go get github.com/phayes/decryptpem
```

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


// It will also work with unencrypted plaintext PEM files
pem, err := decryptpem.DecryptFileWithPrompt("/path/to/plaintext_key.pem") // Will not prompt for pasword.
if err != nil {
  log.Fatal(err)
}
privateKey, err := x509.ParsePKCS1PrivateKey(pem.Bytes());
if err != nil {
  log.Fatal(err)
}
```

## Configuration

There are two configuration variables provided:

```go
// PasswordDelay sets the delay for any password tries and retries as a defence against brute force password guessing
// By default there is no delay
var decryptpem.PasswordDelay time.Duration

// MaxTries sets the maximum number of times a password may be tried before erroring out.
// A MaxTries of 1 means that there is only one try allowed (no retries)
// A MaxTries of 0 means infinite retries are allowed.
// When tries run out, an error of x509.IncorrectPasswordError will be returned.
var decryptpem.MaxTries int
```
