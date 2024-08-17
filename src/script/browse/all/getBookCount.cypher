MATCH (b:Book)-[:HAS_STATUS]->(bs:BookStatus)
WHERE bs.ID <> 2
return count(b)