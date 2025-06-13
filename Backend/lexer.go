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
	"keyword":    `\b(import|from|function|return|const|let|var|if|else|useState|useEffect|export|default|React)\b`,
	"identifier": `\b[a-zA-Z_][a-zA-Z0-9_]*\b`,
	"number":     `\b\d+(\.\d+)?\b`,
	"string":     `"[^"\n]*"|'[^'\n]*'`,
	"symbol":     `[{}\[\]();=.,<>:+\-*/]`,
}

// Lista de palabras clave esperadas (para detectar errores como "imprt")
var expectedKeywords = []string{
	"import", "from", "function", "return", "const", "let", "var", "if", "else",
	"useState", "useEffect", "export", "default", "React",
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

	// Revisión de errores léxicos más robusta
	codeWords := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`).FindAllString(code, -1)
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
			errors = append(errors, "Posible error léxico: '"+word+"' no es una palabra clave válida")
		}
	}

	return tokens, errors
}
