package tests

import (
	"encoding/json"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/restra-social/jypher/models"
	"io/ioutil"
	"os"
	"testing"
	J "github.com/restra-social/jypher"
)

func getGraphConnection(ip string) bolt.Conn{

	driver := bolt.NewDriver()
	conn, _ := driver.OpenNeo(fmt.Sprintf("bolt://neo4j:restra247@%s:7687", ip))
	//defer conn.Close()
	return conn
}

func TestProfile(t *testing.T){


	contents, err := ioutil.ReadFile("/Users/diablo/go/src/gitlab.com/restra-core/document-models/sample-json/user/template/profile.json")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	var data map[string]interface{}
	err = json.Unmarshal(contents, &data)
	if err != nil {
		t.Error(err.Error())
	}

	j := J.Jypher{
		ParentNode: &models.EntityInfo{
			Name: "profile",
			ID:   "1234",
		},
	}
	jsonInfo := models.JSONInfo{
		DecodedJSON: data,
		Rules: &models.Rules{
			SkipField: []string{"name", "work_education"},
		},
	}
	decodedGraph, err := j.GetJypher(jsonInfo)
	if err != nil {
		t.Error(err.Error())
	}

	cypher := j.BuildCypher(decodedGraph)

	//conn := getGraphConnection("66.42.59.213")

	for _, v := range cypher {
		// Start by creating a node
		fmt.Println(v)
		//_, err := conn.ExecNeo(v, nil)
		//if err != nil {
		//	panic(err.Error())
		//}
	}
}