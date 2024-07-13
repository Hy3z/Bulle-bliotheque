MATCH (u:User{UUID:$uuuid})-[r:HAS_LIKED]->(b:Book{UUID:$buuid})
  DELETE r
