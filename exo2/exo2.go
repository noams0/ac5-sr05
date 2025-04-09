package main

import (
	"flag"
	"fmt"
	"math"
	"time"
)

func filtre(nb int, in <-chan int, result chan<- int) {
	out := make(chan int)
	var suivantCree bool = false

	for {
		n := <-in
		fmt.Println("Filtre", nb, "a reçu :", n)

		if n == 0 {
			result <- nb
			close(out)
			return
		}

		if n%nb != 0 {
			if !suivantCree { // si on n'a pas créé le filtre pour le premier nombre premier sur lequelle on tombe
				suivantCree = true
				go filtre(n, out, result) // nouveau filtre avec nouveau nombre premier et chanel dédié
			}
			out <- n
		}
	}
}

func main() {
	start := time.Now()

	source := make(chan int)
	result := make(chan int)

	p_num := flag.Int("n", 30, "limite des nombres")
	flag.Parse()

	go filtre(2, source, result)

	for i := 3; i <= *p_num; i++ {
		source <- i
	}

	// lancement du signal pour récupération des nombres premiers
	source <- 0
	timeout := time.Duration(float64(*p_num)*math.Log(float64(*p_num))) * time.Microsecond * 2 // l'algo a un complexité de O(n log(n)) et on fait *2 pour avoir de la marge
	var primes []int

LOOP:
	for {
		select {
		case p := <-result:
			primes = append(primes, p)
		case <-time.After(timeout):
			break LOOP
		}
	}
	fmt.Println("Liste des nombres premiers jusqu’à", *p_num, ":", primes)

	// Calcul du temps d'exécution
	duration := time.Since(start)
	fmt.Println("Temps d'exécution :", duration)
}
