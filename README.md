### JSON to CYPHER QUERY Generator (Jypher)

[![Go Report Card](https://goreportcard.com/badge/github.com/kite-social/jypher)](https://goreportcard.com/report/github.com/kite-social/jypher)


# Important Components

Rules are what used to avoid creating redundant nodes . E.g. Encounter has object field called `subject` , since every object makes a node in the graph
So a node called `Subject` will be created without the Rules . But we specified rename `subject` with `patient` since we know that this object
refers to a patient object so instead of `Subject` , `Patient` node will be created or merged with existing Node.
```golang
	rules := map[string]map[string]interface{}{
		"Encounter": {
			"subject":         "patient",
			"serviceProvider": "organization",
		},
		"Patient": {
			"generalPractitioner": "practitioner",
		},
	}
```

Sample Usages

```golang

func main() {

	data := []byte(`
	{
        "id": "49c2df26-6fb9-43fb-82a2-3bd413a95934",
        "meta": {
          "profile": [
            "http://standardhealthrecord.org/fhir/StructureDefinition/shr-demographics-PersonOfRecord"
          ]
        },
        "identifier": [
          {
            "type": {
              "coding": [
                {
                  "system": "http://hl7.org/fhir/identifier-type",
                  "code": "SB"
                }
              ]
            },
            "system": "http://hl7.org/fhir/sid/us-ssn",
            "value": "999608294"
          }
        ],
        "name": [
          {
            "use": "official",
            "family": "Abernathy691",
            "given": [
              "Carlotta591"
            ]
          }
        ],
        "telecom": [
          {
            "system": "phone",
            "value": "1-892-670-6907 x741",
            "use": "home"
          }
        ],
        "gender": "female",
        "birthDate": "2004-09-12",
        "address": [
        {
            "line": [
              "775 Jerde Pike"
            ],
            "city": "Sunderland",
            "state": "MA",
            "postalCode": "01375",
            "country": "US"
          }
        ],
        "maritalStatus": {
          "coding": [
            {
              "system": "http://hl7.org/fhir/v3/MaritalStatus",
              "code": "S"
            }
          ],
          "text": "Never Married"
        },
        "multipleBirthBoolean": false,
        "communication": [
          {
            "language": {
              "coding": [
                {
                  "system": "http://hl7.org/fhir/ValueSet/languages",
                  "code": "en-US",
                  "display": "English (United States)"
                }
              ]
            }
          }
        ],
        "generalPractitioner": [
          {
            "reference": "urn:uuid:1a3711a3-1d38-4aeb-b25f-705ea922b1c2"
          }
        ],
        "resourceType": "Patient"
      }
	`)

	/*var graphs []Graph

	master := createMasterNode(&unmarshal)
	graphs = append(graphs, master)*/

	rules := map[string]map[string]interface{}{
		"Patient": {
			"generalPractitioner": "practitioner",
		},
	}

	var unmarshal map[string]interface{}

	err := json.Unmarshal(data, &unmarshal)

	if err != nil {
		panic(err)
	}

	resource := unmarshal["resourceType"].(string)

	jsonInfo := models.JSONInfo{
		DecodedJSON: unmarshal,
		Rules:       getRules(resource, rules),
		Master:      strings.ToLower(resource),
		ID:          unmarshal["id"].(string),
	}

	j := J.Jypher{}
	graph := j.GetJypher(&jsonInfo)

	data, _ = json.MarshalIndent(graph, " ", " ")

	fmt.Println(string(data))

	cypher := j.BuildCypher()

	fmt.Println(cypher)
}

func getRules(resource string, rules map[string]map[string]interface{}) map[string]interface{} {
	if rule, ok := rules[resource]; ok {
		return rule
	} else {
		panic(fmt.Sprintf("Rules provided but rules for [%s] not found, check the rules dictionary", resource))
	}
	return nil
}

```

Decoded Graph model out of the JSON

```json
{
   "address0":{
      "nodes":{
         "lebel":"address0",
         "properties":[
            {
               "city":"Sunderland"
            },
            {
               "state":"MA"
            },
            {
               "postalCode":"01375"
            },
            {
               "country":"US"
            }
         ]
      },
      "edges":{
         "source":"patient",
         "target":"address0"
      }
   },
   "coding0":{
      "nodes":{
         "lebel":"coding0",
         "properties":[
            {
               "system":"http://hl7.org/fhir/identifier-type"
            },
            {
               "code":"SB"
            }
         ]
      },
      "edges":{
         "source":"type",
         "target":"coding0"
      }
   },
   "coding1":{
      "nodes":{
         "lebel":"coding1",
         "properties":[
            {
               "system":"http://hl7.org/fhir/v3/MaritalStatus"
            },
            {
               "code":"S"
            }
         ]
      },
      "edges":{
         "source":"maritalStatus",
         "target":"coding1"
      }
   },
   "coding4":{
      "nodes":{
         "lebel":"coding4",
         "properties":[
            {
               "system":"http://hl7.org/fhir/ValueSet/languages"
            },
            {
               "code":"en-US"
            },
            {
               "display":"English (United States)"
            }
         ]
      },
      "edges":{
         "source":"language",
         "target":"coding4"
      }
   },
   "communication4":{
      "nodes":{
         "lebel":"communication4"
      },
      "edges":{
         "source":"patient",
         "target":"communication4"
      }
   },
   "identifier0":{
      "nodes":{
         "lebel":"identifier0",
         "properties":[
            {
               "system":"http://hl7.org/fhir/sid/us-ssn"
            },
            {
               "value":"999608294"
            }
         ]
      },
      "edges":{
         "source":"patient",
         "target":"identifier0"
      }
   },
   "language":{
      "nodes":{
         "lebel":"language"
      },
      "edges":{
         "source":"communication4",
         "target":"language"
      }
   },
   "maritalStatus":{
      "nodes":{
         "lebel":"maritalStatus",
         "properties":[
            {
               "text":"Never Married"
            }
         ]
      },
      "edges":{
         "source":"patient",
         "target":"maritalStatus"
      }
   },
   "meta":{
      "nodes":{
         "lebel":"meta"
      },
      "edges":{
         "source":"patient",
         "target":"meta"
      }
   },
   "name1":{
      "nodes":{
         "lebel":"name1",
         "properties":[
            {
               "use":"official"
            },
            {
               "family":"Abernathy691"
            }
         ]
      },
      "edges":{
         "source":"patient",
         "target":"name1"
      }
   },
   "patient":{
      "nodes":{
         "id":"49c2df26-6fb9-43fb-82a2-3bd413a95934",
         "lebel":"patient",
         "properties":[
            {
               "resourceType":"Patient"
            },
            {
               "id":"49c2df26-6fb9-43fb-82a2-3bd413a95934"
            },
            {
               "gender":"female"
            },
            {
               "birthDate":"2004-09-12"
            }
         ]
      },
      "edges":{
         "source":"patient",
         "target":"patient"
      }
   },
   "practitioner3":{
      "nodes":{
         "id":"urn:uuid:1a3711a3-1d38-4aeb-b25f-705ea922b1c2",
         "lebel":"practitioner3",
         "properties":[
            {
               "reference":"urn:uuid:1a3711a3-1d38-4aeb-b25f-705ea922b1c2"
            }
         ]
      },
      "edges":{
         "source":"patient",
         "target":"practitioner3"
      }
   },
   "telecom2":{
      "nodes":{
         "lebel":"telecom2",
         "properties":[
            {
               "value":"1-892-670-6907 x741"
            },
            {
               "use":"home"
            },
            {
               "system":"phone"
            }
         ]
      },
      "edges":{
         "source":"patient",
         "target":"telecom2"
      }
   },
   "type":{
      "nodes":{
         "lebel":"type"
      },
      "edges":{
         "source":"identifier0",
         "target":"type"
      }
   }
}
```

Generated CYPHER QUERY

```sql

MERGE (patient:Patient {id:'49c2df26-6fb9-43fb-82a2-3bd413a95934'}) ON CREATE SET patient.resourceType = 'Patient', patient.id = '49c2df26-6fb9-43fb-82a2-3bd413a95934', patient.gender = 'female', patient.birthDate = '2004-09-12'
CREATE (identifier0:Identifier) SET identifier0.system = 'http://hl7.org/fhir/sid/us-ssn', identifier0.value = '999608294', identifier0._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (patient)-[:PATIENT_IDENTIFIER]->(identifier0)
CREATE (type:Type) SET type._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (identifier0)-[:IDENTIFIER0_TYPE]->(type)
CREATE (coding0:Coding) SET coding0.system = 'http://hl7.org/fhir/identifier-type', coding0.code = 'SB', coding0._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (type)-[:TYPE_CODING]->(coding0)
CREATE (name1:Name) SET name1.use = 'official', name1.family = 'Abernathy691', name1._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (patient)-[:PATIENT_NAME]->(name1)
CREATE (telecom2:Telecom) SET telecom2.value = '1-892-670-6907 x741', telecom2.use = 'home', telecom2.system = 'phone', telecom2._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (patient)-[:PATIENT_TELECOM]->(telecom2)
MERGE (practitioner3:Practitioner {id:'urn:uuid:1a3711a3-1d38-4aeb-b25f-705ea922b1c2'}) ON CREATE SET practitioner3.reference = 'urn:uuid:1a3711a3-1d38-4aeb-b25f-705ea922b1c2'
MERGE (patient)-[:PATIENT_PRACTITIONER]->(practitioner3)
CREATE (communication4:Communication) SET communication4._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (patient)-[:PATIENT_COMMUNICATION]->(communication4)
CREATE (language:Language) SET language._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (communication4)-[:COMMUNICATION4_LANGUAGE]->(language)
CREATE (coding4:Coding) SET coding4.system = 'http://hl7.org/fhir/ValueSet/languages', coding4.code = 'en-US', coding4.display = 'English (United States)', coding4._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (language)-[:LANGUAGE_CODING]->(coding4)
CREATE (meta:Meta) SET meta._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (patient)-[:PATIENT_META]->(meta) CREATE (address0:Address) SET address0.city = 'Sunderland', address0.state = 'MA', address0.postalCode = '01375', address0.country = 'US', address0._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (patient)-[:PATIENT_ADDRESS]->(address0)
CREATE (maritalStatus:MaritalStatus) SET maritalStatus.text = 'Never Married', maritalStatus._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (patient)-[:PATIENT_MARITALSTATUS]->(maritalStatus) CREATE (coding1:Coding) SET coding1.system = 'http://hl7.org/fhir/v3/MaritalStatus', coding1.code = 'S', coding1._id = '49c2df26-6fb9-43fb-82a2-3bd413a95934'
MERGE (maritalStatus)-[:MARITALSTATUS_CODING]->(coding1)

```

Sample Output Graph

![Graph](https://cdn.pbrd.co/images/GGNOXyt.png "Neo4j Graph")

#### Todo

* Create separate Node with Array of String
* Serialization of Query Generation for better Execution of Cypher
* Remove id field from properties of master node


#### Limitation

* multiple identifier in a array list with same type of field name , right now it skips and takes only one
