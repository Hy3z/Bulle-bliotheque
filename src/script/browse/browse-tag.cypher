MATCH (b:Book)-[:HAS_TAG]->(t:Tag{name:$tag})
RETURN b.UUID, b.title
  SKIP $skip
  LIMIT $limit