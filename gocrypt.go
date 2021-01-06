package main

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

func main() {

	for {
		var path string
		var password string
		var typeOfCrypt string
		var sure string

		fmt.Print("Path to encrypt/decrypt Data -> ")
		fmt.Scan(&path)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("path not found")
		} else {
			fmt.Print("password -> ")
			fmt.Scan(&password)
			bytes := genertateSecurePassword(password)
			key := hex.EncodeToString(bytes)
			fmt.Printf("key to encrypt/decrypt : %s\n", key)
			fmt.Print("encrypt or decrypt -> ")
			fmt.Scan(&typeOfCrypt)
			allFiles := fillFiles(path, typeOfCrypt)

			for _, f := range allFiles {
				fmt.Println(f)
			}

			fmt.Print("are you sure to do this action y/n -> ")
			fmt.Scan(&sure)

			if sure == "y" {
				if typeOfCrypt == "encrypt" {
					for _, v := range allFiles {
						encrypt(v, key)
					}
				} else if typeOfCrypt == "decrypt" {
					for _, v := range allFiles {
						decrypt(v, key)
					}
				} else {
					fmt.Println("Unknown command")
				}
			} else {
				fmt.Println("Cancelled")
			}
		}
	}
}

func fillFiles(path string, crypttype string) []string {
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

func encrypt(file string, keyString string) {
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

func decrypt(file string, keyString string) {
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

func genertateSecurePassword(pw string) []byte {
	s := strings.Split(pw, "")

	if len(pw) > 32 {
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

		if counter >= len(s) {
			counter = 0
		}
	}
	fmt.Println("internal PW used: ", pw)
	return []byte(pw)
}
