MATCH (b:Book)
  WHERE b.title =~ '.*$filter*.'

RETURN elementId(b), b.title, b.cover SKIP $skip LIMIT $limit