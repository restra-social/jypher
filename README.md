### JSON to CYPHER QUERY Generator (Jypher)

[![Go Report Card](https://goreportcard.com/badge/github.com/kite-social/jypher)](https://goreportcard.com/report/github.com/kite-social/jypher)

Sample JSON INPUT

```json
{
   "id":"123456",
   "session":{
      "title":"12 Years of Spring: An Open Source Journey",
      "abstract":"Spring emerged as a core open source project in early 2003 and evolved to a broad portfolio of open source projects up until 2015.",
      "conference":{
         "city":"London"
      }
   },
   "topics":[
      "keynote",
      "spring"
   ],
   "tracks":[
      {
         "main":"Java"
      },
      {
         "second":"Languages"
      },
      {
         "third":"Golang"
      }
   ],
   "room":"Auditorium",
   "timeslot":"Wed 29th, 09:30-10:30",
   "speaker":{
      "name":"Juergen Hoeller",
      "bio":"Juergen Hoeller is co-founder of the Spring Framework open source project.",
      "twitter":"https://twitter.com/springjuergen",
      "picture":"http://www.springio.net/wp-content/uploads/2014/11/juergen_hoeller-220x220.jpeg"
   }
}
```

Decoded model of the JSON

```json
{
   "conference":{
      "nodes":{
         "id":"123456",
         "lebel":"conference",
         "properties":[
            {
               "city":"London"
            }
         ]
      },
      "edges":{
         "source":"session",
         "target":"conference"
      }
   },
   "session":{
      "nodes":{
         "id":"123456",
         "lebel":"session",
         "properties":[
            {
               "title":"12 Years of Spring: An Open Source Journey"
            },
            {
               "abstract":"Spring emerged as a core open source project in early 2003 and evolved to a broad portfolio of open source projects up until 2015."
            }
         ]
      },
      "edges":{
         "source":"talk",
         "target":"session"
      }
   },
   "speaker":{
      "nodes":{
         "id":"123456",
         "lebel":"speaker",
         "properties":[
            {
               "name":"Juergen Hoeller"
            },
            {
               "bio":"Juergen Hoeller is co-founder of the Spring Framework open source project."
            },
            {
               "twitter":"https://twitter.com/springjuergen"
            },
            {
               "picture":"http://www.springio.net/wp-content/uploads/2014/11/juergen_hoeller-220x220.jpeg"
            }
         ]
      },
      "edges":{
         "source":"talk",
         "target":"speaker"
      }
   },
   "talk":{
      "nodes":{
         "id":"123456",
         "lebel":"talk",
         "properties":[
            {
               "room":"Auditorium"
            },
            {
               "timeslot":"Wed 29th, 09:30-10:30"
            },
            {
               "id":"123456"
            }
         ]
      },
      "edges":{
         "source":"talk",
         "target":"talk"
      }
   },
   "tracks0":{
      "nodes":{
         "id":"123456",
         "lebel":"tracks0",
         "properties":[
            {
               "main":"Java"
            }
         ]
      },
      "edges":{
         "source":"talk",
         "target":"tracks0"
      }
   },
   "tracks1":{
      "nodes":{
         "id":"123456",
         "lebel":"tracks1",
         "properties":[
            {
               "second":"Languages"
            }
         ]
      },
      "edges":{
         "source":"talk",
         "target":"tracks1"
      }
   },
   "tracks2":{
      "nodes":{
         "id":"123456",
         "lebel":"tracks2",
         "properties":[
            {
               "third":"Golang"
            }
         ]
      },
      "edges":{
         "source":"talk",
         "target":"tracks2"
      }
   }
}
```

Generated CYPHER QUERY

```sql

MERGE (talk:Talk {id:'123456'}) ON CREATE SET talk.room = 'Auditorium', talk.timeslot = 'Wed 29th, 09:30-10:30', talk.id = '123456'
CREATE (session:Session) SET session.title = '12 Years of Spring: An Open Source Journey', session.abstract = 'Spring emerged as a core open source project in early 2003 and evolved to a broad portfolio of open source projects up until 2015.', session._id = '123456' MERGE (talk)-[:TALK_SESSION]->(session)
CREATE (conference:Conference) SET conference.city = 'London', conference._id = '123456' MERGE (session)-[:SESSION_CONFERENCE]->(conference)
CREATE (tracks0:Tracks) SET tracks0.main = 'Java', tracks0._id = '123456' MERGE (talk)-[:TALK_TRACKS]->(tracks0)
CREATE (tracks1:Tracks) SET tracks1.second = 'Languages', tracks1._id = '123456' MERGE (talk)-[:TALK_TRACKS]->(tracks1)
CREATE (tracks2:Tracks) SET tracks2.third = 'Golang', tracks2._id = '123456' MERGE (talk)-[:TALK_TRACKS]->(tracks2)
CREATE (speaker:Speaker) SET speaker.name = 'Juergen Hoeller', speaker.bio = 'Juergen Hoeller is co-founder of the Spring Framework open source project.', speaker.twitter = 'https://twitter.com/springjuergen', speaker.picture = 'http://www.springio.net/wp-content/uploads/2014/11/juergen_hoeller-220x220.jpeg', speaker._id = '123456' MERGE (talk)-[:TALK_SPEAKER]->(speaker)
```

Sample Output Graph

![Graph](https://cdn.pbrd.co/images/GGrvUTv.png "Neo4j Graph")

#### Todo

* Create separate Node with Array of String
* Serialization of Query Generation for better Execution of Cypher