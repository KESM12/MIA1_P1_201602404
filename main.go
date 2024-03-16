package main

import (
	analizar "P1/Analizador"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n***********Kevin Esuardo Secaida Molina ***********")
		fmt.Print("Comando: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			continue
		}
		input = strings.TrimSpace(input)
		// Aquí puedes procesar el comando de Linux ingresado
		analizar.Comandos(input)

	}
}
