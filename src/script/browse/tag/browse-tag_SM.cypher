MATCH (b:Book)-[:HAS_TAG]->(t:Tag{name:$tag})
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)<-[:PART_OF]-(ob:Book)
RETURN distinct s.name, s.UUID, count(ob),
                CASE WHEN s IS null THEN b.UUID ELSE null END AS uuid,
                CASE WHEN s IS null THEN b.title ELSE null END AS title
  SKIP $skip LIMIT $limit