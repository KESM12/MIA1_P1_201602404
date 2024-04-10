package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

//estructuras

type Particion struct {
	PartStatus byte
	PartType   byte
	PartFit    [2]byte
	PartStart  int64
	PartSize   int64
	PartName   [16]byte
}

type MBR struct {
	Mbrtamano int64
	Mbrfecha  [20]byte
	Mbrdisk   int64
	Diskfit   [2]byte
	Particion [4]Particion
	Part      [4]bool
	Extend    bool
}

type EBR struct {
	EbrStatus byte
	EbrtFit   [2]byte
	EbrtStart int64
	EbrtSize  int64
	EbrtNext  int64
	EbrtName  [16]byte
}
type Gmount struct {
	ID           [10]string
	NameD        string
	Druta        string
	letra        byte
	Mparticiones [10]string
	numero       int64
}

//variables globales
var DMount [27]Gmount

//inicio	fmt.Println(buffer.Bytes)
func main() {
	Menu()
}

func Menu() {
	fmt.Println("ingrese comando a ejecutar")
	var comandoentrante string
	leer := bufio.NewReader(os.Stdin)
	entrada, _ := leer.ReadString('\n')
	comandoentrante = strings.TrimRight(entrada, "\r\n")
	if comandoentrante == "salir" {
		fmt.Println("usted salio exitosamente")
	} else {
		leercomando(comandoentrante)
		Menu()
	}
}

func leercomando(linea string) {
	comando := strings.Split(linea, " ")
	comparador := strings.ToLower(comando[0])
	switch comparador {
	case "exec":
		fmt.Println("-------------comando exec---------------")
		comandoexec(linea)
		fmt.Println("------------fin comando exec -------------")
		fmt.Println("")
	case "pause":
		fmt.Println("-------------comando pause--------------")
		fmt.Print("Presione enter para continuar")
		lector := bufio.NewReader(os.Stdin)
		entrada, _ := lector.ReadString('\n')
		fmt.Print(entrada)
		fmt.Println("------------fin comando pause-----------")
		fmt.Println("")
	case "mkdisk":
		fmt.Println("------------comando mkdisk-------------")
		comando_mksdisk(linea)
		fmt.Println("-------------fin comando mkdisk-------------")
		fmt.Println("")
	case "rmdisk":
		fmt.Println("------------comando rmdisk--------------")
		comando_rmdisk(linea)
		fmt.Println("-----------fin comando rmdisk------------")
		fmt.Println("")
	case "fdisk":
		fmt.Println("------------comando fdisk--------------")
		comando_fsdisk(linea)
		fmt.Println("----------fin comando disk------------")
		fmt.Println("")
	case "mount":
		fmt.Println("------------comando mount---------------")
		comando_mount(linea)
		fmt.Println("-----------fin comando mount------------")
		fmt.Println("")
	case "unmount":
		fmt.Println("------------comando unmount-------------")
		comando_unmount(linea)
		fmt.Println("------------fin comando unmount-------------")
		fmt.Println("")
	case "rep":
		fmt.Println("------------comando reportes-------------")
		comando_rep(linea)
		fmt.Println("-----------fin comando reportes")
		fmt.Println("")
	default:
		if strings.Contains(linea, "#") {
			fmt.Println("")
			fmt.Println("")
		} else {
			fmt.Println("comando erroneo")
		}
	}
}

//comandos

func comandoexec(linea string) {
	comando := strings.Split(linea, " ")
	ruta := strings.Split(strings.ToLower(comando[1]), "-path=")
	path := strings.ReplaceAll(ruta[1], "\"", "")
	leerarchivo(path)
}

func leerarchivo(ruta string) {
	file, err := os.Open(ruta)
	if err != nil {
		fmt.Println("ruta no existe")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		comando := strings.Split(scanner.Text(), "\n")
		for i := 0; i < len(comando); i++ {
			leercomando(comando[i])
			fmt.Println(comando[i])
		}
	}
}

func comando_rep(linea string) {
	var nameRep string = ""
	var rutaRep string = ""
	var id string = ""
	var comillas bool = false
	if strings.Contains(linea, "\"") {
		comillas = true
		aux := strings.Split(linea, "\"")
		for i := 1; i < len(aux)-1; i++ {
			rutaRep += aux[i]
		}
	}
	comando := strings.Split(linea, " ")
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-path=") {
			if comillas == true {
				path := strings.ReplaceAll(rutaRep, "\"", "")
				fmt.Println("ruta", path)
			} else {
				aux := strings.Split(comando[i], "=")
				rutaRep = aux[1]
				path := strings.ReplaceAll(rutaRep, "\"", "")
				fmt.Println("ruta", path)
			}
		} else if strings.Contains(strings.ToLower(comando[i]), "-id=") {
			aux := strings.Split(comando[i], "=")
			id = aux[1]
			fmt.Println("id", id)
		} else if strings.Contains(strings.ToLower(comando[i]), "-name=") {
			aux := strings.Split(comando[i], "=")
			nameRep = aux[1]
			fmt.Println("tipo", nameRep)
		}
	}
	rep(id, rutaRep, nameRep)
}
func rep(id string, rutaImagen string, tipo string) {
	if strings.ToLower(tipo) == "mbr" {
		ReportMBR(rutaImagen, id)
	} else if strings.ToLower(tipo) == "disk" {
		ReportDisk(rutaImagen, id)
	}
}

func ReportMBR(ruta string, id string) {
	var reporte string
	var nombre string
	var fit string
	DiscoM := obtenerDisco(id)
	file, err := os.OpenFile(DiscoM.Druta, os.O_RDWR, 0777)
	check(err)
	defer file.Close()
	file.Seek(0, 0)
	mbraux := obtenerMBR(file)

	fechambr := string(mbraux.Mbrfecha[:19])
	reporte += "digraph MBR {\n"
	reporte += "\tgraph[ label= \"Reporte MBR\"];\n"
	reporte += "\t  node [shape=plain]\n\n"
	reporte += "\t randir = TB; \n\n"

	reporte += "\tmbr[label=<\n"
	reporte += "\t\t<table border=\"1\" cellborder=\"1\" cellspacing=\"0\">\n "
	reporte += "\t\t\t <tr> <td colspan='2'> MBR " + DiscoM.NameD + "</td> </tr>\n"
	reporte += "\t\t\t <tr> <td> mbr_asignature </td> <td>" + strconv.Itoa(int(mbraux.Mbrdisk)) + "</td> </tr>\n"
	reporte += "\t\t\t <tr> <td> Nombre </td> <td> default </td> </tr>\n"
	reporte += "\t\t\t <tr> <td> mbr_Fecha</td> <td>" + fechambr + " </td> </tr>\n"
	reporte += "\t\t\t <tr> <td> mbr_size</td> <td>" + strconv.Itoa(int(mbraux.Mbrtamano)) + "</td> </tr>\n"
	reporte += "\t\t\t <tr> <td> mbr_</td> <td>" + strconv.Itoa(int(mbraux.Mbrtamano)) + "</td> </tr>\n"

	for j, i := range mbraux.Particion {
		aux := j + 1
		indice := strconv.Itoa(aux)
		if i.PartStart > -1 && i.PartStatus != '1' {
			for _, z := range i.PartName {
				if z != 0 {
					nombre = nombre + string(z)
				}
			}
			for _, nameaux := range i.PartFit {
				if nameaux != 0 {
					fit = fit + string(nameaux)
				}
			}

			reporte += "\t\t\t <tr> <td>part_status_" + indice + "</td> <td> " + string(i.PartStatus) + " </td> </tr>\n"
			reporte += "\t\t\t <tr> <td>part_type_" + indice + "</td> <td> " + "a" + " " + " </td> </tr>\n"
			reporte += "\t\t\t <tr> <td>part_fit_" + indice + "</td> <td> " + fit + " </td> </tr>\n"
			reporte += "\t\t\t <tr> <td>part_start_" + indice + "</td> <td> " + strconv.Itoa(int(i.PartStart)) + " </td> </tr>\n"
			reporte += "\t\t\t <tr> <td>part_size_" + indice + "</td> <td> " + strconv.Itoa(int(i.PartSize)) + " </td> </tr>\n"
			reporte += "\t\t\t <tr> <td>part_name" + indice + "</td> <td> " + nombre + " </td> </tr>\n"
			fit = ""
			nombre = ""

		}
	}

	reporte += "\t\t </table>\n"
	reporte += "\t >];\n\n"
	reporte += "\n}\n"

	GenerateDot(ruta, "mbrReporte.dot", reporte)

}

func GenerateDot(ruta string, name string, contenido string) {
	nombre := strings.Split(ruta, "/")
	var path string
	for i := 0; i < (len(nombre) - 1); i++ {
		path = path + nombre[i] + "/"
	}
	err := os.MkdirAll(path, 0777)
	file, err := os.Create(path + name)
	check(err)
	defer file.Close()
	files, err1 := os.OpenFile(path+name, os.O_RDWR|os.O_TRUNC, 0777)
	check(err1)
	defer files.Close()
	_, err = files.WriteString(contenido)
	err = files.Sync()
	check(err)
	archivo := directorio(path + name)
	var tipo string = "-Tpng"
	if strings.Contains(strings.ToLower(ruta), ".jpg") {
		tipo = "-Tjpg"
	} else if strings.Contains(strings.ToLower(ruta), ".png") == true {
		tipo = "-Tpng"
	}
	fmt.Println(tipo, nombre[len(nombre)-1], archivo)
	Procesoexec(tipo, ruta, archivo)

	fmt.Println("se genero el reporte MBR con exito")

}
func abrirReporte(ruta string) string {
	indice := strings.LastIndex(ruta, ".")
	if indice > -1 {
		ruta = ruta[indice:]
		return ruta
	}
	return "-Tpng"
}
func directorio(path string) string {

	if strings.Contains(path, "\"") {
		path = strings.ReplaceAll(path, "\"", "")
		return path
	}
	return path
}

func Procesoexec(tipos string, ruta string, nombre string) int {
	comn := "dot"
	arg0 := tipos
	arg1 := nombre
	arg2 := "-o"
	arg3 := ruta

	cmd := exec.Command(comn, arg0, arg1, arg2, arg3)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		return 0
	}
	if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		if ws.Exited() {
			return ws.ExitStatus()
		}

		if ws.Signaled() {
			return -int(ws.Signal())
		}
	}
	return -1
}

func ReportDisk(ruta string, id string) {
	var reporte string
	DiscoM := obtenerDisco(id)
	file, err := os.OpenFile(DiscoM.Druta, os.O_RDWR, 0777)
	check(err)
	defer file.Close()
	file.Seek(0, io.SeekStart)
	mbraux := obtenerMBR(file)

	reporte += "digraph Disco {\n"
	reporte += "\tgraph[ label= \"Reporte Disco\"];\n"
	reporte += "\t  node [shape=plain]\n\n"
	reporte += "\t randir = TB; \n\n"

	var totaldisco float64 = float64(mbraux.Mbrtamano)
	var totalmbr float64
	totalmbr = float64(unsafe.Sizeof(mbraux)) * 100 / totaldisco

	reporte += "\tdisk[label=<\n"
	reporte += "\t\t<table border='1' cellborder='1'  cellspacing='0' cellpadding='4'>\n "
	reporte += "\t\t\t <tr> <td colspan='6' >  " + DiscoM.NameD + "</td> </tr> <tr>\n"
	reporte += "<td> MBR <br/>" + strconv.FormatFloat(totalmbr, 'f', 6, 64) + "% del disco</td>"

	for j, i := range mbraux.Particion {
		var tampart float64 = float64(i.PartSize)
		if i.PartSize > -1 {
			total := tampart * 100 / totaldisco
			if i.PartStatus != '1' {
				if i.PartType == 'e' {
					file.Seek(i.PartStart, io.SeekStart)
					ebraux := obtenerEBR(file)
					var totalebr float64
					totalebr = float64(unsafe.Sizeof(ebraux)) * 100 / tampart
					var totallibre float64 = (tampart - totalebr) * 100 / tampart
					reporte += "\t\t\t\t<td>\n\t\t\t\t\t<table border='1' cellspacing = '0' cellborder='1'>\n"
					reporte += "\t\t\t\t\t\t<tr>\n"
					reporte += "\t\t\t\t\t\t\t<td colspan='2' > Extendida <br/>" + strconv.FormatFloat(total, 'f', 6, 64) + "% del disco </td>\n"
					reporte += "\t\t\t\t\t\t</tr>\n\t\t\t\t\t\t\t\n"
					reporte += "\t\t\t\t\t\t<tr>\n"
					reporte += "\t\t\t\t\t\t\t<td> EBR <br/>" + strconv.FormatFloat(totalebr, 'f', 6, 64) + "% del disco </td>\n"
					reporte += "\t\t\t\t\t\t\t<td> libre <br/>" + strconv.FormatFloat(totallibre, 'f', 6, 64) + "% del disco </td>\n"
					reporte += "\t\t\t\t\t\t</tr>\n\t\t\t\t\t\t\t\n"
					reporte += "\t\t\t\t </table>\n\t\t\t</td>"
					var next int64
					if j != 3 {
						next += mbraux.Particion[j+1].PartStart
					}
					reporte += siguiente(i.PartStart, i.PartSize, next, j, totaldisco)
				} else {
					if total > 0 {
						reporte += "\t\t\t<td> Primaria <br/> " + strconv.FormatFloat(total, 'f', 6, 64) + "% del disco </td>\n"
						var next int64 = 0
						if j != 3 {
							next += mbraux.Particion[j+1].PartStart
						}
						reporte += siguiente(i.PartStart, i.PartSize, next, j, totaldisco)
					}
				}
			} else {
				if total > 0 {
					reporte += "<td> Libre <br/> " + strconv.FormatFloat(total, 'f', 6, 64) + "% del disco</td>\n"
				}
			}
		}
	}

	reporte += "</tr></table>\n\t >];\n\n"
	reporte += "\n}\n"

	GenerateDot(ruta, "ReporteDisco.dot", reporte)

}

func siguiente(start int64, size int64, sig int64, indice int, total float64) string {
	var reporte string
	var mbraux MBR
	if indice != 3 {
		var Part1 float64 = float64(start + size)
		var tam float64 = total + float64(int64(unsafe.Sizeof(mbraux)))

		if sig == -1 && tam != Part1 {
			var libre float64 = tam - Part1 + float64(unsafe.Sizeof(mbraux))
			var por float64 = libre * 100 / total
			if por == 0 {
				reporte = ""
			} else {
				reporte += "\t\t\t<td> Libre <br/>"
				reporte += strconv.FormatFloat(por, 'f', 6, 64)
				reporte += "% del disco </td>"
			}
		}
	} else {
		var Part1 float64 = float64(start + size)
		var tam float64 = total + float64(int64(unsafe.Sizeof(mbraux)))
		if tam != Part1 {
			var libre float64 = tam - Part1 + float64(unsafe.Sizeof(mbraux))
			var por float64 = libre * 100 / total
			if por == 0 {
				reporte = ""
			} else {
				reporte += "\t\t\t<td> Libre <br/>"
				reporte += strconv.FormatFloat(por, 'f', 6, 64)
				reporte += "% del disco </td>"
			}
		}
	}
	return reporte
}

func obtenerDisco(id string) Gmount {
	var existe bool
	contador := 0
	for i := 0; i < len(DMount); i++ {
		for j := 0; j < len(DMount[i].Mparticiones); j++ {
			if id == DMount[i].ID[j] {
				existe = true
				contador = i
				break
			}
		}
	}
	if existe == true {
		return DMount[contador]
	}
	return DMount[contador]
}
func comando_unmount(linea string) {
	var id string = "vacio"

	comando := strings.Split(linea, " ")
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-id=") {
			aux := strings.Split(comando[i], "=")
			id = aux[1]
			fmt.Println("id", id)
		}
	}
	unomunt(id)
}

func unomunt(identificador string) {
	var existe = false
	var Dcont = 0
	var Pcont = 0
	for i := 0; i < len(DMount); i++ {
		for j := 0; j < len(DMount[i].Mparticiones); j++ {
			if identificador == DMount[i].ID[j] {
				existe = true
				Dcont = i
				Pcont = j
				break
			}
		}
	}
	fmt.Println(Dcont, Pcont, existe)
	if Dcont >= 0 && Pcont >= 0 && existe == true {
		DMount[Dcont].ID[Pcont] = ""
		DMount[Dcont].Mparticiones[Pcont] = ""
		fmt.Println("se desmonto particion con id", identificador)
	} else {
		fmt.Println("ID no existe", identificador)
	}
}

func comando_mount(linea string) {
	var rutaparticion string = "vacio"
	var nameparticion string = "vacio"
	var comillas bool = false
	if strings.Contains(linea, "\"") {
		comillas = true
		aux := strings.Split(linea, "\"")
		for i := 1; i < len(aux)-1; i++ {
			rutaparticion += aux[i]
		}
	}
	comando := strings.Split(linea, " ")
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-path=") {
			if comillas == true {
				path := strings.ReplaceAll(rutaparticion, "\"", "")
				fmt.Println("ruta", path)
			} else {
				aux := strings.Split(comando[i], "=")
				rutaparticion = aux[1]
				path := strings.ReplaceAll(rutaparticion, "\"", "")
				fmt.Println("ruta", path)
			}
		} else if strings.Contains(strings.ToLower(comando[i]), "-name=") {
			aux := strings.Split(comando[i], "=")
			nameparticion = aux[1]
			fmt.Println("nombre", nameparticion)
		}
	}
	mount(rutaparticion, nameparticion)
}

func mount(ruta string, name string) {
	nombre := strings.Split(ruta, "/")
	fmt.Println(nombre[len(nombre)-1])
	var auxName [16]byte
	for i, j := range []byte(name) {
		auxName[i] = byte(j)
	}

	if ruta != "vacio" && name != "vacio" {
		file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
		check(err)
		defer file.Close()
		file.Seek(0, 0)
		mbraux := obtenerMBR(file)
		var indice = -1
		var verficar = false
		for i, j := range mbraux.Particion {
			if auxName == j.PartName {
				indice = i
				verficar = true
				break
			}
		}
		if verficar == true {
			if indice > -1 {
				var contador = 0
				var existe = false
				if DMount[0].ID[0] == "" {
					DMount[0].Druta = ruta
					DMount[0].NameD = nombre[len(nombre)-1]
					DMount[0].letra = 'a'
					DMount[0].Mparticiones[0] = name
					DMount[0].numero = int64(1)
					conv := strconv.FormatInt(DMount[0].numero, 10)
					DMount[0].ID[0] = "vd" + string(DMount[0].letra) + conv
					fmt.Println("se monto particion con id", DMount[contador].ID[0])
				} else {
					for i := 0; i < len(DMount); i++ {
						if nombre[len(nombre)-1] == DMount[i].NameD {
							contador = i
							existe = true
							break
						} else {
							if DMount[i].NameD == "" {
								contador = i
								break
							}
						}
					}
				}
				if contador > 0 && existe == false {
					DMount[contador].Druta = ruta
					DMount[contador].NameD = nombre[len(nombre)-1]
					convb := (97 + contador)
					character := rune(convb)
					DMount[contador].letra = byte(character)
					DMount[contador].Mparticiones[0] = name
					DMount[contador].numero = int64(1)
					conv := strconv.FormatInt(DMount[contador].numero, 10)
					DMount[contador].ID[0] = "vd" + string(DMount[contador].letra) + conv
					fmt.Println("sse monto particion con id", DMount[contador].ID[0])

				} else if existe == true {
					var contador2 = 0
					for i := 0; i < len(DMount[contador].Mparticiones); i++ {
						if DMount[contador].Mparticiones[i] == name {
							fmt.Println("ya monto una particion con este nombre")
							break
						}
						if DMount[contador].Mparticiones[i] == "" {
							contador2 = i
							break
						}
					}
					conv := strconv.FormatInt(int64(contador2+1), 10)
					DMount[contador].Mparticiones[contador2] = name
					DMount[contador].ID[contador2] = "vd" + string(DMount[contador].letra) + conv
					fmt.Println("se monto particion con id", DMount[contador].ID[contador2])
				}
			}
		} else {
			fmt.Println("no existe particion con ese nombre")
		}
	} else {
		fmt.Println("Particion no existe")
	}
}

func comando_fsdisk(linea string) {
	var banderas = [8]bool{false, false, false, false, false, false, false, false}

	var sizeparticion string = "vacio"
	var rutaparticion string = "vacio"
	var unitparticion string = "vacio"
	var fitparticion string = "vacio"
	var typeparticion string = "vacio"
	var deleteparticion string = "vacio"
	var nameparticion string = "vacio"
	var addparticion string = "vacio"
	var comillas bool = false
	if strings.Contains(linea, "\"") {
		comillas = true
		aux := strings.Split(linea, "\"")
		for i := 1; i < len(aux)-1; i++ {
			rutaparticion += aux[i]
		}
	}
	comando := strings.Split(linea, " ")
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-size=") {
			banderas[0] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-unit=") {
			banderas[1] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-path=") {
			banderas[2] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-type=") {
			banderas[3] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-fit=") {
			banderas[4] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-delete=") {
			banderas[5] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-name=") {
			banderas[6] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-add=") {
			banderas[7] = true
		}
	}
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-size=") {
			aux := strings.Split(comando[i], "=")
			sizeparticion = aux[1]
			fmt.Println("tam", sizeparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-unit=") {
			aux := strings.Split(comando[i], "=")
			unitparticion = aux[1]
			fmt.Println("unit", unitparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-path=") {
			if comillas == true {
				path := strings.ReplaceAll(rutaparticion, "\"", "")
				fmt.Println("ruta", path)
			} else {
				aux := strings.Split(comando[i], "=")
				rutaparticion = aux[1]
				path := strings.ReplaceAll(rutaparticion, "\"", "")
				fmt.Println("ruta", path)
			}
		} else if strings.Contains(strings.ToLower(comando[i]), "-type=") {
			aux := strings.Split(comando[i], "=")
			typeparticion = aux[1]
			fmt.Println("tipo", typeparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-fit=") {
			aux := strings.Split(comando[i], "=")
			fitparticion = aux[1]
			fmt.Println("fit", fitparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-delete=") {
			aux := strings.Split(comando[i], "=")
			deleteparticion = aux[1]
			fmt.Println("eliminar", deleteparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-name=") {
			aux := strings.Split(comando[i], "=")
			nameparticion = aux[1]
			fmt.Println("nombre", nameparticion)
		} else if strings.Contains(strings.ToLower(comando[i]), "-add=") {
			aux := strings.Split(comando[i], "=")
			addparticion = aux[1]
			fmt.Println("agregar", addparticion)
		}
	}
	if banderas[5] == false {
		if banderas[1] == false {
			unitparticion = "kb"
			fmt.Println("unidad KB")
		}
		if banderas[3] == false {
			typeparticion = "p"
			fmt.Println("tipo Primaria")
		}
		if banderas[4] == false {
			fitparticion = "wf"
			fmt.Println("fit peor")
		}
	}

	fdisk(sizeparticion, rutaparticion, unitparticion, typeparticion, fitparticion, deleteparticion, nameparticion, addparticion)
}
func fdisk(size string, ruta string, unit string, tipo string, fit string, eliminar string, nombre string, add string) {

	var sizeparticion int64
	var rutaparticion string = ruta
	var unitparticion string = unit
	var typeparticion string = tipo
	var nameparticion string = nombre
	var fitparticion string = fit
	var deleteparticion string = eliminar
	var addparticion int64

	sizeparticion, err := strconv.ParseInt(size, 10, 64)
	addparticion, err1 := strconv.ParseInt(add, 10, 64)
	check(err1)
	files, err := os.OpenFile(rutaparticion, os.O_RDWR, 0777)
	defer files.Close()
	if err != nil {
		fmt.Println("ruta no existe")
	}
	if strings.ToLower(unitparticion) == "k" {
		sizeparticion = sizeparticion * 1024
	} else if strings.ToLower(unitparticion) == "m" {
		sizeparticion = sizeparticion * 1024 * 1024
	} else {
		sizeparticion = sizeparticion * 1
	}

	mbrauxiliar := MBR{}
	//var espacio_particion = false
	mbrauxiliar = obtenerMBR(files)

	var contador_particion = 0
	for i := 0; i < 4; i++ {
		if mbrauxiliar.Part[i] == false {
			contador_particion++
		}
	}
	if eliminar == "vacio" && add == "vacio" {

		if contador_particion >= 1 {
			if strings.ToLower(typeparticion) == "p" {
				ParticionPrimaria(rutaparticion, nameparticion, 'p', fitparticion, unitparticion, sizeparticion)
				if mbrauxiliar.Extend == false {
					fmt.Println("Queda 1 extendidas y principales", (contador_particion - 1))
				} else {
					fmt.Println("Queda 0 extendidas y principales", (contador_particion - 1))
				}

			} else if strings.ToLower(typeparticion) == "e" {
				if mbrauxiliar.Extend == false {
					fmt.Println("Queda 0 extendidas y principales", (contador_particion - 1))
					ParticionExtendida(rutaparticion, nameparticion, 'e', fitparticion, unitparticion, sizeparticion)
				} else {
					fmt.Println("ya existe una particion extendida")
				}
			} else if strings.ToLower(typeparticion) == "l" {
				fmt.Println("No hay particiones logicas ")
			}
		}
	} else if eliminar != "vacio" {
		eliminar_Particion(deleteparticion, nameparticion, rutaparticion)
	} else if add != "vacio" {
		AEpraticion(addparticion, unitparticion, rutaparticion, nameparticion)
	}
}

func eliminar_Particion(_ string, nombre string, ruta string) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	check(err)
	defer file.Close()
	mbraux := obtenerMBR(file)
	var auxName [16]byte
	for i, j := range []byte(nombre) {
		auxName[i] = byte(j)
	}
	var verificar = false
	var indice = 0
	for i := 0; i < 4; i++ {
		if auxName == mbraux.Particion[i].PartName {
			verificar = true
			indice = 0
			break
		}
	}

	if verificar == true && indice > -1 {
		file.Seek(0, io.SeekStart)
		if mbraux.Particion[indice].PartType == 'p' {
			aux := make([]byte, mbraux.Particion[indice].PartSize)
			file.Seek(mbraux.Particion[indice].PartStart, io.SeekStart)
			var bin bytes.Buffer
			binary.Write(&bin, binary.BigEndian, &aux)
			escribirbinario(file, bin.Bytes())

			file.Seek(0, io.SeekStart)
			copy(mbraux.Particion[indice].PartFit[:], "ff")
			for i := 0; i < 16; i++ {
				mbraux.Particion[indice].PartName[i] = '0'
			}
			//mbraux.Particion[indice].PartSize = 0
			mbraux.Particion[indice].PartStart = -1
			mbraux.Particion[indice].PartStatus = '1'
			mbraux.Particion[indice].PartType = 'p'

			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, &mbraux)
			escribirbinario(file, binario.Bytes())
			fmt.Println("particion eliminada")

		} else if mbraux.Particion[indice].PartType == 'e' {
			aux := make([]byte, mbraux.Particion[indice].PartSize)
			file.Seek(mbraux.Particion[indice].PartStart, io.SeekStart)
			var bin bytes.Buffer
			binary.Write(&bin, binary.BigEndian, &aux)
			escribirbinario(file, bin.Bytes())

			file.Seek(0, io.SeekStart)
			copy(mbraux.Particion[indice].PartFit[:], "ff")
			for i := 0; i < 16; i++ {
				mbraux.Particion[indice].PartName[i] = '0'
			}
			//mbraux.Particion[indice].PartSize = 0
			mbraux.Particion[indice].PartStart = -1
			mbraux.Particion[indice].PartStatus = '1'
			mbraux.Particion[indice].PartType = 'p'

			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, &mbraux)
			escribirbinario(file, binario.Bytes())
			fmt.Println("particion eliminada")
		}
	} else {
		fmt.Println("no existe la particion con ese nombre")
	}

}
func AEpraticion(add int64, unit string, ruta string, name string) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	check(err)
	defer file.Close()
	mbraux := obtenerMBR(file)

	var auxName [16]byte
	for i, j := range []byte(name) {
		auxName[i] = byte(j)
	}
	var sizeparticion int64 = add
	if strings.ToLower(unit) == "k" {
		sizeparticion = add * 1024
	} else if strings.ToLower(unit) == "m" {
		sizeparticion = add * 1024 * 1024
	} else {
		sizeparticion = add * 1
	}
	var verificar = false
	var indice = 0

	for i := 0; i < 4; i++ {
		if auxName == mbraux.Particion[i].PartName {
			verificar = true
			indice = 0
			break
		}
	}

	if verificar == true && indice > -1 {
		if sizeparticion > 0 {
			if math.Abs(float64(sizeparticion)) >= float64(mbraux.Mbrtamano) {
				fmt.Println("ya no hay espacio libre")
			} else {
				mbraux.Particion[indice].PartSize = mbraux.Particion[indice].PartSize + int64(math.Abs(float64(sizeparticion)))
				file.Seek(0, 0)
				var bin bytes.Buffer
				binary.Write(&bin, binary.BigEndian, &mbraux)
				escribirbinario(file, bin.Bytes())
				fmt.Println("se agrego es espacio exitosamente")
			}
		} else {
			fmt.Println(mbraux.Particion[indice].PartSize)
			if math.Abs(float64(sizeparticion)) >= float64(mbraux.Particion[indice].PartSize) {
				fmt.Println("el espacio a reducir es mayor al la particion")
			} else {
				mbraux.Particion[indice].PartSize = mbraux.Particion[indice].PartSize - int64(math.Abs(float64(sizeparticion)))
				file.Seek(0, 0)
				var bin bytes.Buffer
				binary.Write(&bin, binary.BigEndian, &mbraux)
				escribirbinario(file, bin.Bytes())
				fmt.Println("se redujo es espacio exitosamente")
			}
		}
	} else {
		fmt.Println("la particion no existe")
	}

}

func ParticionExtendida(ruta string, name string, tipo byte, fit string, unit string, size int64) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	check(err)
	defer file.Close()
	ebraux := EBR{}
	var auxName [16]byte
	for i, j := range []byte(name) {
		auxName[i] = byte(j)
	}
	mbraux := obtenerMBR(file)
	var contador int64 = 0
	for i := 0; i < 4; i++ {
		if mbraux.Particion[i].PartStatus != '1' {
			contador += mbraux.Particion[i].PartSize
		}
	}

	if (mbraux.Mbrtamano - contador) >= size {
		var verficar bool = false
		for i := 0; i < 4; i++ {
			if mbraux.Particion[i].PartName == auxName {
				verficar = true
				break
			}
		}
		if !verficar {
			var indice int = -1
			if strings.ToLower(fit) == "ff" {
				indice = primerAjuste(mbraux, size)
			} else if strings.ToLower(fit) == "bf" {
				indice = mejorAjuste(mbraux, size)
			} else if strings.ToLower(fit) == "wf" {
				indice = peorAjuste(mbraux, size)
			}
			if indice != -1 {

				var auxfit [2]byte
				for i, j := range []byte(fit) {
					auxfit[i] = byte(j)
				}

				mbraux.Particion[indice].PartFit = auxfit
				mbraux.Particion[indice].PartType = tipo
				mbraux.Particion[indice].PartSize = size
				mbraux.Particion[indice].PartStatus = '0'
				mbraux.Part[indice] = false
				mbraux.Extend = true
				copy(mbraux.Particion[indice].PartName[:], name)

				if indice == 0 {
					mbraux.Particion[indice].PartStart = int64(unsafe.Sizeof(mbraux)) + 1
				} else {
					mbraux.Particion[indice].PartStart = mbraux.Particion[indice-1].PartStart + mbraux.Particion[indice-1].PartSize
				}

				ebraux.EbrStatus = '0'
				ebraux.EbrtFit = auxfit
				ebraux.EbrtNext = -1
				ebraux.EbrtSize = 0
				ebraux.EbrtStart = mbraux.Particion[indice].PartStart
				copy(ebraux.EbrtName[:], name)
				CrearEBR(ruta, mbraux, ebraux, indice)

				fmt.Println("Se creo particion extendida")
			} else {
				fmt.Println("Ya se ah creado el maximo de particiones")
			}
		} else {
			fmt.Println("ya se ah creado un partacion con este nombre")
		}
	} else {
		fmt.Println("ya no tiene espacio en el disco ")
	}
}

func ParticionPrimaria(ruta string, name string, tipo byte, fit string, unit string, size int64) {
	fmt.Println(size)
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	check(err)
	defer file.Close()
	var auxName [16]byte
	for i, j := range []byte(name) {
		auxName[i] = byte(j)
	}
	mbraux := obtenerMBR(file)

	var contador int64 = 0
	for i := 0; i < 4; i++ {
		if mbraux.Particion[i].PartStatus != '1' {
			contador += mbraux.Particion[i].PartSize
		}
	}

	if (mbraux.Mbrtamano - contador) >= size {
		var verficar bool = false
		for i := 0; i < 4; i++ {
			if mbraux.Particion[i].PartName == auxName {
				verficar = true
				break
			}
		}

		if !verficar {
			var indice int = 0
			if strings.ToLower(fit) == "ff" {
				indice = primerAjuste(mbraux, size)
			} else if strings.ToLower(fit) == "bf" {
				indice = mejorAjuste(mbraux, size)
			} else if strings.ToLower(fit) == "wf" {
				indice = peorAjuste(mbraux, size)
			}

			if indice != -1 {
				var auxfit [2]byte
				for i, j := range []byte(fit) {
					auxfit[i] = byte(j)
				}
				mbraux.Particion[indice].PartFit = auxfit
				mbraux.Particion[indice].PartType = tipo
				mbraux.Particion[indice].PartSize = size
				mbraux.Particion[indice].PartStatus = '0'
				mbraux.Part[indice] = true
				copy(mbraux.Particion[indice].PartName[:], name)
				if indice == 0 {
					mbraux.Particion[indice].PartStart = int64(unsafe.Sizeof(mbraux)) + 1
				} else {
					mbraux.Particion[indice].PartStart = mbraux.Particion[indice-1].PartStart + mbraux.Particion[indice-1].PartSize
				}

				file.Seek(0, 0)
				var binario bytes.Buffer
				binary.Write(&binario, binary.BigEndian, &mbraux)
				escribirbinario(file, binario.Bytes())
				fmt.Println("se creo la particion")

			} else {
				fmt.Println("Ya se ah creado el maximo de particiones")
			}
		} else {
			fmt.Println("ya se ah creado un partacion con este nombre")
		}
	} else {
		fmt.Println("ya no tiene espacio en el disco ")
	}
}
func mejorAjuste(mbraux MBR, size int64) int {
	var indice int = 0
	var verficar bool = false
	for i := 0; i < 4; i++ {
		if mbraux.Particion[i].PartStart == -1 || (mbraux.Particion[i].PartStatus == '1' && mbraux.Particion[i].PartSize >= size) {
			verficar = true
			if i != indice && mbraux.Particion[indice].PartSize > mbraux.Particion[i].PartSize {
				indice = i
			}
		}
	}
	if verficar {
		return indice
	}
	return -1
}
func peorAjuste(mbraux MBR, size int64) int {
	var indice int = 0
	var verficar bool = false
	for i := 0; i < 4; i++ {
		if mbraux.Particion[i].PartStart == -1 || (mbraux.Particion[i].PartStatus == '1' && mbraux.Particion[i].PartSize >= size) {
			verficar = true
			if i != indice && mbraux.Particion[indice].PartSize < mbraux.Particion[i].PartSize {
				indice = i
			}
		}
	}
	if verficar {
		return indice
	}
	return -1
	//return -1
}

func primerAjuste(mbraux MBR, size int64) int {
	var indice int = 0
	var verficar bool = false
	for ; indice < 4; indice++ {
		if mbraux.Particion[indice].PartStart == -1 || (mbraux.Particion[indice].PartStatus == '1' && mbraux.Particion[indice].PartSize >= size) {
			verficar = true
			break
		}
	}
	if verficar {
		return indice
	}
	return -1
}

func comando_mksdisk(linea string) {
	var banderas = [4]bool{false, false, false, false}
	var sizearchivo string
	var rutaarchivo string
	var unitarchivo string
	var fitarchivo string
	var comillas bool = false
	if strings.Contains(linea, "\"") {
		comillas = true
		aux := strings.Split(linea, "\"")
		for i := 1; i < len(aux)-1; i++ {
			rutaarchivo += aux[i]
		}
	}
	comando := strings.Split(linea, " ")
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-size=") {
			banderas[0] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-path=") {
			banderas[2] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-fit=") {
			banderas[3] = true
		} else if strings.Contains(strings.ToLower(comando[i]), "-unit=") {
			banderas[1] = true
		}
	}
	for i := 1; i < len(comando); i++ {
		if strings.Contains(strings.ToLower(comando[i]), "-size=") && banderas[0] == true {
			aux := strings.Split(comando[i], "=")
			sizearchivo = aux[1]
			fmt.Println("tam", sizearchivo)
		} else if strings.Contains(strings.ToLower(comando[i]), "-path=") && banderas[2] == true {
			if comillas == true {
				path := strings.ReplaceAll(rutaarchivo, "\"", "")
				fmt.Println("ruta", path)
			} else {
				aux := strings.Split(comando[i], "=")
				rutaarchivo = aux[1]
				path := strings.ReplaceAll(rutaarchivo, "\"", "")
				fmt.Println("ruta", path)
			}
		} else if strings.Contains(strings.ToLower(comando[i]), "-fit=") && banderas[3] == true {
			aux := strings.Split(comando[i], "=")
			fitarchivo = aux[1]
			fmt.Println("fit", fitarchivo)
		} else if strings.Contains(strings.ToLower(comando[i]), "-unit=") && banderas[1] == true {
			aux := strings.Split(comando[i], "=")
			unitarchivo = aux[1]
			fmt.Println("unit", unitarchivo)
		}
	}
	if banderas[1] == false {
		unitarchivo = "MB"
		fmt.Println("unidad MB")
	}
	if banderas[3] == false {
		fitarchivo = "FF"
		fmt.Println("fit FF")
	}

	mkdisk(sizearchivo, rutaarchivo, unitarchivo, fitarchivo)
}
func comando_rmdisk(linea string) {

	comando := strings.Split(linea, " ")
	var rutaarchivo string = ""
	var comillas bool = false
	if strings.Contains(linea, "\"") {
		comillas = true
		aux := strings.Split(linea, "\"")
		for i := 1; i < len(aux)-1; i++ {
			rutaarchivo += aux[i]
		}
	}
	if strings.Contains(strings.ToLower(comando[1]), "-path=") {
		if comillas == true {
			path := strings.ReplaceAll(rutaarchivo, "\"", "")
			fmt.Println("ruta", path)
		} else {
			aux := strings.Split(comando[1], "=")
			rutaarchivo = aux[1]
			path := strings.ReplaceAll(rutaarchivo, "\"", "")
			fmt.Println("ruta", path)
		}
	}

	err := os.Remove(rutaarchivo)
	if err != nil {
		fmt.Println("ruta no existe")
	} else {
		fmt.Println("eliminado correctamente")
	}
	//rmdisk -path=/home/alex/disco/Disco1.dsk
}

func mkdisk(size string, ruta string, unidad string, ajuste string) {
	sizeArchivo, err := strconv.ParseInt(size, 10, 64)
	var rutaArchivo string = ruta
	var unitArchivo string = unidad
	var fitArchivo string = ajuste
	if err != nil {
		fmt.Println("ruta no existe")
	}
	if sizeArchivo > 0 {
		crearArchivo(sizeArchivo, unitArchivo, rutaArchivo, fitArchivo)
	}

}

func crearArchivo(size int64, unit string, ruta string, ajuste string) {
	var path = ""
	nombre := strings.Split(ruta, "/")
	for i := 0; i < (len(nombre) - 1); i++ {
		path = path + nombre[i] + "/"
	}
	err := os.MkdirAll(path, 0777)
	file, err := os.Create(ruta)
	check(err)
	var BT int = int(size)
	var KB int = BT * 1024
	var MB int = BT * 1024 * 1024

	var binario int8 = 0
	aux := &binario
	var bin bytes.Buffer
	binary.Write(&bin, binary.BigEndian, aux)

	if file != nil {
		if strings.ToLower(unit) == "k" {
			for i := 0; i < KB; i++ {
				escribirbinario(file, bin.Bytes())
			}
			size = int64(KB)
		} else if strings.ToLower(unit) == "m" {
			for i := 0; i < MB; i++ {
				escribirbinario(file, bin.Bytes())
			}
			size = int64(MB)
		} else {
			fmt.Println("No se pudo crear el disco")
		}
	} else {
		fmt.Println("no se pudo crear el Archivo")
	}

	numero := rand.NewSource(time.Now().UnixNano())
	random := rand.New(numero)
	temporal := MBR{Mbrtamano: size, Mbrdisk: int64(random.Intn(100))}
	fecha := time.Now()
	temporal.Extend = false
	copy(temporal.Mbrfecha[:], fecha.Format("2006-01-02 15:04:05"))
	copy(temporal.Diskfit[:], ajuste)
	for i := 0; i < 4; i++ {
		temporal.Part[i] = false
		temporal.Particion[i].PartStatus = '0'
		temporal.Particion[i].PartStart = -1
		temporal.Particion[i].PartSize = 0
		copy(temporal.Particion[i].PartName[:], "")
	}
	EscribirMBR(file, temporal)

	/*var nuevobuffer bytes.Buffer
	enc := gob.NewEncoder(&nuevobuffer)
	enc.Encode(temporal)
	files, err := os.OpenFile(ruta, os.O_RDWR, 0777)

	files.Seek(0, 0)
	escribirbinario(files, nuevobuffer.Bytes())
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	*/
}

func EscribirMBR(file *os.File, mbraux MBR) {
	file.Seek(0, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &mbraux)
	escribirbinario(file, binario.Bytes())
}
func CrearEBR(ruta string, mbraux MBR, ebraux EBR, indice int) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	check(err)
	defer file.Close()

	file.Seek(0, 0)
	var bin bytes.Buffer
	binary.Write(&bin, binary.BigEndian, &mbraux)
	escribirbinario(file, bin.Bytes())

	file.Seek(mbraux.Particion[indice].PartStart, 0)
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, &ebraux)
	escribirbinario(file, buf.Bytes())
}
func escribirbinario(file *os.File, binario []byte) {
	_, err := file.Write(binario)
	if err != nil {
		fmt.Println("ruta no existe")
	}
}

func check(err error) {
	if err != nil {
		fmt.Println("ruta no existe")
	}
}
func obtenerMBR(file *os.File) MBR {
	mbrActual := MBR{}
	var size int = int(unsafe.Sizeof(mbrActual))
	data := leersiguientebyte(file, size)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &mbrActual)
	check(err)
	return mbrActual
}
func obtenerEBR(file *os.File) EBR {
	ebraux := EBR{}
	sizeRead := binary.Size(ebraux)
	data := leersiguientebyte(file, sizeRead)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &ebraux)
	check(err)
	return ebraux
}
func leersiguientebyte(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		fmt.Println(err)
	}
	return bytes
}