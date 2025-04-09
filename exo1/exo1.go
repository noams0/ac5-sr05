package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func sommer(tab []int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	somme := 0
	for _, v := range tab {
		somme += v
	}
	ch <- somme
}

func main() {
	start := time.Now()

	p_num := flag.Int("n", 10, "taille du tableau")
	flag.Parse()

	tab := make([]int, *p_num+1)
	for i := 0; i <= *p_num; i++ {
		tab[i] = i
	}

	nbCPU := runtime.NumCPU()
	taille := len(tab) / nbCPU
	if taille == 0 {
		taille = 1
	}

	ch := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < len(tab); i += taille {
		fin := i + taille
		if fin > len(tab) {
			fin = len(tab)
		}
		wg.Add(1)
		go sommer(tab[i:fin], ch, &wg)
	}

	// Ferme le canal quand toutes les goroutines ont fini
	go func() {
		wg.Wait()
		close(ch)
	}()

	sommeTotale := 0
	for somme := range ch {
		sommeTotale += somme
	}

	fmt.Println("Somme totale :", sommeTotale)
	fmt.Println("Temps d'exécution :", time.Since(start))
}
