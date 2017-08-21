package main

import (
	"encoding/json"
	"fmt"
	"github.com/kite-social/jypher/generator"
	"github.com/kite-social/jypher/models"
	"reflect"
	"strings"
)

var master string
var id string

var loop int

var parents []string

func main() {

	data := []byte(`
	{
        "id": "0c5d2041-bed9-4451-87ee-7008286fb8ee",
        "meta": {
          "profile": [
            "http://standardhealthrecord.org/fhir/StructureDefinition/shr-encounter-Encounter"
          ]
        },
        "status": "finished",
        "class": {
          "code": "ambulatory"
        },
        "type": [
          {
            "coding": [
              {
                "system": "http://snomed.info/sct",
                "code": "424619006"
              }
            ],
            "text": "Prenatal visit"
          }
        ],
        "subject": {
          "reference": "urn:uuid:9f3a8f66-c58e-4f3a-8192-9c4403ffda47"
        },
        "period": {
          "start": "2007-09-12T23:38:43+06:00",
          "end": "2007-09-12T23:38:43+06:00"
        },
        "reason": {
          "coding": [
            {
              "system": "http://snomed.info/sct",
              "code": "72892002",
              "display": "Normal pregnancy"
            }
          ]
        },
        "serviceProvider": {
          "reference": "urn:uuid:ab9d733f-1b0a-4e12-89e4-320620e5c2eb"
        },
        "resourceType": "Encounter"
      }
	`)

	var unmarshal map[string]interface{}

	err := json.Unmarshal(data, &unmarshal)

	if err != nil {
		panic(err)
	}

	/*var graphs []Graph

	master := createMasterNode(&unmarshal)
	graphs = append(graphs, master)*/

	graph := map[string]models.Graph{}

	master = strings.ToLower(unmarshal["resourceType"].(string))
	id = unmarshal["id"].(string)

	loop = 0

	//fmt.Println(unmarshal)

	generateGraph(id, master, unmarshal, graph)

	data, _ = json.MarshalIndent(graph, " ", " ")

	fmt.Println(string(data))

	generator := generator.CypherGenerator{}

	cypher := generator.Generate(id, graph, parents)

	fmt.Println(cypher)
}

func generateGraph(id string, node string, unmarshal map[string]interface{}, graph map[string]models.Graph) models.Graph {

	// match node in the map if not exists then create
	// if exists then skip.
	// Added to avoid duplicate entry in the model
	if _, ok := graph[node]; !ok {

		graph[node] = models.Graph{}

		var g models.Graph
		g.Nodes.Lebel = node

		g.Edges.Source = master
		g.Edges.Target = node

		parents = append(parents, node)
		//fmt.Println(parents)

		if master == node {
			// only the master node has the id
			g.Nodes.ID = id
		}

		for k, v := range unmarshal {

			var data map[string]interface{}

			t := reflect.ValueOf(v)
			switch t.Kind() {
			case reflect.String:
				pro := map[string]interface{}{
					k: v,
				}

				// if there exists a field called reference then
				// there should be existing node that it is referring to
				// so add the reference id to node id
				if len(pro) > 0 {
					fmt.Println(pro)
					if ref, ok := pro["reference"]; ok {
						g.Nodes.ID = ref.(string)
					}
				}
				g.Nodes.Properties = append(g.Nodes.Properties, pro)

				graph[node] = g

			case reflect.Map:
				data, _ = t.Interface().(map[string]interface{})

				master = node
				generateGraph(id, k, data, graph)
				// loop should reset if we found any object
				loop = 0

			case reflect.Slice:
				slice := reflect.ValueOf(t.Interface())
				length := slice.Len()
				for i := 0; i < length; i++ {
					val := reflect.ValueOf(slice.Index(i).Interface())
					switch val.Kind() {
					case reflect.String:
						// make a saperate node with array of string !!

					case reflect.Map:
						data, _ = slice.Index(i).Interface().(map[string]interface{})
						master = node
						generateGraph(id, fmt.Sprintf("%s%d", k, loop), data, graph)
						loop++
						//parents = parents[:len(parents)-1]
					}
				}
			}

			//node = fmt.Sprintf("%s%d", node, loop)
			graph[node] = g
		}

		return g
	} else {
		// If duplicate found skip the parent but process the child

		//generateGraph(id, master, unmarshal, graph)
	}

	return models.Graph{}
}
