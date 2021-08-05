package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendImage(res http.ResponseWriter, path string, filetype string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	res.Header().Set("Content-Type", fmt.Sprintf("image/%s", filetype))
	res.Write(file)
	return nil
}