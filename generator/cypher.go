package generator

import (
	"fmt"
	"github.com/kite-social/jypher/models"
	"regexp"
	"strings"
)

// CypherGenerator : Cypher Query Generator
type CypherGenerator struct{}

// Generate : This method takes a graph model and generates cypher query
func (c *CypherGenerator) Generate(models map[string]models.Graph, serial []string) (cypher string) {

	// loop through the serial
	for _, term := range serial {

		// search for key in the model in ascending order to generate the query
		if k , ok := models[term]; ok {

			level := k.Nodes.Lebel

			pl := len(k.Nodes.Properties)

			node := regexp.MustCompile(`[A-za-z]+`).FindAllString(strings.Title(level), -1)[0]

			relation := fmt.Sprintf("%s_%s", strings.ToUpper(k.Edges.Source), strings.ToUpper(node))

			if k.Edges.Source == k.Edges.Target {
				cypher += fmt.Sprintf("MERGE (%s:%s {id:'%s'}) ON CREATE SET ", level, node, k.Nodes.ID)

				for i, property := range k.Nodes.Properties {
					for key, val := range property {
						cypher += fmt.Sprintf("%s.%s = '%s'", k.Nodes.Lebel, key, val)
					}
					if pl > 1 {
						if i < pl-1 {
							cypher += ", "
						}
					}

				}
				if k.Edges.Source != k.Edges.Target { // avoids self loop
					cypher += fmt.Sprintf("MERGE (%s)-[:%s]->(%s)", k.Edges.Source, relation, k.Edges.Target)
				}

				cypher += " "
			} else {

				cypher += fmt.Sprintf("CREATE (%s:%s) SET ", level, node)

				for _, property := range k.Nodes.Properties {
					for key, val := range property {
						cypher += fmt.Sprintf("%s.%s = '%s', ", k.Nodes.Lebel, key, val)
					}
				}

				// append the id to each node
				cypher += fmt.Sprintf("%s._id = '%s' ", k.Nodes.Lebel, k.Nodes.ID)

				cypher += fmt.Sprintf("MERGE (%s)-[:%s]->(%s)", k.Edges.Source, relation, k.Edges.Target)

				cypher += " "
			}

		}
	}

	return cypher
}
