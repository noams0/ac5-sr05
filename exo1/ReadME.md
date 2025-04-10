# Somme parallèle d’un tableau en Go

## Objectif

Ce programme Go calcule la somme des éléments d’un tableau d'entiers, en exploitant le **parallélisme de données** : le tableau est découpé en segments traités en parallèle dans des goroutines. Chaque goroutine calcule une somme partielle sur une portion du tableau, et le résultat final est obtenu en les combinant.

---

## Fonctionnement du code

1. **Récupération de la taille du tableau via la ligne de commande** :
   ```bash
   go run main.go -n 1000000
   ```
   L’option `-n` détermine jusqu’à quel entier le tableau va compter.

2. **Initialisation du tableau** :
   ```go
   tab := make([]int, *p_num+1)
   for i := 0; i <= *p_num; i++ {
       tab[i] = i
   }
   ```

3. **Détermination du nombre de cœurs disponibles** :
   ```go
   nbCPU := runtime.NumCPU()
   taille := len(tab) / nbCPU
   if taille == 0 {
       taille = 1 // pour éviter que ça bloque si tableau petit
   }
   ```

4. **Découpage du tableau + lancement des goroutines** :
   Pour chaque portion du tableau, une goroutine est lancée avec la fonction `sommer`, qui calcule la somme locale et l’envoie sur un canal. On utilise un `WaitGroup` pour attendre la fin de toutes les goroutines.
   ```go
   for i := 0; i < len(tab); i += taille {
       fin := i + taille
       if fin > len(tab) {
           fin = len(tab)
       }
       wg.Add(1)
       go sommer(tab[i:fin], ch, &wg)
   }
   ```

5. **Récupération des résultats** :
   Une fois toutes les goroutines terminées, le canal est fermé, ce qui permet de le parcourir avec `range`:
   ```go
   go func() {
       wg.Wait()
       close(ch)
   }()
   for somme := range ch {
       sommeTotale += somme
   }
   ```

6. **Affichage du résultat et du temps d'exécution** :
   ```go
   fmt.Println("Somme totale :", sommeTotale)
   fmt.Println("Temps d'exécution :", time.Since(start))
   ```

---

## Exemple de sortie

```
Somme totale : 500000500000
Temps d'exécution : 3.14ms
```

---

## Conclusion 

### Liens avec le parallélisme de données

Ce programme est un exemple simple mais parlant du **parallélisme de données** : une même opération (la somme) est appliquée à des **données découpées** (les tranches du tableau) **en parallèle**. Chaque cpu disponible peut exécuter une goroutine, ce qui permet d’accélérer le traitement pour de grands tableaux.



### Pour aller plus loin

- Comparer les performances avec une version séquentielle (cf. partie commentée dans le code).

