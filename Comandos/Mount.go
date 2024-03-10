package Comandos

import (
	"MIA1_P1_201602404/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
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

var alfabeto = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

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
	mount(val_path, name, driveletter)
	listaMount(driveletter)
}

func mount(p string, n string, d string) {
	file, error_ := os.Open(p)
	if error_ != nil {
		fmt.Println("No se ha podido abrir el archivo.")
		return
	}

	disk := Structs.NewMBR()
	file.Seek(0, 0)

	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		fmt.Println("Error al leer el archivo")
		return
	}
	file.Close()

	particion := BuscarParticiones(disk, n, p)
	if particion.Part_type == 'E' || particion.Part_type == 'L' {
		var nombre [16]byte
		copy(nombre[:], n)
		if particion.Part_name == nombre && particion.Part_type == 'E' {
			fmt.Println("No se puede montar una partición extendida.")
			return
		} else {
			ebrs := GetLogicas(*particion, p)
			encontrada := false
			if len(ebrs) != 0 {
				for i := 0; i < len(ebrs); i++ {
					ebr := ebrs[i]
					nombreebr := ""
					for j := 0; j < len(ebr.Part_name); j++ {
						if ebr.Part_name[j] != 0 {
							nombreebr += string(ebr.Part_name[j])
						}
					}

					if Comparar(nombreebr, n) && ebr.Part_status == '1' {
						encontrada = true
						n = nombreebr
						break
					} else if nombreebr == n && ebr.Part_status == '0' {
						fmt.Println("No se puede montar una partición Lógica eliminada.")
						return
					}
				}
				if !encontrada {
					fmt.Println("No se encontró la partición Lógica.")
					return
				}
			}
		}
	}
	for i := 0; i < 99; i++ {
		var ruta [150]byte
		copy(ruta[:], p)
		if DiscMont[i].Path == ruta {
			for j := 0; j < 26; j++ {
				var nombre [20]byte
				copy(nombre[:], n)
				if DiscMont[i].Particiones[j].Nombre == nombre {
					fmt.Println("Ya se ha montado la partición " + n)
					return
				}
				if DiscMont[i].Particiones[j].Estado == 0 {
					DiscMont[i].Particiones[j].Estado = 1
					DiscMont[i].Particiones[j].Letra = alfabeto[j]
					copy(DiscMont[i].Particiones[j].Nombre[:], n)
					re := string(d) + strconv.Itoa(i+1) + string("04")
					fmt.Println("Se ha realizado correctamente el mount -id = " + re)
					return
				}
			}
		}
	}
	for i := 0; i < 99; i++ {
		if DiscMont[i].Estado == 0 {
			DiscMont[i].Estado = 1
			copy(DiscMont[i].Path[:], p)
			for j := 0; j < 26; j++ {
				if DiscMont[i].Particiones[j].Estado == 0 {
					DiscMont[i].Particiones[j].Estado = 1
					DiscMont[i].Particiones[j].Letra = alfabeto[j]
					copy(DiscMont[i].Particiones[j].Nombre[:], n)

					re := string(d) + strconv.Itoa(i+1) + string("04")
					fmt.Println("se ha realizado correctamente el mount -id= " + re)
					return
				}
			}
		}
	}
}

func GetMount(comando string, id string, p *string) Structs.Particion {
	if !(id[0] == '7' && id[1] == '9') {
		fmt.Println("El primer identificador no es válido.")
		return Structs.Particion{}
	}
	letra := id[len(id)-1]
	id = strings.ReplaceAll(id, "79", "")
	i, _ := strconv.Atoi(string(id[0] - 1))
	if i < 0 {
		fmt.Println("El primer identificador no es válido.")
		return Structs.Particion{}
	}
	for j := 0; j < 26; j++ {
		if DiscMont[i].Particiones[j].Estado == 1 {
			if DiscMont[i].Particiones[j].Letra == letra {

				path := ""
				for k := 0; k < len(DiscMont[i].Path); k++ {
					if DiscMont[i].Path[k] != 0 {
						path += string(DiscMont[i].Path[k])
					}
				}

				file, error := os.Open(strings.ReplaceAll(path, "\"", ""))
				if error != nil {
					fmt.Println("No se ha encontrado el disco")
					return Structs.Particion{}
				}
				disk := Structs.NewMBR()
				file.Seek(0, 0)

				data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
				buffer := bytes.NewBuffer(data)
				err_ := binary.Read(buffer, binary.BigEndian, &disk)

				if err_ != nil {
					fmt.Println("Error al leer el archivo")
					return Structs.Particion{}
				}
				file.Close()

				nombreParticion := ""
				for k := 0; k < len(DiscMont[i].Particiones[j].Nombre); k++ {
					if DiscMont[i].Particiones[j].Nombre[k] != 0 {
						nombreParticion += string(DiscMont[i].Particiones[j].Nombre[k])
					}
				}
				*p = path
				return *BuscarParticiones(disk, nombreParticion, path)
			}
		}
	}
	return Structs.Particion{}
}

func listaMount(d string) {
	fmt.Println("LISTADO DE MOUNTS")
	for i := 0; i < 99; i++ {
		for j := 0; j < 26; j++ {
			if DiscMont[i].Particiones[j].Estado == 1 {
				nombre := ""
				for k := 0; k < len(DiscMont[i].Particiones[j].Nombre); k++ {
					if DiscMont[i].Particiones[j].Nombre[k] != 0 {
						nombre += string(DiscMont[i].Particiones[j].Nombre[k])
					}
				}
				fmt.Println(" id:" + string(d) + strconv.Itoa(i+1) + string("04") + ", Nombre: " + nombre)
			}
		}
	}
}
