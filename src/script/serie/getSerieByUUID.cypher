MATCH (s:Serie {UUID: $uuid})<-[r:PART_OF]-(b:Book)-[:HAS_STATUS]->(bs:BookStatus)
WHERE bs.ID <> 2
RETURN s.name, b.title, b.UUID, bs.ID
  ORDER BY r.opus ASC