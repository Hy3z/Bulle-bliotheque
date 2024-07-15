MATCH (u:User{UUID:$uuid})-[:REVIEWED]->(b:Book)-[r:HAS_STATUS]->(bs:BookStatus)
return b.UUID, b.title, bs.ID