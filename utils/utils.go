package utils

import log "github.com/sirupsen/logrus"

// CheckErr ... makes fatal if error occured
func CheckErr(err error, text string) {
	if err != nil {
		log.WithField("error", err).Fatal(text)
	}
}
