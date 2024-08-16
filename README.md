# Bullothèque
Application web de l'association Bulle de Telecom Paris !

# Base de données
Le site utilise **neo4j** pour sa base de données. C'est une base de données **orientée graphe**, qui utilise le langage de requêtes **cypher**
Le structure de cette base comme suit:
![](http://github.com/Hy3z/Bulle-bliotheque/blob/main/resources/image/bdd.png)

Comme on est dans une base de données graphe, on fait des requêtes avec du pattern matching
> <ins>Exemple</ins>: ``MATCH (book:Book)-[r:PART_OF]->(s:Serie) WHERE book.title <> "Eldorado" RETURN r.opus`` renvoit tous les numéros d'opus de la relation `PART_OF` lorsque le titre du livre n'est pas `Eldorado`
>
Consulter https://neo4j.com/docs/cypher-manual/current/introduction/ pour plus d'informations

# HTMX
Toute l'architecture du site tourne autour d'un mini-framework Javascript appellé HTMX. Au lieu d'envoyer du JSON depuis le serveur et de l'interpréter en HTML dans le navigateur, HTMX permet de renvoyer directement un fichier HTML qui est automatiquement inséré dans le naviguateur.
HTMX s'intègre directement dans la page HTML, ce qui rend plus lisible les interactions
> <ins>Exemple</ins>: ``<a href="http://monlien" hx-boost="true" hx-target=".main" hx-swap="innerHTML">CLIC</a>`` modifie le comportement de la balise `a` qui, lorsque cliquée,  va remplacer l'**intérieur** de la balise portant l'ID `main` par la réponse HTTP de l'url `http://monlien`
>
Le site utilise donc HTMX pour toujours uniquement renvoyer la partie d'HTML qui nous interesse
> <ins>Exemple</ins>: Lorsqu'on clic depuis le site sur le bouton `Contact`, on ne renvoit depuis le serveur **que** la partie **en dessous** du bandeau du site (le header, avec le logo, la barre de recherche, etc.), car le bandeau n'a pas de raisons de changer quand on change de page
>
Cette méthode pose pourtant un problème: Quand je clic sur le lien http://bulle.rezel.net/contact depuis un nouvel onglet, ou que je recharge la page, il faut que le serveur HTTP me renvoit aussi cette fois-ci le bandeau du haut, car il n'est pas encore dans mon navigateur (contrairement au cas où je clic sur le bouton Contact depuis le site).
Pour répondre à ce problème, on ajoute toujours un paramètre dans l'entête HTTP de la requête pour préciser la nature de l'HTML qu'on attend en retour. Si on vend la page entière, on ne met pas de paramètre. Si on veut seulement la partie en dessous du bandeau, on met le paramètre `Tmpl` à la valeur `main`. Si on veut autre chose, on mettra autre chose.
><ins>Exemple</ins>: ``<a href="http://monlien" hx-boost="true" hx-target=".main" hx-swap="innerHTML" hx-headers='{"Tmpl":"main"}'>CLIC</a>`` possède le même comportement que précedemment, mais avec ce fameux paramètre en plus, qui sera lu par le serveur HTTP
>
Consulter https://htmx.org/ pour plus d'informations

# Template HTML
Pour créer et renvoyer efficacement les bouts d'HTML que recquiert HTMX, on utilise des templates HTML. Ce sont des fichiers HTML classiques, mais qui peuvent recevoir des paramètres ou éxecuter des fonctions afin de remplir le fichier. Comme le serveur est en **Go**, les templates HTML utilisent la syntaxe `{{...}}` pour intéragir avec ces variables
><ins>Exemple</ins>: ``{{template "ma_template"}} Mon nom: {{.Nom}} {{end}}`` définie une template HTML qui prend en entrée une **struct** (comme en C) avec le champ `Nom`, et qui écrit ce champ dans le fichier HTML
>On peut aussi utiliser des tableaux: ``{{template "ma_template"}} {{range .mon_tableau}} Contenu: {{.}} {{end}} {{end}}`` itère sur un tableau passé en entrée et affiche à chaque fois son contenu

# Structure du projet
Le fichier principal est  `/main.go`, c'est l'entrée du programme (comme `__main__` en python). Il regroupe les appels aux différents modules qui sont:
![](http://github.com/Hy3z/Bulle-bliotheque/blob/main/resources/image/files.png)

- ## api
L'api fait le lien entre les différentes routes du site (https://bulle.rezel.net/contact, https://bulle.rezel.net/login, etc.), et les fonctions qui traitent les requêtes (elles se trouvent dans le dossier `service`). L'api peut aussi faire le lien avec un fichier, pour renvoyer le CSS du site par exemple.
- ## auth
L'auth gère tous les problèmes d'authentification sur le site: gérer les cookies d'authentification, garder ces cookies à jour, rediriger vers la page de connection quand on essaie d'accéder à une page nécessitant une authentification, vérifier que l'utilisateur possède un rôle admin, etc.
- ## database
La database gère la connection avec la base de données, et l'envoie de requêtes à celle-ci
- ## logger
Le logger gère la trace écrite de l'éxécution du serveur. Lorsqu'il y a une erreur avec le backend, on retrouve normalement toujours un trace écrite dans les logs
- ## model
Le dossier `model`comporte toutes les structures associées aux templates HTML utilisées par le site.
><ins>Exemple</ins>: Pour visualiser les informations de son compte, on a besoin de l'UUID de l'utilisateur, son nom, et les livres qu'il a commentés/likés/empruntés. Le fichier `model/Account.go` possède la structure `Account` qui contient ces informations, et on peux directement passer cette structure à la template HTML `account` pour visualiser les informations de l'utilisateur
- ## script
Le dossier script contient toutes les requêtes cypher utilisées par le site. On aurait pû techniquement écrire ces requêtes directement dans les fichiers .go (qui sont ensuite compilés), mais faire comme cela permet de modifier/améliorer les requêtes sans redemarrer le site (en théorie en tout cas)
- ## service
Le gros de la logique du site est dans ce fichier. On y retrouve les fonctions qui traitent la requête HTTP entrante, et renvoient la template HTML correspondante, avec les informations demandées, que ça soit un livre en particulier ou les informations de l'utilisateur, qu'on obtient en lisant les cookies d'authentification
- ## util
Juste des variables pour garder le site propre. C'est là qu'on définit vraiment les routes disponibles sur le site, qu'on définit les noms des variables d'entête HTTP qu'on va utiliser, etc.
- ## view
Enfin, on retrouve içi tous les fichiers que va manipuler le site (hormis la base d'images de couverture des livres et séries, qui sera dans `/data`). On y retrouve les templates HTML dans `/view/html`, les fichiers Javascript et CSS dans `/view/script` et `/view/css`, ainsi que les images utilisées dans le bandeau et la page d'acceuil dans `/view/image`
# Exemple d'éxécution
![](http://github.com/Hy3z/Bulle-bliotheque/blob/main/resources/image/example.png)
