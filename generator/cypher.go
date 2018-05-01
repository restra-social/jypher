package generator

import (
	"fmt"
	"github.com/restra-social/jypher/models"
	"regexp"
	"strings"
)

// CypherGenerator : Cypher Query Generator
type CypherGenerator struct{}

// Generate : This method takes a graph model and generates cypher query
func (c *CypherGenerator) Generate(id string, models map[string]models.Graph, serial []string) []string {

	var queries []string

	// loop through the serial
	for _, term := range serial {

		// search for key in the model in ascending order to generate the query
		if k, ok := models[term]; ok {

			level := k.Nodes.Lebel
			nodeRelName := regexp.MustCompile(`[A-Za-z]+`).FindString(level)

			// Filter Special Label like type which is similar to many resource but has different meaning
			// like Organization Type is different than Claim Type

			if strings.HasPrefix(level, "type") {
				// append source
				level = fmt.Sprintf("%s%s", k.Edges.Source, level)
			}

			pl := len(k.Nodes.Properties)

			node := regexp.MustCompile(`[A-Za-z]+`).FindString(strings.Title(level))
			source := regexp.MustCompile(`[A-Za-z]+`).FindString(k.Edges.Source.Name)

			relation := fmt.Sprintf("%s_%s", strings.ToUpper(source), strings.ToUpper(nodeRelName))

			if k.Nodes.ID != "" {
				cypher := fmt.Sprintf("MERGE (%s:%s {id:'%s'}) SET ", level, node, k.Nodes.ID)

				for i, property := range k.Nodes.Properties {
					for key, val := range property {

						switch val.(type) {
						case string:
							// Using ' ' in value assignment so filter for text contains ''
							filteredVal := strings.Replace(val.(string), "'", "", -1)
							cypher += fmt.Sprintf("%s.%s = '%s'", k.Nodes.Lebel, key, filteredVal)
						case float64:
							cypher += fmt.Sprintf("%s.%s = %d", k.Nodes.Lebel, key, int(val.(float64)))
						case bool:
							cypher += fmt.Sprintf("%s.%s = %v", k.Nodes.Lebel, key, val.(bool))
						}

					}
					if pl > 1 {
						if i < pl-1 {
							cypher += ", "
						}
					}

				}
				if k.Edges.Source.Name != k.Edges.Target { // avoids self loop
					cypher += fmt.Sprintf(" WITH %s MATCH (%s:%s {id : '%s'}) WITH %s, %s ", level, k.Edges.Source.Name, strings.Title(source), k.Edges.Source.ID, level, k.Edges.Source.Name)
					// Add Relation
					cypher += fmt.Sprintf("MERGE (%s)-[:%s]->(%s)", k.Edges.Source.Name, relation, level)

				}

				queries = append(queries, cypher)

			} else {

				node := regexp.MustCompile(`[A-Za-z]+`).FindString(strings.Title(level))

				len := len(k.Nodes.Properties)

				// If property found then take them for full merge
				if len > 0 {

					cypher := fmt.Sprintf("MERGE (%s:%s { ", level, node)

					for i, property := range k.Nodes.Properties {
						for key, val := range property {
							// Using ' ' in value assignment so filter for text contains ''
							filteredVal := strings.Replace(val.(string), "'", "", -1)
							cypher += fmt.Sprintf("%s:'%s'", key, filteredVal)

							switch val.(type) {
							case string:
								// Using ' ' in value assignment so filter for text contains ''
								filteredVal := strings.Replace(val.(string), "'", "", -1)
								cypher += fmt.Sprintf("%s: '%s'", key, filteredVal)
							case float64:
								cypher += fmt.Sprintf("%s: %d", key, int(val.(float64)))
							case bool:
								cypher += fmt.Sprintf("%s: %v", key, val.(bool))
							}

							// skip comma for last property
							if i != len-1 {
								cypher += fmt.Sprint(",")
							}
						}

						if i == len-1 {
							cypher += fmt.Sprint(" }) ")
						}
					}

					// append the id to each node
					cypher += fmt.Sprintf(" WITH %s MATCH (%s:%s {id : '%s'}) WITH %s, %s ", level, k.Edges.Source.Name, strings.Title(k.Edges.Source.Name), k.Edges.Source.ID, level, k.Edges.Source.Name)

					// Add Relation
					cypher += fmt.Sprintf("MERGE (%s)-[:%s]->(%s)", k.Edges.Source.Name, relation, level)

					queries = append(queries, cypher)

				} else {

					// for those nodes who doesn't have any properties
					cypher := fmt.Sprintf("MERGE (%s:%s)", level, node)

					/*cypher += fmt.Sprintf("MERGE (%s:%s) SET ", level, node)

					for _, property := range k.Nodes.Properties {
						for key, val := range property {
							cypher += fmt.Sprintf("%s.%s = '%s', ", k.Nodes.Lebel, key, val)
						}
					}

					// append the id to each node
					cypher += fmt.Sprintf("%s._id = '%s' ", k.Nodes.Lebel, id)*/

					cypher += fmt.Sprintf("MERGE (%s)-[:%s]->(%s)", k.Edges.Source.Name, relation, k.Edges.Target)

					queries = append(queries, cypher)
				}
			}

		}
	}

	return queries
}

// #todo #fix
// fix MERGE (patient:Patient {id:'34876259-35cd-497c-a932-94baaaeb555c'}) SET patient.reference = 'urn:uuid:34876259-35cd-497c-a932-94baaaeb555c'
// fix
