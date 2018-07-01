# Byte Golf

Byte Golf is a code solving problem game that awards players with the shortest code that solves the task at hand, just like golf! I have had this idea for a while but I couldn't find a good time to start. So I thought today (6.30.18) is as good as time as any  This is going to be my first commit which is going to include a quick drawing I did for my own reference, pardon my drawing skills. I will try to update this project as often as I can but my 21st birthday and work occupy most of my time.

## Using Locally

To use this app locally, just clone it to a folder, compile it, and Go! It runs on port `:6017` because I thought it looked like "golf"

## Design Idea

### Origin

My roommates and I came up with the idea of making Byte Golf a web app for us to see who can solve code in the shortest solution. I know this already exists (sort of) on Stack Overflow. However there is no solution and any API that compiles and verifies the code.

### Thoughts

My current idea is to create a simple web app that players can clone and compile locally, and host a LAN server. The LAN server the user runs on their local machine will use an AWS Lambda that I will setup to deal with a global leaderboard system as well as the code compiler using DynamoDB and S3. Different langauges can give players more points because they aren't the best programming language Bash.

![Design](https://i.imgur.com/SVmaRN6.jpg "Design")

## Calender

### Days

* 6.30.2018 Day 1 I guess, I bought the byte.golf domain so ill probably do something with that at a later point. Also started the html/css design and I think it matches the style. I hate css alot so my site might not look good on bigger screens so you probably have to reduce the window size