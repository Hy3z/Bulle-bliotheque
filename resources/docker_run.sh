#!/bin/bash

exec docker run \
    --restart always \
    --publish=7474:7474 --publish=7687:7687 \
    --env NEO4J_AUTH=neo4j/$1 \
    --name bbdb-local \
    -e NEO4J_PLUGINS=\[\"apoc\"\] \
    --volume=$PWD/data:/data \
    --volume=$PWD/logs:/logs \
    --volume=$PWD/plugins:/plugins \
    --volume=$PWD/conf:/conf \
    -e NEO4J_apoc_export_file_enabled=true \
    -e NEO4J_apoc_import_file_enabled=true \
    -e NEO4J_apoc_import_file_use__neo4j__config=true \
    neo4j:5.20-community
