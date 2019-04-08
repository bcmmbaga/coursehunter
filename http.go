package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"

	"golang.org/x/net/publicsuffix"

	"github.com/PuerkitoBio/goquery"
)

var (
	lessonClass = ".lessons-list__li"
	descMeta    = "meta[itemprop=description]"
	content     = "content"
	linkAttr    = "link[itemprop=url]"
	linkTag     = "href"
)

type hunter struct {
	path   string
	client *http.Client
	videos []video
}

type video struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	RangeFrom int    `json:"rangeFrom"`
}

var (
	baseURL = "https://coursehunters.net/course"
	authURL = "https://coursehunters.net/sign-in"
)

func newHunter(course, email, password string) (*hunter, error) {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	client := &http.Client{
		Jar: jar,
	}

	if ok := authenticateUser(email, password, client); !ok {
		return nil, errors.New("Authorization failed")
	}

	get := baseURL + "/" + course
	resp, err := client.Get(get)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &hunter{
		path:   course,
		client: client,
		videos: parseCourseContent(resp.Body),
	}, nil
}

func (h *hunter) download() {
	for _, video := range h.videos {
		resp, err := h.client.Get(video.URL)
		if err != nil {
			log.Fatalf("Error while Getting video: %v", err)
		}
		defer resp.Body.Close()

		file := fmt.Sprintf("%s/%s.mp4", h.path, video.Name)
		f, err := os.Create(file)
		if err != nil {
			log.Fatalf("Error Creating file: %v", err)
		}
		defer f.Close()

		done := make(chan int64)

		go showProgress(file, resp.ContentLength, done)
		done <- saveToDisk(resp.Body, f)
	}
}

func authenticateUser(email, password string, client *http.Client) bool {
	resp, err := client.PostForm(authURL, url.Values{
		"e_mail":   {email},
		"password": {password},
	})

	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func parseCourseContent(reader io.Reader) []video {
	contents := make([]video, 0)
	doc, _ := goquery.NewDocumentFromReader(reader)

	doc.Find(lessonClass).Each(func(i int, selector *goquery.Selection) {
		title, _ := selector.Find(descMeta).Attr(content)
		url, _ := selector.Find(linkAttr).Attr(linkTag)
		contents = append(contents, video{
			Name:      renameFileName(title),
			URL:       url,
			RangeFrom: 0,
		})
	})

	return contents
}

func renameFileName(fileName string) string {
	reg := regexp.MustCompile(`[\\\/*:\-|?"\<\>]`)
	return reg.ReplaceAllString(fileName, "")
}
