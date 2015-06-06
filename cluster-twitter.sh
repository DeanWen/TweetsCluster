#!/bin/bash
# How to run:
# cluster-twitter.sh -q 'amazon' > amazon-cluster.txt

go get github.com/bugra/kmeans
go get github.com/srom/tokenizer
go get github.com/ChimeraCoder/anaconda

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
