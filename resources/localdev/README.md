Pour développer en local, vous avez besoin de 5 choses :
  - une instance docker locale de Keycloak pour l'authentification
  - une instance docker locale de Neo4j pour la base de données
  - un fichier d'environnement [`.env`](./.env) pour les configurations
  - un dossier [``data``](./data) contenant les pages de couvertures des livres
  - une copie du dossier [``script``](../../src/script)
  - une copie du dossier [``view``](../../src/view)

Les dossiers et le fichier doivent être placés au même endroit que l'exécutable du site

# Instances docker
Pour mettre en place les instances docker locales, [c'est par ici !](./docker/README.md)

# Fichier d'environnement
Le programme utilise le fichier d'environnement ``.env`` contenant :
  - l'url du site de la Bulle (localhost:[PORT] quand on développe en local, https://bulle.rezel.net quand on est en production)
  - les identifiants de connexion à la base de données neo4j
  - les identifiants de connexion à Keycloak

Ce fichier doit être placé au même endroit que l'exécutable (le mettre dans ``src/`` permet d'utiliser l'utilitaire de live-reload ``air`` pour plus de confort quand on développe en local)

## Conseil
Il est conseillé de laisser la structure du projet intacte et de 
    - créer un lien symbolique entre le fichier ``ressources/localdev/.env`` et ``src/.env`` 
    - créer un lien symbolique entre le dossier ``ressources/localdev/data`` et ``src/data``
afin de ne pas avoir à copier-coller les dossiers ``script`` et ``view`` à chaque modification