package decryptpem

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"
	"syscall"
	"time"

	"github.com/phayes/errors"

	"golang.org/x/crypto/ssh/terminal"
)

// Configuration
var (

	// PasswordDelay sets the delay for any password tries and retries as a defence against brute force password guessing
	// By default there is no delay
	PasswordDelay time.Duration

	// MaxTries sets the maximum number of times a password may be tried before erroring out.
	// A MaxTries of 1 means that there is only one try allowed (no retries)
	// A MaxTries of 0 means infinite retries are allowed.
	// When tries run out, an error of x509.IncorrectPasswordError will be returned.
	MaxTries int
)

// Errors
var (
	ErrReadFile     = errors.New("decryptpem: Cannot read and decrypt file")
	ErrDecryptBlock = errors.New("decryptpem: Cannot decrypt pem block")
	ErrNoBlockFound = errors.New("decryptpem: No PEM block found")
)

// DecryptFileWithPassword retrieives the pem file and decrypts it with the provided password
func DecryptFileWithPassword(filename string, password string) (*pem.Block, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, ErrReadFile)
	}

	block, _, err := DecryptBytesWithPassword(content, password)

	if err != nil {
		return nil, errors.Wrap(err, ErrReadFile)
	}
	if block == nil {
		return nil, errors.Wrap(ErrNoBlockFound, ErrReadFile)
	}

	return block, nil
}

// DecryptFileWithPrompt retrieives the pem file and decrypts it using a prompt from the user
// When password retries run out, a x509.IncorrectPasswordError error will be returned.
func DecryptFileWithPrompt(filename string) (*pem.Block, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, ErrReadFile)
	}

	prompt := "Enter password for " + filename + ": "
	incorrectMessage := "Incorrect password, please try again"

	block, _, err := DecryptBytesWithPrompt(content, prompt, incorrectMessage)

	if err != nil {
		return nil, errors.Wrap(err, ErrReadFile)
	}
	if block == nil {
		return nil, errors.Wrap(ErrNoBlockFound, ErrReadFile)
	}

	return block, nil
}

// DecryptBytesWithPassword will find the next PEM formatted block (certificate, private key etc) in the input.
// It returns that block decrypted and the remainder of the input.
// If no PEM data is found, block is nil and the whole of the input is returned in rest.
func DecryptBytesWithPassword(pembytes []byte, password string) (block *pem.Block, rest []byte, err error) {
	block, rest = pem.Decode(pembytes)
	if block == nil {
		return nil, rest, errors.Wrap(ErrNoBlockFound, ErrDecryptBlock)
	}

	// If it's not encrypted, just return it
	if !x509.IsEncryptedPEMBlock(block) {
		return block, rest, nil
	}

	// It's encrypted, decrypt it
	der, err := x509.DecryptPEMBlock(block, []byte(password))
	if err != nil {
		if errors.IsA(err, x509.IncorrectPasswordError) {
			return nil, rest, err
		}
		return nil, rest, errors.Wrap(err, ErrDecryptBlock)
	}

	// Decryption OK, strip encryption headers and return the decrypted PEMBlock
	delete(block.Headers, "Proc-Type")
	delete(block.Headers, "DEK-Info")
	block.Bytes = der

	return block, rest, nil
}

// DecryptBytesWithPrompt is the same as DecryptBytesWithPassword, but if the
// pem block is password protected, it will prompt stdout / stdin for a password
// When password retries run out, a x509.IncorrectPasswordError error will be returned.
func DecryptBytesWithPrompt(pembytes []byte, prompt string, incorrectMessage string) (block *pem.Block, rest []byte, err error) {
	block, rest = pem.Decode(pembytes)
	if block == nil {
		return nil, rest, errors.Wrap(ErrNoBlockFound, ErrDecryptBlock)
	}

	// If it's not encrypted, just return it
	if !x509.IsEncryptedPEMBlock(block) {
		return block, rest, nil
	}

	// It's encrypted, prompt for password, retrying as needed
	tries := 1
	for {

		// Password prompt delay if configured
		if PasswordDelay != 0 {
			time.Sleep(PasswordDelay)
		}

		// Get the password
		fmt.Print(prompt)
		password, err := getPassword()
		if err != nil {
			return nil, rest, err
		}
		// Print a linebreak to make the password return feel natural
		fmt.Println("")

		block, rest, err = DecryptBytesWithPassword(pembytes, password)
		if err != nil {
			if err == x509.IncorrectPasswordError {
				// If the password is incorrect, either error out or try again depending on configuration
				if MaxTries != 0 && tries >= MaxTries {
					return nil, rest, err
				} else {
					fmt.Println(incorrectMessage)
					tries++
					continue
				}
			}
			return nil, rest, errors.Wrap(ErrNoBlockFound, ErrDecryptBlock)
		}

		// Decryption OK, return block
		return block, rest, nil
	}

}

func getPassword() (string, error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	password := string(bytePassword)

	return strings.TrimSpace(password), nil
}
