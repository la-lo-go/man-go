package general_functions

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetHtmlFromPost(url string, payload *strings.Reader) (*goquery.Document, error) {
	method := "POST"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return doc, nil
}

func GetHtmlFromGet(url string) (*goquery.Document, error) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return doc, nil
}

func GetJsonFromGet(url string) (string, error) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil
}
