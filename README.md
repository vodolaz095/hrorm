Hunt-Rest-ORM
==================================

[![GoDoc](https://godoc.org/github.com/vodolaz095/hrorm?status.svg)](https://godoc.org/github.com/vodolaz095/hrorm)

Go client for [HuntJS](https://huntjs.herokuapp.com) models exported by means of [REST interface
](http://huntjs.herokuapp.com/documentation/ExportModelToRestParameters.html).
The exported interface is fully compatible with [best practises](http://www.restapitutorial.com/lessons/httpmethods.html)


Example
==================================

```go

	package main

	import (
		orm "bitbucket.org/vodolaz095/hrorm"
		"fmt"
	)

	/*
		We use the same data structure as the Trophy model structure in api
	*/
	type Trophy struct {
		Id       string `json:"id,omitempty"`
		Name     string `json:"name"`
		Scored   bool   `json:"scored"`
		Priority int    `json:"priority"`
	}

	func main() {
		var doDebug bool = true
		var apiUrl string = "https://huntjs.herokuapp.com/api/v1/trophy/"
		//var apiUrl string = "http://localhost:3000/api/v1/trophy/"

		var huntKey string = "i_am_game_master_grr"

		//authorization is done by header of `HuntKey`
		//on HuntJS application it can be enabled by setting
		//`huntKeyHeader` key to true
		//in config dictionary object
		//see for details
		//https://huntjs.herokuapp.com/documentation/config.html
		hr := orm.New(apiUrl, huntKey, doDebug)

		/*
			Query a list of items
		*/
		fmt.Println("Query...")
		var trophies []Trophy
		var id string
		parameters := make(map[string]string)
		//Set parameters for to limit the collection
		parameters["page"] = "1"
		parameters["itemsPerPage"] = "10"
		parameters["sort"] = "+priority"
		//we can use the MongoDB query operators in this way
		//see for inspiration http://docs.mongodb.org/manual/reference/operator/query/
		parameters["priority[$gte]"] = "5"

		metadata, err := hr.Query(parameters, &trophies)
		if err != nil {
			panic(err)
		}

		fmt.Printf("We got metadata for modelName %v with  %v items count \n", metadata.ModelName, metadata.Count)
		for _, v := range trophies {
			fmt.Printf("Trophy #%v\nName: %v. Priority %v.  Scored:%v\n", v.Id, v.Name, v.Priority, v.Scored)
			id = v.Id
		}

		/*
			Get one item by id
		*/
		var trophy Trophy
		metadata, err1 := hr.GetOne(id, &trophy)
		if err1 != nil {
			panic(err)
		}
		fmt.Println("Get one...")
		fmt.Printf("Trophy #%v\nName: %v. Priority %v.  Scored:%v\n", trophy.Id, trophy.Name, trophy.Priority, trophy.Scored)

		/*
			Creating new trophy
			Notice: if HuntJS server has `disableCsrf` config key set to true
			(default behaviour), we need to perform at least one `Query`,`GetOne`
			requests before, so we have the `hunt.sid` and `csrf` cookies set properly
			Enabling CSRF protection greatly increases security of site
			See https://en.wikipedia.org/wiki/Cross-site_request_forgery
		*/

		newTrophy := Trophy{
			Id:       "", //new entry!
			Name:     "Caleb Doxsey",
			Priority: 100,
			Scored:   false,
		}

		id, err2 := hr.Create(&newTrophy)
		if err2 != nil {
			panic(err2)
		}
		fmt.Printf("Trying to create, got id of %v, so the Caleb Doxsey trophy has id of %v", id, newTrophy.Id)

		/*
			Updating trophy

			Notice: if HuntJS server has `disableCsrf` config key set to true
			(default behaviour), we need to perform at least one `Query`,`GetOne`
			requests before, so we have the `hunt.sid` and `csrf` cookies set properly
			Enabling CSRF protection greatly increases security of site
			See https://en.wikipedia.org/wiki/Cross-site_request_forgery

		*/
		newTrophy.Priority = 200
		err3 := hr.Update(&newTrophy)
		//err3 := hr.Save(&newTrophy) //the same action
		if err3 != nil {
			panic(err3)
		}

		/*
			Deleting trophy created
			Notice: if HuntJS server has `disableCsrf` config key set to true
			(default behaviour), we need to perform at least one `Query`,`GetOne`
			requests before, so we have the `hunt.sid` and `csrf` cookies set properly
			Enabling CSRF protection greatly increases security of site
			See https://en.wikipedia.org/wiki/Cross-site_request_forgery

		*/
		err4 := hr.Delete(&newTrophy)
		if err4 != nil {
			panic(err4)
		} else {
			fmt.Println("Caleb Doxsey wrote excellent book -  http://www.golang-book.com/, so do not query him!")
		}
	}


```
