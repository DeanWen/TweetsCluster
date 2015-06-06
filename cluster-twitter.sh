#!/bin/bash 
# How to run:
# cluster-twitter.sh -q 'amazon' > amazon-cluster.txt

#install all required packages
#Package for K-Means Algorithm 
go get github.com/bugra/kmeans
#Package for NLP Tokenizer
go get github.com/srom/tokenizer
#Package for Twitter APIs
go get github.com/ChimeraCoder/anaconda

#Parse Command Query Parameter 
qflag= 
while getopts 'q:' OPTION
do
    case $OPTION in 
        q) text="$OPTARG"
           qflag=1
           ;;
    esac
done
if [ "$qflag" ]
then
go run goSearch.go "$text"
else
    printf "Usage: ./cluster-twitter.sh -q query\n"
fi
