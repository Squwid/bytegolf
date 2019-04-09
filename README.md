# Byte Golf

Byte Golf is a code solving problem that rewards players for solving code problems in the shortest amount of code! There is an input for each question and the output must be correct. Byte Golf is open source and is available [HERE!](https://bytegolf.io)


## Rules

The rules for Byte Golf are suprisingly simple, just write the least code using the standard library of a language of your choice.

1. The player with the lowest total characters at the end wins, just like Golf!
1. Total score consists of characters or "bytes" in a working code solution (whitespace and comments not included). The scores are totaled at the end, and all holes must be completed to a working solution to be able to win. Hence the name "Byte Golf"
1. All code must be written in one of the supported languages and must only use standard libraries.

## Design

### Origin

My roommates and I came up with the idea of making Byte Golf a web app for us to see who can solve code in the shortest solution. I know this already exists (sort of) on Stack Overflow. I decided to code a web application in Go and bootstrap that allows users to submit code and get permanently stored.

### How it works

I use the Jdoodle compiling API to submit the code, and currently store all of the responses and questions locally, until I decide to put them in a database. The server is run on AWS currently but soon will be switched over to GCP. Questions are able to be created if you have Admin status, and can be archived for the future Leaderboards page.