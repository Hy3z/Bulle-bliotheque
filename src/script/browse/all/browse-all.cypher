MATCH (b:Book)
RETURN null as col1, null as col2, 1, b.UUID, b.title
  SKIP $skip LIMIT $limit