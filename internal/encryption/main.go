package encryption

import (
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/arpanrec/secureserver/internal/common"
	"log"
	"os"
	"sync"
)

type openGPGInfo struct {
	privateKeyString string

	publicKeyString string

	passphraseString []byte
}

var gpgInfo = openGPGInfo{}

var mu = &sync.Mutex{}
var mo = &sync.Once{}

func setGPGInfo() openGPGInfo {
	mu.Lock()
	mo.Do(func() {
		gpgPrivateKeyPath := common.GetConfig()["encryption"].(map[string]interface{})["private_key_path"].(string)
		gpgPublicKeyPath := common.GetConfig()["encryption"].(map[string]interface{})["public_key_path"].(string)
		gpgPassphraseFilePath := common.GetConfig()["encryption"].(map[string]interface{})["private_key_password_path"].(string)
		gpgPrivateKey, err := os.ReadFile(gpgPrivateKeyPath)
		if err != nil {
			log.Panicln("Error reading private key: ", err)
		}

		gpgPublicKey, err1 := os.ReadFile(gpgPublicKeyPath)
		if err1 != nil {
			log.Panicln("Error reading public key: ", err1)
		}

		gpgPassphrase, err2 := os.ReadFile(gpgPassphraseFilePath)
		if err2 != nil {
			log.Panicln("Error reading passphrase: ", err2)
		}
		gpgInfo = openGPGInfo{
			privateKeyString: string(gpgPrivateKey),
			publicKeyString:  string(gpgPublicKey),
			passphraseString: gpgPassphrase,
		}

		deleteKeys := common.GetConfig()["encryption"].(map[string]interface{})["delete_key_files_after_startup"].(bool)
		if deleteKeys {
			log.Println("Deleting keys")
			err3 := os.Remove(gpgPrivateKeyPath)
			if err3 != nil {
				log.Panicln("Error deleting private key: ", err3)
			}
			err4 := os.Remove(gpgPublicKeyPath)
			if err4 != nil {
				log.Panicln("Error deleting public key: ", err4)
			}
			err5 := os.Remove(gpgPassphraseFilePath)
			if err5 != nil {
				log.Panicln("Error deleting passphrase: ", err5)
			}
		}
	})

	mu.Unlock()
	return gpgInfo
}

func EncryptMessage(message *string) {
	setGPGInfo()
	armor, err := helper.EncryptMessageArmored(gpgInfo.publicKeyString, *message)
	if err != nil {
		log.Panicln("Error encrypting message: ", err)
	}
	*message = armor
}

func DecryptMessage(armor *string) {
	setGPGInfo()
	decrypted, err := helper.DecryptMessageArmored(gpgInfo.privateKeyString, gpgInfo.passphraseString, *armor)
	if err != nil {
		log.Panicln("Error decrypting message: ", err)
	}
	*armor = decrypted
}
