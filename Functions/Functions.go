package functions

import (
	structs "P1/Structs"
	utilities "P1/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var fileCounter int = 0
var particionesMontadasListado = "Listado de particiones montadas:\n"

// MKDISK
func ProcesarElMKDISK(input string, size *int, fit *string, unit *string) {
	input = strings.ToLower(input)
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			*size = sizeValue
		case "fit":
			flagValue = flagValue[:1]
			*fit = flagValue
		case "unit":
			*unit = flagValue
		default:
			fmt.Println("Error: Bandera no encontrada: " + flagName)
		}
	}

	if *fit == "" {
		*fit = "f"
	}
	if *unit == "" {
		*unit = "m"
	}
}

func AsignarAlfabeto(size *int, fit *string, unit *string) {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if *unit == "k" {
		*size = *size * 1024
	} else {
		*size = *size * 1024 * 1024
	}

	if err := crear_Archivos(fmt.Sprintf("/home/taro/go/src/MIA1_P1_201602404/MIA/P1/%c.dsk", letters[fileCounter]), *size, *fit); err != nil {
		fmt.Println("Error al crear el disco:", err)
		return
	}

	fileCounter++
}

func crear_Archivos(filename string, size int, fit string) error {
	fmt.Println("Creando archivo")
	err := utilities.Crear_Archivo(filename)
	if err != nil {
		return err
	}

	file, err := utilities.OpenFile(filename)
	if err != nil {
		return nil
	}
	data := make([]byte, size)

	for i := range data {
		data[i] = 0
	}
	err = utilities.WriteObject(file, data, 0)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05")
	var TempMBR structs.MBR
	TempMBR.Mbr_tamano = int32(size)
	copy(TempMBR.Mbr_fecha_creacion[:], []byte(timeString))
	TempMBR.Mbr_dsk_signature = int32(GenerateUniqueID())
	copy(TempMBR.Dsk_fit[:], fit)

	if err := utilities.WriteObject(file, TempMBR, 0); err != nil {
		return nil
	}

	var mbr structs.MBR
	if err := utilities.ReadObject(file, &mbr, 0); err != nil {
		return nil
	}

	defer file.Close()
	fmt.Println("Archivo creado con exito")
	return nil
}

func ProcessRMDISK(input string, driveletter *string) {
	input = strings.ToLower(input)
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "driveletter":
			*driveletter = flagValue
		default:
			fmt.Println("[Error] Bandera no encontrada: " + flagName)
		}
	}
}

func Eliminar_ArchivoBin(driveletter *string) {
	*driveletter = strings.ToUpper(*driveletter)
	filename := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + *driveletter + ".dsk"
	if _, err := os.Stat(filename); err == nil {
		fmt.Print("¿Desea eliminar el archovo: ?" + *driveletter + ".dsk")
		var input string
		fmt.Print("Ingrese 'y' para continuar o 'n' para cancelar: ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			return
		}

		if input == "y" {
			if err := os.Remove(filename); err != nil {
				fmt.Println("Error al eliminar el archivo:", err)
				return
			}
			fileCounter--
		} else {
			return
		}

	} else if os.IsNotExist(err) {
		fmt.Printf("El archivo %s.dsk no existe.\n", filename)
	} else {
		fmt.Println("Error al verificar la existencia del archivo:", err)
	}
}
func GestionarFDISK(input string, size *int, driveletter *string, name *string, unit *string, type_ *string, fit *string, delete *string, add *int, path *string) {
	input = strings.ToLower(input)
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			*size = sizeValue
		case "driveletter":
			*driveletter = flagValue
		case "name":
			*name = flagValue
		case "unit":
			*unit = flagValue
		case "type":
			*type_ = flagValue
		case "fit":
			flagValue = flagValue[:1]
			*fit = flagValue
		case "delete":
			*delete = flagValue
		case "add":
			addValue := 0
			fmt.Sscanf(flagValue, "%d", &addValue)
			*add = addValue
		case "path":
			*path = flagValue
		default:
			fmt.Println("[Error] Bandera no encontrada: " + flagName)
		}
		if *unit == "" {
			*unit = "k"
		}
		if *fit == "" {
			*fit = "w"
		}
		if *type_ == "" {
			*type_ = "p"
		}
	}
}

func CRUDdeParticiones(size *int, driveletter *string, name *string, unit *string, type_ *string, fit *string, delete *string, add *int, path *string) {

	if *unit == "k" {
		*size = *size * 1024
	} else if *unit == "m" {
		*size = *size * 1024 * 1024
	}
	if *unit == "k" {
		*add = *add * 1024
	} else if *unit == "m" {
		*add = *add * 1024 * 1024
	}

	*driveletter = strings.ToUpper(*driveletter)
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + *driveletter + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var compareMBR structs.MBR
	copy(compareMBR.Mbr_particion[0].Part_name[:], *name)
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")
	var TempMBR structs.MBR

	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) && *delete == "" && *add == 0 {
			fmt.Println("[Error] El nombre de la partición ya está en uso!")
			return
		}
	}

	var EPartition = false
	var EPartitionStart int
	var ELimit int32
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
			ELimit = partition.Part_start + partition.Part_size
		}
	}

	if *delete == "full" {
		encontrada := false
		for i := range TempMBR.Mbr_particion {
			if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
				if bytes.Equal(TempMBR.Mbr_particion[i].Part_type[:], compareMBR.Mbr_particion[0].Part_type[:]) {
					TempMBR.Mbr_particion[i].Part_correlative = 0
					copy(TempMBR.Mbr_particion[i].Part_fit[:], "")
					copy(TempMBR.Mbr_particion[i].Part_id[:], "")
					copy(TempMBR.Mbr_particion[i].Part_name[:], "")
					copy(TempMBR.Mbr_particion[i].Part_type[:], "")
					copy(TempMBR.Mbr_particion[i].Part_status[:], "")
					encontrada = true
				}
				if bytes.Equal(TempMBR.Mbr_particion[i].Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
					end := TempMBR.Mbr_particion[i].Part_start + TempMBR.Mbr_particion[i].Part_size
					utilities.ConvertToZeros(filepath, int64(TempMBR.Mbr_particion[i].Part_start), int64(end))
					TempMBR.Mbr_particion[i].Part_correlative = 0
					copy(TempMBR.Mbr_particion[i].Part_fit[:], "")
					copy(TempMBR.Mbr_particion[i].Part_id[:], "")
					copy(TempMBR.Mbr_particion[i].Part_name[:], "")
					copy(TempMBR.Mbr_particion[i].Part_type[:], "")
					copy(TempMBR.Mbr_particion[i].Part_status[:], "")
					encontrada = true
				}
				break
			}

		}
		if !encontrada && EPartition {
			var x = 0
			for x < 1 {
				var TempEBR structs.EBR
				if err := utilities.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_s != 0 {
					if bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						copy(TempEBR.Part_mount[:], "0")
						copy(TempEBR.Part_fit[:], "")
						TempEBR.Part_s = 0
						copy(TempEBR.Part_name[:], "")
						if err := utilities.WriteObject(file, TempEBR, int64(EPartitionStart)); err != nil {
							return
						}
						encontrada = true
						break
					}
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					x = 1
				}
			}
		}

		if encontrada {
			fmt.Printf("Particion: %s eliminada\n", *name)
		} else {
			println("Error no se encontro la partción.")
		}

	} else if *add != 0 {
		for i := range TempMBR.Mbr_particion {
			if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
				if TempMBR.Mbr_particion[i].Part_size+int32(*add) < 0 {
					fmt.Println("El espacio de la partición no puede ser negativo")
					return
				}
				if i < len(TempMBR.Mbr_particion)-1 && TempMBR.Mbr_particion[i+1].Part_start < TempMBR.Mbr_particion[i].Part_start+TempMBR.Mbr_particion[i].Part_size+int32(*add) {
					if TempMBR.Mbr_particion[i+1].Part_start != 0 {
						fmt.Println("Al añadir espacio, se sobrepasa el start de la siguiente partición")
						return
					}
				}
				TempMBR.Mbr_particion[i].Part_size += int32(*add)
				if TempMBR.Mbr_particion[i].Part_size > TempMBR.Mbr_tamano {
					println("Supera el tamaño del disco")
					return
				}
				break
			}
		}

	} else {
		var count = 0
		var gap = int32(0)
		for i := 0; i < 4; i++ {

			if TempMBR.Mbr_particion[i].Part_size != 0 {
				count++
				gap = TempMBR.Mbr_particion[i].Part_start + TempMBR.Mbr_particion[i].Part_size
			}
		}

		for i := 0; i < 4; i++ {

			if TempMBR.Mbr_particion[i].Part_size == 0 {
				TempMBR.Mbr_particion[i].Part_size = int32(*size)

				if count == 0 {
					TempMBR.Mbr_particion[i].Part_start = int32(binary.Size(TempMBR))
				} else {
					TempMBR.Mbr_particion[i].Part_start = gap
				}

				suma := int32(*size) + int32(binary.Size(TempMBR))
				if suma > TempMBR.Mbr_tamano {
					println("La particion excede el tamaño del disco.")
					return
				}

				copy(TempMBR.Mbr_particion[i].Part_name[:], *name)
				copy(TempMBR.Mbr_particion[i].Part_fit[:], *fit)
				copy(TempMBR.Mbr_particion[i].Part_status[:], "0")
				copy(TempMBR.Mbr_particion[i].Part_type[:], *type_)
				TempMBR.Mbr_particion[i].Part_correlative = int32(count + 1)
				break
			}
		}

		if EPartition && *type_ == "l" {
			var x = 0
			for x < 1 {
				var TempEBR structs.EBR
				if err := utilities.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_s != 0 {
					var newEBR structs.EBR
					copy(newEBR.Part_mount[:], "0")
					copy(newEBR.Part_fit[:], *fit)
					newEBR.Part_start = int32(EPartitionStart) + 1
					newEBR.Part_s = TempEBR.Part_s
					newEBR.Part_next = int32(EPartitionStart) + int32(TempEBR.Part_s)
					copy(newEBR.Part_name[:], TempEBR.Part_name[:])

					if err := utilities.WriteObject(file, newEBR, int64(EPartitionStart)); err != nil {
						return
					}
					EPartitionStart = EPartitionStart + int(TempEBR.Part_s)
					structs.PrintEBR(newEBR)
				} else {
					var newEBR structs.EBR
					copy(newEBR.Part_mount[:], "0")
					copy(newEBR.Part_fit[:], *fit)
					newEBR.Part_start = int32(EPartitionStart) + 1
					newEBR.Part_s = int32(*size)
					newEBR.Part_next = -1
					copy(newEBR.Part_name[:], *name)

					if err := utilities.WriteObject(file, newEBR, int64(EPartitionStart)); err != nil {
						return
					}
					structs.PrintEBR(newEBR)
					suma := newEBR.Part_start + newEBR.Part_s
					if suma > ELimit {
						println("La particion logica supera el tamaño de la particion extendida")
						return
					}
					x = 1
				}
			}
			return
		}
		var Ecount = 0
		for _, partition := range TempMBR.Mbr_particion {
			if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
				if EPartition {
					Ecount += 1
				}
				if Ecount > 1 {
					println("No se puede tener mas de 1 particion extendida por disco!")
					return
				}
			}
		}

	}

	if err := utilities.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 structs.MBR
	if err := utilities.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}
	defer file.Close()
}

func ProcesarElMOUNT(input string, driveletter *string, name *string) {
	input = strings.ToLower(input)
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "driveletter":
			*driveletter = flagValue
		case "name":
			*name = flagValue
		default:
			fmt.Println("Bandera no encontrada: " + flagName)
		}
	}
}

func PaticionesDelMount(driveletter *string, name *string) {
	*driveletter = strings.ToUpper(*driveletter)
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + *driveletter + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR structs.MBR
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	encontrada := false

	var compareMBR structs.MBR
	copy(compareMBR.Mbr_particion[0].Part_name[:], *name)
	copy(compareMBR.Mbr_particion[0].Part_status[:], "1")
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")

	for i := 0; i < 4; i++ {
		if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
			if bytes.Equal(TempMBR.Mbr_particion[i].Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
				println("No es se pudo montar la particion extendida")
				return
			}
			if bytes.Equal(TempMBR.Mbr_particion[i].Part_status[:], compareMBR.Mbr_particion[0].Part_status[:]) {
				println("La particion ya esta montada")
				return
			}
			encontrada = true
			copy(TempMBR.Mbr_particion[i].Part_status[:], "1")
			ID := fmt.Sprintf("%s%d%s", *driveletter, TempMBR.Mbr_particion[i].Part_correlative, "04")
			copy(TempMBR.Mbr_particion[i].Part_id[:], ID)
			particionesMontadasListado += structs.GetPartition(TempMBR.Mbr_particion[i]) + "\n"
			break
		}
	}

	var EPartition = false
	var EPartitionStart int
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	if !encontrada && EPartition {
		for i := 0; i < 4; i++ {
			var x = 0
			for x < 1 {
				var TempEBR structs.EBR
				if err := utilities.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_s != 0 {
					if bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						if bytes.Equal(TempEBR.Part_mount[:], compareMBR.Mbr_particion[0].Part_status[:]) {
							println("La particion ya esta montada")
							return
						}
						copy(TempEBR.Part_mount[:], "1")
						encontrada = true
						if err := utilities.WriteObject(file, TempEBR, int64(EPartitionStart)); err != nil {
							return
						}
						particionesMontadasListado += structs.GetEBR(TempEBR) + "\n"
					}
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					x = 1
				}
			}
		}
	}
	if encontrada {
		println("Particion montada con exito")
		if err := utilities.WriteObject(file, TempMBR, 0); err != nil {
			return
		}

	} else {
		println("No se encontro la particion")
	}
	println(particionesMontadasListado)
}

func ProcesarElUNMOUNT(input string, id *string) {
	input = strings.ToLower(input)
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "id":
			*id = flagValue
		default:
			fmt.Println("Error archivo no encontrado: " + flagName)
		}
	}
}

func ParticionesDelUnMOINT(id *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	correlativo, err := strconv.ParseInt(string((*id)[len(*id)-3]), 10, 32)
	if err != nil {
		fmt.Println("Error al convertir la cadena a int32:", err)
		return
	}
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR structs.MBR
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var compareMBR structs.MBR
	compareMBR.Mbr_particion[0].Part_correlative = int32(correlativo)

	for i := 0; i < 4; i++ {

		if TempMBR.Mbr_particion[i].Part_correlative == compareMBR.Mbr_particion[0].Part_correlative {
			copy(TempMBR.Mbr_particion[i].Part_status[:], "0")
			break
		}
	}

	if err := utilities.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	structs.PrintMBR(TempMBR)
}

func ProcesarElMKFS(input string, id *string, type_ *string, fs *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "id":
			*id = flagValue
		case "type":
			*type_ = flagValue
		case "fs":
			*fs = flagValue
		default:
			fmt.Println("Archivo no encontrado: " + flagName)
		}

		if *type_ == "" {
			*type_ = "full"
		}
		if *fs == "" {
			*fs = "2fs"
		}
	}
}

func MKFS(id *string, type_ *string, fs *string) {
	fmt.Println("Id:", *id)
	fmt.Println("Type:", *type_)
	fmt.Println("Fs:", *fs)

	driveletter := string((*id)[0])
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR structs.MBR
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	structs.PrintMBR(TempMBR)

	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 {
			if strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), *id) {
				fmt.Println("Partición no encontrada")
				if strings.Contains(string(TempMBR.Mbr_particion[i].Part_status[:]), "1") {
					fmt.Println("Partición montada")
					index = i
				} else {
					fmt.Println("Partición no montada")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		structs.PrintPartition(TempMBR.Mbr_particion[index])
	} else {
		fmt.Println("Particion no encontrada")
		return
	}

	numerador := int32(TempMBR.Mbr_particion[index].Part_size - int32(binary.Size(structs.Superblock{})))
	denominador_base := int32(4 + int32(binary.Size(structs.Inode{})) + 3*int32(binary.Size(structs.Fileblock{})))
	var temp int32 = 0
	if *fs == "2fs" {
		temp = 0
	} else {
		temp = int32(binary.Size(structs.Journaling{}))
	}
	denominador := denominador_base + temp
	n := int32(numerador / denominador)

	fmt.Println("N:", n)
	var newSuperblock structs.Superblock
	newSuperblock.S_inodes_count = 0
	newSuperblock.S_blocks_count = 0

	newSuperblock.S_free_blocks_count = 3 * n
	newSuperblock.S_free_inodes_count = n

	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05")

	timeBytes := []byte(timeString)

	copy(newSuperblock.S_mtime[:], timeBytes)
	copy(newSuperblock.S_umtime[:], timeBytes)
	newSuperblock.S_mnt_count = 0

	if *fs == "2fs" {
		create_ext2(n, TempMBR.Mbr_particion[index], newSuperblock, timeString, file)
	} else {
		create_ext3(n, TempMBR.Mbr_particion[index], newSuperblock, timeString, file)
	}

	defer file.Close()

}

func create_ext2(n int32, partition structs.Partition, newSuperblock structs.Superblock, date string, file *os.File) {
	fmt.Println("N:", n)
	fmt.Println("Superblock:", newSuperblock)
	fmt.Println("Date:", date)

	newSuperblock.S_filesystem_type = 2
	newSuperblock.S_bm_inode_start = partition.Part_start + int32(binary.Size(structs.Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*int32(binary.Size(structs.Inode{}))
	newSuperblock.S_magic = 0xEF53
	newSuperblock.S_mnt_count = 1
	newSuperblock.S_inode_size = int32(binary.Size(structs.Inode{}))
	newSuperblock.S_block_size = int32(binary.Size(structs.Folderblock{}))

	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1
	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1

	for i := int32(0); i < n; i++ {
		err := utilities.WriteObject(file, byte(0), int64(newSuperblock.S_bm_inode_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	for i := int32(0); i < 3*n; i++ {
		err := utilities.WriteObject(file, byte(0), int64(newSuperblock.S_bm_block_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newInode structs.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		err := utilities.WriteObject(file, newInode, int64(newSuperblock.S_inode_start+i*int32(binary.Size(structs.Inode{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newFileblock structs.Fileblock
	for i := int32(0); i < 3*n; i++ {
		err := utilities.WriteObject(file, newFileblock, int64(newSuperblock.S_block_start+i*int32(binary.Size(structs.Fileblock{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var Inode0 structs.Inode
	Inode0.I_uid = 1
	Inode0.I_gid = 1
	Inode0.I_size = 0
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	copy(Inode0.I_type[:], "1")
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0

	var Folderblock0 structs.Folderblock
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[2].B_name[:], "users.txt")
	Folderblock0.B_content[2].B_inodo = 1

	var Inode1 structs.Inode
	Inode1.I_uid = 1
	Inode1.I_gid = 1
	Inode1.I_size = int32(binary.Size(structs.Folderblock{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	copy(Inode1.I_type[:], "1")
	copy(Inode1.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}

	Inode1.I_block[0] = 1

	data := "1,G,root\n1,U,root,root,123\n"
	var Fileblock1 structs.Fileblock
	copy(Fileblock1.B_content[:], data)

	newSuperblock.S_inodes_count = int32(2)
	newSuperblock.S_blocks_count = int32(1)
	newSuperblock.S_fist_ino = int32(0)
	newSuperblock.S_first_blo = int32(1)

	err := utilities.WriteObject(file, newSuperblock, int64(partition.Part_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = utilities.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// write bitmap blocks
	err = utilities.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Inode 0:", int64(newSuperblock.S_inode_start))
	fmt.Println("Inode 1:", int64(newSuperblock.S_inode_start+int32(binary.Size(structs.Inode{}))))

	// write inodes
	err = utilities.WriteObject(file, Inode0, int64(newSuperblock.S_inode_start)) //Inode 0
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities.WriteObject(file, Inode1, int64(newSuperblock.S_inode_start+int32(binary.Size(structs.Inode{})))) //Inode 1
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// write blocks
	err = utilities.WriteObject(file, Folderblock0, int64(newSuperblock.S_block_start)) //Bloque 0
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities.WriteObject(file, Fileblock1, int64(newSuperblock.S_block_start+int32(binary.Size(structs.Fileblock{})))) //Bloque 1

	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func create_ext3(n int32, partition structs.Partition, newSuperblock structs.Superblock, date string, file *os.File) {
	fmt.Println("N:", n)
	fmt.Println("Superblock:", newSuperblock)
	fmt.Println("Date:", date)

	newSuperblock.S_filesystem_type = 3
	newSuperblock.S_bm_inode_start = partition.Part_start + int32(binary.Size(structs.Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*int32(binary.Size(structs.Inode{}))

	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1
	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1

	var err error

	for i := int32(0); i < n; i++ {
		err = utilities.WriteObject(file, byte(0), int64(newSuperblock.S_bm_inode_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	for i := int32(0); i < 3*n; i++ {
		err = utilities.WriteObject(file, byte(0), int64(newSuperblock.S_bm_block_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newInode structs.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		err = utilities.WriteObject(file, newInode, int64(newSuperblock.S_inode_start+i*int32(binary.Size(structs.Inode{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newFileblock structs.Fileblock
	for i := int32(0); i < 3*n; i++ {
		err = utilities.WriteObject(file, newFileblock, int64(newSuperblock.S_block_start+i*int32(binary.Size(structs.Fileblock{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var Inode0 structs.Inode
	Inode0.I_uid = 1
	Inode0.I_gid = 1
	Inode0.I_size = 0
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	copy(Inode0.I_type[:], "1")
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0

	var Folderblock0 structs.Folderblock
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[2].B_name[:], "users.txt")
	Folderblock0.B_content[2].B_inodo = 1

	var Inode1 structs.Inode //Inode 1
	Inode1.I_uid = 1
	Inode1.I_gid = 1
	Inode1.I_size = int32(binary.Size(structs.Folderblock{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	copy(Inode1.I_type[:], "1")
	copy(Inode1.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}

	Inode1.I_block[0] = 1

	data := "1,G,root\n1,U,root,root,123\n"
	var Fileblock1 structs.Fileblock
	copy(Fileblock1.B_content[:], data)

	var journal structs.Journaling
	journal.Size = 50
	journal.Ultimo = -1

	err = utilities.WriteObject(file, newSuperblock, int64(partition.Part_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = utilities.WriteObject(file, journal, int64(newSuperblock.S_block_start+(3*n)*int32(binary.Size(structs.Fileblock{}))))
	if err != nil {
		fmt.Println("Error al escribir el jornaling en el disco:", err)
		return
	}

	err = utilities.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Inode 0:", int64(newSuperblock.S_inode_start))
	fmt.Println("Inode 1:", int64(newSuperblock.S_inode_start+int32(binary.Size(structs.Inode{}))))

	err = utilities.WriteObject(file, Inode0, int64(newSuperblock.S_inode_start)) //Inode 0
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities.WriteObject(file, Inode1, int64(newSuperblock.S_inode_start+int32(binary.Size(structs.Inode{})))) //Inode 1
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = utilities.WriteObject(file, Folderblock0, int64(newSuperblock.S_block_start)) //Bloque 0
	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = utilities.WriteObject(file, Fileblock1, int64(newSuperblock.S_block_start+int32(binary.Size(structs.Fileblock{})))) //Bloque 1

	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func ProcessExecute(input string, path *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		*path = match[2]
	}
}

func GenerateUniqueID() int {
	currentTime := time.Now()
	randomNumber := rand.Intn(10000)
	uniqueID := currentTime.UnixNano() * int64(randomNumber) % (1 << 31)
	uniqueID = int64(math.Abs(float64(uniqueID)))
	return int(uniqueID)
}

func ValidarElDriveletter(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z]$`).MatchString(str)
}
