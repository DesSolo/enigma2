package main

import (
	"crypto/rand"
	"enigma/storage"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

// GenRegexpGetView ...
func GenRegexpGetView() *regexp.Regexp {
	p := fmt.Sprintf("^/get/(\\w{%d})$", Config.TokenBytes * 2)
	return regexp.MustCompile(p)

}

// LoadTemplate ...
func LoadTemplate(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("error load template file: %s erro: %s", file, err)
	}
	return data
}

// GenerateUniqToken ...
func GenerateUniqToken(s storage.SecretStorage, retries int) (string, error) {
	for i := 0; ; i++ {
		key, err := GenerateToken(Config.TokenBytes)
		if err != nil {
			return "", err
		}
		uniq, err := s.IsUniq(key)
		if err != nil {
			return "", err
		}
		if uniq == true {
			return key, nil
		}
		if i >= (retries-1) {
			break
		}
	}
	return "", errors.New("maximum retries save")
}

// GenerateToken ...
func GenerateToken(bytes int) (string, error) {
	b := make([]byte, bytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
