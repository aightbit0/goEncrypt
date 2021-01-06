# goEncrypt

A Encryption Programm written in GO to encrypt and decrypt Data 

## ATTENTION IF YOU ENCRYPT THE DATA AND FORGET THE PASSWORD, THE DATA WILL BE GONE FOREVER!
Use at your own risk !
i am not responsible for damage caused by the program

# Install
```
go get -u github.com/aightbit0/goEncrypt
```

# Usage
```golang
//generate a strong password
bytes := gocrypt.GenertateSecurePassword("a strong password")

//convert to hexadecimal
key := hex.EncodeToString(bytes)

//encrypt data with given password
gocrypt.Goencrypt(pathtofile, key)

//decrypt data with given password
gocrypt.Godecrypt(pathtofile, key)
```

