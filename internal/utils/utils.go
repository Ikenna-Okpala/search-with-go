package utils

import (
	"fmt"
	"strings"
	"webcrawler/internal/database"
)

func TokenizeStringSQL(text string) []string {

	wordQuery := make([]string, 0)

	splittedQuery := strings.Split(text, " ")

	for _, word := range splittedQuery {

		wordQuery = append(wordQuery, fmt.Sprintf("'%s'", word))
	}

	return wordQuery
}

func FormatQuery(queryTokens [] string ) string{

	arguments:= strings.Join(queryTokens, ",")

	return fmt.Sprintf(database.SearchWord, arguments)
}

func TokenizeString(text string) [] string {

	return strings.Split(text, " ")
}