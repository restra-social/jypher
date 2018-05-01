package tests

import (
	J "github.com/restra-social/jypher"
	"github.com/restra-social/jypher/models"
	rules2 "github.com/restra-social/jypher/rules"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"testing"
	"encoding/json"
	"fmt"
	"strings"
)

func TestBuildCypher(t *testing.T) {

	data := []byte(`
{
  "id": "men-id",
  "type": "menu",
  "foods": [
    {
      "display": "Burger",
      "code": 1,
      "title": "Special Hot Sauge Burger",
      "tags": [
        "",
        ""
      ],
      "items": [
        {
          "display": "Restaurant Unique Name",
          "code": 123,
          "offers": [
            {
              "display": "Buy One Get One",
              "code": 12
            }
          ],
          "serial": 1,
          "description": "some des",
          "consumable": {
            "display": "1/3",
            "code": 5
          },
          "ingredients": [
            {
              "display": "Bread",
              "code": 85
            }
          ],
          "cuisine": [
            {
              "display": "Thai",
              "code": 23
            }
          ],
          "size": [
            {
              "display": "Small",
              "code": 4
            },
            {
              "display": "Medium",
              "code": 4
            }
          ],
          "price": 45
        }
      ]
    },
    {
      "display": "Chowmin",
      "code": 1,
      "title": "Pasta & Chowmin",
      "tags": [
        "",
        ""
      ],
      "items": [
        {
          "display": "Restaurant Unique Name",
          "code": 1230,
          "offers": [
            {
              "display": "Buy One Get Half",
              "code": 120
            }
          ],
          "serial": 1,
          "description": "some des",
          "consumable": {
            "display": "1/3",
            "code": 50
          },
          "ingredients": [
            {
              "display": "Pasta",
              "code": 850
            }
          ],
          "cuisine": [
            {
              "display": "Chines",
              "code": 230
            }
          ],
          "size": [
            {
              "display": "Small",
              "code": 4
            },
            {
              "display": "Medium",
              "code": 4
            }
          ],
          "price": 450
        }
      ]
    }
  ]
}

	`)

	var unm map[string]interface{}

	err := json.Unmarshal(data, &unm)
	if err != nil {
		t.Error(err.Error())
	}

	conn := getGraphConnection("localhost")

	buildandExecuteCypher(unm, conn)
}

func getRules(resource string, rules map[string]models.Rules) models.Rules {
	resource = strings.Title(resource)

	if rule, ok := rules[resource]; ok {
		return rule
	} else {
		return models.Rules{}
		//panic(fmt.Sprintf("Rules provided but rules for [%s] not found, check the rules dictionary", resource))
	}

	return models.Rules{}
}


func buildandExecuteCypher(data map[string]interface{}, conn bolt.Conn) string {

	jsonInfo := models.JSONInfo{
		DecodedJSON: data,
		Rules:       getRules("menu", rules2.FHIRRules()),
	}

	j := J.Jypher{
		ConnectionNode: models.EntityInfo{
			Name: "restaurant",
			ID: "res-id",
		},
		ParentNode: models.EntityInfo{
			Name: "menu",
			ID:   "men-id",
		},
	}

	decodedGraph := j.GetJypher(jsonInfo)

	//fmt.Println(graph)

	//data, _ := json.MarshalIndent(graph, " ", " ")

	//fmt.Println(string(data))

	cypher := j.BuildCypher(decodedGraph)

	for _, v := range cypher {

		// Start by creating a node
		fmt.Println(v)
		result, err := conn.ExecNeo(v, nil)
		if err != nil {
			panic(err.Error())
		}
		numResult, _ := result.RowsAffected()
		fmt.Println("CREATED ROWS: ", numResult, "\n") // CREATED ROWS: 2 (per each iteration)

	}

	return "DONE" // CREATED ROWS: 1
}

func getGraphConnection(ip string) bolt.Conn{

	driver := bolt.NewDriver()
	conn, _ := driver.OpenNeo(fmt.Sprintf("bolt://neo4j:123456789@%s:7687", ip))
	//defer conn.Close()

	return conn
}
