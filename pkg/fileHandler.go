package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ReadJson(path string) []Post {
	var posts []Post
	if _, err := os.Stat(path); err == nil {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal(data, &posts)

		if err != nil {
			fmt.Println(err)
		}
	}
	return posts
}

func WriteJson(path string, posts []Post) error {
	jsonString, _ := json.Marshal(posts)
	err := ioutil.WriteFile(path, jsonString, os.ModePerm)
	return err
}

func UploadFile(req *http.Request, formKey string, targetPath string, filename string) (string, error) {
	file, handler, err := req.FormFile(formKey)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileext := filepath.Ext(handler.Filename)

	abs, _ := filepath.Abs(fmt.Sprintf("%s/%s%s", targetPath, filename, fileext))
	out, err := os.Create(abs)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return strings.Replace(fileext, ".", "", 1), nil
}

func DeleteMediaDir(path string, p Post) {
	logoPath, _ := filepath.Abs(fmt.Sprintf("%s/%s/Logo.%s", path, p.Title, p.LogoType))
	err := os.Remove(logoPath)
	if err != nil {
		fmt.Println(err)
	}
	coverPath, _ := filepath.Abs(fmt.Sprintf("%s/%s/Banner.%s", path, p.Title, p.BannerType))
	err = os.Remove(coverPath)
	if err != nil {
		fmt.Println(err)
	}

	dirPath, _ := filepath.Abs(fmt.Sprintf("%s/%s", path, p.Title))

	os.Remove(dirPath)
}