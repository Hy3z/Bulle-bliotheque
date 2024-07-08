MATCH (b:Book)-[:HAS_STATUS]->(bs:BookStatus)
RETURN null as col1, null as col2, 1, b.UUID, b.title, bs.ID
  SKIP $skip LIMIT $limit