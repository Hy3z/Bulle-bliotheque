MATCH (b:Book)-[:HAS_STATUS]->(bs:BookStatus)
WHERE bs.ID <> 2
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)
RETURN distinct s.name, s.UUID, count(b),
                CASE WHEN s IS null THEN b.UUID ELSE null END AS uuid,
                CASE WHEN s IS null THEN b.title ELSE null END AS title, bs.ID
  SKIP $skip LIMIT $limit
