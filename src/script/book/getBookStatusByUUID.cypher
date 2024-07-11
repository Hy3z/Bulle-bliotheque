MATCH (b:Book{UUID:$uuid})-[r:HAS_STATUS]->(bs:BookStatus)
RETURN bs.ID