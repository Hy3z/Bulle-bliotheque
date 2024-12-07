MATCH (b:Book)-[:HAS_STATUS]->(bs:BookStatus)
  WHERE bs.ID <> 2
RETURN null as col1,
       null as col2,
       1,
       b.UUID,
       b.title,
       bs.ID
ORDER BY b.title
  SKIP $skip LIMIT $limit