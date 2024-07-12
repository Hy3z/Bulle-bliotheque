MATCH(b:Book)-[r:HAS_STATUS]->(bs:BookStatus{ID:1})
WHERE r.borrowerUUID = $uuid
return b.UUID, b.title, bs.ID