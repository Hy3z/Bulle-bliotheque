MATCH (bs:BookStatus)<-[:HAS_STATUS]-(b:Book)<-[:WROTE]-(a:Author{name:$author})
RETURN null as f1, null as f2, null as f3, b.UUID, b.title, bs.ID
  ORDER BY b.title
  SKIP $skip
  LIMIT $limit