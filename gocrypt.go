package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/purnaresa/bulwark/encryption"
	"github.com/purnaresa/bulwark/utils"
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
			npw := ""
			if len(password) < 32 || len(password) > 32 {
				npw = genertateSecurePassword(password)
			} else {
				npw = password
			}
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
						encryptor([]byte(npw), v)
					}
				} else if typeOfCrypt == "decrypt" {
					for _, v := range allFiles {
						decryptor([]byte(npw), v)
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

func encryptor(secret []byte, file string) {
	image := utils.ReadFile(file)
	encryptionClient := encryption.NewClient()
	cipherImage := encryptionClient.EncryptAES(image, secret)
	if len(cipherImage) == 0 {
		fmt.Println("Wrong password or empty data")
		fmt.Println(file)
		return
	}
	err := utils.WriteFile(cipherImage, file)
	if err != nil {
		fmt.Println("fail encrypt")
		fmt.Println(file)
	}
	os.Rename(file, file+".gocrypt")
	fmt.Println("Encrypt successfull")
}

func decryptor(key []byte, file string) {
	encryptionClient := encryption.NewClient()
	encryptedImage := utils.ReadFile(file)
	plainImage := encryptionClient.DecryptAES(encryptedImage, key)
	if len(plainImage) == 0 {
		fmt.Println("Wrong password or empty data")
		fmt.Println(file)
		return
	}
	err := utils.WriteFile(plainImage, file)
	if err != nil {
		fmt.Println("fail decrypt")
		fmt.Println(file)
	}
	res1 := strings.ReplaceAll(file, ".gocrypt", "")
	os.Rename(file, res1)
	fmt.Println("Decrypt successfull")
}

func genertateSecurePassword(pw string) string {
	s := strings.Split(pw, "")
	if len(pw) > 32 {
		pw = strings.Join(s[:32], "")
		return pw
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
	return pw
}
