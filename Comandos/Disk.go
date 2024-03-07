package Comandos

import (
	"MIA1_P1_201602404/Structs"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

var diskCounter = 0

func ValidarDatosMKDISK(tokens []string) {
	size := ""
	fit := ""
	unit := ""
	band_error := false

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "fit") {
			if fit == "" {
				fit = tk[1]
			} else {
				Error("", "parametro fit repetido ")
				return
			}
		} else if Comparar(tk[0], "size") {
			if size == "" {
				size = tk[1]
			} else {
				Error(" ", "Parametro size repetido.")
				return
			}
		} else if Comparar(tk[0], "unit") {
			if unit == "" {
				unit = tk[1]
			} else {
				Error("MKDISK", "parametro U repetido en el comando: "+tk[0])
				return
			}
		} else {
			Error(" ", "No se esperaba este parametro: "+tk[0])
			return
		}
	}
	if fit == "" {
		fit = "FF"
	}
	if unit == "" {
		unit = "M"
	}
	if band_error {
		return
	}
	if size == "" {
		Error("MKDISK", "se requiere parametro Size para este comando")
		return
	} else if !Comparar(fit, "BF") && !Comparar(fit, "FF") && !Comparar(fit, "WF") {
		Error("MKDISK", "valores en parametro fit no esperados")
		return
	} else if !Comparar(unit, "k") && !Comparar(unit, "m") {
		Error("MKDISK", "valores en parametro unit no esperados")
		return
	} else {
		// Obtenemos la letra del abecedario según el contador global diskCounter
		letter := string(rune('A' + diskCounter))
		// Formateamos dinámicamente la ruta del archivo de disco usando fmt.Sprintf
		val_path := fmt.Sprintf("/home/taro/Escritorio/MIA/P1/%s.dsk", letter)

		// Incrementamos el contador de discos solo si el archivo no existe
		if _, err := os.Stat(val_path); os.IsNotExist(err) {
			diskCounter++
		} else {
			// Si el archivo ya existe, avanzamos al siguiente disco
			diskCounter++
			letter = string(rune('A' + diskCounter))
			val_path = fmt.Sprintf("/home/taro/Escritorio/MIA/P1/%s.dsk", letter)
		}
		crearArchivo(size, fit, unit, val_path)
	}
}

func crearArchivo(s string, f string, u string, path string) { //MAKEFILE
	var disco = Structs.NewMBR()
	size, err := strconv.Atoi(s)

	if err != nil {
		Error("MKDISK", "Size debe ser un número entero")
		return
	}
	if size <= 0 {
		Error("MKDISK", "Size debe ser mayor a 0")
		return
	}
	if Comparar(u, "M") {
		size = 1024 * 1024 * size
	} else if Comparar(u, "k") {
		size = 1024 * size
	}
	f = string(f[0])

	disco.Mbr_tamano = int64(size)
	fecha := time.Now().String()
	copy(disco.Mbr_fecha_creacion[:], fecha)
	aleatorio, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	entero, _ := strconv.Atoi(aleatorio.String())
	disco.Mbr_dsk_signature = int64(entero)
	copy(disco.Dsk_fit[:], string(f[0]))
	disco.Mbr_partition_1 = Structs.NewParticion()
	disco.Mbr_partition_2 = Structs.NewParticion()
	disco.Mbr_partition_3 = Structs.NewParticion()
	disco.Mbr_partition_4 = Structs.NewParticion()

	if ArchivoExiste(path) {
		_ = os.Remove(path)
	}

	if !strings.HasSuffix(path, "dsk") {
		Error("MKDISK", "Extensión de archivo no válida.")
		return
	}
	carpeta := ""
	direccion := strings.Split(path, "/")

	for i := 0; i < len(direccion)-1; i++ {
		carpeta += "/" + direccion[i]
		if _, err_ := os.Stat(carpeta); os.IsNotExist(err_) {
			os.Mkdir(carpeta, 0777)
		}
	}

	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		Error("MKDISK", "No se pudo crear el disco.")
		return
	}
	var vacio int8 = 0
	s1 := &vacio
	var num int64 = 0
	num = int64(size)
	num = num - 1
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s1)
	EscribirBytes(file, binario.Bytes())

	file.Seek(num, 0)

	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	EscribirBytes(file, binario2.Bytes())

	file.Seek(0, 0)
	disco.Mbr_tamano = num + 1

	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, disco)
	EscribirBytes(file, binario3.Bytes())
	file.Close()
	nombreDisco := strings.Split(path, "/")
	Mensaje("MKDISK", "¡Disco \""+nombreDisco[len(nombreDisco)-1]+"\" creado correctamente!")
	diskCounter++
}

func RMDISK(tokens []string) {

	if len(tokens) > 1 {
		Error("RMDISK", "Parametro -driveletter esperado...")
		return
	}

	driveletter := ""

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "driveletter") {
			if driveletter == "" {
				driveletter = tk[1]
			} else {
				Error("RMDISK", "Parametro driveletter repetido en el comando: "+tk[0])
				return
			}
		} else {
			Error("RMDISK", "no se esperaba el parametro "+tk[0])
			return
		}
	}

	driveletter = fmt.Sprintf("/home/taro/Escritorio/MIA/P1/%s.dsk", driveletter)

	if driveletter == "" {
		Error("RMDISK", "se requiere parametro Path para este comando")
		return
	} else {
		if !ArchivoExiste(driveletter) {
			Error("RMDISK", "No se encontró el disco en la ruta indicada.")
			return
		}
		if !strings.HasSuffix(driveletter, "dsk") {
			Error("RMDISK", "Extensión de archivo no válida.")
			return
		}
		if Confirmar("¿Desea eliminar el disco: " + driveletter + " ?") {
			err := os.Remove(driveletter)
			if err != nil {
				Error("RMDISK", "Error al intentar eliminar el archivo. :c")
				return
			}
			Mensaje("RMDISK", "Disco ubicado en "+driveletter+", ha sido eliminado exitosamente. :D")
			return
		} else {
			Mensaje("RMDISK", "Eliminación del disco "+driveletter+", cancelada exitosamente. :D")
			return
		}
	}

}
