package features

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/rchesnokov/tg-bot/service"
	"github.com/rchesnokov/tg-bot/users"
	"github.com/rchesnokov/tg-bot/utils"
	"github.com/ryanuber/columnize"
)

const filterFileName = "resources/filter.txt"

var filter *regexp.Regexp

func init() {
	// Read file
	b, err := ioutil.ReadFile(filterFileName)
	utils.CheckErr(err, "Error while reading file")

	s := strings.Replace(string(b), ", ", "|", -1)

	filter, _ = regexp.Compile(`(?i)(` + s + `)`)
}

// FilterSwearing ... filters message's text for bad words
func FilterSwearing(text string) int {
	count := len(filter.FindAllStringSubmatchIndex(text, -1))
	return count
}

// PrintSwearingRating ... returns string with users' rating
func PrintSwearingRating() string {
	var results []users.User

	db := service.GetDatabase()
	err := db.C("users").Find(nil).Sort("-swearing").All(&results)
	utils.CheckErr(err, "PrintSwearingRating error")

	output := make([]string, len(results))
	config := columnize.DefaultConfig()
	config.Glue = " | "

	for i := 0; i < len(results); i++ {
		best := "🎩"
		if i == 0 {
			best = "🏆"
		} else if i == 1 {
			best = "🥈"
		} else if i == 2 {
			best = "🥉"
		}

		count := fmt.Sprintf("%06d", results[i].Swearing)
		output[i] = best + "  | <code>" + count + "</code> |  " + results[i].Realname
	}

	return columnize.Format(output, config)
}
