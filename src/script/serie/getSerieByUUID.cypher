MATCH (s:Serie {UUID: $uuid})
OPTIONAL MATCH (b:Book)-[:PART_OF]->(s)

OPTIONAL MATCH (u:User)-[:HAS_LIKED]->(b)
OPTIONAL MATCH (b)-[:HAS_TAG]->(t:Tag)
OPTIONAL MATCH (a:Author)-[:WROTE]->(b)

RETURN s.name, count(distinct u), collect(distinct t.name), collect(distinct a.name)
