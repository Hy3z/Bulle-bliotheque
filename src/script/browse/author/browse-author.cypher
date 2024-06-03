MATCH (b:Book)<-[:WROTE]-(a:Author{name:$author})
RETURN null as f1, null as f2, null as f3, b.UUID, b.title
  ORDER BY b.title
  SKIP $skip
  LIMIT $limit