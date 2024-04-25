MATCH (b:Book)
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)
OPTIONAL MATCH (b)<-[:WROTE]-(a:Author)
OPTIONAL MATCH (b)-[:HAS_TAG]->(t:Tag)
UNWIND split(b.title) AS btfield

WITH *,(
       apoc.text.sorensenDiceSimilarity(b.title, $expr) +
       s IS NOT NULL * apoc.text.fuzzyMatch(s.name, $expr) +
       a IS NOT NULL * apoc.text.fuzzyMatch(a.name, $expr) +
       t IS NOT NULL * apoc.text.fuzzyMatch(t.name, $expr)
       ) AS rank
  WHERE rank > 0
RETURN b.ISBN_13, b.title, max(rank)
  ORDER BY max(rank) DESC, b.title
  SKIP $skip LIMIT $limit




MATCH (b:Book)
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)
OPTIONAL MATCH (b)<-[:WROTE]-(a:Author)
OPTIONAL MATCH (b)-[:HAS_TAG]->(t:Tag)
UNWIND split(b.title) AS btfield
WITH *,(
      CASE WHEN s IS NOT NULL THEN 1 ELSE 2

       toInteger(b.title =~ $regex) +

       toInteger(s IS NOT NULL AND s.name =~ $regex) +
       toInteger(a IS NOT NULL AND a.name =~ $regex) +
       toInteger(t IS NOT NULL AND t.name =~ $regex)
       ) AS rank
  WHERE rank > 0
RETURN b.ISBN_13, b.title, max(rank)
  ORDER BY max(rank) DESC, b.title
  SKIP $skip LIMIT $limit





