package functions

import (
	"P1/Global"
	structs "P1/Structs"
	utilities "P1/Utilities"
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

var (
	session      = false
	usuario      = Global.UserInfo{}
	groupCounter = 1
	userCounter  = 1
	blockIndex   = 0
	searchIndex  = 0
	letra        = ""
	ID           = ""
)

// Login
func ProcessLOGIN(input string, user *string, pass *string, id *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user":
			*user = flagValue
		case "pass":
			*pass = flagValue
		case "id":
			*id = flagValue
		default:
			fmt.Println("Error archivo no encontrado: " + flagName)
		}
	}
}

func LOGIN(user *string, pass *string, id *string) {

	letra = string((*id)[0])
	ID = *id

	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	var TempMBR structs.MBR
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error al leer MBR:", err)
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
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock structs.Superblock
	if err := utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error al leer superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error al leer inode:", err)
		return
	}

	var Fileblock structs.Fileblock
	if err := utilities.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error al leer Fileblock:", err)
		return
	}
	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	userFound := false
	for _, line := range lines {
		items := strings.Split(line, ",")
		if len(items) > 3 {
			if *user == items[len(items)-2] {
				userFound = true
				usuario.Nombre = items[len(items)-2]
				usuario.ID = items[0]
				session = true
				break
			}
		}
	}

	if !userFound {
		searchIndex++
		if searchIndex <= blockIndex {
			LOGIN(user, pass, id)
		} else {
			fmt.Println("Usuario no encontrado")
			searchIndex = 0
			return
		}
	} else {
		Global.PrintUser(usuario)
		searchIndex = 0
		return
	}
}

// Logout
func ProcessLOGOUT() {
	if session {
		println("Se ha cerrado la sesion")
		session = false
		searchIndex = 0
		return
	}
	println("Error no hay una sesion activa")
}

// MKGROUP
func ProcessMKGRP(input string, name *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "name":
				*name = flagValue
			default:
				fmt.Println("Error archivo no encontrado: " + flagName)
			}
		}
	} else {
		println("Solo el usuario root puede acceder a este comando")
		return
	}
}

// MKGRP
func MKGRP(name *string) {

	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al ejecutar el disco:", err)
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
		fmt.Println("Partición no encontrada")
		return
	}

	var tempSuperblock structs.Superblock
	if err := utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error al leer el superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error al leer el inode:", err)
		return
	}

	var Fileblock structs.Fileblock
	if err := utilities.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(blockIndex))); err != nil {
		fmt.Println("Error al leer el Fileblock:", err)
		return
	}
	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) == 3 {
			if *name == items[2] {
				println("Nombres repetidos.")
				return
			}
		}
	}

	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	nuevoGrupo := fmt.Sprintf("%d,G,%s\n", groupCounter, *name)
	newContent := currentContent + nuevoGrupo

	if len(newContent) > len(Fileblock.B_content) {
		if blockIndex > int(len(crrInode.I_block)) {
			fmt.Println("No hay mas bloques disponibles")
			return
		}
		blockIndex++
		var NEWFileblock structs.Fileblock
		if err := utilities.WriteObject(file, &NEWFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(blockIndex))); err != nil {
			fmt.Println("Error al leer Fileblock:", err)
			return
		}

		crrInode.I_block[blockIndex] = 1

		if err := utilities.WriteObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
			fmt.Println("Error writing Inode to disk:", err)
			return
		}

		MKGRP(name)
	} else {
		copy(Fileblock.B_content[:], newContent)

		if err := utilities.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(blockIndex))); err != nil {
			fmt.Println("Error writing Fileblock to disk:", err)
			return
		}

		println("ACTUALIZACION")
		data := string(Fileblock.B_content[:])
		lines := strings.Split(data, "\n")

		for _, line := range lines {
			fmt.Println(line)
		}
	}
}

// RMGRP
func ProcessRMGRP(input string, name *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "name":
				*name = flagValue
			default:
				fmt.Println("Error archivo no encontrado: " + flagName)
			}
		}
	} else {
		println("Solo el usuario root puede acceder a este comando")
		return
	}
}

func RMGRP(name *string) {
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir el disco:", err)
		return
	}
	defer file.Close()

	// Leer el MBR del disco
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
		fmt.Println("Partición no encontrada")
		return
	}

	var tempSuperblock structs.Superblock
	if err := utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error al leer el superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error al leer el inode:", err)
		return
	}

	var Fileblock structs.Fileblock
	if err := utilities.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error al leer el Fileblock:", err)
		return
	}

	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	lines := strings.Split(currentContent, "\n")
	deleted := false
	for i, line := range lines {
		if strings.Contains(line, *name) {
			lines[i] = "0,G," + *name
			deleted = true
			break
		}
	}

	if !deleted {
		searchIndex++
		if searchIndex > blockIndex {
			fmt.Println("Grupo no encontrado.")
			searchIndex = 0
			return
		}
		RMGRP(name)

	}

	newContent := strings.Join(lines, "\n")
	copy(Fileblock.B_content[:], newContent)

	if deleted {
		if err := utilities.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(searchIndex))); err != nil {
			fmt.Println("Error al escribir el Fileblock en el disco:", err)
			return
		}

		currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
		lines := strings.Split(currentContent, "\n")
		for i := range lines {
			println(lines[i])
		}

		searchIndex = 0
	}
}

func ProcessMKUSR(input string, user *string, pass *string, grp *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				*user = flagValue
			case "pass":
				*pass = flagValue
			case "grp":
				*grp = flagValue
			default:
				fmt.Println("Error: " + flagName)
			}
		}
	} else {
		println("Solo el usuario root puede acceder a este comando")
		return
	}
}

func MKUSR(user *string, pass *string, grp *string) {
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir el disco:", err)
		return
	}
	defer file.Close()

	var TempMBR structs.MBR
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error al leer MBR:", err)
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
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock structs.Superblock
	if err := utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error al leer superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error al leer inode:", err)
		return
	}

	var Fileblock structs.Fileblock
	if err := utilities.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(blockIndex))); err != nil {
		fmt.Println("Error al leer Fileblock:", err)
		return
	}

	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	searchIndex = 0
	var nuevoUsuario = BuscarGrupo(user, pass, grp)
	if nuevoUsuario == "" {
		fmt.Println("Error No se encontro el grupo")
		return
	}
	newContent := currentContent + nuevoUsuario

	if len(newContent) > len(Fileblock.B_content) {
		if blockIndex > int(len(crrInode.I_block)) {
			fmt.Println("Error no hay mas bloques disponibles")
			return
		}
		blockIndex++
		var NEWFileblock structs.Fileblock
		copy(NEWFileblock.B_content[:], nuevoUsuario)
		if err := utilities.WriteObject(file, &NEWFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(blockIndex))); err != nil {
			fmt.Println("Error al leer Fileblock:", err)
			return
		}

		println("MKUSR EXITOSO")
		data := string(NEWFileblock.B_content[:])
		lines := strings.Split(data, "\n")

		for _, line := range lines {
			fmt.Println(line)
		}

		crrInode.I_block[blockIndex] = 1

		if err := utilities.WriteObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
			fmt.Println("Error writing Inode to disk:", err)
			return
		}
		searchIndex = 0

	} else {
		println("MKUSR EXITOSO")
		copy(Fileblock.B_content[:], newContent)

		if err := utilities.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(blockIndex))); err != nil {
			fmt.Println("Error writing Fileblock to disk:", err)
			return
		}

		data := string(Fileblock.B_content[:])
		lines := strings.Split(data, "\n")

		for _, line := range lines {
			fmt.Println(line)
		}
		searchIndex = 0
	}
}

func ProcessRMUSR(input string, user *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				*user = flagValue
			default:
				fmt.Println("Error parametro no encontrado: " + flagName)
			}
		}
	} else {
		println("Solo el usuario root puede acceder a este comando")
		return
	}
}

func RMUSR(user *string) {
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir el disco:", err)
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
		fmt.Println("Partición no encontrada")
		return
	}

	var tempSuperblock structs.Superblock
	if err := utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error al leer el superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error al leer el inode:", err)
		return
	}

	var Fileblock structs.Fileblock
	if err := utilities.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error al leer el Fileblock:", err)
		return
	}

	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		items := strings.Split(line, ",")
		if len(items) > 3 {
			if *user == items[len(items)-2] {
				items[0] = "0" // Cambiar el estado del usuario a 0
				newLine := strings.Join(items, ",")
				copy(Fileblock.B_content[:], []byte(newLine))
				if err := utilities.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(searchIndex))); err != nil {
					fmt.Println("Error al escribir el Fileblock en el disco:", err)
					return
				}
				return
			}
		}
	}

	searchIndex++
	if searchIndex <= blockIndex {
		RMUSR(user)
	} else {
		fmt.Println("Usuario no encontrado")
	}
}

func ProcessCHGRP(input string, user *string, grp *string) {
}

func CHGRP(user *string, grp *string) {
}

func ImprimirBloques() {
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error al abrir el disco:", err)
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

	var tempSuperblock structs.Superblock
	if err := utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error al leer superblock:", err)
		return
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error al leer inode:", err)
		return
	}

	var Fileblock structs.Fileblock
	if err := utilities.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error al leer Fileblock:", err)
		return
	}
	fmt.Println("Fileblock " + fmt.Sprint(searchIndex))
	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		fmt.Println(line)
	}

	if searchIndex < blockIndex {
		searchIndex++
		ImprimirBloques()
	} else {
		searchIndex = 0
	}
}

func BuscarGrupo(user *string, pass *string, grp *string) string {
	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return ""
	}
	defer file.Close()

	var TempMBR structs.MBR
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error al leer MBR:", err)
		return ""
	}

	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return ""
	}

	var tempSuperblock structs.Superblock
	if err := utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error al leer superblock:", err)
		return ""
	}

	indexInode := int32(1)
	var crrInode structs.Inode
	if err := utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs.Inode{})))); err != nil {
		fmt.Println("Error al leer inode:", err)
		return ""
	}

	var Fileblock structs.Fileblock
	if err := utilities.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))+crrInode.I_block[0]*int32(binary.Size(structs.Fileblock{}))*int32(searchIndex))); err != nil {
		fmt.Println("Error al leer Fileblock:", err)
		return ""
	}
	data := string(Fileblock.B_content[:])
	lines := strings.Split(data, "\n")

	groupFound := false
	var newUserLine string
	for _, line := range lines {
		items := strings.Split(line, ",")
		if len(items) == 3 {
			if *grp == items[2] {
				groupFound = true
				newUserLine = fmt.Sprintf("%d,G,%s,%s,%s\n", userCounter, *grp, *user, *pass)
				userCounter++
				break
			}
		}
	}

	if !groupFound {
		searchIndex++
		if searchIndex <= blockIndex {
			return BuscarGrupo(user, pass, grp)
		}
	} else {
		return newUserLine
	}
	return ""
}
