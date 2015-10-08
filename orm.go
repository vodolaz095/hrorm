// Package hrorm is a Object Relational Mapper between
// Go structures and REST interface exported by HuntJS
// framework `exportModel` function
// It allows to perform create, read, update and delete functions on
// collectons' entries.
// See http://huntjs.herokuapp.com/documentation/ExportModelToRestParameters.html
// for more  details on HuntJS framework
package hrorm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// ORM is a structure, that represents single HuntJS REST API client
type ORM struct {
	APIURL  string
	HuntKey string
	HuntSid string
	Csrf    string
	Debug   bool
}

// New is a constructor, it builds ready to use ORM client instances
// accepting strings of `apiURL`, `huntKey` and `debug` boolean
func New(apiURL, huntKey string, debug bool) ORM {
	return ORM{
		APIURL:  apiURL,
		HuntKey: huntKey,
		HuntSid: "",
		Csrf:    "",
		Debug:   debug,
	}
}

func (o *ORM) prepareRequest(req *http.Request) {
	req.Header.Set("Cookie", fmt.Sprintf("hunt.sid=%v; XSRF-TOKEN=%v", o.HuntSid, o.Csrf))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("huntKey", o.HuntKey)
	req.Header.Set("User-Agent", "Hunt-Rest-Orm")
	if o.Csrf != "" {
		req.Header.Set("X-XSRF-TOKEN", o.Csrf)
	}
}

func (o *ORM) extractFromResponse(res *http.Response) {
	for _, v := range res.Cookies() {
		if v.Name == "XSRF-TOKEN" {
			o.Csrf = v.Value
		}
		if v.Name == "hunt.sid" {
			o.HuntSid = v.Value
		}
	}
	if o.Debug {
		fmt.Printf("Status code %v \n", res.StatusCode)
		fmt.Printf("CSRF %v\n", o.Csrf)
		fmt.Printf("HuntSid %v\n", o.HuntSid)
		h, m, s := time.Now().Clock()
		fmt.Printf("Response completed at %v:%v:%v\n", h, m, s)
		fmt.Println("---------------------------------------------")
		fmt.Println()
		fmt.Println()
	}
}

func makeError(statusCode int) error {
	var e error
	switch statusCode {
	case 400:
		e = errors.New("Bad Request")
		break
	case 401:
		e = errors.New("Unauthorized")
		break
	case 403:
		e = errors.New("Access denied")
		break
	case 404:
		e = errors.New("Not found")
		break
	case 420:
		e = errors.New("Enhance Your Calm")
		break
	case 500:
		e = errors.New("Internal Server Error")
		break
	default:
		e = errors.New("Badaboom!")
	}
	return e
}

// Query extracts entries that satisfies the `parameters`
// and marshals them into the slice of  interface{}'s of `data`
func (o *ORM) Query(parameters map[string]string, data interface{}) (Metadata, error) {
	//http://stackoverflow.com/q/28329938/1885921
	rP := responseParsed{
		Status: "Ok",
		Mtd:    Metadata{},
		Data:   data,
	}
	client := &http.Client{}

	var queryString string

	for k, v := range parameters {
		queryString = queryString + fmt.Sprintf("%v=%v&", url.QueryEscape(k), url.QueryEscape(v))
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%v?%v", o.APIURL, queryString), nil)
	o.prepareRequest(req)
	if o.Debug {
		fmt.Println("---------------------------------------------")
		fmt.Println(fmt.Sprintf("HuntJS ORM - [GET] %v?%v ...", o.APIURL, queryString))
	}

	if err != nil {
		return Metadata{}, err
	}
	res, err1 := client.Do(req)
	defer res.Body.Close()
	if err1 != nil {
		return Metadata{}, err1
	}
	o.extractFromResponse(res)
	if res.StatusCode == 200 {
		raw, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			return Metadata{}, err2
		}
		err2 = json.Unmarshal(raw, &rP)
		data = rP.Data
		return rP.Mtd, nil
	}
	return Metadata{}, makeError(res.StatusCode)
}

// GetOne item by `id` provided and marshal it into `data`
func (o *ORM) GetOne(id string, data interface{}) (Metadata, error) {
	rP := responseParsed{
		Status: "Ok",
		Mtd:    Metadata{},
		Data:   data,
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/%v", o.APIURL, id), nil)
	if o.Debug {
		fmt.Println("---------------------------------------------")
		fmt.Println(fmt.Sprintf("HuntJS ORM - [GET] %v/%v ...", o.APIURL, id))
	}
	o.prepareRequest(req)
	if err != nil {
		return Metadata{}, err
	}
	res, err1 := client.Do(req)
	defer res.Body.Close()
	if err1 != nil {
		return Metadata{}, err1
	}
	o.extractFromResponse(res)
	if res.StatusCode == 200 {
		raw, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			return Metadata{}, err2
		}
		err2 = json.Unmarshal(raw, &rP)
		data = rP.Data
		return rP.Mtd, nil
	}
	return Metadata{}, makeError(res.StatusCode)
}

// Create - function to create new item by marshaling `data` provided
// Notice: if HuntJS server has `disableCsrf` config key set to true
// (default behaviour), we need to perform at least one `Query`,`GetOne`
// requests before, so we have the `hunt.sid` and `csrf` cookies set properly
// Enabling CSRF protection greatly increases security of site
// See https://en.wikipedia.org/wiki/Cross-site_request_forgery
func (o *ORM) Create(data interface{}) (string, error) {
	var id string
	client := &http.Client{}
	body, err0 := json.Marshal(data)
	if err0 != nil {
		return "", err0
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%v", o.APIURL), bytes.NewBuffer(body))
	if o.Debug {
		fmt.Println("---------------------------------------------")
		fmt.Println(fmt.Sprintf("HuntJS ORM - [POST] %v ...", o.APIURL))
	}
	o.prepareRequest(req)
	if err != nil {
		return "", err
	}
	res, err1 := client.Do(req)
	defer res.Body.Close()
	if err1 != nil {
		return "", err1
	}
	o.extractFromResponse(res)
	if res.StatusCode == 201 {
		location := res.Header["Location"][0]
		parts := strings.Split(location, "/")
		id = parts[len(parts)-1]
		_, err2 := o.GetOne(id, &data)
		return id, err2
	}
	return "", makeError(res.StatusCode)
}

// Update the current entry in database by marshaling the `data` provided.
// The entry `id` is extracted from `data.Id`
// Notice: if HuntJS server has `disableCsrf` config key set to true
// (default behaviour), we need to perform at least one `Query`,`GetOne`
// requests before, so we have the `hunt.sid` and `csrf` cookies set properly
// Enabling CSRF protection greatly increases security of site
// See https://en.wikipedia.org/wiki/Cross-site_request_forgery
func (o *ORM) Update(data interface{}) error {
	//https://stackoverflow.com/questions/27992821/how-get-pointer-of-structs-member-from-interface-reflection-golang
	id := string(reflect.ValueOf(data).Elem().FieldByName("ID").String())
	if id == "" {
		return errors.New("Object does not have the `ID` field! It cannot be saved!")
	}
	client := &http.Client{}
	body, err0 := json.Marshal(data)
	if err0 != nil {
		return err0
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%v/%v", o.APIURL, id), bytes.NewBuffer(body))
	if o.Debug {
		fmt.Println("---------------------------------------------")
		fmt.Println(fmt.Sprintf("HuntJS ORM - [PUT] %v/%v ...", o.APIURL, id))
	}
	o.prepareRequest(req)
	if err != nil {
		return err
	}
	res, err1 := client.Do(req)
	defer res.Body.Close()
	if err1 != nil {
		return err1
	}
	o.extractFromResponse(res)
	if res.StatusCode == 200 {
		return nil
	}
	return makeError(res.StatusCode)
}

// Delete entity
// Notice: if HuntJS server has `disableCsrf` config key set to true
// The entry `id` is extracted from `data.Id`
// (default behaviour), we need to perform at least one `Query`,`GetOne`
// requests before, so we have the `hunt.sid` and `csrf` cookies set properly
// Enabling CSRF protection greatly increases security of site
// See https://en.wikipedia.org/wiki/Cross-site_request_forgery
func (o *ORM) Delete(data interface{}) error {
	id := string(reflect.ValueOf(data).Elem().FieldByName("ID").String())
	if id == "" {
		return errors.New("Object does not have the `ID` field! It cannot be saved!")
	}
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/%v", o.APIURL, id), nil)
	if o.Debug {
		fmt.Println("---------------------------------------------")
		fmt.Println(fmt.Sprintf("HuntJS ORM - [DELETE] %v/%v ...", o.APIURL, id))
	}
	o.prepareRequest(req)
	if err != nil {
		return err
	}
	res, err1 := client.Do(req)
	defer res.Body.Close()
	if err1 != nil {
		return err1
	}
	o.extractFromResponse(res)
	if res.StatusCode == 200 {
		return nil
	}
	return makeError(res.StatusCode)
}
