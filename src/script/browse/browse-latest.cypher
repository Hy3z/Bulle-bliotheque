MATCH (b:Book)
  WHERE b.date IS NOT NULL
RETURN b.UUID, b.title
  ORDER BY b.date DESC
  LIMIT $limit