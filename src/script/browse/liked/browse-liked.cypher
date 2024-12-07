MATCH (u:User)-[:HAS_LIKED]->(b:Book)-[:HAS_STATUS]->(bs:BookStatus)
  WHERE bs.ID <> 2
RETURN null as col1,
       null as col2,
       1,
       b.UUID,
       b.title,
       bs.ID,
       count(u) as likes
ORDER BY likes DESC, b.title
  SKIP $skip LIMIT $limit
