MATCH (bs:BookStatus)<-[:HAS_STATUS]-(b:Book)-[:HAS_TAG]->(t:Tag{name:$tag})
WHERE bs.ID <> 2
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)<-[:PART_OF]-(ob:Book)

RETURN distinct s.name,
                s.UUID,
                count(ob),
                CASE WHEN s IS null THEN b.UUID ELSE null END AS uuid,
                CASE WHEN s IS null THEN b.title ELSE null END AS title,
                bs.ID,
                CASE WHEN s IS null then b.title ELSE s.name END AS name
ORDER BY name
  SKIP $skip LIMIT $limit
