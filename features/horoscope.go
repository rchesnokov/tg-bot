package features

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"time"

	"github.com/rchesnokov/tg-bot/service"
	"github.com/rchesnokov/tg-bot/utils"

	"github.com/djherbis/times"
	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"
)

const serviceURL = "http://ignio.com/r/export/utf/xml/daily/bus.xml"
const horoscopeFileName = "resources/horoscope.xml"

var signTrans = map[string]string{
	"Aries":       "Овен ♈️",
	"Taurus":      "Телец ♉️",
	"Gemini":      "Близнецы ♊️",
	"Cancer":      "Рак ♋️",
	"Leo":         "Лев ♌️",
	"Virgo":       "Дева ♍️",
	"Libra":       "Весы ♎️",
	"Scorpio":     "Скорпион ♏️",
	"Sagittarius": "Стрелец ♐️",
	"Capricorn":   "Козерог ♑️",
	"Aquarius":    "Водолей ♒️",
	"Pisces":      "Рыбы ♓️",
}

type date struct {
	Yesterday  string `xml:"yesterday,attr"`
	Today      string `xml:"today,attr"`
	Tomorrow   string `xml:"tomorrow,attr"`
	Tomorrow02 string `xml:"tomorrow02,attr"`
}

type signDates struct {
	Yesterday  string `xml:"yesterday"`
	Today      string `xml:"today"`
	Tomorrow   string `xml:"tomorrow"`
	Tomorrow02 string `xml:"tomorrow02"`
}

type forecast struct {
	Date        *date      `xml:"date"`
	Aries       *signDates `xml:"aries"`
	Taurus      *signDates `xml:"taurus"`
	Gemini      *signDates `xml:"gemini"`
	Cancer      *signDates `xml:"cancer"`
	Leo         *signDates `xml:"leo"`
	Virgo       *signDates `xml:"virgo"`
	Libra       *signDates `xml:"libra"`
	Scorpio     *signDates `xml:"scorpio"`
	Sagittarius *signDates `xml:"sagittarius"`
	Capricorn   *signDates `xml:"capricorn"`
	Aquarius    *signDates `xml:"aquarius"`
	Pisces      *signDates `xml:"pisces"`
}

// ProvideHoroscope ... returns today's horoscope for given user
func ProvideHoroscope(birthdate string) string {
	date, _ := time.Parse("2006-01-02", birthdate)
	_, month, day := date.Date()

	log.WithFields(log.Fields{
		"string": birthdate,
		"month":  month,
		"day":    day,
	}).Debug("User's birthday")

	// Download new if file doesn't exist
	if _, err := os.Stat(horoscopeFileName); os.IsNotExist(err) {
		err := service.DownloadFile(serviceURL, horoscopeFileName)
		utils.CheckErr(err, "Error while downloading horoscope")
	}

	// Get file's creation date
	t, err := times.Stat(horoscopeFileName)
	if err != nil {
		utils.CheckErr(err, "Error while reading xml stats")
	}

	// Download new file if current one has expired
	if !t.ChangeTime().Truncate(time.Hour * 24).Equal(time.Now().Truncate(time.Hour * 24)) {
		err := service.DownloadFile(serviceURL, horoscopeFileName)
		utils.CheckErr(err, "Error while downloading horoscope")
	}

	// Read file
	xmlFile, err := service.ReadFile(horoscopeFileName)
	utils.CheckErr(err, "Error while reading horoscope")
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	var f forecast
	xml.Unmarshal(b, &f)

	userSign := getSign(month.String(), day)
	fstruct := structs.New(f)
	messageField := fstruct.Field(userSign).Field("Today")
	messageValue := messageField.Value().(string)
	message := signTrans[userSign] + "\n-------------------" + messageValue

	return message
}

func getSign(month string, day int) string {
	switch month {
	case "January":
		if day < 21 {
			return "Capricorn"
		}
		return "Aquarius"

	case "February":
		if day < 20 {
			return "Aquarius"
		}
		return "Pisces"

	case "March":
		if day < 21 {
			return "Pisces"
		}
		return "Aries"

	case "April":
		if day < 21 {
			return "Aries"
		}
		return "Taurus"

	case "May":
		if day < 21 {
			return "Taurus"
		}
		return "Gemini"

	case "June":
		if day < 22 {
			return "Gemini"
		}
		return "Cancer"

	case "July":
		if day < 24 {
			return "Cancer"
		}
		return "Leo"

	case "August":
		if day < 24 {
			return "Leo"
		}
		return "Virgo"

	case "September":
		if day < 24 {
			return "Virgo"
		}
		return "Libra"

	case "October":
		if day < 24 {
			return "Libra"
		}
		return "Scorpio"

	case "November":
		if day < 23 {
			return "Scorpio"
		}
		return "Sagittarius"

	case "December":
		if day < 22 {
			return "Sagittarius"
		}
		return "Capricorn"
	}

	return ""
}
