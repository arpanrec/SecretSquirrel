package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/arpanrec/secureserver/internal/serverconfig"
)

var (
	pkiConfigVar serverconfig.PkiConfig
	mu           = &sync.Mutex{}
	oncePki      = &sync.Once{}
)

func getPkiConfig() serverconfig.PkiConfig {
	mu.Lock()
	oncePki.Do(func() {
		pkiConfigVar = serverconfig.GetConfig().PkiConfig
		log.Println("Removing password from CA key: ", pkiConfigVar.CaPrivateKeyFile)
		removePassCmd := exec.Command("openssl",
			"rsa",
			"-in", pkiConfigVar.CaPrivateKeyFile,
			"-passin", "file:"+pkiConfigVar.CaPrivateKeyPasswordFile,
			"-passout", "pass:\"\"",
			"-out", pkiConfigVar.CaPrivateKeyNoPasswordFile)
		if err := removePassCmd.Run(); err != nil {
			log.Fatal("Error removing password from CA key: ", err)
		}
		caCertBytes, err := os.ReadFile(pkiConfigVar.CaCertFile)
		if err != nil {
			log.Fatalln("Error reading ca cert file", err)
		}
		caCertBlock, _ := pem.Decode(caCertBytes)
		CaCert, errParseCert := x509.ParseCertificate(caCertBlock.Bytes)
		if errParseCert != nil {
			log.Fatalln("Error parsing ca cert", errParseCert)
		}
		pkiConfigVar.CaCert = CaCert
		CaPrivateKeyNoPasswordFile := pkiConfigVar.CaPrivateKeyNoPasswordFile
		caPrivKeyBytes, errReadNpPk := os.ReadFile(CaPrivateKeyNoPasswordFile)
		if errReadNpPk != nil {
			log.Fatalln("Error reading ca private key", errReadNpPk)
		}
		caPrivKeyBlock, _ := pem.Decode(caPrivKeyBytes)
		caPrivKey, errParsePKCS8 := x509.ParsePKCS8PrivateKey(caPrivKeyBlock.Bytes)
		if errParsePKCS8 != nil {
			log.Fatalln("Error parsing ca private key", errParsePKCS8)
		}
		pkiConfigVar.CaPrivateNoPasswordKey = caPrivKey.(*rsa.PrivateKey)

		if pkiConfigVar.CaDeleteKeys {
			log.Println("Deleting CA key files")
			errRemoveKey := os.Remove(pkiConfigVar.CaPrivateKeyFile)
			if errRemoveKey != nil {
				log.Fatalln("Error deleting CA key file: ", errRemoveKey)
			}
			errRemoveKey = os.Remove(pkiConfigVar.CaPrivateKeyNoPasswordFile)
			if errRemoveKey != nil {
				log.Fatalln("Error deleting CA key file: ", errRemoveKey)
			}
			errRemoveKey = os.Remove(pkiConfigVar.CaPrivateKeyPasswordFile)
			if errRemoveKey != nil {
				log.Fatalln("Error deleting CA key file: ", errRemoveKey)
			}
			errRemoveKey = os.Remove(pkiConfigVar.CaCertFile)
			if errRemoveKey != nil {
				log.Fatalln("Error deleting CA key file: ", errRemoveKey)
			}
		}
	})
	mu.Unlock()
	return pkiConfigVar
}

func generateSerialNumber() (*big.Int, error) {
	// Generate a random 128-bit number
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	// Convert the number to a byte array and set it as the serial number
	serialNumberBytes := serialNumber.Bytes()
	certSerialNumber := &big.Int{}
	certSerialNumber.SetBytes(serialNumberBytes)
	return certSerialNumber, nil
}

func cetCert(dnsAltNames []string, extKeyUsage []x509.ExtKeyUsage, isCA bool) (string, string, error) {
	pkiCurrentConfig := getPkiConfig()
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Println("Error generating cert private key: ", err)
		return "", "", err
	}
	certSerialNumber, err := generateSerialNumber()
	if err != nil {
		log.Println("Error generating cert serial number: ", err)
		return "", "", err
	}
	cert := &x509.Certificate{
		SerialNumber:   certSerialNumber, // TODO: Convert to upper case hex
		NotBefore:      time.Now(),
		NotAfter:       time.Now().AddDate(0, 0, 30),
		SubjectKeyId:   []byte{1, 2, 3, 4, 6}, // TODO: Create proper subject key id from public key
		ExtKeyUsage:    extKeyUsage,
		KeyUsage:       x509.KeyUsageDigitalSignature,
		DNSNames:       dnsAltNames,
		IsCA:           isCA,
		AuthorityKeyId: pkiCurrentConfig.CaCert.SubjectKeyId,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader,
		cert,
		pkiCurrentConfig.CaCert,
		&certPrivKey.PublicKey,
		pkiCurrentConfig.CaPrivateNoPasswordKey)
	if err != nil {
		log.Println("Error creating cert: ", err)
		return "", "", err
	}
	certPEM := new(bytes.Buffer)
	errPemEncode := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if errPemEncode != nil {
		log.Println("Error encoding cert: ", errPemEncode)
		return "", "", errPemEncode
	}
	certPrivKeyPEM := new(bytes.Buffer)
	errPemEncodePriv := pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if errPemEncodePriv != nil {
		log.Println("Error encoding cert private key: ", errPemEncodePriv)
		return "", "", errPemEncodePriv
	}
	return certPEM.String(), certPrivKeyPEM.String(), nil
}

func GetServerCert(dnsNames []string) (string, string, error) {
	return cetCert(dnsNames,
		[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		false)
}

func GetClientCert(dnsNames []string) (string, string, error) {
	return cetCert(dnsNames, []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, false)
}
