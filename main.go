package main

import (
	"MIA1_P1_201602404/Comandos"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var logued = false

func main() {
	for {
		fmt.Println("***** KEVIN ESTUARDO SECAIDA MOLINA ***** ")
		fmt.Println("Ingrese un comando: ")

		reader := bufio.NewReader(os.Stdin)
		entrada, _ := reader.ReadString('\n')
		eleccion := strings.TrimRight(entrada, "\r\n")
		if eleccion == "exit" {
			break
		}
		comando := Comando(eleccion)
		eleccion = strings.TrimSpace(eleccion)
		eleccion = strings.TrimLeft(eleccion, comando)
		tokens := SepararTokens(eleccion)
		funciones(comando, tokens)
		fmt.Println("\tPresione Enter para continuar....")
		fmt.Scanln()
	}
}

func Comando(text string) string {
	var tkn string
	terminar := false
	for i := 0; i < len(text); i++ {
		if terminar {
			if string(text[i]) == " " || string(text[i]) == "-" {
				break
			}
			tkn += string(text[i])
		} else if string(text[i]) != " " && !terminar {
			if string(text[i]) == "#" {
				tkn = text
			} else {
				tkn += string(text[i])
				terminar = true
			}
		}
	}
	return tkn
}

func SepararTokens(texto string) []string {
	var tokens []string
	if texto == "" {
		return tokens
	}
	texto += " "
	var token string
	estado := 0
	for i := 0; i < len(texto); i++ {
		c := string(texto[i])
		if estado == 0 && c == "-" {
			estado = 1
		} else if estado == 0 && c == "#" {
			continue
		} else if estado != 0 {
			if estado == 1 {
				if c == "=" {
					estado = 2
				} else if c == " " {
					continue
				} else if (c == "P" || c == "p") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
					continue
				} else if (c == "R" || c == "r") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
					continue
				}
			} else if estado == 2 {
				if c == " " {
					continue
				}
				if c == "\"" {
					estado = 3
					continue
				} else {
					estado = 4
				}
			} else if estado == 3 {
				if c == "\"" {
					estado = 4
					continue
				}
			} else if estado == 4 && c == "\"" {
				tokens = []string{}
				continue
			} else if estado == 4 && c == " " {
				estado = 0
				tokens = append(tokens, token)
				token = ""
				continue
			}
			token += c
		}
	}
	return tokens
}

func funciones(token string, tks []string) {
	if token != "" {
		if Comandos.Comparar(token, "EXEC") {
			FuncionExec(tks)
		} else if Comandos.Comparar(token, "MKDISK") {
			Comandos.ValidarDatosMKDISK(tks)
		} else if Comandos.Comparar(token, "RMDISK") {
			Comandos.RMDISK(tks)
		} else if Comandos.Comparar(token, "FDISK") {
			//Comandos.ValidarDatosFDISK(tks)
		} else if Comandos.Comparar(token, "MOUNT") {
			//Comandos.ValidarDatosMOUNT(tks)
		} else if Comandos.Comparar(token, "MKFS") {
			//	Comandos.ValidarDatosMKFS(tks)
		} else if Comandos.Comparar(token, "LOGIN") {
			if logued {
				Comandos.Error("LOGIN", "Ya hay un usuario en línea.")
				return
			} else {
				fmt.Print("hola1")
				//		logued = Comandos.ValidarDatosLOGIN(tks)
			}
		} else if Comandos.Comparar(token, "LOGOUT") {
			if !logued {
				Comandos.Error("LOGOUT", "Inicie sesión, por favor.")
				return
			} else {
				fmt.Print("hola1")
				//logued = Comandos.CerrarSesion()
			}
		} else if Comandos.Comparar(token, "MKGRP") {
			if !logued {
				Comandos.Error("MKGRP", "Inicie sesión, por favor.")
				return
			} else {
				fmt.Print("hola1")
				//Comandos.ValidarDatosGrupos(tks, "MK")
			}
		} else if Comandos.Comparar(token, "RMGRP") {
			if !logued {
				Comandos.Error("RMGRP", "Inicie sesión, por favor.")
				return
			} else {
				fmt.Print("hola1")
				//Comandos.ValidarDatosGrupos(tks, "RM")
			}
		} else if Comandos.Comparar(token, "MKUSER") {
			if !logued {
				Comandos.Error("MKUSER", "Inicie sesión, por favor.")
				return
			} else {
				fmt.Print("hola1")
				//	Comandos.ValidarDatosUsers(tks, "MK")
			}
		} else if Comandos.Comparar(token, "RMUSR") {
			if !logued {
				Comandos.Error("RMUSER", "Inicie sesión, por favor.")
				return
			} else {
				fmt.Print("hola1")
				//Comandos.ValidarDatosUsers(tks, "RM")
			}
		} else {
			Comandos.Error("ANALIZADOR", "No se reconoce el comando \""+token+"\"")
		}
	}
}

func FuncionExec(tokens []string) {
	path := ""
	for i := 0; i < len(tokens); i++ {
		datos := strings.Split(tokens[i], "=")
		if Comandos.Comparar(datos[0], "path") {
			path = datos[1]
		}
	}
	if path == "" {
		Comandos.Error("EXEC", "Se requiere el parámetro \"path\" para este comando")
		return
	}
	Exec(path)
}

func Exec(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		texto := fileScanner.Text()
		texto = strings.TrimSpace(texto)
		tk := Comando(texto)
		if texto != "" {
			if Comandos.Comparar(tk, "pause") {
				var pause string
				Comandos.Mensaje("Presione \"enter\" para continuar...", " ")
				fmt.Scanln(&pause)
				continue
			} else if string(texto[0]) == "#" {
				Comandos.Mensaje("COMENTARIO", texto)
				continue
			}
			texto = strings.TrimLeft(texto, tk)
			tokens := SepararTokens(texto)
			funciones(tk, tokens)
		}
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error al leer el archivo: %s", err)
	}
	file.Close()
}
