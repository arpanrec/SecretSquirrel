package encryption

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/arpanrec/secureserver/internal/common"
	"github.com/arpanrec/secureserver/internal/serverconfig"

	"github.com/ProtonMail/gopenpgp/v2/helper"
)

var encryptionConfig serverconfig.EncryptionConfig

var mu = &sync.Mutex{}
var mo = &sync.Once{}

func setGPGInfo() serverconfig.EncryptionConfig {
	mu.Lock()
	mo.Do(func() {
		encryptionConfig = serverconfig.GetConfig().Encryption
		encryptionConfig.GPGPrivateKey = common.ReadFileStringSureOrStop(&encryptionConfig.GPGPrivateKeyFile)
		encryptionConfig.GPGPublicKey = common.ReadFileStringSureOrStop(&encryptionConfig.GPGPublicKeyFile)

		gpgPassphrase, err2 := os.ReadFile(encryptionConfig.GPGPassphraseFile)
		if err2 != nil {
			log.Fatalln("Error reading passphrase: ", err2)
		}
		gpgPassphraseSanitized := strings.Split(string(gpgPassphrase), "\n")[0]
		encryptionConfig.GPGPrivateKeyPassphrase = []byte(gpgPassphraseSanitized)

		if encryptionConfig.GPGDeleteKeys {
			log.Println("Deleting keys")
			common.DeleteFileSureOrStop(&encryptionConfig.GPGPrivateKeyFile)
			common.DeleteFileSureOrStop(&encryptionConfig.GPGPublicKeyFile)
			common.DeleteFileSureOrStop(&encryptionConfig.GPGPassphraseFile)
		}
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
	decrypted, err := helper.DecryptMessageArmored(encryptionConfig.GPGPrivateKey, encryptionConfig.GPGPrivateKeyPassphrase, *armor)
	if err != nil {
		log.Println("Error decrypting message: ", err)
	}
	*armor = decrypted
	return err
}
