package functions

import (
	structs "P1/Structs"
	utilities "P1/Utilities"
	"fmt"
	"regexp"
	"strings"
)

// MKDIR
func ProcessMKDIR(input string, path *string, r *bool) {
	flags := strings.Split(input, "-")
	for _, i := range flags {
		if i == "r" {
			*r = true
		}
		f := strings.Split(i, "=")
		if f[0] == "path" {
			*path = f[1]
			if strings.Contains(f[1], " ") {
				*path = `"` + f[1] + `"`
			}
		}
	}

	re := regexp.MustCompile(`-(\w+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]

		switch flagName {
		case "r":
			*r = true
		case "path":
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func MKDIR(path *string, r *bool) {
	// Buscamos la letra para crear el disco /home/taro/go/src/MIA1_P1_201602404/MIA/P1
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("No se encontro:", err)
		return
	}
	defer file.Close()

	var TempMBR structs.MBR
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}

	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partición no encontrada.")
		return
	}

	// Leer el superbloque
	var tempSuperblock structs.Superblock
	if err := utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error al leer superblock:", err)
		return
	}
}

// MKFILE
func ProcessMKFILE(input string, path *string, r *bool, size *int, cont *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "r":
			*r = true
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			*size = sizeValue
		case "cont":
			*cont = flagValue
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func MKFILE(path *string, r *bool) {
}

// Cat
func ProcessCAT(input string, file *string) string {
	re := regexp.MustCompile(`-file(\d*)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagIndex := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		// Generar y asignar el nombre de la clave para el mapa
		*file = flagValue
		return flagIndex
	}
	return ""
}

func CAT(file *string) {
}

// Remove
func ProcessREMOVE(input string, path *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func REMOVE(path *string) {
}

// Edit
func ProcessEDIT(input string, path *string, cont *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "cont":
			*cont = flagValue
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func EDIT(path *string, cont *string) {
}

// Rename
func ProcessRENAME(input string, path *string, name *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "name":
			*name = flagValue
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func RENAME(path *string, name *string) {
}

// Copy
func ProcessCOPY(input string, path *string, destino *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "destino":
			*destino = flagValue
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func COPY(path *string, destino *string) {
}

// Move
func ProcessMOVE(input string, path *string, destino *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "destino":
			*destino = flagValue
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func MOVE(path *string, destino *string) {
}

// Find
func ProcessFIND(input string, path *string, destino *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "destino":
			*destino = flagValue
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func FIND(path *string, destino *string) {
}

// CHOWN
func ProcessCHOWN(input string, path *string, user *string, r *bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "user":
			*user = flagValue
		case "r":
			*r = true
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func CHOWN(path *string, user *string, r *bool) {
}

// CHMOD
func ProcessCHMOD(input string, path *string, ugo *string, r *bool) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			*path = flagValue
		case "ugo":
			*ugo = flagValue
		case "r":
			*r = true
		default:
			fmt.Println("Error parametro no encontrado: " + flagName)
		}
	}
}

func CHMOD(path *string, ugo *string, r *bool) {
}
