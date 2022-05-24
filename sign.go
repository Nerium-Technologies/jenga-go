package jenga

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
	"io/ioutil"
	"strings"
)

func (j *Jenga) getPrivateKey() (*rsa.PrivateKey, error) {
	j.CheckPreconditions()

	key, err := ioutil.ReadFile(j.JengaAuthCredentials.PrivateKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read private key file: %w", err)
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("crypto: no key found")
	}

	rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	return rsa, nil
}

func (j *Jenga) generateSignature(params ...string) (string, error) {
	j.CheckPreconditions()

	sigString := strings.Join(params, "")
	signer, err := j.getPrivateKey()
	if err != nil {
		return "", fmt.Errorf("could not get private key: %w", err)
	}

	signed, err := j.signSHA256([]byte(sigString), signer)
	if err != nil {
		return "", fmt.Errorf("could not sign param string: %w", err)
	}

	signature := base64.StdEncoding.EncodeToString(signed)

	return signature, nil

}

func (j *Jenga) signSHA256(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, d)
}
