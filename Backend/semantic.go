package main

import (
	"regexp"
	"strings"
)

type ResultadoSemantico struct {
	Errores      []string `json:"errores"`
	Advertencias []string `json:"advertencias"`
}

func AnalisisSemantico(codigo string) ResultadoSemantico {
	lineas := strings.Split(codigo, "\n")
	variables := map[string]bool{}
	errores := []string{}
	advertencias := []string{}
	usadas := map[string]bool{}

	for i, linea := range lineas {
		l := strings.TrimSpace(linea)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}

		// Asignación: VAR=valor
		asig := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)=(.*)$`)
		if asig.MatchString(l) {
			m := asig.FindStringSubmatch(l)
			varName := m[1]
			if _, existe := variables[varName]; existe {
				advertencias = append(advertencias, "Variable '"+varName+"' reasignada (línea "+itoa(i+1)+")")
			}
			variables[varName] = true
			continue
		}

		// Uso de variable: $VAR o ${VAR}
		usoVar := regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)|\$\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)
		usos := usoVar.FindAllStringSubmatch(l, -1)
		for _, u := range usos {
			varName := u[1]
			if varName == "" {
				varName = u[2]
			}
			usadas[varName] = true
			if _, existe := variables[varName]; !existe {
				errores = append(errores, "Variable '"+varName+"' usada sin haber sido asignada (línea "+itoa(i+1)+")")
			}
		}

		// Comando mal formado (ejemplo simple: línea vacía o solo símbolo)
		cmd := regexp.MustCompile(`^[a-zA-Z0-9_\-\.\$\{\}\[\]=/\":' ]+$`)
		if !asig.MatchString(l) && !cmd.MatchString(l) && !strings.HasPrefix(l, "if") && !strings.HasPrefix(l, "for") && !strings.HasPrefix(l, "while") && l != "then" && l != "fi" && l != "do" && l != "done" && l != "{" && l != "}" {
			errores = append(errores, "Línea "+itoa(i+1)+": Comando o sintaxis no reconocida: '"+l+"'")
		}
	}

	// Variables declaradas y no usadas
	for v := range variables {
		if !usadas[v] {
			advertencias = append(advertencias, "Variable '"+v+"' asignada pero nunca usada")
		}
	}

	return ResultadoSemantico{
		Errores:      errores,
		Advertencias: advertencias,
	}
}
