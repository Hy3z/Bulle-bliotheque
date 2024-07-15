MATCH (u:User{UUID:$uuuid})-[r:REVIEWED]->(b:Book{UUID:$buuid})
DELETE r