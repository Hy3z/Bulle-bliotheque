MATCH (b:Book)-[:HAS_STATUS]->(bs:BookStatus)
  WHERE b.date IS NOT NULL
RETURN b.UUID, b.title, bs.ID
  ORDER BY b.date DESC
  LIMIT $limit