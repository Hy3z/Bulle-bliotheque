MATCH (u:User{UUID:$uuuid})-[:HAS_LIKED]->(b:Book{UUID:$buuid})
  DELETE r
