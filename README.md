# Byte Golf

Byte Golf is a code solving problem game that awards players with the shortest code that solves the task at hand, just like golf! I have had this idea for a while but I couldn't find a good time to start. So I thought today (6.30.18) is as good as time as any  This is going to be my first commit which is going to include a quick drawing I did for my own reference, pardon my drawing skills.

## CURRENT UPDATE (OCTOBER 2018)

I currently have started a company with my roommate and am in my senior year at college, so this project will continue to be updated at a reduced speed. However I want to change the scope of the project to host externally instead of local on an ec2 instance. I already own the domain byte.golf so all I have to do is make a site for it and upgrade security. This project is open source so I want to find a good balance between security and open source to show what I can accomplish.

## Rules

The rules for Byte Golf are suprisingly simple, just write the least code using the standard library of a language of your choice.

1. The player with the lowest total characters at the end wins, just like Golf!
1. Total score consists of characters or "bytes" in a working code solution (whitespace not included). The scores are totaled at the end, and all holes must be completed to a working solution to be able to win. Hence the name "Byte Golf"
1. Golfing in any order is allowed, as long as all hole solutions are completed by the conclusion of the game.
1. All code must be written in one of the supported languages and must only use standard libraries.

## Setup

The setup for this app to run locally is pretty simple, just get api credits on JDoodle and get started.

1. Navigate to [JDoodle](https://www.jdoodle.com/compiler-api) and subscribe to obtain tokens to compile the solutions.
1. Set the environmental variables of both `RUNNER_ID` and `RUNNER_SECRET` to the Client ID and Client Secret on the JDoodle API
1. Find the config.yml file and configure it to the settings you would like.
1. The last and final step is to use the `make` command which will compile and start the application.

## Design

### Origin

My roommates and I came up with the idea of making Byte Golf a web app for us to see who can solve code in the shortest solution. I know this already exists (sort of) on Stack Overflow. However there is no solution and any API that compiles and verifies the code.

### Thoughts

My current idea is to create a simple web app that players can clone and compile locally, and host a LAN server. The LAN server the user runs on their local machine will use an AWS Lambda that I will setup to deal with a global leaderboard system as well as the code compiler using DynamoDB and S3. Different langauges can give players more points because they aren't the best programming language Bash.

![Design](https://i.imgur.com/SVmaRN6.jpg "Design")