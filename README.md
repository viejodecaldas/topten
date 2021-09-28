#TOP 10 Processor

##Introduction

This repo process files and prints out the Top 10 repositories
sorted by amount of watch events, top 10 active users sorted
by amount of PRs created and commits pushed and Top 10 repositories
sorted by amount of commits pushed.

Given the files input it will get info from files through go routines
process the info and prints out the result.

##Instructions

You can run the app and pass it parameters indicating the name of
each file so the program can look for them in the right place.

To achieve this you can simply 

    go run main.go --actors=actors.csv --events=events.csv --repos=repos.csv --commits=commits.csv

You can omit the flags and some default values will be assumed to process the files
included in the project.

##Assumptions

Based on the information given in _events.csv_ there might not need to use _commits.csv_.
This is because _events.csv_ already have the event type inside, and you can know who
did an action from just that file and search for those actors that did an action