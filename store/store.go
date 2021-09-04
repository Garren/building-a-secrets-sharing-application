package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/Garren/building-a-secrets-sharing-application/types"
	"golang.org/x/crypto/scrypt"
)

type fileStore struct {
	Store map[string]string
	mutex sync.Mutex
}

var FileStoreConfig struct {
	DataFilePath string
	Fs           fileStore
	Gcm          cipher.AEAD
	Nonce        []byte
}

func initCrypto(password, salt string) (cipher.AEAD, []byte, error) {
	key, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	return gcm, nonce, err
}

func Init(dataFilePath, password, salt string) error {
	_, err := os.Stat(dataFilePath)

	if err != nil {
		_, err := os.Create(dataFilePath)
		if err != nil {
			return err
		}
	}

	gcm, nonce, err := initCrypto(password, salt)
	if err != nil {
		return err
	}

	FileStoreConfig.Fs = fileStore{
		mutex: sync.Mutex{},
		Store: make(map[string]string),
	}

	FileStoreConfig.Gcm = gcm
	FileStoreConfig.Nonce = nonce
	FileStoreConfig.DataFilePath = dataFilePath

	return nil
}

func encrypt(plainText string) (cipherText []byte, err error) {
	if _, err := io.ReadFull(rand.Reader, FileStoreConfig.Nonce); err != nil {
		log.Fatal(err)
	}
	cipherText = FileStoreConfig.Gcm.Seal(
		FileStoreConfig.Nonce, FileStoreConfig.Nonce, []byte(plainText), nil)

	return
}

func decrypt(cipherText []byte) (plainText []byte, err error) {
	nonce := cipherText[:FileStoreConfig.Gcm.NonceSize()]
	cipherText = cipherText[FileStoreConfig.Gcm.NonceSize():]
	plainText, err = FileStoreConfig.Gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return
}

func (j *fileStore) ReadFromFile() error {
	f, err := os.Open(FileStoreConfig.DataFilePath)
	if err != nil {
		return err
	}
	jsonData, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	if len(jsonData) != 0 {
		return json.Unmarshal(jsonData, &j.Store)
	}
	return nil
}

func (j *fileStore) WriteToFile() error {
	jsonData, err := json.Marshal(j.Store)
	if err != nil {
		return err
	}
	ioutil.WriteFile(FileStoreConfig.DataFilePath, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if len(jsonData) != 0 {
		return json.Unmarshal(jsonData, &j.Store)
	}
	return nil
}

func (j *fileStore) Write(data types.SecretData) (err error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	err = j.ReadFromFile()
	if err != nil {
		return
	}

	bytes, err := encrypt(data.Secret)
	if err != nil {
		log.Fatal(err)
	}
	j.Store[data.Id] = base64.StdEncoding.EncodeToString(bytes)
	return j.WriteToFile()
}

func (j *fileStore) Read(id string) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	err := j.ReadFromFile()
	if err != nil {
		return "", err
	}

	if plainText, ok := j.Store[id]; ok {
		bytes, err := base64.StdEncoding.DecodeString(plainText)
		if err != nil {
			return "", err
		}

		delete(j.Store, id)
		j.WriteToFile()

		bytes, err = decrypt(bytes)
		if err != nil {
			return "", err
		}

		plainText = string(bytes)
		return plainText, nil
	}

	return "", err
}
