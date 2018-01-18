package service

import (
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// DownloadFile ... downloads file from url
func DownloadFile(url string, fileName string) error {
	output, err := os.Create(fileName)
	if err != nil {
		log.Warnln("Error while creating", fileName, "-", err)
		return err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		log.Warnln("Error while downloading", url, "-", err)
		return err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		log.Warnln("Error while downloading", url, "-", err)
		return err
	}

	log.WithFields(log.Fields{
		"file":  fileName,
		"bytes": n,
	}).Debug("File has been successfully downloaded")

	return nil
}

// ReadFile ... reads content from file
func ReadFile(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Warnln("Error opening file:", err)
		return nil, err
	}

	return file, nil
}
