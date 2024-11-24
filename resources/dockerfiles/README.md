On range ici les dockerfiles utilisées pour créer les images docker de Keycloak et Neo4j pour le développement en local.
Le dossier  ``keycloak`` contient le dockerfile pour l'image Keycloak ainsi que les fichiers de configuration pour recréer un realm similaire à celui de Rezel
Le dossier ``neo4j`` contient le dockerfile pour l'image Neo4j, mais les dossier ``data``, ``plugins``, ``logs`` et le fichier ``conf`` doivent être récupérées depuis la base de données de production, ou une backup.
Pour build une image, par exemple neo4j: ``docker build -t bullotheque-local-neo4j .`` depuis le dossier ``neo4j``
Pour exporter une image, par exemple neo4j: ``docker save bullotheque-local-neo4j > bullotheque-local-neo4j.tar``
