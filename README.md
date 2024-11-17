# Hangman-WEB

Bienvenue dans Le Supplice du Pendard !

Le Supplice du Pendard est un jeu du pendu sur le thème du fantastique et du Moyen Age développé majoritairement en Golang et HTML par
Johan Daumas et Yolan Chiotti dans le cadre d'un projet à Ynov.

Comment acceder au jeu sur son navigateur: 
- Il faut lancer l'executable Launcher.exe trouvable dans le dossier "src" du projet.

Liste des roots du site web:

    http://localhost:8080/mainMenu : Menu principal du jeu.

    http://localhost:8080/mainMenu/userForm : Choix du mode de jeu et du pseudonyme.

    http://localhost:8080/game : Jeu du pendu.

    http://localhost:8080/game/resultat : Ecran des résultats de fin de partie (gagné ou perdu).

    http://localhost:8080/leaderboard : Ecran du tableau des scores.


    *TRAITEMENT DES DONNEES*

    http://localhost:8080/game/initialisation/first : Rénitialisation des variables utilisateur (pseudonyme, score, word streak) et du jeu (mot secret, nombre de vies, historique des lettres ...etc) lors d'une nouvelle partie.

    http://localhost:8080/game/initialisation : Rénitialisation des variables du jeu lorsque le joueur souhaite rejouer après une victoire.

    http://localhost:8080/game/treatment : Vérification de la lettre entrée par le joueur.

    http://localhost:8080/game/treatment/joker : Utilisation d'un joker.



