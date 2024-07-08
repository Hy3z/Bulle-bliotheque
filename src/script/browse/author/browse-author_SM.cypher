MATCH (bs:BookStatus)<-[:HAS_STATUS]-(b:Book)<-[:WROTE]-(a:Author{name:$author})
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)<-[:PART_OF]-(ob:Book)
RETURN distinct s.name, s.UUID, count(ob),
                CASE WHEN s IS null THEN b.UUID ELSE null END AS uuid,
                CASE WHEN s IS null THEN b.title ELSE null END AS title, bs.ID
  ORDER BY s.name
  SKIP $skip
  LIMIT $limit