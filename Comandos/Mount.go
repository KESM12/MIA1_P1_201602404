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
	Id     [4]byte
}

var alfabeto = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

func ValidarDatosMOUNT(context []string) {
	//fmt.Println(context)

	name := ""
	driveletter := ""
	val_path := ""
	fmt.Println(val_path)
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
	if driveletter == "" || name == "" {
		fmt.Println("El comando MOUNT requiere parámetros obligatorios")
		return
	}
	mount(val_path, name, driveletter)
	listaMount(driveletter)
}

func mount(path string, name string, driveletter string) {
	file, error_ := os.Open(path)
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

	particion := BuscarParticiones(disk, name, path)
	if particion.Part_type == 'E' || particion.Part_type == 'L' {
		var nombre [16]byte
		copy(nombre[:], []byte(name))
		if particion.Part_name == nombre && particion.Part_type == 'E' || particion.Part_type == 'L' {
			fmt.Println("No se puede montar una partición extendida, ni lógica.")
			return
		}
		encontrada := false
		if !encontrada {
			fmt.Println("No se encontró la partición Lógica.")
			return
		}
	}
	for i := 0; i < 99; i++ {
		var ruta [150]byte
		copy(ruta[:], path)
		if DiscMont[i].Path == ruta {
			for j := 0; j < 26; j++ {
				var nombre [20]byte
				copy(nombre[:], name)
				if DiscMont[i].Particiones[j].Nombre == nombre {
					fmt.Println("Ya se ha montado la partición " + name)
					return
				}
				if DiscMont[i].Particiones[j].Estado == 0 {
					DiscMont[i].Particiones[j].Estado = 1
					DiscMont[i].Particiones[j].Letra = alfabeto[j]
					copy(DiscMont[i].Particiones[j].Nombre[:], name)
					re := string(driveletter) + strconv.Itoa(i+1)
					fmt.Println("Se ha realizado correctamente el mount -id = " + re + string("04"))
					return
				}
			}
		}
	}
	for i := 0; i < 99; i++ {
		if DiscMont[i].Estado == 0 {
			DiscMont[i].Estado = 1
			copy(DiscMont[i].Path[:], path)
			for j := 0; j < 26; j++ {
				if DiscMont[i].Particiones[j].Estado == 0 {
					DiscMont[i].Particiones[j].Estado = 1
					DiscMont[i].Particiones[j].Letra = alfabeto[2]
					copy(DiscMont[i].Particiones[j].Nombre[:], name)
					re := string(driveletter) + strconv.Itoa(i+1)
					fmt.Println("Se ha realizado correctamente el mount -id = " + re + string("04"))
					return
				}
			}
		}
	}
}

func GetMount(comando string, id string, p *string) Structs.Particion {

	if !(len(id) == 4 && id[0] >= 'A' && id[0] <= 'Z' && id[1] >= '0' && id[1] <= '9' && id[2] == '0' && id[3] == '4') {
		fmt.Println("El primer identificador no es válido.")
		return Structs.Particion{}
	}
	letra := id[len(id)-1]
	id = strings.ReplaceAll(id, "04", "")
	i, _ := strconv.Atoi(string(id[0] - 1))
	if i < 0 {
		//fmt.Println("El primer identificador no es válido.")
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
					fmt.Println("No se ha encontrado el disco4")
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
