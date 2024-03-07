package Comandos

import (
	"fmt"
	"strings"
	// "Ejemplos_Proyecto/Structs"
	// "bytes"
	// "context"
	// "encoding/binary"
	// "fmt"
	// "os"
	// "strconv"
	// "strings"
	// "unsafe"
)

var DiscMont [99]DiscoMontado

type DiscoMontado struct {
	Path        [150]byte
	Estado      byte
	Particiones [26]ParticionMontada
}

type ParticionMontada struct {
	Letra  byte
	Estado byte
	Nombre [20]byte
}

func ValidarDatosMOUNT(context []string) {
	name := ""
	driveletter := ""
	val_path := ""

	for i := 0; i < len(context); i++ {
		current := context[i]
		comando := strings.Split(current, "=")
		if Comparar(comando[0], "name") {
			name = comando[1]
		} else if Comparar(comando[0], "driveletter") {
			driveletter = strings.ReplaceAll(comando[1], "\"", "")
			val_path = fmt.Sprintf("/home/taro/go/src/MIA1_P1_201602404/MIA/P1/%s.dsk", driveletter)
		}
	}
	if val_path == "" || name == "" {
		fmt.Println("El comando MOUNT requiere parámetros obligatorios")
		return
	}
	mount(val_path, name)
	listaMount()
}
