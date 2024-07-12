MATCH (b:Book)-[:HAS_STATUS]->(bs:BookStatus)
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)
  WHERE bs.ID <> 2
RETURN distinct s.name, s.UUID, count(b),
                CASE WHEN s IS null THEN b.UUID ELSE null END AS uuid,
                CASE WHEN s IS null THEN b.title ELSE null END AS title, bs.ID
  SKIP $skip LIMIT $limit