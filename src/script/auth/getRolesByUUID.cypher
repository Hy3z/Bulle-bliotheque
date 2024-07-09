MATCH (u:User{UUID:$uuid})-[:HAS_ROLE]->(r:Role)
return r.name
