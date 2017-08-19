package main

import (
	"encoding/json"
	"fmt"
	"github.com/kite-social/jypher-old/generator"
	"github.com/kite-social/jypher-old/models"
	"reflect"
)

var master string
var id string

func main() {

	data := []byte(`
	{
	  "id" : "123456",
	  "session" : {
		"title" : "12 Years of Spring: An Open Source Journey",
		"abstract" : "Spring emerged as a core open source project in early 2003 and evolved to a broad portfolio of open source projects up until 2015.",
		"conference" : {
		  "city" : "London"
		}
	  },
	  "topics" : [
		"keynote",
		"spring"
	  ],
	  "tracks": [{ "main":"Java" }, { "second":"Languages" }, { "third":"Golang" }],
	  "room" : "Auditorium",
	  "timeslot" : "Wed 29th, 09:30-10:30",
	  "speaker" : {
		"name" : "Juergen Hoeller",
		"bio" : "Juergen Hoeller is co-founder of the Spring Framework open source project.",
		"twitter" : "https://twitter.com/springjuergen",
		"picture" : "http://www.springio.net/wp-content/uploads/2014/11/juergen_hoeller-220x220.jpeg"
	  }
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

	master = "talk"
	id = "123456"

	generateGraph(id, master, unmarshal, graph)

	data, _ = json.MarshalIndent(graph, " ", " ")

	fmt.Println(string(data))

	generator := generator.CypherGenerator{}

	fmt.Println(generator.Generate(graph))
}

func generateGraph(id string, node string, unmarshal map[string]interface{}, graph map[string]models.Graph) models.Graph {

	graph[node] = models.Graph{}

	var g models.Graph
	g.Nodes.Lebel = node

	g.Edges.Source = master
	g.Edges.Target = node

	// only the master node has the id
	g.Nodes.ID = id

	for k, v := range unmarshal {
		t := reflect.ValueOf(v)
		switch t.Kind() {
		case reflect.String:
			pro := map[string]interface{}{
				k: v,
			}
			g.Nodes.Properties = append(g.Nodes.Properties, pro)

			graph[node] = g

		case reflect.Map:
			data, _ := t.Interface().(map[string]interface{})

			master = node
			generateGraph(id, k, data, graph)

		case reflect.Slice:
			slice := reflect.ValueOf(t.Interface())
			len := slice.Len()
			for i := 0; i < len; i++ {
				val := reflect.ValueOf(slice.Index(i).Interface())
				switch val.Kind() {
				case reflect.String:
					// make a saperate node with array of string !!

				case reflect.Map:
					data, _ := slice.Index(i).Interface().(map[string]interface{})
					master = node
					generateGraph(id, fmt.Sprintf("%s%d", k, i), data, graph)
				}
			}
		}
	}

	return models.Graph{}
}
