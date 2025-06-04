Pour développer en local, vous avez besoin d'une instance locale de Keycloak et de neo4j

# Images docker
L'[image docker de Keycloak](./bullotheque-local-keycloak.tar) est basée sur l'image officielle de Keycloak, et contient un realm similaire à celui de Rezel.
L'[image docker de neo4j](./bullotheque-local-neo4j.tar) est basée sur la base de donnée de production, et contient une partie de ces données.
Pour les utiliser, vous pouvez les importer avec la commande ``docker load -i bullotheque-local-keycloak.tar`` et ``docker load -i bullotheque-local-neo4j.tar``
Pour lancer les deux images, vous pouvez utiliser le fichier [docker-compose.yml](./docker-compose.yml) avec la commande ``docker-compose up -d``
Au bout de quelques secondes, vous devriez pouvoir vous connecter à Keycloak sur [http://localhost:8080](http://localhost:8080) et à neo4j sur [http://localhost:7474](http://localhost:7474)

## Utilisateurs Keycloak
- username: admin , password: password
- username: user  , password: user

L'utilisateur ``admin`` l'est aussi dans la base de données, donc celui-ci a accès au pannel admin

## Base de données neo4j
  - username: neo4j , password: YOUR_PASSWORD
  (c'est réellement le mot de passe)

# ATTENTION
Les images docker ne fournissent pas de volumes pour sauvegarder les données. Si vous arrêtez les containers, vous perdrez les données.
