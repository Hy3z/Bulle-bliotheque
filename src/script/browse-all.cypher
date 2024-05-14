MATCH (b:Book)
OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie)
RETURN distinct s.name AS sname, count(b) AS bcount,

                CASE
                  WHEN s IS null THEN b.ISBN_13
                  ELSE null
                  END AS bisbn,

                CASE
                  WHEN s IS null THEN b.title
                  ELSE null
                  END AS btitle