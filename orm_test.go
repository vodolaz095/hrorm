package hrorm

import (
	"testing"
)

var someId string

type Trophy struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Scored   bool   `json:"scored"`
	Priority int    `json:"priority"`
}

var someTrophy Trophy

var apiUrl string = "https://huntjs.herokuapp.com/api/v1/trophy"

//var apiUrl string = "http://localhost:3000/api/v1/trophy"

var huntKey string = "i_am_game_master_grr"

func TestQueryAll(t *testing.T) {
	hr := New(apiUrl, huntKey, true)
	var trophies []Trophy
	parameters := make(map[string]string)
	metadata, err := hr.Query(parameters, &trophies)
	if err != nil {
		t.Error("We have issues contacting API " + err.Error())
	} else {
		if len(trophies) > 0 {
			for _, v := range trophies {
				if v.Id == "" {
					t.Error("Id is not recieved!")
				}
				someId = v.Id
				if v.Name == "" {
					t.Error("Name is not recieved!")
				}
				if !(v.Scored == true || v.Scored == false) {
					t.Error("Scored is not recieved!")
				}
				if v.Priority < 0 {
					t.Error("Priority is not recieved!")
				}
				someTrophy = v
			}
		} else {
			t.Error("API returned 0 items!")
		}
		/*
			"metadata":{
				"modelName":"Trophy",
				"fieldsAccessible":["id","name","scored","priority"],
				"filter":{},
				"page":1,
				"sort":"-_id",
				"itemsPerPage":10,
				"numberOfPages":1,
				"count":6
				}
		*/
		if metadata.ModelName != "Trophy" {
			t.Error("We recieved model of " + metadata.ModelName + " not `Trophy`")
		}
		if metadata.Sort != "-_id" {
			t.Error("We are sorting by " + metadata.Sort + " while we need to sort by  `-_id`")
		}

		if int(metadata.ItemsPerPage) < len(trophies) {
			t.Error("We recieved wrong itemsPerPage!")
		}

		if metadata.Count < 0 {
			t.Error("We recieved  wrong number of items!")
		}
	}
}

func TestQuerySorted(t *testing.T) {
	hr := New(apiUrl, huntKey, true)
	var trophies []Trophy
	parameters := make(map[string]string)
	parameters["itemsPerPage"] = "2"
	parameters["sort"] = "+name"
	metadata, err := hr.Query(parameters, &trophies)
	if err != nil {
		t.Error("We have issues contacting API " + err.Error())
	} else {
		if len(trophies) > 0 {
			for _, v := range trophies {
				if v.Id == "" {
					t.Error("Id is not recieved!")
				}
				someId = v.Id
				if v.Name == "" {
					t.Error("Name is not recieved!")
				}
				if !(v.Scored == true || v.Scored == false) {
					t.Error("Scored is not recieved!")
				}
				if v.Priority < 0 {
					t.Error("Priority is not recieved!")
				}
				someTrophy = v
			}
		} else {
			t.Error("API returned 0 items!")
		}
		/*
			"metadata":{
				"modelName":"Trophy",
				"fieldsAccessible":["id","name","scored","priority"],
				"filter":{},
				"page":1,
				"sort":"-_id",
				"itemsPerPage":10,
				"numberOfPages":1,
				"count":6
				}
		*/
		if metadata.ModelName != "Trophy" {
			t.Error("We recieved model of " + metadata.ModelName + " not `Trophy`")
		}
		if metadata.Sort != "+name" {
			t.Error("We are sorting by " + metadata.Sort + " while we need to sort by  `+name`")
		}

		if int(metadata.ItemsPerPage) < len(trophies) {
			t.Error("We recieved wrong itemsPerPage!")
		}

		if metadata.ItemsPerPage != 2 {
			t.Error("We recieved wrong itemsPerPage!")
		}

		if metadata.Count < 0 {
			t.Error("We recieved  wrong number of items!")
		}
	}
}

func TestQueryFilteredById(t *testing.T) {
	hr := New(apiUrl, huntKey, true)
	var trophies []Trophy
	parameters := make(map[string]string)
	parameters["_id"] = someTrophy.Id
	metadata, err := hr.Query(parameters, &trophies)
	if err != nil {
		t.Error("We have issues contacting API " + err.Error())
	} else {
		if len(trophies) > 0 {
			for _, v := range trophies {
				if v.Id == "" {
					t.Error("Id is not recieved!")
				}
				someId = v.Id
				if v.Name == "" {
					t.Error("Name is not recieved!")
				}
				if !(v.Scored == true || v.Scored == false) {
					t.Error("Scored is not recieved!")
				}
				if v.Priority < 0 {
					t.Error("Priority is not recieved!")
				}
				someTrophy = v
			}
		} else {
			t.Error("API returned 0 items!")
		}
		/*
			"metadata":{
				"modelName":"Trophy",
				"fieldsAccessible":["id","name","scored","priority"],
				"filter":{},
				"page":1,
				"sort":"-_id",
				"itemsPerPage":10,
				"numberOfPages":1,
				"count":6
				}
		*/
		if metadata.Filter["_id"] != someTrophy.Id {
			t.Error("We were unable to parse filter!")
		}

		if metadata.ModelName != "Trophy" {
			t.Error("We recieved model of " + metadata.ModelName + " not `Trophy`")
		}
		if metadata.Sort != "-_id" {
			t.Error("We are sorting by " + metadata.Sort + " while we need to sort by  `-_id`")
		}

		if int(metadata.ItemsPerPage) < len(trophies) {
			t.Error("We recieved wrong itemsPerPage!")
		}

		if metadata.ItemsPerPage != 10 {
			t.Error("We recieved wrong itemsPerPage!")
		}

		if metadata.Count != 1 {
			t.Error("We recieved  wrong number of items!")
		}
	}

}

func TestGetOneById(t *testing.T) {
	hr := New(apiUrl, huntKey, true)
	var oneTrophy Trophy
	metadata, err := hr.GetOne(someTrophy.Id, &oneTrophy)
	if err != nil {
		t.Error("We have issues contacting API " + err.Error())
	} else {
		if oneTrophy.Id != someTrophy.Id {
			t.Error("We get wrong Trophy")
		}
		if oneTrophy.Name != someTrophy.Name {
			t.Error("We get wrong Trophy")
		}
		if oneTrophy.Scored != someTrophy.Scored {
			t.Error("We get wrong Trophy")
		}
		if oneTrophy.Priority != someTrophy.Priority {
			t.Error("We get wrong Trophy")
		}
		if &metadata == nil {
			t.Error("We got emtpy metadata!")
		}

	}
}

func TestCreateUpdateDelete(t *testing.T) {
	hr := New(apiUrl, huntKey, true)
	var trophies []Trophy
	parameters := make(map[string]string)
	parameters["name"] = "John Doe"
	metadata, err := hr.Query(parameters, &trophies)
	if metadata.Count > 0 {
		t.Error("We found wrong trophy!")
	}

	newTrophy := Trophy{
		Id:       "", //new entry!
		Name:     "John Doe",
		Priority: 100,
		Scored:   false,
	}
	id, err := hr.Create(&newTrophy)
	if err != nil {
		t.Error("We have error creating - " + err.Error())
	} else {
		if id == "" {
			t.Error("We haven't recieved the id!")
		} else {
			var nt Trophy
			_, err1 := hr.GetOne(id, &nt)
			if err1 != nil {
				t.Error("We have error creating - " + err1.Error())
			} else {
				if nt.Id != id {
					t.Error("We recieved wrong id")
				}
				if newTrophy.Id != id {
					t.Error("The id is not updated!")
				}

				if nt.Name != newTrophy.Name {
					t.Error("We recieved wrong name")
				}
				//update
				newTrophy.Priority = 10
				err2 := hr.Update(&newTrophy)
				if err2 != nil {
					t.Error("We have error updating - " + err2.Error())
				}
				//get
				_, err3 := hr.GetOne(id, &nt)
				if err3 != nil {
					t.Error("We have error updating - " + err3.Error())
				}
				if nt.Id != id {
					t.Error("We recieved wrong id")
				}
				if nt.Name != newTrophy.Name {
					t.Error("We recieved wrong name")
				}
				if nt.Priority != 10 {
					t.Error("The priority is not updated!")
				}

				//delete
				err4 := hr.Delete(&newTrophy)
				if err4 != nil {
					t.Error("We have error deleting - " + err4.Error())
				}
				_, err5 := hr.GetOne(id, &nt)
				if err5.Error() != "Not found" {
					t.Error("We got item, so the item is not deleted!")
				}
			}

		}
	}
}

func TestUpdateObjectWithoutId(t *testing.T) {
	newTrophy := Trophy{
		Id:       "", //new entry!
		Name:     "John Doe",
		Priority: 100,
		Scored:   false,
	}
	hr := New(apiUrl, huntKey, true)
	err := hr.Update(&newTrophy)
	if err == nil {
		t.Error("Error is not thrown!")
	}
	if err.Error() != "Object does not have the `Id` field! It cannot be saved!" {
		t.Error("We recieved a bad error!")
	}
}

func TestDeleteObjectWithoutId(t *testing.T) {
	newTrophy := Trophy{
		Id:       "", //new entry!
		Name:     "John Doe",
		Priority: 100,
		Scored:   false,
	}
	hr := New(apiUrl, huntKey, true)
	err := hr.Delete(&newTrophy)
	if err == nil {
		t.Error("Error is not thrown!")
	}
	if err.Error() != "Object does not have the `Id` field! It cannot be saved!" {
		t.Error("We recieved a bad error!")
	}

}
