MATCH (b:Book)<-[:WROTE]-(a:Author{name:$author})
RETURN b.UUID, b.title
  ORDER BY b.title
  SKIP $skip
  LIMIT $limit