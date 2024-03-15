package analyzer

import (
	functions "P1/Functions"
	utilities "P1/Utilities"
	"bufio"
	"flag"
	"fmt"
	"strings"
)

func Command(input string) {

	// Verificar si el input está vacío
	if input == "" {
		return // No hacer nada si el input está vacío
	}

	comando := input
	input = strings.ToLower(input)
	switch {
	case strings.HasPrefix(input, "mkdisk"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoMKDISK(comando)

	case strings.HasPrefix(input, "rmdisk"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoRMDISK(comando)

	case strings.HasPrefix(input, "fdisk"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoFDISKC(comando)

	case strings.HasPrefix(input, "mount"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoMOUNT(comando)

	case strings.HasPrefix(input, "unmount"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoUNMOUNT(comando)

	case strings.HasPrefix(input, "mkfs"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoMKFS(comando)

	case strings.HasPrefix(input, "login"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoLOGIN(comando)

	case strings.HasPrefix(input, "logout"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoLOGOUT()

	case strings.HasPrefix(input, "mkgrp"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoMKGRP(comando)

	case strings.HasPrefix(input, "rmgrp"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoRMGRP(comando)

	case strings.HasPrefix(input, "mkusr"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoMKUSR(comando)

	case strings.HasPrefix(input, "rmusr"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoRMUSR(comando)

	case strings.HasPrefix(input, "mkfile"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoMKFILE(comando)

	case strings.HasPrefix(input, "cat"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoCAT(comando)

	case strings.HasPrefix(input, "remove"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoREMOVE(comando)

	case strings.HasPrefix(input, "edit"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoEDIT(comando)

	case strings.HasPrefix(input, "rename"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoRENAME(comando)

	case strings.HasPrefix(input, "mkdir"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoMKDIR(comando)

	case strings.HasPrefix(input, "copy"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoCOPY(comando)

	case strings.HasPrefix(input, "move"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoMOVE(comando)

	case strings.HasPrefix(input, "find"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoFIND(comando)

	case strings.HasPrefix(input, "chown"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoCHOWN(comando)

	case strings.HasPrefix(input, "chgrp"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoCHGRP(comando)

	case strings.HasPrefix(input, "chmod"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoCHMOD(comando)

	case strings.HasPrefix(input, "pause"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoPAUSE()

	case strings.HasPrefix(input, "exec"):
		comandoEXECUTE(comando)

	case strings.HasPrefix(input, "rep"):
		fmt.Println("******Kevin Secaida******")
		fmt.Println("Ingrese un comando" + comando)
		comandoREP(comando)

	case strings.HasPrefix(input, "#"):
		// Print the comment and continue with the program execution
		fmt.Println(input)
	default:
		fmt.Println("Comando no reconocido:", comando)
	}
}

var (
	size        = flag.Int("size", 0, "Tamaño")
	fit         = flag.String("fit", "", "Ajuste")
	unit        = flag.String("unit", "", "Unidad")
	type_       = flag.String("type", "", "Tipo")
	driveletter = flag.String("driveletter", "", "Busqueda")
	name        = flag.String("name", "", "Nombre")
	delete      = flag.String("delete", "", "Eliminar")
	add         = flag.Int("add", 0, "Añadir/Quitar")
	path        = flag.String("path", "", "Directorio")
	id          = flag.String("id", "", "ID")
	fs          = flag.String("fs", "", "FDISK")
	ruta        = flag.String("ruta", "", "Ruta")
	user        = flag.String("user", "", "Usuario")
	pass        = flag.String("pass", "", "Password")
	grp         = flag.String("grp", "", "Group")
	r           = flag.Bool("r", false, "Rewrite")
	cont        = flag.String("cont", "", "Cont")
	destino     = flag.String("destino", "", "Destino")
	ugo         = flag.String("ugo", "", "UGO")
	file        = flag.String("file", "", "File to process")
)

// Comandos.

func comandoMKDISK(input string) {

	flag.Parse()
	functions.MKDISK(input, size, fit, unit)

	// validate size > 0
	if *size <= 0 {
		fmt.Println("Error Size deben ser mayor 0.")
		return
	}

	// validate fit equals to b/w/f
	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error Fit debe ser bf|ff|wf")
		return
	}

	// validate unit equals to k/m
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error Unit debe ser k|m")
		return
	}

	// Create the file
	functions.CreateBinFile(size, fit, unit)
	*size = 0
	*fit = ""
	*unit = ""
}

func comandoRMDISK(input string) {
	flag.Parse()
	functions.RMDISK(input, driveletter)
	// validate driveletter be una letra and not empty
	if !functions.ValidDriveLetter(*driveletter) {
		fmt.Println("Error DriveLetter deben ser una letra")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error DriveLetter no puede estar vacio..")
		return
	}

	functions.DeleteBinFile(driveletter)
	*driveletter = ""
}

func comandoFDISKC(input string) {
	flag.Parse()
	functions.FDISK(input, size, driveletter, name, unit, type_, fit, delete, add, path)

	//Obligatorio cuando no existe la particion
	// validate size > 0
	if *size <= 0 && *delete != "full" && *add == 0 {
		fmt.Println("Error Size deben ser mayor 0.")
		return
	}

	// validate driveletter be una letra and not empty
	if !functions.ValidDriveLetter(*driveletter) {
		fmt.Println("Error DriveLetter deben ser una letra")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error DriveLetter no puede estar vacio..")
		return
	}

	// validate fit equals to b/w/f
	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error Fit deben ser (BF/FF/WF)")
		return
	}

	// validate unit equals to b/k/m
	if *unit != "b" && *unit != "k" && *unit != "m" {
		fmt.Println("Error Unit deben ser (B/K/M)")
		return
	}

	//println("ADD", *add)
	// validate type equals to P/E/L
	if *type_ != "p" && *type_ != "e" && *type_ != "l" && *delete != "full" && *add == 0 {
		fmt.Println("Error Type deben ser (P/E/L)")
		return
	}

	if *delete != "" {
		if *delete != "full" {
			fmt.Println("Error Delete deben ser full")
			return
		}
		if *name == "" && *path == "" {
			println("Error you need path and name to delete")
			return
		}
	}

	functions.CRUD_Partitions(size, driveletter, name, unit, type_, fit, delete, add, path)
	*size = 0
	*driveletter = ""
	*name = ""
	*unit = ""
	*type_ = ""
	*fit = ""
	*delete = ""
	*add = 0
	*path = ""
}

func comandoMOUNT(input string) {
	flag.Parse()
	functions.MOUNT(input, driveletter, name)

	// validate driveletter be una letra and not empty
	if !functions.ValidDriveLetter(*driveletter) {
		fmt.Println("Error DriveLetter deben ser una letra.")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error DriveLetter no puede estar vacio...")
		return
	}

	functions.MountPartition(driveletter, name)
	*driveletter = ""
	*name = ""
}

func comandoUNMOUNT(input string) {
	flag.Parse()
	functions.UNMOUNT(input, id)

	if *id == "" {
		println("Error Id es un campo obligatorio.")
	}

	functions.UNMOUNT_Partition(id)
	*id = ""
}

func comandoMKFS(input string) {
	flag.Parse()
	functions.MKFSS(input, id, type_, fs)

	if *id == "" {
		println("Error id no puede estar vacio..")
	}

	if *fs != "2fs" && *fs != "3fs" {
		println("Error fs deben ser 2fs o 3fs")
	}

	functions.MKFS(id, type_, fs)
	*id = ""
	*type_ = ""
	*fs = ""
}

//ADMINISTRACION DE USUARIOS

func comandoLOGIN(input string) {
	flag.Parse()
	functions.LLOGIN(input, user, pass, id)

	if *user == "" || *pass == "" || *id == "" {
		println("Error parametros incompletos")
	}

	functions.LOGIN(user, pass, id)

	*user = ""
	*pass = ""
	*id = ""
}

func comandoLOGOUT() {
	functions.LLOGOUT()
}

func comandoMKGRP(input string) {
	flag.Parse()
	functions.MMKGRP(input, name)

	if *name == "" {
		println("Error el campo name no puede estar vacio.")
		return
	}

	functions.MKGRP(name)
	*name = ""
}

func comandoRMGRP(input string) {
	flag.Parse()
	functions.RRMGRP(input, name)

	if *name == "" {
		println("Error el campo name no puede estar vacio.")
		return
	}

	functions.RMGRP(name)
	*name = ""
}

func comandoMKUSR(input string) {
	flag.Parse()
	functions.MMKUSR(input, user, pass, grp)

	if len(*user) > 10 {
		println("Error user no puede ser mayor a 10 caracteres")
		return
	}
	if len(*pass) > 10 {
		println("Error password no puede ser mayor a 10 caracteres")
		return
	}
	if len(*grp) > 10 {
		println("Error grupo no puede ser mayor a 10 caracteres")
		return
	}

	if *user == "" || *pass == "" || *grp == "" {
		println("Error parametros incompletos")
		return
	}

	functions.MKUSR(user, pass, grp)

	*user = ""
	*pass = ""
	*grp = ""
}

func comandoRMUSR(input string) {
	flag.Parse()
	functions.RRMUSR(input, user)

	if *user == "" {
		println("Error user no puede estar vacio..")
		return
	}

	functions.RMUSR(user)

	*user = ""
}

func comandoCHGRP(input string) {
	flag.Parse()
	functions.CCHGRP(input, user, grp)
}

//ADMINISTRACION DE CARPETAS

func comandoMKDIR(input string) {
	flag.Parse()
	functions.MMKDIR(input, path, r)

	if *path == "" {
		println("Error path no puede estar vacio..")
		return
	}

	fmt.Println("Path: " + *path)
	fmt.Print("r: ")
	fmt.Println(*r)

	*path = ""
	*r = false
}

func comandoMKFILE(input string) {
	flag.Parse()
	functions.ProcessMKFILE(input, path, r, size, cont)
}

func comandoCAT(input string) {
	flag.Parse()
	functions.ProcessCAT(input, file)
}

func comandoREMOVE(input string) {
	flag.Parse()
	functions.ProcessREMOVE(input, path)
}

func comandoEDIT(input string) {
	flag.Parse()
	functions.ProcessEDIT(input, path, cont)
}

func comandoRENAME(input string) {
	flag.Parse()
	functions.ProcessRENAME(input, path, name)
}

func comandoCOPY(input string) {
	flag.Parse()
	functions.ProcessCOPY(input, path, destino)
}

func comandoMOVE(input string) {
	flag.Parse()
	functions.ProcessMOVE(input, path, destino)
}

func comandoFIND(input string) {
	flag.Parse()
	functions.ProcessFIND(input, path, destino)
}

func comandoCHOWN(input string) {
	flag.Parse()
	functions.ProcessCHOWN(input, path, user, r)
}

func comandoCHMOD(input string) {
	flag.Parse()
	functions.ProcessCHMOD(input, path, ugo, r)
}

//   COMANDOS AUXILIARES

func comandoPAUSE() {
	fmt.Println("Presione ENTER tecla para continuar...")
	fmt.Scanln()
}

func comandoEXECUTE(input string) {
	flag.Parse()
	functions.Exxecute(input, path)
	if *path == "" {
		fmt.Println("Error Path no puede estar vacio..")
		return
	}
	// Open bin file
	file, err := utilities.OpenFile(*path)
	if err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	// Crea un nuevo scanner para leer el archivo
	scanner := bufio.NewScanner(file)

	// Itera sobre cada línea del archivo
	for scanner.Scan() {
		linea := scanner.Text() // Lee la línea actual
		//fmt.Println(linea)
		Command(linea)
	}

	// Verifica si hubo algún error durante la lectura
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
	*path = ""
}

//         REPORTES

func comandoREP(input string) {
	flag.Parse()
	functions.RREP(input, name, path, id, ruta)

	if *name == "" || *path == "" || *id == "" {
		println("Error parametros incompletos. ")
		return
	}

	functions.GenerateReports(name, path, id, ruta)
}
