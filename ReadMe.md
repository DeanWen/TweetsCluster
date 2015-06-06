ReadMe

This Program to cluster the search result of tweets 

Required Files - put into same directory
	1. goSearch.go (main program)
	2.stop_words.txt (using for NLP)

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


Please refer the main program goSearch.go for more details. It is highly architected and object oriented with fully commentation.

If you have any questions, feel free to address me at

Dian (Dean) Wen
dwen@cmu.edu
github.com/DeanWen

Carnegie Mellon University
Jun-5-2015 