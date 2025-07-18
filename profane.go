package main

import "strings"

func profaneCleaner(body string) string {

	profane := []string{"kerfuffle", "sharbert", "fornax"}

	bodySplit := strings.Split(body, " ")

	for i := range bodySplit {
		for j := range profane {
			if strings.ToLower(bodySplit[i]) == profane[j] {
				bodySplit[i] = strings.Repeat("*", 4)
			}
		}
	}

	return strings.Join(bodySplit, " ")
}
