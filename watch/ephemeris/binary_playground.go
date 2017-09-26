package ephemeris

import "fmt"

func main() {
	v := uint8(9)
	fmt.Printf("%08b\n", v)
	//v >>= 1
	fmt.Printf("%08b\n", 255-3) //3 peut s'écrire sur 2 bits = 11
	fmt.Printf("%d\n", v)
	v &^= 255-3   //met tous les bits à zéro sauf les 2 derniers qui seront inchangés: 255-3 = 11111100
	v &^= 255-1 // de même on peut faire ceci pour récupérer uniquement le dernier bit
	v &^= 255-7 // de même on peut faire ceci pour récupérer les 3 derniers bits
	fmt.Printf("%08b\n", v)
	fmt.Printf("%d\n", v)
}

