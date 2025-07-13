package main

import (
	"regexp"
	"strings"
)

type Token struct {
	Tipo  string `json:"tipo"`
	Valor string `json:"valor"`
}

var patrones = map[string]string{
	"palabra_clave": `\b(if|then|else|elif|fi|for|while|do|done|function|case|esac|in|select|until|time|coproc|break|continue|return|exit|echo|read|declare|local|export|let|test|shift|unset|trap|source|exec|set|eval|wait|bg|fg|kill)\b`,
	"identificador": `\$?[a-zA-Z_][a-zA-Z0-9_]*|\${[a-zA-Z_][a-zA-Z0-9_]*}`,
	"numero":        `\b\d+(\.\d+)?\b`,
	"cadena":        `'[^'\\n\\r]*'|"([^"\\n\\r]*)"`,
	"operador":      `==|!=|=|\+=|-=|\*=|/=|%=|\|\||&&|!|\||;|&|\+|-|\*|/|%|<|>|<=|>=|-eq|-ne|-lt|-gt|-le|-ge`,
	"simbolo":       `[{}\[\]()<>]`,
	"comentario":    `#[^\n]*`,
}

// Palabras clave de Bash
var palabrasClave = []string{
	"if", "then", "else", "elif", "fi", "for", "while", "do", "done", "function", "case", "esac", "in", "select", "until", "time", "coproc", "break", "continue", "return", "exit", "echo", "read", "declare", "local", "export", "let", "test", "shift", "unset", "trap", "source", "exec", "set", "eval", "wait", "bg", "fg", "kill",
}

func esPalabraClave(palabra string) bool {
	for _, kw := range palabrasClave {
		if palabra == kw {
			return true
		}
	}
	return false
}

func eliminarCadenasDelCodigo(codigo string) string {
	// Elimina los literales de texto entre comillas simples y dobles
	codigo = regexp.MustCompile(`'[^']*'`).ReplaceAllString(codigo, "")
	codigo = regexp.MustCompile(`"[^"]*"`).ReplaceAllString(codigo, "")
	return codigo
}

func AnalisisLexico(codigo string) ([]Token, []string) {
	var tokens []Token
	var errores []string

	patronMaestro := ""
	for _, p := range patrones {
		patronMaestro += "(" + p + ")|"
	}
	patronMaestro = strings.TrimSuffix(patronMaestro, "|")

	re := regexp.MustCompile(patronMaestro)
	coincidencias := re.FindAllString(codigo, -1)

	for _, m := range coincidencias {
		tipoToken := "desconocido"
		for tTipo, pat := range patrones {
			if matched, _ := regexp.MatchString("^"+pat+"$", m); matched {
				tipoToken = tTipo
				break
			}
		}
		tokens = append(tokens, Token{Tipo: tipoToken, Valor: m})
	}

	// Revisión de errores léxicos (ignorando palabras dentro de cadenas)
	codigoSinCadenas := eliminarCadenasDelCodigo(codigo)
	palabras := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`).FindAllString(codigoSinCadenas, -1)
	for _, palabra := range palabras {
		esTokenizado := false
		for _, tok := range tokens {
			if tok.Valor == palabra {
				esTokenizado = true
				break
			}
		}
		if !esPalabraClave(palabra) && !esTokenizado {
			if !regexp.MustCompile(`^\d+$`).MatchString(palabra) {
				errores = append(errores, "Posible error léxico: '"+palabra+"' no es una palabra clave válida ni identificador reconocido")
			}
		}
	}

	return tokens, errores
}

// Clasificación para la tabla
func ClasificarToken(t Token) string {
	switch t.Tipo {
	case "palabra_clave":
		return "Palabra Clave"
	case "identificador":
		return "Identificador"
	case "numero":
		return "Número"
	case "cadena":
		return "Cadena"
	case "operador":
		return "Operador"
	case "simbolo":
		return "Símbolo"
	case "comentario":
		return "Comentario"
	default:
		return "Error"
	}
}

func ResumenLexico(codigo string) map[string]interface{} {
	tokens, errores := AnalisisLexico(codigo)

	tabla := []map[string]interface{}{}
	conteos := map[string]int{"Palabra Clave": 0, "Identificador": 0, "Número": 0, "Cadena": 0, "Operador": 0, "Símbolo": 0, "Comentario": 0, "Error": 0}

	for _, t := range tokens {
		cat := ClasificarToken(t)
		fila := map[string]interface{}{
			"valor":   t.Valor,
			"tipo":    cat,
			"esError": cat == "Error",
		}
		tabla = append(tabla, fila)
		if cat != "Error" {
			conteos[cat]++
		}
	}
	conteos["Error"] += len(errores)

	return map[string]interface{}{
		"filas":           tabla,
		"totales":         conteos,
		"errores_lexicos": errores,
	}
}
