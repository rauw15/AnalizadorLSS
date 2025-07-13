package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type NodoSintactico struct {
	Nodo  string                 `json:"nodo"`
	Tipo  string                 `json:"tipo,omitempty"`
	Valor string                 `json:"valor,omitempty"`
	Hijos []NodoSintactico       `json:"hijos,omitempty"`
	Extra map[string]interface{} `json:"extra,omitempty"`
}

var palabrasReservadasBash = []string{
	"if", "then", "else", "elif", "fi", "for", "while", "do", "done", "function", "case", "esac", "in", "select", "until", "time", "coproc", "break", "continue", "return", "exit", "echo", "read", "declare", "local", "export", "let", "test", "shift", "unset", "trap", "source", "exec", "set", "eval", "wait", "bg", "fg", "kill",
}

// Operadores de test de Bash
var operadoresTestBash = []string{
	"eq", "ne", "lt", "le", "gt", "ge",
}

func distanciaLevenshtein(a, b string) int {
	la := len(a)
	lb := len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	dp := make([][]int, la+1)
	for i := range dp {
		dp[i] = make([]int, lb+1)
	}
	for i := 0; i <= la; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			coste := 0
			if a[i-1] != b[j-1] {
				coste = 1
			}
			dp[i][j] = min(
				dp[i-1][j]+1,
				dp[i][j-1]+1,
				dp[i-1][j-1]+coste,
			)
		}
	}
	return dp[la][lb]
}

func min(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

func sugerirPalabraReservada(palabra string) string {
	minDist := math.MaxInt32
	mejor := ""
	for _, pr := range palabrasReservadasBash {
		d := distanciaLevenshtein(palabra, pr)
		if d < minDist {
			minDist = d
			mejor = pr
		}
	}
	if minDist <= 2 { // Solo sugerir si es muy parecido
		return mejor
	}
	return ""
}

// Elimina texto dentro de comillas simples y dobles
func eliminarStrings(linea string) string {
	// Elimina comillas dobles
	linea = regexp.MustCompile(`"([^"\\]|\\.)*"`).ReplaceAllString(linea, "")
	// Elimina comillas simples
	linea = regexp.MustCompile(`'([^'\\]|\\.)*'`).ReplaceAllString(linea, "")
	return linea
}

func AnalisisSintactico(codigo string) (NodoSintactico, []string) {
	errores := []string{}
	lineas := strings.Split(codigo, "\n")
	arbol := NodoSintactico{Nodo: "script", Hijos: []NodoSintactico{}}

	// Pilas para bloques
	pilaIf := 0
	pilaWhile := 0
	pilaFor := 0
	pilaLlaves := 0

	for i, linea := range lineas {
		l := strings.TrimSpace(linea)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}

		// Comillas balanceadas
		numComillasSimples := strings.Count(l, "'")
		numComillasDobles := strings.Count(l, "\"")
		if numComillasSimples%2 != 0 {
			errores = append(errores, "Línea "+itoa(i+1)+": Comillas simples no balanceadas")
		}
		if numComillasDobles%2 != 0 {
			errores = append(errores, "Línea "+itoa(i+1)+": Comillas dobles no balanceadas")
		}

		// Sugerencia de palabra reservada mal escrita (ignorando strings y operadores de test)
		sinStrings := eliminarStrings(l)
		palabras := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`).FindAllString(sinStrings, -1)
		for _, palabra := range palabras {
			if contiene(palabrasReservadasBash, palabra) || contiene(operadoresTestBash, palabra) {
				continue
			}
			// No sugerir para palabras que están dentro de $(( ... )) (expresiones aritméticas)
			if strings.HasPrefix(palabra, "$((") && strings.HasSuffix(palabra, "))") {
				continue
			}
			sugerencia := sugerirPalabraReservada(palabra)
			if sugerencia != "" {
				errores = append(errores, "Línea "+itoa(i+1)+": Palabra reservada desconocida '"+palabra+"'. ¿Quizás quisiste escribir '"+sugerencia+"'?")
			}
		}

		// Asignación: VAR=valor
		asig := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)=(.*)$`)
		if asig.MatchString(l) {
			m := asig.FindStringSubmatch(l)
			arbol.Hijos = append(arbol.Hijos, NodoSintactico{
				Nodo:  "asignacion",
				Tipo:  "asignacion",
				Valor: m[1],
				Extra: map[string]interface{}{"valor": strings.TrimSpace(m[2])},
			})
			continue
		}

		// If
		ifif := regexp.MustCompile(`^if\s+(.+)\s*;?\s*then$`)
		if ifif.MatchString(l) {
			m := ifif.FindStringSubmatch(l)
			arbol.Hijos = append(arbol.Hijos, NodoSintactico{
				Nodo:  "if",
				Tipo:  "condicional",
				Valor: m[1],
			})
			pilaIf++
			continue
		}
		if l == "fi" {
			if pilaIf > 0 {
				pilaIf--
			} else {
				errores = append(errores, "Línea "+itoa(i+1)+": 'fi' sin 'if' correspondiente")
			}
			arbol.Hijos = append(arbol.Hijos, NodoSintactico{Nodo: "fin_if"})
			continue
		}

		// While
		wh := regexp.MustCompile(`^while\s+(.+)\s*;?\s*do$`)
		if wh.MatchString(l) {
			m := wh.FindStringSubmatch(l)
			arbol.Hijos = append(arbol.Hijos, NodoSintactico{
				Nodo:  "while",
				Tipo:  "bucle",
				Valor: m[1],
			})
			pilaWhile++
			continue
		}
		if l == "done" {
			if pilaWhile > 0 {
				pilaWhile--
			} else if pilaFor > 0 {
				pilaFor--
			} else {
				errores = append(errores, "Línea "+itoa(i+1)+": 'done' sin 'while' o 'for' correspondiente")
			}
			arbol.Hijos = append(arbol.Hijos, NodoSintactico{Nodo: "fin_while"})
			continue
		}

		// For (permitir guion bajo en nombre de variable)
		forr := regexp.MustCompile(`^for\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+in\s+(.+)\s*;?\s*do$`)
		if forr.MatchString(l) {
			m := forr.FindStringSubmatch(l)
			arbol.Hijos = append(arbol.Hijos, NodoSintactico{
				Nodo:  "for",
				Tipo:  "bucle",
				Valor: m[1],
				Extra: map[string]interface{}{"en": m[2]},
			})
			pilaFor++
			continue
		}

		// Función (permitir guion bajo en nombre)
		fun := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*\(\)\s*\{$`)
		if fun.MatchString(l) {
			m := fun.FindStringSubmatch(l)
			arbol.Hijos = append(arbol.Hijos, NodoSintactico{
				Nodo:  "funcion",
				Tipo:  "definicion_funcion",
				Valor: m[1],
			})
			pilaLlaves++
			continue
		}
		if l == "{" {
			pilaLlaves++
			continue
		}
		if l == "}" {
			if pilaLlaves > 0 {
				pilaLlaves--
			} else {
				errores = append(errores, "Línea "+itoa(i+1)+": '}' sin '{' correspondiente")
			}
			arbol.Hijos = append(arbol.Hijos, NodoSintactico{Nodo: "fin_bloque"})
			continue
		}

		// Comando
		cmd := regexp.MustCompile(`^[a-zA-Z0-9_\-\.\$\{\}\[\]=/\":' ]+$`)
		if cmd.MatchString(l) {
			arbol.Hijos = append(arbol.Hijos, NodoSintactico{
				Nodo:  "comando",
				Tipo:  "comando",
				Valor: l,
			})
			continue
		}

		errores = append(errores, "Línea "+itoa(i+1)+": Sintaxis no reconocida: '"+l+"'")
	}

	// Al final, verificar bloques abiertos
	if pilaIf > 0 {
		errores = append(errores, "Bloque 'if' sin 'fi' de cierre")
	}
	if pilaWhile > 0 {
		errores = append(errores, "Bloque 'while' sin 'done' de cierre")
	}
	if pilaFor > 0 {
		errores = append(errores, "Bloque 'for' sin 'done' de cierre")
	}
	if pilaLlaves > 0 {
		errores = append(errores, "Bloque '{' sin '}' de cierre")
	}

	return arbol, errores
}

func contiene(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
