MATCH (bs:BookStatus)<-[:HAS_STATUS]-(b:Book)-[:HAS_TAG]->(t:Tag{name:$tag})
RETURN null as f1, null as f2, null as f3, b.UUID, b.title, bs.ID
  SKIP $skip
  LIMIT $limit