package utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	logs "pickfighter.com/pkg/logger"
	"github.com/spf13/viper"
)

var ErrAuthCertsPathRequired = fmt.Errorf("authentication certificates path is required")

// Capitalize returns a capitalized string.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

// LoadJwtCerts loads the JWT certificates required for authentication from the specified paths.
// It expects paths to the X.509 certificate (certPath) and private key (keyPath) in the configuration.
// The loaded keypair is used for signing JWT tokens, and the public key is used for token verification.
// If the certificate or key cannot be loaded or parsed, an error is returned.
// The loaded keys are set in the configuration for later use in JWT signing and parsing.
func LoadJwtCerts() error {
	certPath := viper.GetString("auth.jwt.cert")
	keyPath := viper.GetString("auth.jwt.key")

	hasRsaKeys := len(certPath) > 0 && len(keyPath) > 0

	if !hasRsaKeys {
		return ErrAuthCertsPathRequired
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		fmt.Println("ERR", err)
		logs.Errorf("Unable to load client keypair: %s", err)
		return err
	}

	viper.Set("auth.jwt.signing_key", cert.PrivateKey)

	clientCert, err := os.ReadFile(certPath)
	if err != nil {
		logs.Errorf("Unable to read key file bytes: %s", err)
		return err
	}

	block, _ := pem.Decode(clientCert)
	readCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		logs.Errorf("Unable to parse certificate: %s", err)
		return err
	}

	viper.Set("auth.jwt.parse_key", readCert.PublicKey)

	logs.Debugw("Loaded jwt certs",
		"cert_path", viper.GetString("auth.jwt.cert"),
		"key_path", viper.GetString("auth.jwt.key"),
	)

	return nil
}
