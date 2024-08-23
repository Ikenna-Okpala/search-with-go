package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"slices"
	"webcrawler/internal/database"
	"webcrawler/internal/processor"
	"webcrawler/internal/utils"
)

type Handler struct {
	DB *sql.DB
}

func (h Handler) Search (w http.ResponseWriter, r * http.Request){

	searchQuery:= r.URL.Query().Get("query")

	processor.Process(searchQuery)

	queryTokensSQL := utils.TokenizeStringSQL(searchQuery)

	queryForSQL := utils.FormatQuery(queryTokensSQL)

	rows, err:= h.DB.Query(queryForSQL)

	if err != nil {
		http.Error(w, "Server Failure", http.StatusInternalServerError)
	}

	tokens:= utils.TokenizeString(searchQuery)

	keywords:= make([] database.Keyword, 0)

	for rows.Next() {

		var keyword database.Keyword

		rows.Scan(&keyword.Word, &keyword.Url, &keyword.Idf, &keyword.TfIdf)

		keywords = append(keywords, keyword)
	}

	tfIdfMap, idfMap:= processor.ConstructMaps(keywords)

	queryMap := processor.ConstructQueryMap(tokens, idfMap)

	queryVectorLength:= processor.ComputeQueryVectorLength(queryMap)

	docVectorLengthMap:= processor.ComputeDocVectorLength(tfIdfMap)

	docRankList := processor.ComputeCosineSimilarity(tfIdfMap, docVectorLengthMap, queryMap, queryVectorLength)

	
	slices.SortFunc(docRankList, func(doc1 processor.DocRank, doc2 processor.DocRank) int {
		
		return int((doc2.Score - doc1.Score) * 10000)
	})

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(docRankList)
}