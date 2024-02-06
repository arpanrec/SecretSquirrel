package encryption

import (
	"log"
	"sync"

	"github.com/arpanrec/secretsquirrel/internal/appconfig"
	"github.com/arpanrec/secretsquirrel/internal/common"

	"github.com/ProtonMail/gopenpgp/v2/helper"
)

var encryptionConfig *appconfig.ApplicationEncryptionConfig

var mu = &sync.Mutex{}
var mo = &sync.Once{}

func setGPGInfo() *appconfig.ApplicationEncryptionConfig {
	mu.Lock()
	mo.Do(func() {
		newEncryptionConfig := appconfig.GetConfig().Encryption
		err := common.IfFileExistsReplaceStringVal(&newEncryptionConfig.GPGPrivateKey)
		if err != nil {
			log.Fatalln("Error reading private key: ", err)
		}
		errPublicKey := common.IfFileExistsReplaceStringVal(&newEncryptionConfig.GPGPublicKey)
		if errPublicKey != nil {
			log.Fatalln("Error reading public key: ", errPublicKey)
		}

		if newEncryptionConfig.GPGPrivateKeyPassphrase != "" {
			errPassphrase := common.IfFileExistsReplaceStringVal(&newEncryptionConfig.GPGPrivateKeyPassphrase)
			if errPassphrase != nil {
				log.Fatalln("Error reading passphrase: ", errPassphrase)
			}
		}
		encryptionConfig = &newEncryptionConfig
	})

	mu.Unlock()
	return encryptionConfig
}

func EncryptMessage(message *string) error {
	setGPGInfo()
	armor, err := helper.EncryptMessageArmored(encryptionConfig.GPGPublicKey, *message)
	if err != nil {
		log.Println("Error encrypting message: ", err)
	}
	*message = armor
	return err
}

func DecryptMessage(armor *string) error {
	setGPGInfo()
	decrypted, err := helper.DecryptMessageArmored(encryptionConfig.GPGPrivateKey, []byte(encryptionConfig.GPGPrivateKeyPassphrase), *armor)
	if err != nil {
		log.Println("Error decrypting message: ", err)
	}
	*armor = decrypted
	return err
}
