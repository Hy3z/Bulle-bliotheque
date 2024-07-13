MATCH (b:Book{UUID:$buuid})
MATCH (u:User{UUID:$uuuid})
WHERE NOT (u)-[:HAS_LIKED]->(b)
CREATE (u)-[:HAS_LIKED]->(b)
