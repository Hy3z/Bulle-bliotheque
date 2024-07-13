MATCH(u:User{UUID:$uuid})-[:HAS_LIKED]->(b:Book)-[:HAS_STATUS]->(bs:BookStatus)
RETURN b.UUID, b.title, bs.ID