/*
 * Tweets Automatic Clustering Program
 *
 * Author: Dian (Dean) Wen
 * Github: github.com/DeanWen
 * Carnegie Mellon University
 * Jun-5-2015 
 * 
 * All Rights Reserved
 */
package main

import (
	"io"
	"os"
	"fmt"
	"math"
	"time"
	"strconv"
	"net/url"
	"reflect"
	"path/filepath"
	"github.com/bugra/kmeans"//K-Means Library
	"github.com/srom/tokenizer"//NLP Tokenizer for English
	"github.com/ChimeraCoder/anaconda"//Twitter Library for Golang
)

/*
 *Twitter APIs
 *@consumerKey,@consumerSecret
 *@accessToken,@accessSecret
 *
 *K-Means Constants
 *@K cluster by k groups
 *@threshold the maximum clustering times
 */
const consumerKey = "47m4XBT9qogkUr1wyJv5sNiOi"
const consumerSecret = "gz7c2zNkBPanG2AdR8MvlLqoi16AveGsSneOe05N9DkBiwonnY"
const accessToken = "532932305-82LoqwU604eVUb8RkMIIWN5lHGLJMl3czqKJ8KMf"
const accessSecret = "qf0NmAK9f6otfYHBneYKwe6dPQOz8DTn1RWlQvzeE3zXr"
const TWEETS_AMOUNT = "100"
const K = 10
const threshold = 100


/*
 * Customer Data Type
 * @time tweets create time
 * @post tweets text
 * @token tokenizes and stemmed tweets text
 * @vector tf-idf vector
 * @clusterID K-Means clustered ID
 * @distance the distance to the cluster center
 */
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
	 * Read Query Paramter from comander line
	 */
	args := os.Args
	query := string(args[1])
	/*
	 *@posts : the list of all tweets.Text
	 *@nodes : the list of all nodes
	 */
	posts, nodes := searchTweets(query)
	
	/*
	 *@tokenizedList : the list of tokenized posts (tweet.text)
	 *@bagOfWords : the list of all unique words
	 */
	tokenizedList, bagOfWords := tokenize(posts)
	
	/*
	 *@dictionary : a global index
	 */
	dictionary := indexingTerms(bagOfWords)
	
	/*
	 *@data : the list of all vectors
	 */
	data := buildVector(tokenizedList, dictionary)
	
	/*
	 *@clusters : the list of cluster IDs
	 */
	clusters, _ := kmeans.Kmeans(data, K, kmeans.EuclideanDistance, threshold);

	/*
	 *complete attributes to each node
	 */
	for i, _ := range nodes {
		nodes[i].token = tokenizedList[i]
		nodes[i].vector = data[i]
		nodes[i].clusterID = clusters[i]
	}

	/*
	 *put node that is in same cluster to the same group
	 *map <key = groupID, value = node>
	 */
	grouped := make(map[int][]Node)
	for _, node := range nodes {
		grouped[node.clusterID] = append (grouped[node.clusterID], node)
	}

	/*
	 *assign Euclidean distance to each node vector from group center
	 */
	finalClusters := setDistance(grouped)

	/*
	 *print in terminal and write to file
	 */
	filename := query + "-cluster.txt"
	printRes(finalClusters, filename)
}

/*
 *Tokenize Method in English
 * 1. to remove all stopwords 
 * 2. stemming the terms
 * 3. atomize the unique terms
 *@posts all tweets text
 * Return
 * @HashMap<term, appering_times> tokenized token for each tweets text
 * @bagOfUniqueWords[]string a global unique bag of words
 */
func tokenize(posts []string) ([]map[string]int, []string) {
	pwd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    /*
     * get the stop-words list from txt file
     * should in the same directory
     */
    txtPath := string(pwd) + "/stop_words.txt"
	absPath, _ := filepath.Abs(txtPath)
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

/*
 * Calculate term frequency(tf)
 * @doc
 * @term
 */
func getTF(doc map[string]int, term string) float64 {
	var tf float64 = 0.0
	_, present := doc[term]
    if present {
        tf = float64(doc[term])  
    }
    return float64(tf / float64(len(doc)))
}

/*
 * Calculate inverse document frequency(idf)
 * @documents the corpus 
 * @term
 */
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

/*
 * build Vector method
 * @tokenizedList all tonkenized post collections
 * @dictionary the global index to look up terms
 * Return 
 * 		@[][]float64 a N-Dimension Vector Collections
 */
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

/*
 * indexing method to create the dictionary
 * @bagOfWords all unique term collections
 * Return 
 *	@indexList map[string]int a dictionary to look up all terms
 */
func indexingTerms(bagOfWords []string) map[string]int {
	indexList := make(map[string]int)
	
	for index, term := range bagOfWords {
		indexList[term] = index
	}
	
	return indexList
}

/*
 * calculate the Euclidean distance of each node to group center 
 */
func setDistance(groups map[int][]Node) map[int][]Node {
	res := make(map[int][]Node)
	for gid, eachGroup := range groups {		
		sum := make([]float64, len(eachGroup[0].vector))
		for _, node := range eachGroup {
			sum = add(node.vector, sum)
		}

		/*
		 *center vector
		 */
		center := getCenter(sum, len(eachGroup))

		/*
		 * Using Euclidean Distance
		 */
		for _, node := range eachGroup {
			dis, _ := kmeans.EuclideanDistance(node.vector, center)
			node.distance = dis
			res[gid] = append(res[gid], node)
		}
	}
	return res
}

/*
 * two vector sum method
 * @vector1
 * @vector2
 */
func add(vector1, vector2 []float64) []float64 {
	res := make([]float64, len(vector1))
	for ii, _ := range vector1 {
		res[ii] = vector1[ii] + vector2[ii]
	}
	return res
}

/*
 * vector divide method
 * @vector
 * @divider
 */
func div(vector []float64, op float64) []float64 {
	res := make([]float64, len(vector))
	for ii, jj := range vector {
		res[ii] = jj / op
	}
	return res
}

/*
 * calculate the cluster center vector
 * formula: sum[0...0] / total # of nodes
 */
func getCenter(sum []float64, num int) []float64{
	center := div(sum, float64(num))
	return center
}

/*
 * search tweets
 * @query the content want to retrieval
 * return
 *	 @[]string all tweets text list
 *   @[]Node all nodes with tweets text and timestamp
 */
func searchTweets(query string) ([]string, []Node){
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)
	//set up the optional parameter[@count amount of tweets]
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

/*
 * Word Count method
 * @s post text
 * Return
 * map<key = term, value = appearing times>
 */
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

/*
 * Bubble Sort for sort by time
 * Since amount is limited, BB sort is quite effetive as well
 */
func bubbleSort(arrayzor []Node) {
	swapped := true;
	for swapped {
		swapped = false
		for i := 0; i < len(arrayzor) - 1; i++ {
			if arrayzor[i + 1].time.After(arrayzor[i].time) {
				swap(arrayzor, i, i + 1)
				swapped = true
			}
		}
	}	
}

/*
 * swap method for bubble sort
 */
func swap(arrayzor []Node, i, j int) {
	tmp := arrayzor[j]
	arrayzor[j] = arrayzor[i]
	arrayzor[i] = tmp
}


func printRes(finalClusters map[int][]Node, filename string) {
	/*
	 * Get current directory path
	 */
	pwd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    /*
     * create result file in same dir
     * name format query-cluster
     */
    newFile := string(pwd) + "/" + filename
	f, err := os.Create(newFile)
	if err != nil {
		fmt.Println(err)
	}

	for cid, clusters := range finalClusters {
		fmt.Println("Cluster :", cid)
		n, err := io.WriteString(f, "Cluster : " + strconv.Itoa(cid) + "\n")
		
		/*
		 * Find the closest node to be BEST
		 * Find the earilest node to be FIRST
		 */
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
		
		/*
		 * Print in Terminal
		 */
		fmt.Println("Best Result")
		fmt.Println(best.time.String(), best.post)
		fmt.Println("First Result")
		fmt.Println(first.time.String(), first.post)

		/*
		 * Write to the File
		 */
		n, err = io.WriteString(f, "Best Result \n")
		n, err = io.WriteString(f, best.time.String() + best.post + "\n")
		n, err = io.WriteString(f, "First Result \n")
		n, err = io.WriteString(f, first.time.String() + first.post + "\n")

		/*
		 *check error
		 */
		if err != nil {
    		fmt.Println(n, err)
  		}

		fmt.Println("Other Result - newest to oldest")
		n, err = io.WriteString(f, "Other Result - newest to oldest \n")

		/*
		 *split rest post
		 */
		var rest []Node
		for _, node := range clusters {
			if !reflect.DeepEqual(node, first) && 
			   !reflect.DeepEqual (node, best) {
			   	rest = append(rest, node)
			}
		}

		/*
		 * Sort by tweets.time
		 * make sure newest comes first
		 */
		bubbleSort(rest)

		for _, node := range rest {
			fmt.Println(node.time, node.post)
			n, err = io.WriteString(f, node.time.String() + node.post + "\n")
		}

		n, err = io.WriteString(f, "\n")
		fmt.Println()
	}

	f.Close()
}
