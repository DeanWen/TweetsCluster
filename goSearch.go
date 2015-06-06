package main

import (
	"io"
	"os"
	"fmt"
	"math"
	"time"
	"net/url"
	"reflect"
	"path/filepath"
	"github.com/bugra/kmeans"
	"github.com/srom/tokenizer"
	"github.com/ChimeraCoder/anaconda"
)

const consumerKey = "47m4XBT9qogkUr1wyJv5sNiOi"
const consumerSecret = "gz7c2zNkBPanG2AdR8MvlLqoi16AveGsSneOe05N9DkBiwonnY"
const accessToken = "532932305-82LoqwU604eVUb8RkMIIWN5lHGLJMl3czqKJ8KMf"
const accessSecret = "qf0NmAK9f6otfYHBneYKwe6dPQOz8DTn1RWlQvzeE3zXr"
const TWEETS_AMOUNT = "100"
const K = 10
const threshold = 100

type Node struct {
	time time.Time
	post string
	token map[string]int
	vector []float64
	clusterID int
	distance float64
}

func main() {
	/*
	@posts : the list of all tweets.Text
	@nodes : the list of all nodes
	*/
	posts, nodes := searchTweets("Obama")
	
	/*
	@tokenizedList : the list of tokenized posts (tweet.text)
	@bagOfWords : the list of all unique words
	*/
	tokenizedList, bagOfWords := tokenize(posts)
	
	/*
	@dictionary : a global index
	*/
	dictionary := indexingTerms(bagOfWords)
	
	/*
	@data : the list of all vectors
	*/
	data := buildVector(tokenizedList, dictionary)
	
	/*
	@clusters : the list of cluster IDs
	*/
	clusters, _ := kmeans.Kmeans(data, K, kmeans.EuclideanDistance, threshold);

	/*
	complete attributes to each node
	*/
	for i, _ := range nodes {
		nodes[i].token = tokenizedList[i]
		nodes[i].vector = data[i]
		nodes[i].clusterID = clusters[i]
	}

	/*
	put node that is in same cluster to the same group
	map <key = groupID, value = node>
	*/
	grouped := make(map[int][]Node)
	for _, node := range nodes {
		grouped[node.clusterID] = append (grouped[node.clusterID], node)
	}

	finalClusters := setDistance(grouped)

	printRes(finalClusters)
	writeFile()
}

func printRes(finalClusters map[int][]Node) {
	for cid, clusters := range finalClusters {
		fmt.Println("Cluster :", cid)

		best := clusters[0]
		first := clusters[0]
		for _, node := range clusters {
			if node.distance < best.distance {
				best = node
			}
			if node.time.Before(first.time) {
				first = node
			}
		}
		
		fmt.Println("Best Result")
		fmt.Println(best.post)
		fmt.Println("First Result")
		fmt.Println(first.post)

		fmt.Println("Other Result - newest to oldest")
		for _, node := range clusters {
			if !reflect.DeepEqual(node, first) && 
			   !reflect.DeepEqual (node, best) {
				fmt.Println(node.post)
			}
		}

		fmt.Println()
	}
}

func writeFile() {
	filename := "Documents/GO/searchTweets/test.txt"

	fmt.Println("writing: " + filename)
	
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	
	n, err := io.WriteString(f, "blahblahblah")
	if err != nil {
		fmt.Println(n, err)
	}
	
	f.Close()
}

func setDistance(groups map[int][]Node) map[int][]Node {
	res := make(map[int][]Node)

	for gid, eachGroup := range groups {
		//fmt.Println(gid)
		
		sum := make([]float64, len(eachGroup[0].vector))
		for _, node := range eachGroup {
			//fmt.Println(node.vector)
			sum = add(node.vector, sum)
		}
		
		center := getCenter(sum, len(eachGroup))

		for _, node := range eachGroup {
			dis, _ := kmeans.EuclideanDistance(node.vector, center)
			node.distance = dis
			res[gid] = append(res[gid], node)
		}
	}

	return res
}

func add(vector1, vector2 []float64) []float64 {
	res := make([]float64, len(vector1))
	for ii, _ := range vector1 {
		res[ii] = vector1[ii] + vector2[ii]
	}
	return res
}

func div(vector []float64, op float64) []float64 {
	res := make([]float64, len(vector))
	for ii, jj := range vector {
		res[ii] = jj / op
	}
	return res
}

func getCenter(sum []float64, num int) []float64{
	center := div(sum, float64(num))
	return center
}

func searchTweets(query string) ([]string, []Node){
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
	var nodes []Node
	for _, tweet := range search_result.Statuses {
		var node Node
		node.post = tweet.Text
		node.time, _ = time.Parse(time.RubyDate,tweet.CreatedAt)
		
		result = append(result, tweet.Text)
		nodes = append(nodes, node)
	}

	return result, nodes
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
			//calculate tf (tdoc, doc)
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