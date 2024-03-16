package utilities

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

func Crear_Archivo(name string) error {
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err crear_Archivo dir==", err)
		return err
	}

	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err crear_Archivo create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

func ConvertToZeros(filename string, start int64, end int64) error {
	file, err := OpenFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	if start >= fileSize || end >= fileSize || start > end {
		return fmt.Errorf("posición final o inicial, invalida. Tamaño: %d", fileSize)
	}

	zeroBytes := make([]byte, end-start+1)

	_, err = file.WriteAt(zeroBytes, start)
	if err != nil {
		return err
	}

	return nil
}

func CalcularPorcentaje(tamanoParticion int64, tamanoDisco int64) int64 {
	return (int64(tamanoParticion) * 100 / int64(tamanoDisco))
}
