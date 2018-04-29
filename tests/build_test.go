package tests

import (
	J "github.com/restra-social/jypher"
	"github.com/restra-social/jypher/models"
	rules2 "github.com/restra-social/jypher/rules"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"strings"
	"testing"
	"encoding/json"
	"fmt"
)

func TestBuildCypher(t *testing.T) {

	data := []byte(`
{
  "title": "The New KFC",
  "type": "restaurant",
  "social": {
    "email": "fahm",
    "website": "sitom",
    "facebook": "",
    "phone": "01256983"
  },
  "picture": {
    "logo": "logo",
    "cover": "cover"
  },
  "cuisine": [
    {
      "code": 12,
      "display": "Italian"
    }
  ],
  "tags": ["party", "birthdy", "treat"],
  "delivery": {
    "status": false,
    "area": [
      {
        "code": 1,
        "display": "Modhubag",
        "charge": 45
      }
    ]
  },
  "address": {
    "division": {
      "code": 25,
      "display": "Dhaka"
    },
    "district": {
      "code": 1,
      "display": "Dhaka Zila"
    },
    "area": {
      "code": 1,
      "display": "Banani"
    },
    "postal": 1210,
    "street": "Houde",
    "longitude": 28.02525525,
    "latitude": 58.32656996
  },
  "description": "Some description about this restaurant ..... ",
  "additional": [
    {
      "display": "wifi",
      "code": 1
    },
    {
      "display": "smoking",
      "code": 2
    },
    {
      "display": "restaurant",
      "code": 3
    },
    {
      "display": "pub",
      "code": 4
    },
    {
      "display": "fuck",
      "code": 5
    }
  ],
  "rating": {
    "quality": 0,
    "service": 0,
    "value": 2,
    "place": 0
  },
  "time": {
    "saturday": {
      "open": "time",
      "close": "time"
    },
    "sunday": {
      "open": "time",
      "close": "time"
    },
    "monday": {
      "open": "time",
      "close": "time"
    },
    "tuesday" : {
      "open": "time",
      "close": "time"
    },
    "wednesday": {
      "open": "time",
      "close": "time"
    },
    "thursday": {
      "open": "time",
      "close": "time"
    },
    "friday": {
      "open": "time",
      "close": "time"
    }
  },
  "items": 59,
  "verified": true,
  "status": "pending",
  "created_at": "2013",
  "updated_at": "2016"
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
	if rule, ok := rules[resource]; ok {
		return rule
	} else {
		return models.Rules{}
		//panic(fmt.Sprintf("Rules provided but rules for [%s] not found, check the rules dictionary", resource))
	}

	return models.Rules{}
}


func buildandExecuteCypher(data map[string]interface{}, conn bolt.Conn) string {

	resource := "Restaurant"

	jsonInfo := models.JSONInfo{
		DecodedJSON: data,
		Rules:       getRules(resource, rules2.FHIRRules()),
		Master:      strings.ToLower(resource),
		ID:          "res-id",
	}

	j := J.Jypher{}
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
