# Byte Golf

Byte Golf is a code solving problem that rewards players for solving code problems in the shortest amount of code! Some questions require parsing of Stdin while others have to do with printing a simple output. 

## Rules

The rules for Byte Golf are surprisingly simple, just write the least code using the standard library of a language of your choice.

1. The player with the lowest total bytes at the end wins, just like Golf!
1. Total score consists of bytes in a working code solution (whitespace and comments included). The test output is NOT shown to the user, *this will probably change in the future*. 
1. All code must be written in one of the supported languages and must only use standard libraries.

## Design

### Origin

A few years ago, my roommates and I decided that it would be cool to have a game where code problems could be solved in interesting and short ways. I wanted to open source the entire solution as a way to show my progress of backend, frontend, and cloud hosting stack. Take note that I am new to frontend development, so making custom components, icons and loading screens to implement rather than design kits was exciting.

### Tech Stack

The tech stack consists of a Go backend, React Typescript frontend, and GCP cloud hosting using the Firestore database in Native Mode. Both the frontend and backend are hosted in Google Cloud Run in Docker containers to keep the bill low as well as leveraging a completely serverless solution. I *would* host the frontend using static storage behind a load balancer with Cloud Armor, but it starts to get expensive quick, this might be revisited in the future.

### Running Submissions Remotely

Currently, I am using the [Jdoodle API](https://docs.jdoodle.com/compiler-api/compiler-api) to compile code remotely, but am working on another non-open sourced way to run code remotely at scale.

## Future

There are still many features to implement before this is a complete working product at scale. Those will be listed below

### TODOs:

- Profiles
- Output of some test cases, also some hidden test cases
    - It would be good to see exactly what test cases failed and which ones passed
- Pagination of Submissions
- Logout
- Footer
- Automated build pipelines
- ~~Submitted Time~~
- Colored text editor in submissions modal (based on language)

## FAQ

- How were the logos and icons made?
    - They were made in Photoshop and After Effects