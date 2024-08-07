MATCH (b:Book)-[:HAS_STATUS]->(bs:BookStatus)
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)
OPTIONAL MATCH (b)<-[:WROTE]-(a:Author)  
OPTIONAL MATCH (b)-[:HAS_TAG]->(t:Tag)
WITH *,(
       $titleCoeff * apoc.text.sorensenDiceSimilarity(b.title, $expr) +
       $serieCoeff * CASE WHEN s IS NOT NULL THEN apoc.text.sorensenDiceSimilarity(s.name, $expr) ELSE 0 END +
       $authorCoeff * CASE WHEN a IS NOT NULL THEN apoc.text.sorensenDiceSimilarity(a.name, $expr) ELSE 0 END +
       $tagCoeff * CASE WHEN t IS NOT NULL THEN apoc.text.sorensenDiceSimilarity(t.name, $expr) ELSE 0 END
       ) AS rank
  WHERE rank > $minRank and bs.ID <> 2
RETURN null as f1, null as f2, null as f3, max(rank), b.UUID, b.title, bs.ID
  ORDER BY max(rank) DESC, b.title
  SKIP $skip
  LIMIT $limit