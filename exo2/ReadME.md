# Nombres premiers en Go — Crible de Hoare

## Objectif

Ce programme implémente le **crible de Hoare**, une variante concurrente du crible d'Ératosthène, pour identifier les **nombres premiers jusqu'à une limite donnée**. L'idée principale est d'utiliser une **chaîne dynamique de goroutines** agissant comme des filtres, chacun bloquant les multiples d'un nombre premier découvert.

---

## Fonctionnement du programme

1. **Lancement** avec un flag `-n` indiquant la borne supérieure :
   ```bash
   go run main.go -n 30
   ```

2. **Génération des entiers** de 3 à `n` :
   ```go
   for i := 3; i <= *p_num; i++ {
       source <- i
   }
   ```
   Le `2` est traité à part, car c’est le premier nombre premier utilisé pour initialiser la première goroutine de filtre.

3. **Mise en place des filtres en cascade** :
    - Le premier filtre est créé sur la base de 2.
    - Chaque filtre reçoit les nombres qui ne sont pas multiples de son propre nombre.
    - Dès qu’un filtre reçoit un nombre qu’il considère premier (non divisible par ses prédécesseurs), il lance une **nouvelle goroutine** avec un canal dédié.

   ```go
   go filtre(2, source, result)
   ```

   Puis, dans la fonction `filtre` :
   ```go
   if !suivantCree {
       suivantCree = true
       go filtre(n, out, result)
   }
   ```

4. **Fin du pipeline** :
    - L’envoi d’un `0` dans le canal source est utilisé comme **signal de fin** pour déclencher la remontée des résultats vers le canal `result`.

5. **Récupération des résultats** :
    - Les nombres premiers sont récupérés depuis le canal `result` dans un `select`, avec un **timeout dynamique** basé sur la complexité du crible :
      ```go
      timeout := time.Duration(float64(*p_num)*math.Log(float64(*p_num))) * time.Microsecond * 2
      ```

    - Ce timeout empêche un blocage si la remontée des résultats n’est pas complète (ou en cas de bug).

6. **Affichage des résultats** :
   Le programme imprime la liste des nombres premiers et le temps d’exécution.

---

## Exemple de sortie

```
Filtre 2 a reçu : 3
Filtre 3 a reçu : 5
Filtre 2 a reçu : 4
...
Liste des nombres premiers jusqu’à 30 : [2 3 5 7 11 13 17 19 23 29]
Temps d'exécution : 2.18ms
```

---

## Conclusion

### Liens avec la programmation répartie

Ce programme illustre un **parallélisme de pipeline**, où chaque unité de traitement (goroutine) effectue une opération simple (filtrage) sur un flux de données, et passe les valeurs non filtrées à la suite du pipeline.


### Points intéressants

- Le programme ne stocke pas tous les entiers en mémoire : les données circulent dans un **pipeline**.
- Le nombre de goroutines dépend du nombre de premiers trouvés (et donc de `n`).

Pour observer le comportement du pipeline, laisser les `fmt.Println()` dans `filtre()` ; pour plus de performance, les commenter.

---
