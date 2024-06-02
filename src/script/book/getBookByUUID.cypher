MATCH (b:Book {UUID: $uuid})
OPTIONAL MATCH (a:Author)-[:WROTE]->(b)
OPTIONAL MATCH (b)-[:HAS_TAG]->(t:Tag)
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie) 
RETURN b.title, b.UUID, b.description, b.publishedDate, b.publisher, b.cote, b.pageCount, collect(distinct(a.name)) as authors, collect(distinct(t.name)) as tags, s.name, s.UUID
  LIMIT 1