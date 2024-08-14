MATCH (u:User{UUID:$uuuid})
MATCH (b:Book{UUID:$buuid})
OPTIONAL MATCH (u)-[r:REVIEWED]->(b)
DELETE r
CREATE (u)-[:REVIEWED{date:datetime({timezone: 'Europe/Paris'}), message:$message}]->(b)