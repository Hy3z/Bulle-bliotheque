MATCH (b:Book)-[:HAS_STATUS]->(bs:BookStatus)
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)<-[:PART_OF]-(ob:Book)
OPTIONAL MATCH (b)<-[:WROTE]-(a:Author)
OPTIONAL MATCH (b)-[:HAS_TAG]->(t:Tag)
WITH *,(
       $titleCoeff * apoc.text.sorensenDiceSimilarity(b.title, $expr) +
       $serieCoeff * CASE WHEN s IS NOT NULL THEN apoc.text.sorensenDiceSimilarity(s.name, $expr) ELSE 0 END +
       $authorCoeff * CASE WHEN a IS NOT NULL THEN apoc.text.sorensenDiceSimilarity(a.name, $expr) ELSE 0 END +
       $tagCoeff * CASE WHEN t IS NOT NULL THEN apoc.text.sorensenDiceSimilarity(t.name, $expr) ELSE 0 END
       ) AS rank
  WHERE rank > $minRank and bs.ID <> 2
RETURN distinct s.name, s.UUID, count(ob), max(rank),
                CASE WHEN s IS null THEN b.UUID ELSE null END AS uuid,
                CASE WHEN s IS null THEN b.title ELSE null END AS title, bs.ID
  ORDER BY max(rank) DESC
  SKIP $skip
  LIMIT $limit
