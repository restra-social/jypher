package jypher

import (
	"fmt"
	"github.com/restra-social/jypher/generator"
	"github.com/restra-social/jypher/helper"
	"github.com/restra-social/jypher/models"
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
	MainNode string
	NodeHead       string
	ObjIteration int
}

// GetJypher build Object , Init method mostly
func (j *Jypher) GetJypher(jsonInfo *models.JSONInfo) map[string]models.Graph {

	j.Graph = map[string]models.Graph{}
	j.NodeHead = strings.ToLower(jsonInfo.Master)
	j.ID = jsonInfo.ID
	j.ObjIteration = 0
	j.MainNode = jsonInfo.Master
	j.generateGraph(j.NodeHead, jsonInfo.DecodedJSON, jsonInfo.Rules)

	// Append Rules Connection Nodes

	if len(jsonInfo.Rules.Connections) > 0 {
		for _, connectTo := range jsonInfo.Rules.Connections {

			sp := strings.Split(connectTo, "#")
			connectNode := sp[0]
			key := sp[1]

			var conn models.Graph
			conn.Nodes = models.Node{
				ID:    jsonInfo.DecodedJSON[key].(string),
				Lebel: connectNode,
			}

			conn.Edges = models.Edges{
				Source:   j.MainNode,
				Target:   connectNode,
				Relation: "l",
			}

			cn := strings.ToLower(connectNode)
			j.Graph[cn] = conn

			j.Tree = append(j.Tree, cn)

		}
	}

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

			if strings.HasPrefix(j.NodeHead, "type") {
				g.Edges.Source = fmt.Sprintf("%s%s", j.Tree[0], j.NodeHead)
			} else {
				g.Edges.Source = j.NodeHead
			}

			g.Edges.Target = node

			j.Tree = append(j.Tree, node)
			//fmt.Println(parents)

			if j.NodeHead == node {
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
					if nodeName == "coding" && field == "code" {
						g.Nodes.ID = value.(string)
					}

					j.Graph[node] = g

				case reflect.Map:
					data, _ = fieldValue.Interface().(map[string]interface{})

					j.NodeHead = node
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
							j.NodeHead = node
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
	} else {
		fmt.Sprintf("Skipping %s", nodeName)
	}

	return j.Graph
}
