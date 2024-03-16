package analyzer

import (
	functions "P1/Functions"
	utilities "P1/Utilities"
	"bufio"
	"flag"
	"fmt"
	"strings"
)

func Comandos(input string) {

	if input == "" {
		return
	}

	comando := input
	input = strings.ToLower(input)
	switch {
	case strings.HasPrefix(input, "mkdisk"):
		fmt.Printf("\nComando: %s\n", comando)
		MMkdisk(comando)
	case strings.HasPrefix(input, "rmdisk"):
		fmt.Printf("\nComando: %s\n", comando)
		RRmdisk(comando)
	case strings.HasPrefix(input, "fdisk"):
		fmt.Printf("\nComando: %s\n", comando)
		FFdisk(comando)
	case strings.HasPrefix(input, "mount"):
		fmt.Printf("\nComando: %s\n", comando)
		MMount(comando)
	case strings.HasPrefix(input, "unmount"):
		fmt.Printf("\nComando: %s\n", comando)
		UUnMount(comando)
	case strings.HasPrefix(input, "mkfs"):
		fmt.Printf("\nComando: %s\n", comando)
		MMkfs(comando)
	case strings.HasPrefix(input, "login"):
		fmt.Printf("\nComando: %s\n", comando)
		LLogin(comando)
	case strings.HasPrefix(input, "logout"):
		fmt.Printf("\nComando: %s\n", comando)
		LLogut()
	case strings.HasPrefix(input, "mkgrp"):
		fmt.Printf("\nComando: %s\n", comando)
		MMkgrp(comando)
	case strings.HasPrefix(input, "rmgrp"):
		fmt.Printf("\nComando: %s\n", comando)
		RRmgrp(comando)
	case strings.HasPrefix(input, "mkusr"):
		fmt.Printf("\nComando: %s\n", comando)
		MMkusr(comando)
	case strings.HasPrefix(input, "rmusr"):
		fmt.Printf("\nComando: %s\n", comando)
		RRmusr(comando)
	case strings.HasPrefix(input, "mkfile"):
		fmt.Printf("\nComando: %s\n", comando)
		MMkfile(comando)
	case strings.HasPrefix(input, "cat"):
		fmt.Printf("\nComando: %s\n", comando)
		CCat(comando)
	case strings.HasPrefix(input, "remove"):
		fmt.Printf("\nComando: %s\n", comando)
		Remove(comando)
	case strings.HasPrefix(input, "edit"):
		fmt.Printf("\nComando: %s\n", comando)
		edit(comando)
	case strings.HasPrefix(input, "rename"):
		fmt.Printf("\nComando: %s\n", comando)
		rename(comando)
	case strings.HasPrefix(input, "mkdir"):
		fmt.Printf("\nComando: %s\n", comando)
		MMkdir(comando)
	case strings.HasPrefix(input, "copy"):
		fmt.Printf("\nComando: %s\n", comando)
		copy(comando)
	case strings.HasPrefix(input, "move"):
		fmt.Printf("\nComando: %s\n", comando)
		move(comando)
	case strings.HasPrefix(input, "find"):
		fmt.Printf("\nComando: %s\n", comando)
		Ffind(comando)
	case strings.HasPrefix(input, "chown"):
		fmt.Printf("\nComando: %s\n", comando)
		comando_chown(comando)
	case strings.HasPrefix(input, "chgrp"):
		fmt.Printf("\nComando: %s\n", comando)
		comando_chgrp(comando)
	case strings.HasPrefix(input, "chmod"):
		fmt.Printf("\nComando: %s\n", comando)
		comando_Chmod(comando)
	case strings.HasPrefix(input, "pause"):
		fmt.Printf("\nComando: %s\n", comando)
		pause()
	case strings.HasPrefix(input, "execute"):
		execute(comando)
	case strings.HasPrefix(input, "rep"):
		fmt.Printf("\nComando: %s\n", comando)
		rep(comando)
	case strings.HasPrefix(input, "#"):
		fmt.Printf("\nComentario: %s\n", comando)
	default:
		fmt.Println("Comando no valido: ", comando)
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

// Administración de discos.
func MMkdisk(input string) {
	fmt.Println("************ MKDISK ************")
	flag.Parse()
	functions.ProcesarElMKDISK(input, size, fit, unit)
	// Para size
	if *size <= 0 {
		fmt.Println("Error El tamaño debe ser mayor que 0")
		return
	}
	// Para fit
	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error El ajuste debe ser (bf/ff/wf)")
		return
	}

	// Validar si es kilobytes o megabytes
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error La unidad debe ser kilobytes o megabytes (k/m)")
		return
	}

	// Crear el archivo
	functions.AsignarAlfabeto(size, fit, unit)
	*size = 0
	*fit = ""
	*unit = ""
}

func RRmdisk(input string) {
	fmt.Println("************ RMDISK ************")
	flag.Parse()
	functions.ProcessRMDISK(input, driveletter)
	if !functions.ValidarElDriveletter(*driveletter) {
		fmt.Println("Error DriveLetter debe ser una letra")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error DriveLetter no puede estar vacía")
		return
	}
	functions.Eliminar_ArchivoBin(driveletter)
	*driveletter = ""
}

func FFdisk(input string) {
	fmt.Println("************ FDISK ************")
	flag.Parse()
	functions.GestionarFDISK(input, size, driveletter, name, unit, type_, fit, delete, add, path)

	if *size <= 0 && *delete != "full" && *add == 0 {
		fmt.Println("Error El tamaño debe ser mayor que 0.")
		return
	}

	if !functions.ValidarElDriveletter(*driveletter) {
		fmt.Println("Error DriveLetter debe ser una letra")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error DriveLetter no puede estar vacío")
		return
	}

	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error El ajuste debe ser (BF/FF/WF)")
		return
	}

	if *unit != "b" && *unit != "k" && *unit != "m" {
		fmt.Println("Error La unidad debe ser (B/K/M)")
		return
	}

	if *type_ != "p" && *type_ != "e" && *type_ != "l" && *delete != "full" && *add == 0 {
		fmt.Println("Error El tipo debe ser (P/E/L)")
		return
	}

	if *delete != "" {
		if *delete != "full" {
			fmt.Println("Error Eliminar debe estar lleno")
			return
		}
		if *name == "" && *path == "" {
			println("Error ecesitas ruta y nombre para eliminar")
			return
		}
	}

	functions.CRUDdeParticiones(size, driveletter, name, unit, type_, fit, delete, add, path)
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

func MMount(input string) {
	fmt.Println("************ FDISK ************")
	flag.Parse()
	functions.ProcesarElMOUNT(input, driveletter, name)

	if !functions.ValidarElDriveletter(*driveletter) {
		fmt.Println("Error DriveLetter debe ser una letra")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error DriveLetter no puede estar vacío")
		return
	}
	fmt.Println("Ejecutando MOUNT...")
	functions.PaticionesDelMount(driveletter, name)
	*driveletter = ""
	*name = ""
}

func UUnMount(input string) {
	fmt.Println("************ FDISK ************")
	flag.Parse()
	functions.ProcesarElUNMOUNT(input, id)

	if *id == "" {
		println("Error Id es un campo obligatorio")
	}
	fmt.Println("Ejecutando UNMOUNT...")
	functions.ParticionesDelUnMOINT(id)
	*id = ""
}

func MMkfs(input string) {
	fmt.Println("************ MKFS ************")
	flag.Parse()
	functions.ProcesarElMKFS(input, id, type_, fs)

	if *id == "" {
		println("Error la identificación no puede estar vacía")
	}

	if *fs != "2fs" && *fs != "3fs" {
		println("Error fs debe ser 2fs o 3fs")
	}
	fmt.Println("Ejecutando MKFS...")
	functions.MKFS(id, type_, fs)
	*id = ""
	*type_ = ""
	*fs = ""
}

// Administración de usuarios y grupos.
func LLogin(input string) {
	flag.Parse()
	functions.ProcessLOGIN(input, user, pass, id)

	if *user == "" || *pass == "" || *id == "" {
		println("Error campos incompletos")
	}

	functions.LOGIN(user, pass, id)

	*user = ""
	*pass = ""
	*id = ""
}

func LLogut() {
	functions.ProcessLOGOUT()
}

func MMkgrp(input string) {
	flag.Parse()
	functions.ProcessMKGRP(input, name)

	if *name == "" {
		println("Error el campo name no puede estar vacio")
		return
	}

	functions.MKGRP(name)
	*name = ""
}

func RRmgrp(input string) {
	flag.Parse()
	functions.ProcessMKGRP(input, name)

	if *name == "" {
		println("Error el campo name no puede estar vacio")
		return
	}

	functions.RMGRP(name)
	*name = ""
}

func MMkusr(input string) {
	flag.Parse()
	functions.ProcessMKUSR(input, user, pass, grp)

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
		println("Error campos incompletos")
		return
	}

	functions.MKUSR(user, pass, grp)

	*user = ""
	*pass = ""
	*grp = ""
}

func RRmusr(input string) {
	flag.Parse()
	functions.ProcessRMUSR(input, user)

	if *user == "" {
		println("Error user no puede estar vacio")
		return
	}

	functions.RMUSR(user)

	*user = ""
}

func comando_chgrp(input string) {
	flag.Parse()
	functions.ProcessCHGRP(input, user, grp)
}

// Adminstración de carpetas.
func MMkdir(input string) {
	flag.Parse()
	functions.ProcessMKDIR(input, path, r)

	if *path == "" {
		println("Error path no puede estar vacio")
		return
	}

	fmt.Println("Path: " + *path)
	fmt.Print("r: ")
	fmt.Println(*r)

	*path = ""
	*r = false
}

func MMkfile(input string) {
	flag.Parse()
	functions.ProcessMKFILE(input, path, r, size, cont)
}

func CCat(input string) {
	flag.Parse()
	functions.ProcessCAT(input, file)
}

func Remove(input string) {
	flag.Parse()
	functions.ProcessREMOVE(input, path)
}

func edit(input string) {
	flag.Parse()
	functions.ProcessEDIT(input, path, cont)
}

func rename(input string) {
	flag.Parse()
	functions.ProcessRENAME(input, path, name)
}

func copy(input string) {
	flag.Parse()
	functions.ProcessCOPY(input, path, destino)
}

func move(input string) {
	flag.Parse()
	functions.ProcessMOVE(input, path, destino)
}

func Ffind(input string) {
	flag.Parse()
	functions.ProcessFIND(input, path, destino)
}

func comando_chown(input string) {
	flag.Parse()
	functions.ProcessCHOWN(input, path, user, r)
}

func comando_Chmod(input string) {
	flag.Parse()
	functions.ProcessCHMOD(input, path, ugo, r)
}

func pause() {
	fmt.Println("Presione ENTER tecla para continuar...")
	fmt.Scanln()
}

func execute(input string) {
	flag.Parse()
	functions.ProcessExecute(input, path)
	if *path == "" {
		fmt.Println("Error Path no puede estar vacío.")
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
		Comandos(linea)
	}

	// Verifica si hubo algún error durante la lectura
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
	*path = ""
}

// Reportes
func rep(input string) {
	flag.Parse()
	fmt.Println("El comando es: " + input)
	functions.ProcessREP(input, name, path, id, ruta)
	fmt.Println("Namere: ", *name)
	fmt.Println("Pathre: ", *path)
	fmt.Println("Idre: ", *id)
	fmt.Println("Rutare: ", *ruta)
	if *name == "" || *path == "" || *id == "" {
		println("Faltan parametros para el REP.")
		return
	}
	fmt.Println("Ejecutando REP...")
	functions.GenerateReports(name, path, id, ruta)
}
