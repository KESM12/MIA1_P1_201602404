package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var diskCounter = 0

type mbr = struct {
	Mbr_tamano         [100]byte
	Mbr_fecha_creacion [100]byte
	Mbr_dsk_signature  [100]byte
	Dsk_fit            [100]byte
	Mbr_partition      [4]partition
	// Mbr_partition_2    partition
	// Mbr_partition_3    partition
	// Mbr_partition_4    partition
}

type partition = struct {
	Part_status      [100]byte
	Part_type        [100]byte
	Part_fit         [100]byte
	Part_start       [100]byte
	Part_size        [100]byte
	Part_name        [100]byte
	Part_correlative [100]byte
	Part_id          [100]byte
}

type Ebr = struct {
	Part_mount [100]byte //Estado
	Part_fit   [100]byte
	Part_start [100]byte
	Part_s     [100]byte
	Part_next  [100]byte
	Part_name  [100]byte
}

type super_bloque = struct {
	S_filesystem_type   [100]byte
	S_inodes_count      [100]byte
	S_blocks_count      [100]byte
	S_free_blocks_count [100]byte
	S_free_inodes_count [100]byte
	S_mtime             [100]byte
	S_mnt_count         [100]byte
	S_magic             [100]byte
	S_inode_size        [100]byte
	S_block_size        [100]byte
	S_firts_ino         [100]byte
	S_first_blo         [100]byte
	S_bm_inode_start    [100]byte
	S_bm_block_start    [100]byte
	S_inode_start       [100]byte
	S_block_start       [100]byte
}

type inodo = struct {
	I_uid   [100]byte
	I_gid   [100]byte
	I_size  [100]byte
	I_atime [100]byte
	I_ctime [100]byte
	I_mtime [100]byte
	I_block [100]byte
	I_type  [100]byte
	I_perm  [100]byte
}

type bloque_archivo = struct {
	B_content [100]byte
}

type content = struct {
	B_name  [100]byte
	B_inodo [100]byte
}

// Bloque de apuntadores.
type apuntadores = struct {
	B_pointers [100]byte
}

func main() {
	analizar()
}

func analizar() {
	finalizar := false
	fmt.Println()
	fmt.Println("***** KEVIN ESTUARDO SECAIDA MOLINA ***** ")
	reader := bufio.NewReader(os.Stdin)
	for !finalizar {
		fmt.Print("Ingrese un comando: ")
		comando, _ := reader.ReadString('\n')
		if strings.Contains(comando, "exit") {
			finalizar = true
			fmt.Println("Saliendo...")
		} else if strings.Contains(comando, "EXIT") {
			finalizar = true
			fmt.Println("Saliendo...")
		} else {
			if comando != "" && comando != "exit\n" && comando != "EXIT\n" {
				split_comando(comando)
			}
		}
	}
}

func split_comando(comando string) {
	var commandArray []string
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)

	band_comentario := false

	if strings.Contains(comando, "pause") {
		commandArray = append(commandArray, comando)
	} else if strings.Contains(comando, "#") {
		band_comentario = true
		fmt.Println(comando)
	} else {
		commandArray = strings.Split(comando, " -")
	}

	if !band_comentario {
		ejecutar_comando(commandArray)
	}
}

func ejecutar_comando(commandArray []string) {
	data := strings.ToLower(commandArray[0])

	if data == "mkdisk" {
		mkdisk(commandArray)
		fmt.Println()
	} else if data == "fdisk" {
		fdisk(commandArray)
		fmt.Println()
	} else if data == "rmdisk" {
		rmdisk(commandArray)
		fmt.Println()
	} else if data == "mount" {
		//mount()
		fmt.Println("Mount.")
	} else if data == "rep" {
		mostrar_mkdisk()
		fmt.Println()
	} else if data == "execute" {
		execute(commandArray)
		fmt.Println()
	} else if data == "pause" {
		pause()
		fmt.Println()
	} else {
		fmt.Println("Comando no fue reconocido...")
	}
}

func mkdisk(commandArray []string) {
	val_size := 0
	val_fit := ""
	val_unit := ""
	//val_path := ""

	band_size := false
	band_fit := false
	band_unit := false
	band_path := true
	band_error := false

	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "size="):
			if band_size {
				fmt.Println("Parametro -size ya fue ingresado...")
				band_error = true
				break
			}

			band_size = true

			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size

			if err != nil {
				msg_error(err)
			}

			if val_size < 0 {
				band_error = true
				fmt.Println("Parametro -size es negativo...")
				break
			}

		case strings.Contains(data, "fit="):
			if band_fit {
				fmt.Println("Parametro -fit ya fue ingresado...")
				band_error = true
				break
			}

			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" { // Best Fit
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" { // First Fit
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" { // Worst Fit
				band_fit = true
				val_fit = "w"
			} else {
				fmt.Println("Valor del parametro -fit no es valido...")
				band_error = true
				break
			}

		case strings.Contains(data, "unit="):
			if band_unit {
				fmt.Println("Parametro -unit ya fue ingresado...")
				band_error = true
				break
			}

			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			if val_unit == "k" || val_unit == "m" { // Kilobytes o Megabytes
				band_unit = true
			} else {
				fmt.Println("Valor del parametro -unit no es valido...")
				band_error = true
				break
			}

		default:
			fmt.Println("Parametro no valido...")
		}
	}

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

	if !band_error {
		if band_path {
			if band_size {
				total_size := 1024
				master_boot_record := mbr{}

				crear_disco(val_path)

				fecha := time.Now()
				str_fecha := fecha.Format("02/01/2006 15:04:05")

				copy(master_boot_record.Mbr_fecha_creacion[:], str_fecha)

				rand.Seed(time.Now().UnixNano())
				min := 0
				max := 100
				num_random := rand.Intn(max-min+1) + min

				copy(master_boot_record.Mbr_dsk_signature[:], strconv.Itoa(int(num_random)))

				if band_fit {
					copy(master_boot_record.Dsk_fit[:], val_fit)
				} else {
					copy(master_boot_record.Dsk_fit[:], "f")
				}

				if band_unit {
					if val_unit == "m" {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
						total_size = val_size * 1024
					} else {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024)))
						total_size = val_size
					}
				} else {
					copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
					total_size = val_size * 1024
				}

				// Inicializar Parcticiones
				for i := 0; i < 4; i++ {
					copy(master_boot_record.Mbr_partition[i].Part_status[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_type[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_fit[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_start[:], "-1")
					copy(master_boot_record.Mbr_partition[i].Part_size[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_name[:], "")
				}

				str_total_size := strconv.Itoa(total_size)

				cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=\""+val_path+"\" bs=1024 count="+str_total_size)
				cmd.Dir = "/"
				_, err := cmd.Output()

				if err != nil {
					msg_error(err)
				}

				disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

				if err != nil {
					msg_error(err)
				}

				mbr_byte := struct_a_bytes(master_boot_record)

				newpos, err := disco.Seek(0, os.SEEK_SET)

				if err != nil {
					msg_error(err)
				}
				_, err = disco.WriteAt(mbr_byte, newpos)

				if err != nil {
					msg_error(err)
				}

				disco.Close()
			}
		}
	}
	fmt.Println("MKDISK creado exitosamente")
}

func rmdisk(commandArray []string) {
	val_driveletter := ""
	//val_path := ""

	band_driveletter := false
	band_error := false

	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]
		switch {
		case strings.Contains(data, "driveletter="):
			if band_driveletter {
				fmt.Println("El parametro -driveletter ya fue ingresado...")
				band_error = true
				break
			}
			band_driveletter = true

			// Obtenemos la letra del abecedario según el contador global diskCounter
			//letter := string(rune('A'))
			// Formateamos dinámicamente la ruta del archivo de disco usando fmt.Sprintf
			val_driveletter = fmt.Sprintf("/home/taro/Escritorio/MIA/P1/%s.dsk", val_data)
		default:
			fmt.Println("Error, parametro no valido...")
		}
	}

	if !band_error {
		for band_driveletter {
			_, e := os.Stat(val_driveletter)

			if e != nil {
				if os.IsNotExist(e) {
					fmt.Println("Error no existe el disco.")
					band_driveletter = false
				}

			} else {
				fmt.Print("¿Esta seguro de eliminar el disco? [S/N]")

				var opcion string
				fmt.Scanln(&opcion)
				if opcion == "s" || opcion == "S" {
					cmd := exec.Command("/bin/sh", "-c", "rm \""+val_driveletter+"\"")
					cmd.Dir = "/"
					_, err := cmd.Output()

					if err != nil {
						msg_error(err)
					} else {
						fmt.Println("El disco fue eliminado satisfactoriamente")
					}
					band_driveletter = false
				} else if opcion == "n" || opcion == "N" {
					fmt.Print("El disco no se eliminara")
					band_driveletter = false
				} else {
					fmt.Println("Opción no valida.")
				}
			}
		}
	}
}

func fdisk(commandArray []string) {
	val_size := 0
	val_driveletter := ""
	val_name := ""
	val_unit := ""
	val_type := ""
	val_fit := ""
	val_delete := "Full"
	val_add := ""

	band_error := false
	band_size := false
	band_driveletter := false
	band_name := false
	band_unit := false
	band_type := false
	band_fit := false
	band_delete := false
	band_add := false

	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "size="): //obligatorio
			if band_size {
				fmt.Println("El parametro -size ya fue ingresado.")
				band_error = true
				break
			}
			band_size = true
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size
			fmt.Println("Size: ", val_size)
			if err != nil {
				msg_error(err)
				band_error = true
			}
			if val_size < 0 {
				band_error = true
				fmt.Println("El parametro -size no puede ser negativo.")
				break
			}
		case strings.Contains(data, "driveletter="): //obligatorio
			if band_driveletter {
				fmt.Println("El parametro -driveletter ya fue ingresado...")
				band_error = true
				break
			}
			band_driveletter = true
			val_driveletter = fmt.Sprintf("/home/taro/Escritorio/MIA/P1/%s.dsk", val_data)
			fmt.Println("driver letter: ", val_driveletter)
		case strings.Contains(data, "name="): //obligatorio
			if band_name {
				fmt.Println("El parametro -name ya fue ingresado.")
				band_error = true
				break
			}
			band_name = true
			val_name = strings.Replace(val_data, "\"", "", 2)
			fmt.Println("Name: ", val_name)

		//problem system validation
		case strings.Contains(data, "unit="):
			if band_unit {
				fmt.Println("El parametro -unit ya fue ingresado.")
				band_error = true
				break
			}
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)
			fmt.Println("Unit: ", val_unit)
			if val_unit == "b" || val_unit == "k" || val_unit == "m" {
				band_unit = true
			} else {
				fmt.Println("El valor del parametro -unit no es valido.")
				band_error = true
				break
			}
		case strings.Contains(data, "type="):
			if band_type {
				fmt.Println("El parametro -type ya fue ingresado.")
				band_error = true
				break
			}
			val_type = strings.Replace(val_data, "\"", "", 2)
			val_type = strings.ToLower(val_type)
			fmt.Println("Type: ", val_type)
			if val_type == "p" || val_type == "e" || val_type == "l" {
				band_type = true
			} else {
				fmt.Println("El valor del parametro -type no es valido.")
				band_error = true
				break
			}

		case strings.Contains(data, "fit="):
			if band_fit {
				fmt.Println("Parametro -fit ya fue ingresado...")
				band_error = true
				break
			}
			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" { // Best Fit
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" { // First Fit
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" { // Worst Fit
				band_fit = true
				val_fit = "w"
			} else {
				fmt.Println("Valor del parametro -fit no es valido...")
				band_error = true
				break
			}

		case strings.Contains(data, "delete="):
			if band_delete {
				fmt.Println("El parametro -delete ya fue ingresado.")
				band_error = true
				break
			}
			val_delete = strings.Replace(val_delete, "\"", "", 2)
			val_delete = strings.ToLower(val_delete)
			if val_delete != "full" {
				fmt.Println("Error: el valor del parametro -delete debe ser 'full'.")
				band_error = true
				break
			}
		case strings.Contains(data, "add="):
			if band_add {
				fmt.Println("El parametro -add ya fue ingresado.")
				band_error = true
				break
			}
			band_add = true
			val_add = strings.Replace(val_add, "\"", "", 2)
			fmt.Println("Add: ", val_add)
		default:
			fmt.Println("Error parametro no valido.")
		}
	}

	if band_delete && band_add {
		fmt.Println("Error: no pueden venir -add y -delete en el mismo comando.")
		band_error = true
	}

	if !band_size || !band_driveletter || !band_name {
		fmt.Println("Error: faltan parámetros obligatorios.")
		band_error = true
	}

	if band_error {
		fmt.Println("Error en tiempo de ejecución.")
		return
	}

	if band_unit || band_type || band_fit {
		if val_type == "e" {
			fmt.Println("Crear partición extendida.")
		} else if val_type == "l" {
			fmt.Println("Crear partición lógica.")
		} else {
			fmt.Println("Crear partición primaria.")
			crear_particion_primaria(val_driveletter, val_name, val_size, val_fit, val_unit)
		}
	} else {
		fmt.Println("Crear partición primaria.")
		crear_particion_primaria(val_driveletter, val_name, val_size, "", "")
	}

	if band_delete {
		fmt.Println("Borrar partición.")
		// Llamar a función para borrar partición
	}

	if band_add {
		fmt.Println("Agregar espacio a la partición.")
		// Llamar a función para agregar espacio a la partición
	}
}

func eliminar_particion(direccion string, name string, size int) { //val_driveletter, val_name, val_size) //full?
	fmt.Println("Eliminando partición en", direccion)
	fmt.Println("Nombre de la partición:", name)
	fmt.Println("Tamaño de la partición:", size)
}

func agregar_espacio_particion(direccion string, name string, size string, unidad string) { //val_driveletter, val_name, val_add, val_unit, val_size)
	fmt.Println("Agregando espacio a partición en", direccion)
	fmt.Println("Nombre de la partición:", name)
	fmt.Println("Tamaño a agregar:", size, unidad)
}

func quitar_espacio_particion(direccion string, name string, size string, unidad string) { //val_driveletter, val_name, val_add, val_unit, val_size)
	fmt.Println("Quitando espacio a partición en", direccion)
	fmt.Println("Nombre de la partición:", name)
	fmt.Println("Tamaño a quitar:", size, unidad)
}

func crear_particion_primaria(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := mbr{}
	var empty [100]byte

	// Verifico si tiene Ajuste
	if fit != "" {
		aux_fit = fit
	} else {
		// Por default es Peor ajuste
		aux_fit = "w"
	}

	// Verifico si tiene Unidad
	if unit != "" {
		aux_unit = unit

		// *Bytes
		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			// *Kilobytes
			size_bytes = size * 1024
		} else {
			// *Megabytes
			size_bytes = size * 1024 * 1024
		}
	} else {
		// Por default Kilobytes
		size_bytes = size * 1024
	}

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		fmt.Println("[ERROR] No existe un disco duro con ese nombre...")
	} else {
		// Bandera para ver si hay una particion disponible
		band_particion := false
		// Valor del numero de particion
		num_particion := 0

		// Calculo del tamaño de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_start := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
				s_part_start = strings.Trim(s_part_start, "\x00")

				// Verifico si en las particiones hay espacio
				if s_part_start == "-1" {
					band_particion = true
					num_particion = i
					break
				}
			}

			// Si hay una particion disponible
			if band_particion {
				espacio_usado := 0
				s_part_size := ""
				i_part_size := 0
				s_part_status := ""

				// Recorro las 4 particiones
				for i := 0; i < 4; i++ {
					// Obtengo el espacio utilizado
					s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
					// Le quito los caracteres null
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ = strconv.Atoi(s_part_size)

					// Obtengo el status de la particion
					s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
					// Le quito los caracteres null
					s_part_status = strings.Trim(s_part_status, "\x00")

					if s_part_status != "1" {
						// Le sumo el valor al espacio
						espacio_usado += i_part_size
					}
				}

				/* Tamaño del disco */

				// Obtengo el tamaño del disco
				s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
				// Le quito los caracteres null
				s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
				i_tamaño_disco, _ := strconv.Atoi(s_tamaño_disco)

				espacio_disponible := i_tamaño_disco - espacio_usado

				fmt.Println("[ESPACIO DISPONIBLE] ", espacio_disponible, " Bytes")
				fmt.Println("[ESPACIO NECESARIO] ", size_bytes, " Bytes")

				// Verifico que haya espacio suficiente
				if espacio_disponible >= size_bytes {
					// Verifico si no existe una particion con ese nombre
					if !existe_particion(direccion, nombre) {
						// Antes de comparar limpio la cadena
						s_dsk_fit := string(master_boot_record.Dsk_fit[:])
						s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

						/*  Primer Ajuste  */
						if s_dsk_fit == "f" {
							copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

							// Si esta iniciando
							if num_particion == 0 {
								// Guardo el inicio de la particion y dejo un espacio de separacion
								mbr_empty_byte := struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
							} else {
								// Obtengo el inicio de la particion anterior
								s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
								// Le quito los caracteres null
								s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
								i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

								// Obtengo el tamaño de la particion anterior
								s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
								// Le quito los caracteres null
								s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
								i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant+1))
							}

							copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], nombre)

							// Se guarda de nuevo el MBR

							// Conversion de struct a bytes
							mbr_byte := struct_a_bytes(master_boot_record)

							// Se posiciona al inicio del archivo para guardar la informacion del disco
							f.Seek(0, os.SEEK_SET)
							f.Write(mbr_byte)

							// Obtengo el inicio de la particion
							s_part_start = string(master_boot_record.Mbr_partition[num_particion].Part_start[:])
							// Le quito los caracteres null
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							// Se posiciona en el inicio de la particion
							f.Seek(int64(i_part_start), os.SEEK_SET)

							// Lo llena de unos
							for i := 1; i < size_bytes; i++ {
								f.Write([]byte{1})
							}

							fmt.Println("[SUCCES] La Particion primaria fue creada con exito!")
						} else if s_dsk_fit == "b" {
							/*  Mejor Ajuste  */
							best_index := num_particion

							// Variables para conversiones
							s_part_start_act := ""
							s_part_status_act := ""
							s_part_size_act := ""
							i_part_size_act := 0
							s_part_start_best := ""
							i_part_start_best := 0
							s_part_start_best_ant := ""
							i_part_start_best_ant := 0
							s_part_size_best := ""
							i_part_size_best := 0
							s_part_size_best_ant := ""
							i_part_size_best_ant := 0

							for i := 0; i < 4; i++ {
								// Obtengo el inicio de la particion actual
								s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start_act = strings.Trim(s_part_start_act, "\x00")

								// Obtengo el size de la particion actual
								s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
								// Le quito los caracteres null
								s_part_status_act = strings.Trim(s_part_status_act, "\x00")

								// Obtengo la posicion de la particion actual
								s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_size_act = strings.Trim(s_part_size_act, "\x00")
								i_part_size_act, _ = strconv.Atoi(s_part_size_act)

								if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
									if i != num_particion {
										// Obtengo el tamaño de la particion del mejor indice
										s_part_size_best = string(master_boot_record.Mbr_partition[best_index].Part_size[:])
										// Le quito los caracteres null
										s_part_size_best = strings.Trim(s_part_size_best, "\x00")
										i_part_size_best, _ = strconv.Atoi(s_part_size_best)

										// Obtengo la posicion de la particion actual
										s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
										// Le quito los caracteres null
										s_part_size_act = strings.Trim(s_part_size_act, "\x00")
										i_part_size_act, _ = strconv.Atoi(s_part_size_act)

										if i_part_size_best > i_part_size_act {
											best_index = i
											break
										}
									}
								}
							}

							// Primaria
							copy(master_boot_record.Mbr_partition[best_index].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[best_index].Part_fit[:], aux_fit)

							// Si esta iniciando
							if best_index == 0 {
								// Guardo el inicio de la particion y dejo un espacio de separacion
								mbr_empty_byte := struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
							} else {
								// Obtengo el inicio de la particion actual
								s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])
								// Le quito los caracteres null
								s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
								i_part_start_best_ant, _ = strconv.Atoi(s_part_start_best_ant)

								// Obtengo el inicio de la particion actual
								s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])
								// Le quito los caracteres null
								s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
								i_part_size_best_ant, _ = strconv.Atoi(s_part_size_best_ant)

								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant))
							}

							copy(master_boot_record.Mbr_partition[best_index].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[best_index].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[best_index].Part_name[:], nombre)

							// Se guarda de nuevo el MBR

							// Conversion de struct a bytes
							mbr_byte := struct_a_bytes(master_boot_record)

							// Se posiciona al inicio del archivo para guardar la informacion del disco
							f.Seek(0, os.SEEK_SET)
							f.Write(mbr_byte)

							// Obtengo el inicio de la particion best
							s_part_start_best = string(master_boot_record.Mbr_partition[best_index].Part_start[:])
							// Le quito los caracteres null
							s_part_start_best = strings.Trim(s_part_start_best, "\x00")
							i_part_start_best, _ = strconv.Atoi(s_part_start_best)

							// Conversion de struct a bytes

							// Se posiciona en el inicio de la particion
							f.Seek(int64(i_part_start_best), os.SEEK_SET)

							// Lo llena de unos
							for i := 1; i < size_bytes; i++ {
								f.Write([]byte{1})
							}

							fmt.Println("[SUCCES] La Particion primaria fue creada con exito!")
						} else {
							/*  Peor ajuste  */
							worst_index := num_particion

							// Variables para conversiones
							s_part_start_act := ""
							s_part_status_act := ""
							s_part_size_act := ""
							i_part_size_act := 0
							s_part_start_worst := ""
							i_part_start_worst := 0
							s_part_start_worst_ant := ""
							i_part_start_worst_ant := 0
							s_part_size_worst := ""
							i_part_size_worst := 0
							s_part_size_worst_ant := ""
							i_part_size_worst_ant := 0

							for i := 0; i < 4; i++ {
								// Obtengo el inicio de la particion actual
								s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start_act = strings.Trim(s_part_start_act, "\x00")

								// Obtengo el size de la particion actual
								s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
								// Le quito los caracteres null
								s_part_status_act = strings.Trim(s_part_status_act, "\x00")

								// Obtengo la posicion de la particion actual
								s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_size_act = strings.Trim(s_part_size_act, "\x00")
								i_part_size_act, _ = strconv.Atoi(s_part_size_act)

								if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
									if i != num_particion {
										// Obtengo el tamaño de la particion del mejor indice
										s_part_size_worst = string(master_boot_record.Mbr_partition[worst_index].Part_size[:])
										// Le quito los caracteres null
										s_part_size_worst = strings.Trim(s_part_size_worst, "\x00")
										i_part_size_worst, _ = strconv.Atoi(s_part_size_worst)

										// Obtengo la posicion de la particion actual
										s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
										// Le quito los caracteres null
										s_part_size_act = strings.Trim(s_part_size_act, "\x00")
										i_part_size_act, _ = strconv.Atoi(s_part_size_act)

										if i_part_size_worst < i_part_size_act {
											worst_index = i
											break
										}
									}
								}
							}

							// Particiones Primarias
							copy(master_boot_record.Mbr_partition[worst_index].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[worst_index].Part_fit[:], aux_fit)

							// Se esta iniciando
							if worst_index == 0 {
								// Guardo el inicio de la particion y dejo un espacio de separacion
								mbr_empty_byte := struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
							} else {
								// Obtengo el inicio de la particion anterior
								s_part_start_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_start[:])
								// Le quito los caracteres null
								s_part_start_worst_ant = strings.Trim(s_part_start_worst_ant, "\x00")
								i_part_start_worst_ant, _ = strconv.Atoi(s_part_start_worst_ant)

								// Obtengo el tamaño de la particion anterior
								s_part_size_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_size[:])
								// Le quito los caracteres null
								s_part_size_worst_ant = strings.Trim(s_part_size_worst_ant, "\x00")
								i_part_size_worst_ant, _ = strconv.Atoi(s_part_size_worst_ant)

								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(i_part_start_worst_ant+i_part_size_worst_ant))
							}

							copy(master_boot_record.Mbr_partition[worst_index].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[worst_index].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[worst_index].Part_name[:], nombre)

							// Se guarda de nuevo el MBR

							// Conversion de struct a bytes
							mbr_byte := struct_a_bytes(master_boot_record)

							// Escribe desde el inicio del archivo
							f.Seek(0, os.SEEK_SET)
							f.Write(mbr_byte)

							// Obtengo el inicio de la particion best
							s_part_start_worst = string(master_boot_record.Mbr_partition[worst_index].Part_start[:])
							// Le quito los caracteres null
							s_part_start_worst = strings.Trim(s_part_start_worst, "\x00")
							i_part_start_worst, _ = strconv.Atoi(s_part_start_worst)

							// Se posiciona en el inicio de la particion
							f.Seek(int64(i_part_start_worst), os.SEEK_SET)

							// Lo llena de unos
							for i := 1; i < size_bytes; i++ {
								f.Write([]byte{1})
							}

							fmt.Println("[SUCCES] La Particion primaria fue creada con exito!")
						}
					} else {
						fmt.Println("[ERROR] Ya existe una particion creada con ese nombre...")
					}
				} else {
					fmt.Println("[ERROR] La particion que desea crear excede el espacio disponible...")
				}
			} else {
				fmt.Println("[ERROR] La suma de particiones primarias y extendidas no debe exceder de 4 particiones...")
				fmt.Println("[MENSAJE] Se recomienda eliminar alguna particion para poder crear otra particion primaria o extendida")
			}
		} else {
			fmt.Println("[ERROR] el disco se encuentra vacio...")
		}
		f.Close()
	}
}

func existe_particion(direccion string, nombre string) bool {
	extendida := -1
	mbr_empty := mbr{}
	ebr_empty := Ebr{}
	var empty [100]byte
	cont := 0
	fin_archivo := false

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		msg_error(err)
	} else {
		// Procedo a leer el archivo

		// Calculo del tamano de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		_, err = f.ReadAt(lectura, 0)

		// ERROR
		if err != nil && err != io.EOF {
			msg_error(err)
		}

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)
		sstruct = len(lectura)

		// ERROR
		if err != nil {
			msg_error(err)
		}

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_name := ""
			s_part_type := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_name = string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")

				// Verifico si ya existe una particion con ese nombre
				if s_part_name == nombre {

				}

				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				// Verifico si de tipo extendida
				if s_part_type == "E" {
					extendida = i
				}
			}

			// Lo busco en las extendidas
			if extendida != -1 {
				// Obtengo el inicio de la particion
				s_part_start := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
				// Le quito los caracteres null
				s_part_start = strings.Trim(s_part_start, "\x00")
				i_part_start, err := strconv.Atoi(s_part_start)

				// ERROR
				if err != nil {
					msg_error(err)
					fin_archivo = true
				}

				// Obtengo el espacio de la partcion
				s_part_size := string(master_boot_record.Mbr_partition[extendida].Part_size[:])
				// Le quito los caracteres null
				s_part_size = strings.Trim(s_part_size, "\x00")
				i_part_size, err := strconv.Atoi(s_part_size)

				// ERROR
				if err != nil {
					msg_error(err)
					fin_archivo = true
				}

				// Calculo del tamano de struct en bytes
				ebr2 := struct_a_bytes(ebr_empty)
				sstruct := len(ebr2)

				// Lectrura de conjunto de bytes desde el inicio de la particion
				for !fin_archivo {
					// Lectrura de conjunto de bytes en archivo binario
					lectura := make([]byte, sstruct)
					n_leidos, err := f.ReadAt(lectura, int64(sstruct*cont+i_part_start))

					// ERROR
					if err != nil {
						msg_error(err)
						fin_archivo = true
					}

					// Posicion actual en el archivo
					pos_actual, err := f.Seek(0, os.SEEK_CUR)

					// ERROR
					if err != nil {
						msg_error(err)
						fin_archivo = true
					}

					// Si no lee nada y ya se paso del tamaño de la particion
					if n_leidos == 0 && pos_actual < int64(i_part_start+i_part_size) {
						fin_archivo = true
						break
					}

					// Conversion de bytes a struct
					extended_boot_record := bytes_a_struct_ebr(lectura)
					sstruct = len(lectura)

					if err != nil {
						msg_error(err)
					}

					if extended_boot_record.Part_s == empty {
						fin_archivo = true
					} else {
						fmt.Print(" Nombre: ")
						fmt.Print(string(extended_boot_record.Part_name[:]))

						// Antes de comparar limpio la cadena
						s_part_name = string(extended_boot_record.Part_name[:])
						s_part_name = strings.Trim(s_part_name, "\x00")

						// Verifico si ya existe una particion con ese nombre
						if s_part_name == nombre {
							f.Close()
							return true
						}

						// Obtengo el espacio utilizado
						s_part_next := string(extended_boot_record.Part_next[:])
						// Le quito los caracteres null
						s_part_next = strings.Trim(s_part_next, "\x00")
						i_part_next, err := strconv.Atoi(s_part_next)

						// ERROR
						if err != nil {
							msg_error(err)
						}

						// Si ya termino
						if i_part_next != -1 {
							f.Close()
							return false
						}
					}
					cont++
				}
			}
		}
	}
	f.Close()
	return false
}

func mount(commandArray []string) {
	//mount
}

func msg_error(err error) {
	fmt.Println(" ", err)
}

func struct_a_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)

	if err != nil && err != io.EOF {
		msg_error(err)
	}

	return buf.Bytes()
}

// Verifica o crea la ruta para el disco duro
func crear_disco(ruta string) {
	aux, err := filepath.Abs(ruta)

	// ERROR
	if err != nil {
		msg_error(err)
	}

	// Crea el directiorio de forma recursiva
	cmd1 := exec.Command("/bin/sh", "-c", "sudo mkdir -p '"+filepath.Dir(aux)+"'")
	cmd1.Dir = "/"
	_, err = cmd1.Output()

	// ERROR
	if err != nil {
		msg_error(err)
	}

	// Da los permisos al directorio
	cmd2 := exec.Command("/bin/sh", "-c", "sudo chmod -R 777 '"+filepath.Dir(aux)+"'")
	cmd2.Dir = "/"
	_, err = cmd2.Output()

	// ERROR
	if err != nil {
		msg_error(err)
	}

	// Verifica si existe la ruta para el archivo
	if _, err := os.Stat(filepath.Dir(aux)); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			fmt.Println("[FAILURE] No se pudo crear el disco...")
		}
	}
}

func mostrar_mkdisk() {
	fin_archivo := false
	var empty [100]byte
	mbr_empty := mbr{}
	cont := 0

	fmt.Println("Reporte de MKDISK:")
	disco, err := os.OpenFile("/home/taro/Escritorio/Tarea2/Disk1.dsk", os.O_RDWR, 0660) //CAMBIAR

	if err != nil {
		msg_error(err)
	}

	mbr2 := struct_a_bytes(mbr_empty)
	sstruct := len(mbr2)

	for !fin_archivo {
		lectura := make([]byte, sstruct)
		_, err = disco.ReadAt(lectura, int64(sstruct*cont))

		if err != nil && err != io.EOF {
			msg_error(err)
		}

		mbr := bytes_a_struct_mbr(lectura)
		sstruct = len(lectura)

		if err != nil {
			msg_error(err)
		}

		if mbr.Mbr_tamano == empty {
			fin_archivo = true
		} else {
			fmt.Print("Tamaño: ")
			fmt.Println(string(mbr.Mbr_tamano[:]))
			fmt.Print("Fecha: ")
			fmt.Println(string(mbr.Mbr_fecha_creacion[:]))
			fmt.Print("Signature: ")
			fmt.Println(string(mbr.Mbr_dsk_signature[:]))
		}

		cont++
	}
	disco.Close()
}

func bytes_a_struct_mbr(s []byte) mbr {
	p := mbr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	if err != nil && err != io.EOF {
		msg_error(err)
	}

	return p
}

// Decodifica de [] Bytes a Struct
func bytes_a_struct_ebr(s []byte) Ebr {
	p := Ebr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		msg_error(err)
	}

	return p
}

func pause() {
	fmt.Print("Presiona enter para continuar...")
	fmt.Scanln()
}

func execute(commandArray []string) {
	val_path := ""
	band_path := false
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]
		switch {
		case strings.Contains(data, "path="):
			if band_path {
				fmt.Println("Parametro -path ya fue ingresado...")
				break
			}

			band_path = true
			val_path = strings.Replace(val_data, "\"", "", 2)
			cargarArchivo(val_path)
		default:
			fmt.Println(" Parametro no valido...")
		}
	}

}

func cargarArchivo(ruta string) {
	file, err := os.Open(ruta)
	if err != nil {
		fmt.Println(err)
		return
	}
	scanner := bufio.NewScanner(file)
	fmt.Println("Analizando comandos....")
	for scanner.Scan() {
		split_comando(scanner.Text())
	}
	file.Close()
}
