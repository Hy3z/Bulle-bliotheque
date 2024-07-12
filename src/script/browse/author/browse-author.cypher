MATCH (bs:BookStatus)<-[:HAS_STATUS]-(b:Book)<-[:WROTE]-(a:Author{name:$author})
  WHERE bs.ID <> 2
RETURN null as f1, null as f2, null as f3, b.UUID, b.title, bs.ID
  ORDER BY b.title
  SKIP $skip
  LIMIT $limit