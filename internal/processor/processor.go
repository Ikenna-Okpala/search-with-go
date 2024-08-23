package processor

import (
	"log"
	"maps"
	"math"
	"os"
	"regexp"
	"slices"
	"strings"
	"webcrawler/internal/database"

	stemmer "github.com/agonopol/go-stem"
)

type DocRank struct {
	URL string `json:"url"`
	Score float64 `json:"score"`
}

func Process(content string) string {

	stopWords, err2:= os.ReadFile("../../internal/processor/stopwords.txt")

	if err2 != nil {
		log.Fatalln(err2)
	}

	stopWordsLineEnding := strings.ReplaceAll(string(stopWords), "\r\n", "\n")

	stopWordsSet := makeStopWordSet((strings.Split(stopWordsLineEnding, "\n")))

	content = removeExtraSpaces(content)
	content = toLowerCase(content)
	content = removePunctuation(content)
	content = removeStopWords(stopWordsSet, content)
	content = stemWords(content)

	return content
}

func makeStopWordSet(stopWords []string) map[string]bool {

	stopWordsSet := make(map[string]bool)

	for _, stopWord := range stopWords {

		stopWordsSet[stopWord] = true
	}

	return stopWordsSet
}

func removeExtraSpaces(webText string) string {

	return strings.Join(strings.Fields(webText), " ")
}

func toLowerCase(webText string) string {
	return strings.ToLower(webText)
}

func removePunctuation(webText string) string {

	regex:= regexp.MustCompile(`['!"#$%&\\'()\*+,\-\.\/:;<=>?@\[\\\]\^_{|}~']`)

	return regex.ReplaceAllString(webText, "")
}

func removeStopWords(stopWordsSet map[string]bool, webText string) string {

	words := strings.Split(webText, " ")

	var builder strings.Builder

	for _, word := range words {

		exists := stopWordsSet[word]

		if !exists {
			builder.WriteString(word)
			builder.WriteString(" ")
		}
	}

	return strings.Trim(builder.String(), " ")
}

func stemWords(webText string) string {

	words := strings.Split(webText, " ")

	var builder strings.Builder

	for _, word := range words {

		stemmedWord := stemmer.Stem([]byte(word))

		builder.WriteString(string(stemmedWord))
		builder.WriteString(" ")
	}

	return builder.String()
}

func ConstructMaps(keywords [] database.Keyword) (map[string]map[string] float64, map[string] float64){

	tfIdf := make(map[string] map[string] float64)
	idfMap := make(map[string] float64)

	for keyword:= range slices.Values(keywords){

		wordMap, ok := tfIdf[keyword.Url]

		if ok {
			wordMap[keyword.Word] = keyword.TfIdf
		}else{
			tfIdf[keyword.Url] = make(map[string] float64)
			tfIdf[keyword.Url][keyword.Word] = keyword.TfIdf
		}

		idfMap[keyword.Word] = keyword.Idf
	}

	return tfIdf, idfMap
}

func ConstructQueryMap(query [] string, idfMap map[string] float64) map[string] float64{

	queryMap:= make(map[string] float64)

	maxFreq := 0.0

	for _,word:= range query{

		queryMap[word] ++

		newFreq:= queryMap[word]

		if newFreq > maxFreq{
			maxFreq = newFreq
		}
	}
	
	for word, freq := range maps.All(queryMap){

		queryMap[word] = (freq / maxFreq) * idfMap[word]
	}

	return queryMap
}

func ComputeDocVectorLength(tfIdfMap map[string]map[string] float64) map[string] float64{

	docLengthMap := make(map[string] float64)

	for doc, wordMap := range tfIdfMap {

		for _, tfIdf:= range wordMap {
			 docLengthMap[doc]+= math.Pow(tfIdf, 2)
		}

		docLengthMap[doc] = math.Sqrt(docLengthMap[doc])
	}

	return docLengthMap
}

func ComputeQueryVectorLength(queryMap map[string] float64) float64 {

	length := 0.0

	for _, tfIdf := range queryMap {

		length += math.Pow(tfIdf, 2)
	}

	return math.Sqrt(length)
}

func ComputeCosineSimilarity(docTfIdfMap map[string] map[string] float64, docLengthMap map[string] float64, queryMap map[string] float64, queryLength float64) [] DocRank{

	docRank := make(map[string] float64)
	for doc, wordMap := range docTfIdfMap{

		docSum := 0.0

		for word, queryTfIdf := range queryMap{

			docSum += wordMap[word] * queryTfIdf
		}
		
		docSum /= docLengthMap[doc] * queryLength
		docRank[doc] = docSum
	}

	docRankList := make([]DocRank, 0)

	for url, score := range docRank {

		docRankList = append(docRankList, DocRank{URL: url, Score: score})
	}

	return docRankList
}