package main

import (
	"fmt"
	"math"
	"net/url"
	"path/filepath"
	"github.com/bugra/kmeans"
	"github.com/srom/tokenizer"
	"github.com/ChimeraCoder/anaconda"
)

var consumerKey = "47m4XBT9qogkUr1wyJv5sNiOi"
var consumerSecret = "gz7c2zNkBPanG2AdR8MvlLqoi16AveGsSneOe05N9DkBiwonnY"
var accessToken = "532932305-82LoqwU604eVUb8RkMIIWN5lHGLJMl3czqKJ8KMf"
var accessSecret = "qf0NmAK9f6otfYHBneYKwe6dPQOz8DTn1RWlQvzeE3zXr"
var TWEETS_AMOUNT = "15"

func main() {
	list := searchTweets("Obama")
	tokenizedList, bagOfWords := tokenize(list)
	globalIndex := indexingTerms(bagOfWords)
	data := buildVector(tokenizedList, globalIndex)
	// for _, item := range data {
	// 	fmt.Println(item)
	// }

	clusters, _ := kmeans.Kmeans(data, 5, kmeans.EuclideanDistance, 100);
	fmt.Println(clusters)
}



func searchTweets(query string) []string {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)
	
	v := url.Values{}
	v.Set("count", TWEETS_AMOUNT)
	search_result, err := api.GetSearch(query, v)
	
	if err != nil {
		panic(err)
	}
	var result []string
	for _, tweet := range search_result.Statuses {
		result = append(result, tweet.Text)
	}

	return result
}

func indexingTerms(bagOfWords []string)map[string]int {
	indexList := make(map[string]int)
	for index, term := range bagOfWords {
		indexList[term] = index
	}
	return indexList
}

//Questions
func buildVector(tokenizedList []map[string]int, globalIndex map[string]int) [][]float64 {
	var data [][]float64
	for _, doc := range tokenizedList {
		post := make([]float64, len(globalIndex))
		for term, _ := range doc {	
			//calculate tf*idf for (q, d)
			tf := getTF(doc, term)
			idf := getIDF(tokenizedList, term)
			score := float64(tf * idf)
			
			//initialize the n-d array
			index := globalIndex[term]
			post[index] = score
		}
		data = append(data, post)
	}
	return data
}
 
func getTF(doc map[string]int, term string) float64 {
	var tf float64 = 0.0
	_, present := doc[term]
    if present {
        tf = float64(doc[term])  
    }
    return float64(tf / float64(len(doc)))
}

func getIDF(documents []map[string]int, term string) float64 {
	var df float64 = 0.0
	for _, doc := range documents {
		_, present := doc[term]
    	if present {
        	df++ 
    	}
	}
	var totalDocs = float64(len(documents))
	return math.Log(totalDocs / df)
}

//tokenize every post as a array of HashMap<term, appering_times>
//tokenize all terms return a global unique bag of words
func tokenize(list []string) ([]map[string]int, []string) {
	absPath, _ := filepath.Abs("Documents/GO/stop_words.txt")
	bwtokenizer := tokenizer.NewBagOfWordsTokenizer(absPath)
	var tokenizedList []map[string]int
	var rawCorpus []string
	for _, text := range list {
		tokens := bwtokenizer.Tokenize(text)
		rawCorpus = append(rawCorpus, tokens...)
		dict := wordCount(tokens)
		tokenizedList = append(tokenizedList, dict)
	}
	
	//find tokenized global unique bag of words corpus
	uniqueCorpusDict := wordCount(rawCorpus)
	var bagOfUniqueWords []string
	for term, _ := range uniqueCorpusDict {
		bagOfUniqueWords = append(bagOfUniqueWords, term)
	}

	return tokenizedList, bagOfUniqueWords
}

func wordCount(s []string) map[string]int {
    dict := make(map[string]int)
    for _, string := range s {
        _, present := dict[string]
        if present {
            dict[string]++  
        } else {
            dict[string] = 1
        }
    }
    return dict
}