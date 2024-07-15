MATCH (u:User{UUID:$uuid})
MATCH (b:Book{UUID:$buuid})
OPTIONAL MATCH (u)-[r:REVIEWED]->(b)
DELETE r
CREATE (u)-[:REVIEWED{date:$date, message:$message}]->(b)