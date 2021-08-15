package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateFromJsonFile(path string) (Account, error) {
	var account Account
	content, err := ioutil.ReadFile(path)
	
	if err != nil {
		return account, nil
	}

	err = json.Unmarshal([]byte(content), &account)
	account.Password, _ = HashPassword(account.Password)

	return account, err
}

func HashPassword(password string) (string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func IsOwnerFromReq(req *http.Request, account Account) bool {
	username, password, ok := req.BasicAuth()
	fmt.Println(username)
	fmt.Println(username)
	if !ok {
		return false
	}

	return username == account.Username && CheckPassword(password, account.Password)
}