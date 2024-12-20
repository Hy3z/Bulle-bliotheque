MATCH (bs:BookStatus)<-[:HAS_STATUS]-(b:Book)-[:HAS_TAG]->(t:Tag{name:$tag})
WHERE bs.ID <> 2
RETURN null as f1, null as f2, null as f3, b.UUID, b.title, bs.ID
ORDER BY b.title
  SKIP $skip
  LIMIT $limit