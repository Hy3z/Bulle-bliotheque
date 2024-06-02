MATCH (b:Book)-[:HAS_TAG]->(t:Tag{name:$tag})
RETURN null as f1, null as f2, null as f3, b.UUID, b.title
  SKIP $skip
  LIMIT $limit