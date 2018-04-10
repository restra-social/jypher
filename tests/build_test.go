package tests

import (
	"testing"
	"github.com/restra-social/jypher/models"
	"strings"
	J "github.com/restra-social/jypher"
	rules2 "github.com/restra-social/jypher/rules"

	"encoding/json"
)

func TestBuildCypher(t *testing.T) {

	data := []byte(`
{
  "foods": [
    {
      "items": [
        {
          "cuisine": "Desserts",
          "description": "Soft, moist, fuzzy chocolaty almond Brownie garnished with chewy chocolate syrap",
          "price": 130,
          "name": "Brownie ",
          "consumable": "1/1"
        },
        {
          "cuisine": "Desserts",
          "description": "Soft, moist, fuzzy chocolaty almond Brownie served hot loaded with a scoope of vanilla icecream  and garnished with chocolate syrap served on sizzling plate",
          "price": 335,
          "name": "Sizzling Brownie ",
          "consumable": "1/1"
        }
      ],
      "code": 68,
      "title": "Sweet sin",
      "tags": [
        "ice",
        "fried"
      ],
      "display": "Sweet"
    }
  ],
  "type": "menu",
  "r_id": "SMCBMU"
}

	`)
	resource := "Menu"

	var unm map[string]interface{}

	err := json.Unmarshal(data, &unm)
	if err != nil {
		t.Error(err.Error())
	}

	jsonInfo := models.JSONInfo{
		DecodedJSON: unm,
		Rules:       getRules(resource, rules2.FHIRRules()),
		Master:      strings.ToLower(resource),
		ID:          "BG1OLE",
	}

	j := J.Jypher{}
	j.GetJypher(&jsonInfo)

	//fmt.Println(graph)

	//data, _ := json.MarshalIndent(graph, " ", " ")

	//fmt.Println(string(data))

	cypher := j.BuildCypher()

	t.Log(cypher)
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
