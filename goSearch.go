package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/srom/tokenizer"
	"net/url"
	"path/filepath"
	"math"
)

var consumerKey = "47m4XBT9qogkUr1wyJv5sNiOi"
var consumerSecret = "gz7c2zNkBPanG2AdR8MvlLqoi16AveGsSneOe05N9DkBiwonnY"
var accessToken = "532932305-82LoqwU604eVUb8RkMIIWN5lHGLJMl3czqKJ8KMf"
var accessSecret = "qf0NmAK9f6otfYHBneYKwe6dPQOz8DTn1RWlQvzeE3zXr"
var TWEETS_AMOUNT = "100"

func main() {
	list := searchTweets("Obama")
	// setUpDict(list);
	tokenizedList := Tokenize(list)
	fmt.Println(tokenizedList)
}

func searchTweets(query string)([]string){
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

func getTF(vector map[string]int, term string)float64 {
	var n int
	_, present := vector[term]
    if present {
        n = vector[term]  
    } else {
    	n = 0
    }
    return float64(n / len(vector))
}

func getIDF(documents []map[string]int, term string)float64 {
	var n = 0
	for _, doc := range documents {
		_, present := doc[term]
    	if present {
        	n++ 
    	}
	}
	return math.Log(float64(len(documents) / n))
}

func Tokenize(list []string)[]map[string]int {
	absPath, _ := filepath.Abs("Documents/GO/stop_words.txt")
	bwtokenizer := tokenizer.NewBagOfWordsTokenizer(absPath)
	var tokenizedList []map[string]int
	for _, text := range list {
		tokens := bwtokenizer.Tokenize(text)
		dict := WordCount(tokens)
		tokenizedList = append(tokenizedList, dict)
	}
	return tokenizedList
}

func WordCount(s []string) map[string]int {
    dict := make(map[string]int)

    splited := s
    for _, string := range splited {
        _, present := dict[string]
        if present {
            dict[string]++  
        } else {
            dict[string] = 1
        }
    }
    return dict
}