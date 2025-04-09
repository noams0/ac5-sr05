package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

func sommer(tabint []int, channel chan int) {
	somme := 0
	for _, valeur := range tabint {
		somme += valeur
	}
	channel <- somme
}

func main() {
	start := time.Now()
	p_num := flag.Int("n", 10, "nombre")
	flag.Parse()

	tabnum := make([]int, *p_num+1)
	for i := 0; i <= *p_num; i++ {
		tabnum[i] = i
	}

	//fmt.Println("Dans main(), tabnum =", tabnum)

	//En utilisant le parallélisme de donnée

	channel_result := make(chan int)
	nbcpu := runtime.NumCPU()
	taille := len(tabnum) / nbcpu

	nbGoroutines := 0
	for i := 0; i < len(tabnum); i += taille {
		fin := i + taille
		if fin > len(tabnum) {
			fin = len(tabnum)
		}
		go sommer(tabnum[i:fin], channel_result)
		nbGoroutines++
	}

	sommeTotale := 0
	for i := 0; i < nbGoroutines; i++ {
		sommePartielle := <-channel_result
		sommeTotale += sommePartielle
	}

	//Sans utiliser le parallélisme de donnée
	//sommeTotale := 0
	//for _, valeur := range tabnum {
	//	sommeTotale += valeur
	//}

	fmt.Println("Somme totale :", sommeTotale)
	duration := time.Since(start)
	fmt.Println("Temps d'exécution :", duration)
}
