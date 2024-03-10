package Comandos

import (
	"fmt"
	"strings"
)

var contadorBloques int
var contadorArchivos int
var bloquesUsados []int64

func ValidarDatosREP(context []string) {
	contadorBloques = 0
	contadorArchivos = 0
	bloquesUsados = []int64{}

	name := ""
	path := ""
	id := ""
	ruta := ""

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "path") {
			path = strings.ReplaceAll(tk[1], "\"", "")
		} else if Comparar(tk[0], "name") {
			name = tk[1]
		} else if Comparar(tk[0], "id") {
			id = tk[1]
		} else if Comparar(tk[0], "ruta") {
			ruta = tk[1]
			fmt.Print(ruta)
		}
	}

	if name == "" || path == "" || id == "" {
		fmt.Println("Error parametros obligatorios faltantes.")
		return
	}

	if Comparar(name, "DISK") {
		//dsk(path, id)
	} else if Comparar(name, "MBR") {
		//rep para mbr
	} else if Comparar(name, "INODE") {

	} else if Comparar(name, "JOURNALIGN") {

	} else if Comparar(name, "BLOCK") {

	} else if Comparar(name, "BM_INODE") {

	} else if Comparar(name, "BM_BLOCK") {

	} else if Comparar(name, "TREE") {

	} else if Comparar(name, "SB") {

	} else if Comparar(name, "FILE") {

	} else if Comparar(name, "LS") {

	} else {
		fmt.Println("Error: parametro name incorrecto.")
	}
}
