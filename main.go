package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var apikey string
var zoneid string
var mail string

type Config struct {
	APIKey string `json:"APIKey"`
	ZONEID string `json:"ZONEID"`
	Mail   string `json:"Mail"`
}

func ConfigLoad() {
	configjson, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	var config Config
	json.Unmarshal(configjson, &config)
	apikey = config.APIKey
	zoneid = config.ZONEID
	mail = config.Mail
	return
}

func DeleteDNS(id string) {
	url := "https://api.cloudflare.com/client/v4/zones/" + zoneid + "/dns_records/" + id
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("X-Auth-Email", mail)
	req.Header.Add("X-Auth-Key", apikey)
	req.Header.Add("Content-Type", "application/json")
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

type DNSs struct {
	Result []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

func main() {
	ConfigLoad()
	url := "https://api.cloudflare.com/client/v4/zones/" + zoneid + "/dns_records"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Auth-Email", mail)
	req.Header.Add("X-Auth-Key", apikey)
	req.Header.Add("Content-Type", "application/json")
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	var DNSs DNSs
	json.Unmarshal(body, &DNSs)
	for _, dns := range DNSs.Result {
		fmt.Println(dns.Name)
		DeleteDNS(dns.ID)
	}
}
