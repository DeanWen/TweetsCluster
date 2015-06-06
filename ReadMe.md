ReadMe

This program to cluster the search result of tweets   

Required Files - put into same directory   

	1. goSearch.go (main program)
	2. stop_words.txt (using for NLP)
	3. Steve Jobs-cluster.txt (Sample Output)

Main Idea  

	1. NLP for English to tokenize each tweet, remove stop words, stemming the tweets

	2. Set up a dictionary for all unique words, as a global index

	3. Calculate the TF-IDF, use this score to build a N Dimension Vector for each tweet

	4. Using K-Means Algorithm to cluster all tweets based on vector similiary (Using Euclidean Distance)

	5. After clustering, calculate the center position and Euclidean Distance

	6. Find the BEST result based on closeat to the center got in step 5

	7. Find the FIRST result based on created timestamp

	8. Sort the rest of results by newest coming first principle

	9. Generate a txt file to store all results (name formate: query-cluster.txt)


How to Run

	Put folder under the $GOPATH directory 
	
		go get github.com/bugra/kmeans
		go get github.com/srom/tokenizer
		go get github.com/ChimeraCoder/anaconda 
		go run goSearch.go 'query'  
 	
 	Output will be in the same directory as well  

 	Please refer Sample output "Steve-Jobs-Cluster.txt"  
 		All Tweets are clustered in 10 groups (0-9)  
 		Group 8 only get BEST & FIRST  
 		Group 9 is good sample, BEST & FIRST & Rest sort by timestamp  

Please refer the main program goSearch.go for more details. It is highly architected and object oriented with fully commentation.  

If you have any questions, feel free to address me at  
Dian (Dean) Wen  
dwen@cmu.edu  
github.com/DeanWen  
  
Carnegie Mellon University  
Jun-5-2015   
