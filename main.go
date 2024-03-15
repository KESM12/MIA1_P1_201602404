package main

import (
	analyzer "P1/Analyzer"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("******Kevin Secaida******")
		fmt.Print("Ingrese un comando: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			continue
		}

		input = strings.TrimSpace(input) //quitamos el salto de linea
		//Para llamar una funcion desde otro archivo este debe ir en mayuscula al inicio
		analyzer.Command(input)
	}
}
