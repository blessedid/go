package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Auth struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type configFile struct {
	AccessToken  string `yaml:"access_token"`
	ExpiresIn    int64  `yaml:"expires_in"`
	Login        string `yaml:"login"`
	Password     string `yaml:"password"`
	RefreshToken string `yaml:"refresh_token"`
	Resume       string `yaml:"resume"`
}

func main() {
	config := readYamlFile()

	if config.Login == "" || config.Password == "" {
		log.Fatalln("Not found 'login' or 'password'")
	}

	if config.ExpiresIn < time.Now().Unix() {
		refreshData := refreshToken(config)
		fmt.Println(refreshData.RefreshToken)
		saveYamlFile(config, refreshData)
		config = readYamlFile()
	}

	if config.AccessToken == "" {
		authData := authRequest(config.Login, config.Password, "K811HJNKQA8V1UN53I6PN1J1CMAD2L1M3LU6LPAU849BCT031KDSSM485FDPJ6UF")
		saveYamlFile(config, authData)
		config = readYamlFile()
	}

	updateDateResume(config.Resume, config.AccessToken)
}

func authRequest(login string, password string, token string) Auth {
	result := Auth{}
	params := map[string]string{
		"app_id":           "ru.hh.android",
		"app_version":      "6.71",
		"app_type":         "applicant",
		"platform":         "android",
		"platform_version": "8.1.0",
		"grant_type":       "password",
		"login":            login,
		"password":         password,
	}

	resp, err := multipartRequest("https://hh.ru/oauth/password_credentials?host=hh.ru&locale=RU", params, token)
	if err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()

		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		json.Unmarshal(body.Bytes(), &result)
	}
	return result
}

func multipartRequest(uri string, params map[string]string, token string) (*http.Response, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "ru.hh.android/6.71.11147, Device: Redmi 5, Android OS: 8.1.0 (UUID: 7c115245-b735-432e-a487-c5192d3438fd)")

	client := &http.Client{}
	return client.Do(req)
}

func refreshToken(config configFile) Auth {
	result := Auth{}
	params := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": config.RefreshToken,
	}

	response, err := postRequest("https://hh.ru/oauth/token", params, config.AccessToken)
	checkError(err)
	defer response.Body.Close()

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body.Bytes(), &result)
	fmt.Println(body)
	return result
}

func postRequest(uri string, params map[string]string, token string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	data := url.Values{}
	for key, val := range params {
		data.Add(key, val)
	}
	encodedData := data.Encode()

	req, err := http.NewRequest("POST", uri, strings.NewReader(encodedData))
	checkError(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", "ru.hh.android/6.71.11147, Device: Redmi 5, Android OS: 8.1.0 (UUID: 7c115245-b735-432e-a487-c5192d3438fd)")

	return client.Do(req)
}

func updateDateResume(resume string, token string) {
	fmt.Println(resume, token)
	params := map[string]string{}
	res, err := postRequest("https://api.hh.ru/resumes/"+resume+"/publish", params, token)
	checkError(err)
	defer res.Body.Close()

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(body)
}

func readYamlFile() configFile {
	file, err := ioutil.ReadFile("config.yaml")
	checkError(err)

	data := configFile{}
	err2 := yaml.Unmarshal(file, &data)
	checkError(err2)

	return data
}

func saveYamlFile(config configFile, result Auth) {
	date := time.Now()

	config.AccessToken = result.AccessToken
	config.ExpiresIn = date.Add(time.Duration(result.ExpiresIn) * time.Second).Unix()
	config.RefreshToken = result.RefreshToken

	byte, err := yaml.Marshal(config)
	checkError(err)
	ioutil.WriteFile("config.yaml", byte, 644)
}

func checkError(err error) {
	if err != nil {
		//fmt.Println(err)
		log.Println(err)
	}
}
