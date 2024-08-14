MATCH (b:Book)-[:HAS_TAG]->(t:Tag{name:$name})
RETURN count(b)