package tfidf

import (
	"math"
	"webcrawler/internal/webcrawler"
)

//TF
// count number of times word appear
// count number of words in doc
// document: list of terms

func TF(websites map[string] webcrawler.Website) map[string]map[string] float64 {

	tf:= make(map[string]map[string] float64)
	for url, content := range websites {

		tf[url] = make(map[string] float64)
		
		for _, word:= range content.Words{

			tf[url][word] ++
		}

		n_words:= len(websites[url].Words)

		for word := range tf[url] {

			// fmt.Println("-------TESTING-------")

			// fmt.Printf("n_words: %v\n", n_words)

			// fmt.Printf("count_word: %s is %v\n", word, tf[url][word])

			tf[url][word] /= float64(n_words)

			
		}

	}

	return tf
}

func IDF(websites map[string] webcrawler.Website) map[string] float64 {

	websiteWordSet := make(map[string]map[string] bool)

	idf := make(map[string] float64)

	nWebsites := len(websites)

	for website, content := range websites {
		websiteWordSet[website] = make(map[string] bool)

		for _,word := range content.Words{

			websiteWordSet[website][word] = true
			idf[word] = 0
		}
	}

	for idfWord := range idf{

		for _, websiteWords := range websiteWordSet{

			_, ok := websiteWords[idfWord]

			if ok {
				idf[idfWord]++
			}
		}
	}

	for idfWord, freq := range idf {

		// fmt.Printf("[%v, %v]", idfWord, freq)

		idf[idfWord] = math.Log2(float64(1 + nWebsites) / (1 + freq)) + 1
	}
	
	return idf
}


func TFIDF(tf map[string]map[string]float64, idf map[string]float64) map[string]map[string]float64{

	tfidf := make(map[string]map[string] float64)

	for website, wordsMap := range tf{

		tfidfWebsite := make(map[string]float64)

		for word, tfWord := range wordsMap {

			tfidfWebsite[word] = tfWord * idf[word]
		}

		tfidf[website] = tfidfWebsite
	}

	return tfidf
}