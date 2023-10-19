package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/arpanrec/secureserver/internal/appconfig"
	"github.com/arpanrec/secureserver/internal/common"
)

var (
	pkiConfigVar appconfig.ApplicationPkiConfig
	mu           = &sync.Mutex{}
	oncePki      = &sync.Once{}
)

type pkiRequest struct {
	DnsNames []string `json:"dns_names"`
}

type pkiResponse struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func getPkiConfig() appconfig.ApplicationPkiConfig {
	mu.Lock()
	oncePki.Do(func() {
		pkiConfigVar = appconfig.GetConfig().PkiConfig
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
		caCertBytes := common.ReadFileSureOrStop(&pkiConfigVar.CaCertFile)
		caCertBlock, _ := pem.Decode(caCertBytes)
		CaCert, errParseCert := x509.ParseCertificate(caCertBlock.Bytes)
		if errParseCert != nil {
			log.Fatalln("Error parsing ca cert", errParseCert)
		}
		pkiConfigVar.CaCert = CaCert
		caPrivKeyBytes := common.ReadFileSureOrStop(&pkiConfigVar.CaPrivateKeyNoPasswordFile)
		caPrivKeyBlock, _ := pem.Decode(caPrivKeyBytes)
		caPrivKey, errParsePKCS8 := x509.ParsePKCS8PrivateKey(caPrivKeyBlock.Bytes)
		if errParsePKCS8 != nil {
			log.Fatalln("Error parsing ca private key", errParsePKCS8)
		}
		pkiConfigVar.CaPrivateNoPasswordKey = caPrivKey.(*rsa.PrivateKey)

		if pkiConfigVar.CaDeleteKeys {
			log.Println("Deleting CA key files")
			common.DeleteFileSureOrStop(&pkiConfigVar.CaPrivateKeyFile)
			common.DeleteFileSureOrStop(&pkiConfigVar.CaPrivateKeyNoPasswordFile)
			common.DeleteFileSureOrStop(&pkiConfigVar.CaPrivateKeyPasswordFile)
			common.DeleteFileSureOrStop(&pkiConfigVar.CaCertFile)
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

	subjectKeyID := sha1.Sum(certPrivKey.PublicKey.N.Bytes())

	cert := &x509.Certificate{
		SerialNumber:          certSerialNumber,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, 30),
		SubjectKeyId:          subjectKeyID[:],
		ExtKeyUsage:           extKeyUsage,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		DNSNames:              dnsAltNames,
		IsCA:                  isCA,
		AuthorityKeyId:        pkiCurrentConfig.CaCert.SubjectKeyId,
		BasicConstraintsValid: true,
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

func getServerCert(dnsNames []string) (string, string, error) {
	return cetCert(dnsNames,
		[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		false)
}

func getClientCert(dnsNames []string) (string, string, error) {
	return cetCert(dnsNames, []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, false)
}

func GetCert(locationPath *string, body *[]byte) (string, error) {
	pkiRequestJson := pkiRequest{}
	pkiResponseJson := pkiResponse{}

	err := json.Unmarshal(*body, &pkiRequestJson)
	if err != nil {
		return "", err
	}
	if strings.HasSuffix(*locationPath, "clientcert") {
		cert, k, e := getClientCert(pkiRequestJson.DnsNames)
		if e != nil {
			return "", e
		}
		pkiResponseJson.Cert = cert
		pkiResponseJson.Key = k
	} else if strings.HasSuffix(*locationPath, "servercert") {
		cert, k, e := getServerCert(pkiRequestJson.DnsNames)
		if e != nil {
			return "", e
		}
		pkiResponseJson.Cert = cert
		pkiResponseJson.Key = k
	} else {
		return "", errors.New("invalid path for pki: " + *locationPath)
	}
	pkiResponseJsonBytes, err := json.Marshal(pkiResponseJson)
	if err != nil {
		return "", err
	}
	return string(pkiResponseJsonBytes), nil
}
