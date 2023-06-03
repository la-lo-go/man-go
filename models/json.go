package models

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

const NYAA_MANGAS_NAMES_JSON_PATH = "json/NYAA_mangasNames.json"
const NYAA_MANGAS_NAMES_JSON_EXPIRE_TIME = 12

type Json struct {
	Path       string `json:"path"`
	ExpireTime int    `json:"expireTime"` // 0 = no expire
}

func NewNyaaSearchJson() *Json {
	return newJson(NYAA_MANGAS_NAMES_JSON_PATH, NYAA_MANGAS_NAMES_JSON_EXPIRE_TIME)
}

func newJson(path string, expireTime int) *Json {
	p := new(Json)
	p.Path = path
	p.ExpireTime = expireTime

	return p
}

func (j Json) Check() (bool, error) {
	fileInfo, err := os.Stat(j.Path)

	if err != nil {
		log.Printf("\n>>>> [" + j.Path + "]: Creating \n\n")
		return false, err
	}

	// The file has been updated in the last expireTime hours or if the expireTime is 0 (infinite)
	if j.ExpireTime == 0 || time.Since(fileInfo.ModTime()) < time.Duration(j.ExpireTime)*time.Hour {
		log.Printf("\n>>>> [" + j.Path + "]: Up to date\n\n")
		return true, nil
	} else {
		log.Printf("\n>>>>[" + j.Path + "]: Updating \n\n")
		return false, nil
	}
}

func (j Json) Read() ([]byte, error) {
	jsonFile, err := os.Open(j.Path)

	if err != nil {
		// Create an empty json file
		err = j.Write(nil)

		if err != nil {
			return []byte{}, err
		}
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	err = jsonFile.Close()
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

/*
Writes the slice 'T' in the json file
*/
func (j Json) Write(T interface{}) error {
	jsonByte, err := json.Marshal(T)

	if err != nil {
		return err
	}

	err = j.save(jsonByte)
	if err != nil {
		return err
	}

	return nil
}

func (j Json) save(jsonByte []byte) error {
	// See if the folder 'json' exists
	_, err := os.ReadDir("json")

	if err != nil {
		// Create the folder 'json'
		err = os.Mkdir("json", os.ModePerm)

		if err != nil {
			return err
		}
	}

	// Write the json file with jsonByte
	err = os.WriteFile(j.Path, jsonByte, 0644)

	if err != nil {
		return err
	}

	return nil
}
