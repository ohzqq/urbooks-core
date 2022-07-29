package audible

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gosimple/slug"
)

func getAsin(path string) string {
	paths := strings.Split(path, "/")
	return paths[len(paths)-1]
}

func DownloadCover(name, u string) {
	response, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatal(response.StatusCode)
	}

	file, err := os.Create(slug.Make(name) + ".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}

var countryCodes = map[string]string{
	"us": ".com",
	"ca": ".ca",
	"uk": ".co.uk",
	"au": ".co.uk",
	"fr": "fr",
}

func countrySuffix(code string) string {
	switch code {
	case "uk", "au", "jp", "in":
		return ".co." + code
	default:
		return "." + code
	}
}

func countryCode(suffix string) string {
	switch suffix {
	case "com":
		return "us"
	default:
		return suffix
	}
}
