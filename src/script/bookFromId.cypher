MATCH (a:Author)-[:WROTE]->(b:Book)-[:HAS_TAG]->(t:Tag)
  WHERE elementId(b) = $id
RETURN
  b.title, b.cover, b.summary, b.date, b.language , collect(distinct(a.name)) as authors, collect(distinct(t.name)) as tags
    LIMIT 1




MATCH (a:Author)-[:WROTE]->(b:Book)-[:HAS_TAG]->(t:Tag)
  WHERE elementId(b) = "4:3b1486ef-c88b-4a6d-9454-af2e8705a336:15"
RETURN
  b.title, b.cover, b.summary, b.date, b.language , collect(distinct(a.name)) as authors, collect(distinct(t.name)) as tags
  LIMIT 1