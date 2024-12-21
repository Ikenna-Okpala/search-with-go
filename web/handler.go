package web

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"slices"
	"webcrawler/internal/database"
	"webcrawler/internal/processor"
	"webcrawler/internal/utils"
	"webcrawler/internal/webcrawler"
)

type Handler struct {
	DB *sql.DB
}


func (h Handler) Search (w http.ResponseWriter, r * http.Request){

	searchQuery:= r.URL.Query().Get("query")

	searchQuery = processor.Process(searchQuery)

	queryTokensSQL := utils.TokenizeStringSQL(searchQuery)

	queryForSQL := utils.FormatQuery(queryTokensSQL)

	rows, err:= h.DB.Query(queryForSQL)

	if err != nil {
		http.Error(w, "Server Failure", http.StatusInternalServerError)
	}

	tokens:= utils.TokenizeString(searchQuery)

	keywords:= make([] database.Keyword, 0)
	websites:= make([] webcrawler.Website, 0)

	for rows.Next() {

		var keyword database.Keyword
		var website webcrawler.Website

		rows.Scan(&keyword.Word,
			 &keyword.Url, &website.Title, &website.Icon,
			  &website.Name, &website.Description,
			   &keyword.Idf, &keyword.TfIdf)

		keywords = append(keywords, keyword)

		website.URL = keyword.Url
		websites = append(websites, website)
	}

	// fmt.Println("websites: ", websites)

	tfIdfMap, idfMap:= processor.ConstructMaps(keywords)

	// fmt.Println("idfMap", idfMap)

	queryMap := processor.ConstructQueryMap(tokens, idfMap)

	queryVectorLength:= processor.ComputeQueryVectorLength(queryMap)

	docVectorLengthMap:= processor.ComputeDocVectorLength(tfIdfMap)
	

	docRankList := processor.ComputeCosineSimilarity(tfIdfMap, docVectorLengthMap, queryMap, queryVectorLength)

	// fmt.Println("queryMap: ", queryMap)

	// fmt.Println("idfMap", idfMap)

	// fmt.Println("tfidfMap: ", tfIdfMap)

	// fmt.Println("idfMap: ", idfMap)

	// fmt.Println("queryVec: ", queryVectorLength)

	// fmt.Println("tfidfVecLength: ", docVectorLengthMap)

	// fmt.Println("docRankList: ", docRankList)

	
	slices.SortStableFunc(docRankList, func(doc1 processor.DocRank, doc2 processor.DocRank) int {

		if doc2.Score > doc1.Score {
			return 1
		}else if doc2.Score == doc1.Score {
			return 0
		}else{
			return -1
		}
	})

	
	for i, doc := range docRankList{

		index:= slices.IndexFunc(websites, func(website webcrawler.Website) bool {
			return doc.URL == website.URL
		})

		// fmt.Println("index: ", index)

		if index < 0 {
			continue
		}

		website:= websites[index]

		// fmt.Println("website: ", website)

		docRankList[i].Title = website.Title
		docRankList[i].Name = website.Name
		docRankList[i].Description = website.Description
		docRankList[i].Icon = website.Icon
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(docRankList)
}