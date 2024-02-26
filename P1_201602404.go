package main

import (
	"bufio"
	"bytes"
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
	Mbr_partition_1    partition
	Mbr_partition_2    partition
	Mbr_partition_3    partition
	Mbr_partition_4    partition
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
	Part_mount [100]byte
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
				copy(master_boot_record.Mbr_partition_1.Part_status[:], "0")
				copy(master_boot_record.Mbr_partition_1.Part_type[:], "0")
				copy(master_boot_record.Mbr_partition_1.Part_fit[:], "0")
				copy(master_boot_record.Mbr_partition_1.Part_start[:], "-1")
				copy(master_boot_record.Mbr_partition_1.Part_size[:], "0")
				copy(master_boot_record.Mbr_partition_1.Part_name[:], "")

				copy(master_boot_record.Mbr_partition_2.Part_status[:], "0")
				copy(master_boot_record.Mbr_partition_2.Part_type[:], "0")
				copy(master_boot_record.Mbr_partition_2.Part_fit[:], "0")
				copy(master_boot_record.Mbr_partition_2.Part_start[:], "-1")
				copy(master_boot_record.Mbr_partition_2.Part_size[:], "0")
				copy(master_boot_record.Mbr_partition_2.Part_name[:], "")

				copy(master_boot_record.Mbr_partition_3.Part_status[:], "0")
				copy(master_boot_record.Mbr_partition_3.Part_type[:], "0")
				copy(master_boot_record.Mbr_partition_3.Part_fit[:], "0")
				copy(master_boot_record.Mbr_partition_3.Part_start[:], "-1")
				copy(master_boot_record.Mbr_partition_3.Part_size[:], "0")
				copy(master_boot_record.Mbr_partition_3.Part_name[:], "")

				copy(master_boot_record.Mbr_partition_4.Part_status[:], "0")
				copy(master_boot_record.Mbr_partition_4.Part_type[:], "0")
				copy(master_boot_record.Mbr_partition_4.Part_fit[:], "0")
				copy(master_boot_record.Mbr_partition_4.Part_start[:], "-1")
				copy(master_boot_record.Mbr_partition_4.Part_size[:], "0")
				copy(master_boot_record.Mbr_partition_4.Part_name[:], "")

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

func crear_disco(ruta string) {
	aux, err := filepath.Abs(ruta)

	if err != nil {
		msg_error(err)
	}

	cmd1 := exec.Command("/bin/sh", "-c", "sudo mkdir -p '"+filepath.Dir(aux)+"'")
	cmd1.Dir = "/"
	_, err1 := cmd1.Output()

	if err1 != nil {
		msg_error(err)
	}

	cmd2 := exec.Command("/bin/sh", "-c", "sudo chmod -R 777 '"+filepath.Dir(aux)+"'")
	cmd2.Dir = "/"
	_, err2 := cmd2.Output()

	if err2 != nil {
		msg_error(err)
	}

	if _, err := os.Stat(filepath.Dir(aux)); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			fmt.Println("No se pudo crear el disco...")
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
