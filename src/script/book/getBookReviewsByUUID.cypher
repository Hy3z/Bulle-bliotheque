MATCH (u:User)-[r:REVIEWED]->(b:Book{UUID:$uuid})
RETURN u.UUID, u.name, r.message, apoc.temporal.format(r.date, 'dd/MM/yyyy (HH:mm)')
ORDER BY r.date ASC