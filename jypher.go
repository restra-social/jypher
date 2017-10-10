package jypher

import (
	"fmt"
	"github.com/kite-social/jypher/generator"
	"github.com/kite-social/jypher/helper"
	"github.com/kite-social/jypher/models"
	"reflect"
	"regexp"
	"strings"
)

// Jypher struct
type Jypher struct {
	Graph    map[string]models.Graph
	GraphObj models.JSONInfo
	ID       string
	Tree     []string

	Master       string
	ObjIteration int
}

// GetJypher build Object , Init method mostly
func (j *Jypher) GetJypher(jsonInfo *models.JSONInfo) map[string]models.Graph {

	j.Graph = map[string]models.Graph{}
	j.Master = strings.ToLower(jsonInfo.Master)
	j.ID = jsonInfo.ID
	j.ObjIteration = 0

	j.generateGraph(j.Master, jsonInfo.DecodedJSON, jsonInfo.Rules)

	return j.Graph
}

// BuildCypher : Builds Cypher Query based on Graph Object
func (j *Jypher) BuildCypher() string {

	generator := generator.CypherGenerator{}
	cypher := generator.Generate(j.ID, j.Graph, j.Tree)
	return cypher
}

func (j *Jypher) generateGraph(node string, decodedJSON map[string]interface{}, rules models.Rules) map[string]models.Graph {

	nodeName := regexp.MustCompile(`[A-za-z]+`).FindAllString(node, -1)[0]

	if rules.Rename != nil {
		// apply rename rules before creating a node
		if name, ok := rules.Rename[nodeName]; ok {
			node = regexp.MustCompile(`[A-za-z]+`).ReplaceAllString(node, name.(string))
		}
	}

	// match node in the map if not exists then create
	// if exists then skip.
	// Added to avoid duplicate entry in the model
	// Also Skip Meta and Text field
	if nodeName != "meta" && nodeName != "text" && nodeName != "extension" {
		if _, ok := j.Graph[node]; !ok {

			j.Graph[node] = models.Graph{}

			var g models.Graph
			g.Nodes.Lebel = node

			g.Edges.Source = j.Master
			g.Edges.Target = node

			j.Tree = append(j.Tree, node)
			//fmt.Println(parents)

			if j.Master == node {
				// only the master node has the id
				g.Nodes.ID = j.ID
			}

			for field, value := range decodedJSON {

				var data map[string]interface{}

				fieldValue := reflect.ValueOf(value)

				switch fieldValue.Kind() {
				case reflect.String:
					pro := map[string]interface{}{
						field: value,
					}
					// if there exists a field called reference then
					// there should be existing node that it is referring to
					// so add the reference id to node id
					if len(pro) > 0 {
						if ref, ok := pro["reference"]; ok {
							g.Nodes.ID = helper.IDfilter("urn", ref.(string))
						}
					}
					g.Nodes.Properties = append(g.Nodes.Properties, pro)

					// If nodeName coding then set code value to ID
					if nodeName == "coding" && field == "code"  {
						g.Nodes.ID = value.(string)
					}

					j.Graph[node] = g

				case reflect.Map:
					data, _ = fieldValue.Interface().(map[string]interface{})

					j.Master = node
					j.generateGraph(field, data, rules)
					// loop should reset if we found any object
					j.ObjIteration = 0

				case reflect.Slice:
					slice := reflect.ValueOf(fieldValue.Interface())
					length := slice.Len()
					for i := 0; i < length; i++ {
						val := reflect.ValueOf(slice.Index(i).Interface())
						switch val.Kind() {
						case reflect.String:
							// make a saperate node with array of string !! #todo

						case reflect.Map:
							data, _ = slice.Index(i).Interface().(map[string]interface{})
							j.Master = node
							j.generateGraph(fmt.Sprintf("%s%d", field, j.ObjIteration), data, rules)
							j.ObjIteration++
						}
					}
				}
				j.Graph[node] = g
			}

		} else {
			// If duplicate found skip the parent but process the child #todo

			//generateGraph(id, master, unmarshal, graph)
		}
	}else{
		fmt.Sprintf("Skipping %s", nodeName)
	}

	return j.Graph
}
