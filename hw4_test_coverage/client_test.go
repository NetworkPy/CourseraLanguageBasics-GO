package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"
)

type TestCase struct {
	Input   SearchRequest
	IsError bool
}

type Persons struct {
	Values []Person `xml:"row"`
}

type Person struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

func (p *Person) getName() string {
	return p.FirstName + " " + p.LastName
}

type PersonsJson struct {
	Values []PersonJson `json:"persons"`
}

type PersonJson struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	About  string `json:"about"`
	Gender string `json:"gender"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	
	query := r.FormValue("query")
	// query := "Everett"
	// Open file
	file, err := os.Open("dataset.xml")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
		// fmt.Errorf("cant read xml file: %s", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
		// fmt.Errorf("cant read data from xml: %s", err)
	}
	var result Persons
	err = xml.Unmarshal(data, &result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
		// fmt.Errorf("cant unmarshal: %s", err)
	}

	var resJson PersonsJson

	switch query {
	case "":
		for _, p := range result.Values {
			resJson.Values = append(resJson.Values, PersonJson{
				Id:     p.Id,
				Name:   p.getName(),
				Age:    p.Age,
				About:  p.About,
				Gender: p.Gender,
			})
		}
	default:
		for _, p := range result.Values {
			if strings.Contains(p.FirstName+" "+p.LastName, query) || strings.Contains(p.About, query) {
				resJson.Values = append(resJson.Values, PersonJson{
					Id:     p.Id,
					Name:   p.getName(),
					Age:    p.Age,
					About:  p.About,
					Gender: p.Gender,
				})
			}
		}
	}

	orderField := r.FormValue("order_field")
	orderBy := r.FormValue("order_by")
	switch orderField {
	case "":
		switch orderBy {
		case "-1":
			sort.Slice(resJson.Values, func(i, j int) bool {
				return resJson.Values[i].Name > resJson.Values[j].Name
			})
		case "1":
			sort.Slice(resJson.Values, func(i, j int) bool {
				return resJson.Values[i].Name < resJson.Values[j].Name
			})
		}
	case "Id":
		switch orderBy {
		case "-1":
			sort.Slice(resJson.Values, func(i, j int) bool {
				return resJson.Values[i].Id > resJson.Values[j].Id
			})
		case "1":
			sort.Slice(resJson.Values, func(i, j int) bool {
				return resJson.Values[i].Id < resJson.Values[j].Id
			})
		}
	case "Age":
		switch orderBy {
		case "-1":
			sort.Slice(resJson.Values, func(i, j int) bool {
				return resJson.Values[i].Age > resJson.Values[j].Age
			})
		case "1":
			sort.Slice(resJson.Values, func(i, j int) bool {
				return resJson.Values[i].Age < resJson.Values[j].Age
			})
		}
	case "Name":
		switch orderBy {
		case "-1":
			sort.Slice(resJson.Values, func(i, j int) bool {
				return resJson.Values[i].Name > resJson.Values[j].Name
			})
		case "1":
			sort.Slice(resJson.Values, func(i, j int) bool {
				return resJson.Values[i].Name < resJson.Values[j].Name
			})
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(resJson.Values)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(j))
}

func TestFindUser(t *testing.T) {
	cases := []TestCase{
		{
			Input: SearchRequest{
				Limit:      -1,
				Offset:     1,
				Query:      "Dil",
				OrderField: "Name",
				OrderBy:    1,
			},
			IsError: true,
		},
		{
			Input: SearchRequest{
				Limit:      27,
				Offset:     -1,
				Query:      "",
				OrderField: "Id",
				OrderBy:    -1,
			},
			IsError: true,
		},
		{
			Input: SearchRequest{
				Limit:      10,
				Offset:     1,
				Query:      "Dil",
				OrderField: "NoName",
				OrderBy:    -1,
			},
			IsError: true,
		},
		{
			Input: SearchRequest{
				Limit:      25,
				Offset:     1,
				Query:      "Dil",
				OrderField: "Name",
				OrderBy:    0,
			},
			IsError: false,
		},
		{
			Input: SearchRequest{
				Limit:      25,
				Offset:     1,
				Query:      "",
				OrderField: "Name",
				OrderBy:    -1,
			},
			IsError: true,
		},
		{
			Input: SearchRequest{
				Limit:      0,
				Offset:     1,
				Query:      "Lynn",
				OrderField: "",
				OrderBy:    -1,
			},
			IsError: true,
		},
		{
			Input: SearchRequest{
				Limit:      0,
				Offset:     1,
				Query:      "Lynn",
				OrderField: "",
				OrderBy:    -1,
			},
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for caseNum, item := range cases {
		c := &SearchClient{
			URL:         ts.URL,
			AccessToken: "cnmcn2527dhalknoiywb",
		}

		_, err := c.FindUsers(item.Input)
		if err != nil && !item.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
	}
	ts.Close()
}
