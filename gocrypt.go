package gocrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//FillFiles ...
func FillFiles(path string, crypttype string) []string {
	var files []string
	root := path
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if crypttype == "encrypt" {
				if filepath.Ext(path) != ".gocrypt" {
					files = append(files, path)
				}
			} else {
				if filepath.Ext(path) == ".gocrypt" {
					files = append(files, path)
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("fail")
	}

	return files
}

//Goencrypt ...
func Goencrypt(file string, keyString string) {
	key, _ := hex.DecodeString(keyString)
	plaintext, fileerror := ioutil.ReadFile(file)

	if fileerror != nil || len(plaintext) == 0 {
		fmt.Println(file, "error empty file")
		return
	}

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)

	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	errorr := ioutil.WriteFile(file, ciphertext, 0644)

	if errorr != nil {
		fmt.Println(file, "ERROR")
		return
	}

	os.Rename(file, file+".gocrypt")
	fmt.Println("Encrypt successfull")
}

//Godecrypt ...
func Godecrypt(file string, keyString string) {
	key, _ := hex.DecodeString(keyString)
	enc, fileerror := ioutil.ReadFile(file)

	if fileerror != nil {
		fmt.Println(file, "error")
	}

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)

	if err != nil {
		panic(err.Error())
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	if len(plaintext) == 0 {
		fmt.Println("ERROR wrong password?", file)
		return
	}

	if err != nil {
		fmt.Println(file, "ERROR")
		return
	}

	rrorr := ioutil.WriteFile(file, plaintext, 0644)
	if rrorr != nil {
		fmt.Println(file, "ERROR")
		return
	}

	res1 := strings.ReplaceAll(file, ".gocrypt", "")
	os.Rename(file, res1)
	fmt.Println("Decrypt successfull")
}

//GenertateSecurePassword ...
func GenertateSecurePassword(pw string) []byte {
	s := strings.Split(pw, "")

	if len(pw) >= 32 {
		pw = strings.Join(s[:32], "")
		return []byte(pw)
	}

	counter := 0
	for len(pw) < 32 {
		if counter%4 == 0 {
			pw += strconv.Itoa(counter)
		} else {
			pw += s[counter]
		}

		counter++

		if counter >= 9 || counter >= len(s) {
			counter = 0
		}
	}

	fmt.Println("internal PW used: ", pw)

	return []byte(pw)
}
