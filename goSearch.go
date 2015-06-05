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

const consumerKey = "47m4XBT9qogkUr1wyJv5sNiOi"
const consumerSecret = "gz7c2zNkBPanG2AdR8MvlLqoi16AveGsSneOe05N9DkBiwonnY"
const accessToken = "532932305-82LoqwU604eVUb8RkMIIWN5lHGLJMl3czqKJ8KMf"
const accessSecret = "qf0NmAK9f6otfYHBneYKwe6dPQOz8DTn1RWlQvzeE3zXr"
const TWEETS_AMOUNT = "3"

type Node struct {
	post string
	token map[string]int
	vector []float64
	clusterID int
}

func main() {
	post := searchTweets("Obama")
	tokenizedList, bagOfWords := tokenize(post)
	nodes := make([]Node, len(post))
	for i, post := range post {
		nodes[i].post = post
		nodes[i].token = tokenizedList[i]
	}
	
	globalIndex := indexingTerms(bagOfWords)
	data := buildVector(tokenizedList, globalIndex)
	clusters, _ := kmeans.Kmeans(data, 5, kmeans.EuclideanDistance, 100);

	printResult(post, clusters)

	for _, node := range nodes {
		fmt.Println(node.post)
		fmt.Println(node.token)
	}
}

func printResult(posts []string, clusters[]int) {
	// res := make(map[int][]string)
	// for index, groupID := range clusters {
	// 	res[groupID] = append(res[groupID], posts[index])
	// }

	// for id, texts :=range res {
	// 	fmt.Println("cluster", id)
	// 	for _ , text :=range texts {
	// 		fmt.Println(text)
	// 	}
	// }
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

func indexingTerms(bagOfWords []string) map[string]int {
	indexList := make(map[string]int)
	for index, term := range bagOfWords {
		indexList[term] = index
	}
	return indexList
}

func buildVector(tokenizedList []map[string]int, dictionary map[string]int) [][]float64 {
	//a list of ND-Vector
	var data [][]float64
	for _, doc := range tokenizedList {
		//set up ND-Vector n is the length of dictionary
		vector := make([]float64, len(dictionary))
		for term, _ := range doc {	
			//calculate tf (term, doc)
			tf := getTF(doc, term)
			//calculate idf (corpus, term)
			idf := getIDF(tokenizedList, term)
			//calculate score of tf*idf
			score := float64(tf * idf)
			
			//loop up the index of term from dictionary
			index := dictionary[term]
			vector[index] = score
		}
		//add to ND-Vector list
		data = append(data, vector)
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
func tokenize(posts []string) ([]map[string]int, []string) {
	absPath, _ := filepath.Abs("Documents/GO/stop_words.txt")
	bwtokenizer := tokenizer.NewBagOfWordsTokenizer(absPath)
	var tokenizedList []map[string]int
	var rawCorpus []string
	for _, text := range posts {
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