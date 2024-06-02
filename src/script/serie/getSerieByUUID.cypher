MATCH (s:Serie {UUID: $uuid})<-[r:PART_OF]-(b:Book)
RETURN s.name, b.title, b.UUID
  ORDER BY r.opus ASC