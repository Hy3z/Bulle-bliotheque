MATCH (b:Book {UUID: $buuid})-[r:HAS_STATUS]->(bs:BookStatus)
OPTIONAL MATCH (a:Author)-[:WROTE]->(b)
OPTIONAL MATCH (b)-[:HAS_TAG]->(t:Tag)
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)
OPTIONAL MATCH (u:User{UUID:$uuuid})-[:HAS_LIKED]->(b)
OPTIONAL MATCH (other:User)-[:HAS_LIKED]->(b)
RETURN b.title, b.UUID, b.description, b.publishedDate, b.publisher, b.cote, b.pageCount, collect(distinct(a.name)) as authors, collect(distinct(t.name)) as tags, s.name, s.UUID, bs.ID, r.borrowerUUID, u IS NOT NULL, count(distinct(other))
  LIMIT 1
