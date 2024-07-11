MATCH (b:Book{UUID:$buuid})-[r:HAS_STATUS]->(bs:BookStatus{ID:3})
MATCH (borrowStatus:BookStatus{ID:1})
CREATE (b)-[rp:HAS_STATUS{borrowerUUID:$uuuid}]->(borrowStatus)
delete r