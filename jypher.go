package jypher

import (
	"fmt"
	"github.com/restra-social/jypher/generator"
	"github.com/restra-social/jypher/models"
	"reflect"
	"regexp"
	"go/types"
	"errors"
)

const (
	IdentifierField = "code"
)

// Jypher struct
type Jypher struct {
	Tree         []string
	ParentNode   models.EntityInfo
	ConnectionNode models.EntityInfo
	ObjIteration int
}

// GetJypher build Object , Init method mostly
func (j *Jypher) GetJypher(jsonInfo models.JSONInfo) map[string]models.Graph {
	decodedGraph := map[string]models.Graph{}
	j.ObjIteration = 0

	// Remove Skipable field
	// #todo implement nested skipable like time.wednesday
	for _, skip := range jsonInfo.Rules.SkipField {
		 delete(jsonInfo.DecodedJSON, skip)
	}

	j.generateGraph(j.ParentNode, jsonInfo.DecodedJSON, jsonInfo.Rules, decodedGraph)

	return decodedGraph
}

// BuildCypher : Builds Cypher Query based on Graph Object
func (j *Jypher) BuildCypher(decodedGraph map[string]models.Graph) []string {

	generator := generator.CypherGenerator{}
	cypher := generator.Generate(j.ParentNode.ID, decodedGraph, j.Tree)
	return cypher
}

func (j *Jypher) generateGraph(currentNode models.EntityInfo, decodedJSON map[string]interface{}, rules models.Rules, decodedGraph map[string]models.Graph) error{

	nodeName := regexp.MustCompile(`[A-za-z]+`).FindAllString(currentNode.Name, -1)[0]

	if rules.Rename != nil {
		// apply rename rules before creating a node
		if name, ok := rules.Rename[nodeName]; ok {
			currentNode.Name = regexp.MustCompile(`[A-za-z]+`).ReplaceAllString(currentNode.Name, name.(string))
		}
	}

	if _, ok := decodedGraph[currentNode.Name]; !ok {

		var g models.Graph
		g.Nodes.Lebel = currentNode.Name

		// check if Connection Node Available

		if j.ConnectionNode.ID != "" && len(j.Tree) == 0 {
			g.Edges.Source = models.EntityInfo{
				ID:   j.ConnectionNode.ID,
				Name: j.ConnectionNode.Name,
			}
		}else {
			g.Edges.Source = models.EntityInfo{
				ID:   j.ParentNode.ID,
				Name: j.ParentNode.Name,
			}
		}

		g.Edges.Target = currentNode.Name

		j.Tree = append(j.Tree, currentNode.Name)
		//fmt.Println(parents)

		if j.ParentNode.Name == currentNode.Name {
			// only the master node has the id
			g.Nodes.ID = j.ParentNode.ID
		}

		for field, value := range decodedJSON {

			var data map[string]interface{}

			fieldValue := reflect.ValueOf(value)

			switch fieldValue.Kind() {
			case reflect.String, reflect.Float64:
				pro := map[string]interface{}{
					field: value,
				}

				g.Nodes.Properties = append(g.Nodes.Properties, pro)
				// If nodeName coding then set code value to ID
				if field == IdentifierField {
					var val string
					switch value.(type) {
					case string:
						val = value.(string)
					case float64:
						val = fmt.Sprintf("%d", int(value.(float64)))
					}
					g.Nodes.ID = val
				}

				decodedGraph[currentNode.Name] = g

			case reflect.Map:
				data, _ = fieldValue.Interface().(map[string]interface{})
				entityInfo := models.EntityInfo{
					Name: field,
				}
				j.generateGraph(entityInfo, data, rules, decodedGraph)
				// loop should reset if we found any object
				j.ObjIteration = 0

			case reflect.Slice:
				slice := reflect.ValueOf(fieldValue.Interface())
				length := slice.Len()
				j.ParentNode = currentNode

				for i := 0; i < length; i++ {
					object := slice.Index(i).Interface()
					val := reflect.ValueOf(object)
					switch val.Kind() {
					case reflect.String:
						// make a saperate node with array of string !! #todo

					case reflect.Map:
						data, _ = object.(map[string]interface{})
						var id string
						switch data[IdentifierField].(type) {
						case float64:
							id = fmt.Sprintf("%d", int(data[IdentifierField].(float64)))
							break
						case types.Nil:
							return errors.New(fmt.Sprintf("Id not found on Data %+v", data))
						}
						entityInfo := models.EntityInfo{
							Name: fmt.Sprintf("%s%d", field, j.ObjIteration),
							ID: id,
						}
						j.ParentNode = currentNode
						j.generateGraph(entityInfo, data, rules, decodedGraph)
						j.ObjIteration++
					}
				}
			}
			decodedGraph[currentNode.Name] = g
		}
	}

	return nil
}
