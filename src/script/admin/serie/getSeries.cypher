MATCH (s:Serie)
OPTIONAL MATCH (s)<-[:PART_OF]-(b:Book)
RETURN s.UUID, s.name, count(b)
ORDER BY s.name ASC
