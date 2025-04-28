package gofile

import (
	"encoding/json"
	"errors"
	"main/utils"
)

const (
	referer   = "https://gofile.io/"
	serverUrl = "https://api.gofile.io/servers"
)

type GetServer struct {
	Status string `json:"status"`
	Data   struct {
		Servers map[string]struct {
			Name string `json:"name"`
		} `json:"servers"`
	} `json:"data"`
}

type Upload struct {
	Status string `json:"status"`
	Data   struct {
		DownloadPage string `json:"downloadPage"`
	} `json:"data"`
}

func getServer() (string, error) {
	respBody, err := utils.DoGet(serverUrl, nil, nil)
	if err != nil {
		return "", err
	}
	defer respBody.Close()

	var obj GetServer
	if err := json.NewDecoder(respBody).Decode(&obj); err != nil {
		return "", err
	}
	if obj.Status != "ok" {
		return "", errors.New("bad response from server")
	}
	if len(obj.Data.Servers) == 0 {
		return "", errors.New("no servers available")
	}
	for _, server := range obj.Data.Servers {
		return server.Name, nil
	}
	return "", errors.New("no servers found")
}

func upload(uploadUrl, path string, size, byteLimit int64, headers map[string]string) (string, error) {
	respBody, err := utils.MultipartUpload(uploadUrl, path, "file", size, byteLimit, nil, nil, headers)
	if err != nil {
		return "", err
	}
	defer respBody.Close()

	var obj Upload
	if err := json.NewDecoder(respBody).Decode(&obj); err != nil {
		return "", err
	}
	if obj.Status != "ok" {
		return "", errors.New("bad response during upload")
	}
	return obj.Data.DownloadPage, nil
}

func Run(args *utils.Args, path string) (string, error) {
	size, err := utils.CheckSize(path, "unlim")
	if err != nil {
		return "", err
	}

	server, err := getServer()
	if err != nil {
		return "", err
	}

	uploadUrl := "https://" + server + ".gofile.io/uploadFile"

	headers := map[string]string{
		"Referer": referer,
	}

	url, err := upload(uploadUrl, path, size, args.ByteLimit, headers)
	if err != nil {
		return "", err
	}

	return url, nil
}
