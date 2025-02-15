  MATCH (b:Book)-[:HAS_STATUS]->(bs:BookStatus)
  WHERE b.date IS NOT NULL and bs.ID <> 2
RETURN b.UUID, b.title, bs.ID
  ORDER BY b.date DESC, b.title
  LIMIT $limit
