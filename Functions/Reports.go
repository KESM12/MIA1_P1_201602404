package functions

import (
	structs "P1/Structs"
	utilities "P1/Utilities"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func ProcessREP(input string, name *string, path *string, id *string, ruta *string) {
	fmt.Println("inputdsafa", input)
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "name":
			*name = flagValue
		case "path":
			*path = flagValue
		case "id":
			*id = flagValue
		case "ruta":
			*ruta = flagValue
		default:
			fmt.Println("Error faltan parametros.")
		}
	}
}

func GenerateReports(name *string, path *string, id *string, ruta *string) {
	fmt.Println("Generando reporte:", *name, "en", *path, "de", *id, "en", *ruta)
	switch *name {
	case "mbr":
		REPORT_MBR(id, path)
	case "disk":
		REPORT_DISK(id, path)
	case "inode":
		REPORT_INODE(id, path)
	case "Journaling":
		REPORT_JOURNALING()
	case "block":
		REPORT_BLOCK(id, path)
	case "bm_inode":
		REPORT_BM_INODE(id, path)
	case "bm_block":
		REPORT_BM_BLOCK(id, path)
	case "tree":
		REPORT_TREE()
	case "sb":
		REPORT_SB(id, path)
	case "file":
		REPORT_FILE(id, path, ruta)
	case "ls":
		REPORT_LS(id, path, ruta)
	default:
		println("Reporte no reconocido:", *name)
	}
}

func REPORT_MBR(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)
	uniqueNumber := time.Now().UnixNano()

	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	var TempMBR structs.MBR
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var EPartition = false
	var EPartitionStart int

	var compareMBR structs.MBR
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")

	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	strP := ""
	strE := ""

	for _, partition := range TempMBR.Mbr_particion {
		partNameClean := strings.Trim(string(partition.Part_name[:]), "\x00")
		if partition.Part_correlative == 0 {
			continue
		} else {
			strP += fmt.Sprintf(`
		|Particion %d
		|{part_status|%s}
		|{part_type|%s}
		|{part_fit|%s}
		|{part_start|%d}
		|{part_size|%d}
		|{part_name|%s}`,
				partition.Part_correlative,
				string(partition.Part_status[:]),
				string(partition.Part_type[:]),
				string(partition.Part_fit[:]),
				partition.Part_start,
				partition.Part_size,
				partNameClean,
			)
		}

		//?EBR verificacion
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) && EPartition {
			var x = 0
			for x < 1 {
				var TempEBR structs.EBR
				if err := utilities.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if EPartitionStart != 0 && TempEBR.Part_next != -1 {
					partNameClean := strings.Trim(string(TempEBR.Part_name[:]), "\x00")
					strE += fmt.Sprintf(`
		|Particion Logica
		|{part_status|%s}
		|{part_next|%d}
		|{part_fit|%s}
		|{part_start|%d}
		|{part_size|%d}
		|{part_name|%s}`,
						string(TempEBR.Part_mount[:]),
						TempEBR.Part_next,
						string(TempEBR.Part_fit[:]),
						TempEBR.Part_start,
						TempEBR.Part_s,
						partNameClean,
					)
					print("fit logica")
					println(string(TempEBR.Part_fit[:]))
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					print("fit logica")
					println(string(TempEBR.Part_fit[:]))
					partNameClean := strings.Trim(string(TempEBR.Part_name[:]), "\x00")
					strE += fmt.Sprintf(`
		|Particion Logica
		|{part_status|%s}
		|{part_next|%d}
		|{part_fit|%s}
		|{part_start|%d}
		|{part_size|%d}
		|{part_name|%s}`,
						string(TempEBR.Part_mount[:]),
						TempEBR.Part_next,
						string(TempEBR.Part_fit[:]),
						TempEBR.Part_start,
						TempEBR.Part_s,
						partNameClean,
					)
					strP += strE
					x = 1
				}
			}

		}

	}

	dotCode := fmt.Sprintf(`
	digraph G {
		fontname="Helvetica,Arial,sans-serif"
		node [fontname="Helvetica,Arial,sans-serif", style="filled", color="lightblue", shape="record"]
		edge [fontname="Helvetica,Arial,sans-serif"]
		concentrate=True;
		rankdir=TB;
	
		title [label="Reporte MBR" shape=plaintext fontname="Helvetica,Arial,sans-serif" color="darkorange1" fontcolor="darkorange4"];
	
		mbr[label="
			{MBR: %s.dsk|
				{mbr_tamaño|%d}
				|{mbr_fecha_creacion|%s}
				|{mbr_disk_signature|%d}
				%s
			}
		"] [color="lightgoldenrod" fontcolor="darkgoldenrod"];
	
		title2 [label="Reporte EBR" shape=plaintext fontname="Helvetica,Arial,sans-serif" color="darkgreen" fontcolor="darkolivegreen1"];
		
		ebr[label="
			{EBR%s}
		"] [color="palegreen1" fontcolor="darkgreen"];
	
		title -> mbr [style=invis];
		mbr -> title2[style=invis];
		title2 -> ebr[style=invis];
	}`,
		letra,
		TempMBR.Mbr_tamano,
		TempMBR.Mbr_fecha_creacion,
		TempMBR.Mbr_dsk_signature,
		strP,
		strE,
	)

	dotFilePath := fmt.Sprintf("/home/taro/go/src/MIA1_P1_201602404/MIA/Reportes/MBR_%d.dot", uniqueNumber)
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	pngFilePath := *path
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("Reporte MBR, EBR generado en", pngFilePath)
}

func REPORT_DISK(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	var TempMBR structs.MBR
	// Read object from bin file
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var EPartition = false
	var EPartitionStart int

	var compareMBR structs.MBR
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")

	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	strP := ""

	for _, partition := range TempMBR.Mbr_particion {
		if partition.Part_correlative == 0 {
			porcentaje := utilities.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			strP += fmt.Sprintf(`|Libre\n%d%%`, porcentaje)
		}

		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[0].Part_type[:]) {
			porcentaje := utilities.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			strP += fmt.Sprintf(`|Primaria\n%d%%`, porcentaje)
		}

		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) && EPartition {
			porcentaje := utilities.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			println("extendida size")
			println(partition.Part_size)
			strP += fmt.Sprintf(`|{Extendida %d%%|{`, porcentaje)
			var x = 0
			for x < 1 {
				var TempEBR structs.EBR
				if err := utilities.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_next != -1 {
					if !bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						porcentaje := utilities.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|EBR|Particion logica %d%%`, porcentaje)
					} else {
						porcentaje := utilities.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
					}
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					if !bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						porcentaje := utilities.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|EBR|Particion logica %d%%`, porcentaje)
					} else {
						porcentaje := utilities.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
					}
					x = 1
				}
			}
			strP += "}}"
		}

	}

	dotCode := fmt.Sprintf(`
    digraph G {
        graph [bgcolor="#F8F8F8"]; // Fondo gris claro
        node [fontname="Helvetica,Arial,sans-serif", shape=record, style=filled, fillcolor="#FFFFFF", color="#000000", penwidth=2]; // Estilo de los nodos
        edge [fontname="Helvetica,Arial,sans-serif", color="#555555", penwidth=2]; // Estilo de las aristas
        concentrate=true; // Mejora la legibilidad
        rankdir=TB; // Orientación de arriba hacia abajo
        fontsize=12; // Tamaño de fuente

        // Definir nodos
        title [label="Reporte DISK %s" fontsize=16 fontcolor="#333333"];
        dsk[label=<
            <table border="0" cellborder="1" cellspacing="0" cellpadding="4" bgcolor="#FFFFFF">
                <tr><td colspan="2" bgcolor="#FFD700"><b>MBR</b></td></tr>
                %s
            </table>
        >];

        // Definir relaciones
        title -> dsk [style=invis];
    }
`,
		letra,
		strP,
	)

	// Escribir el contenido en el archivo DOT
	dotFilePath := "/home/taro/go/src/MIA1_P1_201602404/MIA/Reportes/DISK.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	// Llamar a Graphviz para generar el gráfico
	pngFilePath := *path
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	structs.PrintMBR(TempMBR)
}

func REPORT_INODE(id *string, path *string) {

}

func REPORT_BLOCK(id *string, path *string) {

}

func REPORT_BM_INODE(id *string, path *string) {

}

func REPORT_BM_BLOCK(id *string, path *string) {

}

func REPORT_TREE() {

}

func REPORT_SB(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	filepath := "/home/taro/go/src/MIA1_P1_201602404/MIA/P1/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	var TempMBR structs.MBR
	if err := utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
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
		fmt.Println("Error reading superblock:", err)
		return
	}

	dotCode := fmt.Sprintf(`
    digraph G {
        graph [bgcolor="#E6E6FA"];
        node [fontname="Helvetica,Arial,sans-serif", shape=record, style=filled, fillcolor="#FFFFFF", color="#000000", penwidth=2];
        edge [fontname="Helvetica,Arial,sans-serif", color="#000000", penwidth=2];
        concentrate=true;
        rankdir=TB;

        title [label="Reporte SUPERBLOCK" shape=plaintext fontname="Helvetica,Arial,sans-serif" fontsize=16 fontcolor="#333333"];

        sb[label=<
            <table border="0" cellborder="1" cellspacing="0" cellpadding="4" bgcolor="#FFFFFF">
                <tr><td colspan="2" bgcolor="#87CEFA"><b>Superblock</b></td></tr>
                <tr><td><b>S_filesystem_type</b></td><td>%d</td></tr>
                <tr><td><b>S_inodes_count</b></td><td>%d</td></tr>
                <tr><td><b>S_blocks_count</b></td><td>%d</td></tr>
                <tr><td><b>S_free_blocks_count</b></td><td>%d</td></tr>
                <tr><td><b>S_free_inodes_count</b></td><td>%d</td></tr>
                <tr><td><b>S_mtime</b></td><td>%s</td></tr>
                <tr><td><b>S_umtime</b></td><td>%s</td></tr>
                <tr><td><b>S_mnt_count</b></td><td>%d</td></tr>
                <tr><td><b>S_magic</b></td><td>%d</td></tr>
                <tr><td><b>S_inode_size</b></td><td>%d</td></tr>
                <tr><td><b>S_block_size</b></td><td>%d</td></tr>
                <tr><td><b>S_fist_ino</b></td><td>%d</td></tr>
                <tr><td><b>S_first_blo</b></td><td>%d</td></tr>
                <tr><td><b>S_bm_inode_start</b></td><td>%d</td></tr>
                <tr><td><b>S_bm_block_start</b></td><td>%d</td></tr>
                <tr><td><b>S_inode_start</b></td><td>%d</td></tr>
                <tr><td><b>S_block_start</b></td><td>%d</td></tr>
            </table>
        >];

        title -> sb [style=invis];
    }
`,
		int(tempSuperblock.S_filesystem_type),
		int(tempSuperblock.S_inodes_count),
		int(tempSuperblock.S_blocks_count),
		int(tempSuperblock.S_free_blocks_count),
		int(tempSuperblock.S_free_inodes_count),
		tempSuperblock.S_mtime[:],
		tempSuperblock.S_umtime[:],
		int(tempSuperblock.S_mnt_count),
		int(tempSuperblock.S_magic),
		int(tempSuperblock.S_inode_size),
		int(tempSuperblock.S_block_size),
		int(tempSuperblock.S_fist_ino),
		int(tempSuperblock.S_first_blo),
		int(tempSuperblock.S_bm_inode_start),
		int(tempSuperblock.S_bm_block_start),
		int(tempSuperblock.S_inode_start),
		int(tempSuperblock.S_block_start),
	)

	// Escribir el contenido en el archivo DOT
	dotFilePath := "/home/taro/go/src/MIA1_P1_201602404/MIA/Reportes/SB.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	pngFilePath := *path
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico del SB:", err)
		return
	}

	fmt.Println("Reporte MBR, EBR generado en", pngFilePath)
}

func REPORT_FILE(id *string, path *string, ruta *string) {

}

func REPORT_LS(id *string, path *string, ruta *string) {

}

func REPORT_JOURNALING() {

}
