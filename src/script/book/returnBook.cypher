MATCH (b:Book{UUID:$buuid})-[r:HAS_STATUS{borrowerUUID:$uuuid}]->(bs:BookStatus{ID:1})
MATCH (availableStatus:BookStatus{ID:3})
CREATE (b)-[rp:HAS_STATUS]->(availableStatus)
delete r