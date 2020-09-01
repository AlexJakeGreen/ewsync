package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Account struct {
	Login    string     `yaml:"login"`
	Password string     `yaml:"password"`
	Maildir  string     `yaml:"maildir"`
	Folders  []*folders `yaml:"folders"`
}

type folders struct {
	Remote string `yaml:"remote"`
	Local  string `yaml:"local"`
}

type config struct {
	Accounts []*Account `yaml:"accounts"`
}

func readConfig() *config {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("Can't read config file %s\n", err.Error())
	}

	c := &config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %s\n", err.Error())
	}

	return c
}

func readState(account *Account, path string) *folderState {
	log.Println("Reading state")
	state := &folderState{}
	if _, err := os.Stat(account.Maildir + "/" + path + "/.state.json"); os.IsNotExist(err) {
		return state
	}

	data, err := ioutil.ReadFile(account.Maildir + "/" + path + "/.state.json")
	if err != nil {
		return state
	}

	result := &folderState{}
	if err := json.Unmarshal(data, result); err != nil {
		return state
	}

	return result
}

func saveState(account *Account, path string, state *folderState) error {
	log.Println("Saving state")
	b, _ := json.Marshal(state)
	err := ioutil.WriteFile(account.Maildir+"/"+path+"/.state.json", b, 0644)
	if err != nil {
		return err
	}

	return nil
}
