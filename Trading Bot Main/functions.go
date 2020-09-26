package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func pathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func createFiles(){
	var fileName = "config.txt"
	if pathExists(fileName) {
		os.Remove(fileName)
	}
	os.Create("config.txt")
	var insertData = "Binance API Key{123456789}\nBinance Secret Key{987654321}\nTelegram Bot Token{abcdefghiklmnopqrstuvwxyz}"
	fileWriteErr := ioutil.WriteFile(fileName, []byte(insertData), 0644)
	if fileWriteErr != nil{
		panic(fileWriteErr)
	}
	getData(1, fileName)
}

func getData(line int, file string) string{
	data, err := ioutil.ReadFile(file)
	if err != nil{
		panic(err)
	}
	var split []string = strings.Split(string(data), "\n")
	return strings.Replace(strings.Split(split[line], "{")[1], "}", "", 1)
}