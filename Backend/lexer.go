package main

import (
	"regexp"
	"strings"
)

type Token struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

var patterns = map[string]string{
	"keyword":    `\b(public|class|static|void|main|String|int|double|float|char|boolean|if|else|for|while|return|System|out|println|new)\b`,
	"identifier": `\b[a-zA-Z_][a-zA-Z0-9_]*\b`,
	"number":     `\b\d+(\.\d+)?\b`,
	"string":     `"([^"\\n\\r]*)"`,
	"symbol":     `[{}\[\]();=.,<>:+\-*/!]`,
	"dot":        `[.]`,
}

// Lista de palabras clave esperadas (para detectar errores como "imprt")
var expectedKeywords = []string{
	"public", "class", "static", "void", "main", "String", "int", "double", "float", "char", "boolean",
	"if", "else", "for", "while", "return", "System", "out", "println", "new",
}

// Palabras reservadas para Java
var reservedWords = []string{
	"public", "class", "static", "void", "main", "String", "int", "double", "float", "char", "boolean",
	"if", "else", "for", "while", "return", "System", "out", "println", "new",
}

// Verifica si un string es una palabra clave esperada
func isExpectedKeyword(word string) bool {
	for _, kw := range expectedKeywords {
		if word == kw {
			return true
		}
	}
	return false
}

// Verifica si un string es una palabra reservada
func isReservedWord(word string) bool {
	for _, kw := range reservedWords {
		if word == kw {
			return true
		}
	}
	return false
}

func removeStringsFromCode(code string) string {
	// Elimina los literales de texto entre comillas dobles
	return regexp.MustCompile(`"[^"]*"`).ReplaceAllString(code, "")
}

func LexicalAnalysis(code string) ([]Token, []string) {
	var tokens []Token
	var errors []string

	masterPattern := ""
	for _, p := range patterns {
		masterPattern += "(" + p + ")|"
	}
	masterPattern = strings.TrimSuffix(masterPattern, "|")

	re := regexp.MustCompile(masterPattern)
	matches := re.FindAllString(code, -1)

	for _, m := range matches {
		tokenType := "unknown"
		for tType, pat := range patterns {
			if matched, _ := regexp.MatchString("^"+pat+"$", m); matched {
				tokenType = tType
				break
			}
		}
		tokens = append(tokens, Token{Type: tokenType, Value: m})
	}

	// Revisión de errores léxicos más robusta (ignorando palabras dentro de strings)
	codeNoStrings := removeStringsFromCode(code)
	codeWords := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`).FindAllString(codeNoStrings, -1)
	for _, word := range codeWords {
		isTokenized := false
		for _, tok := range tokens {
			if tok.Value == word {
				isTokenized = true
				break
			}
		}
		// Si no es una palabra clave válida, pero parece ser una... es un error
		if !isExpectedKeyword(word) && !isTokenized {
			// Solo marcar error si no es un número
			if !regexp.MustCompile(`^\d+$`).MatchString(word) {
				errors = append(errors, "Posible error léxico: '"+word+"' no es una palabra clave válida")
			}
		}
	}

	return tokens, errors
}

// Clasificación para la tabla
func ClassifyToken(t Token) string {
	switch t.Type {
	case "keyword":
		if isReservedWord(t.Value) {
			return "PR"
		}
		return "Otro"
	case "identifier":
		if isReservedWord(t.Value) {
			return "PR"
		}
		return "ID"
	case "number":
		return "Numeros"
	case "symbol":
		return "Simbolos"
	default:
		return "Error"
	}
}

// Nueva función para el resumen léxico
func LexicalSummary(code string) map[string]interface{} {
	tokens, errors := LexicalAnalysis(code)

	// Inicializar estructura para la tabla
	table := []map[string]interface{}{}
	counts := map[string]int{"PR": 0, "ID": 0, "Numeros": 0, "Simbolos": 0, "Error": 0}

	for _, t := range tokens {
		cat := ClassifyToken(t)
		if cat == "Otro" {
			continue
		}
		row := map[string]interface{}{
			"value":    t.Value,
			"PR":       cat == "PR",
			"ID":       cat == "ID",
			"Numeros":  cat == "Numeros",
			"Simbolos": cat == "Simbolos",
			"Error":    cat == "Error",
		}
		table = append(table, row)
		if cat != "Error" {
			counts[cat]++
		}
	}
	// Contar errores léxicos
	counts["Error"] += len(errors)

	return map[string]interface{}{
		"rows":           table,
		"totals":         counts,
		"lexical_errors": errors,
	}
}
