package pkg

import (
	"encoding/json"
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
	username := req.FormValue("username")
	password := req.FormValue("password")

	return username == account.Username && CheckPassword(password, account.Password)
}