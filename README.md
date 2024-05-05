# Byte Golf

Play Byte Golf - https://byte.golf/

Byte Golf is a code solving problem that rewards players for solving code problems in the shortest amount of code! Some questions require parsing of Stdin while others have to do with printing a simple output. Byte Golf is completely open sourced [here](https://github.com/Squwid/bytegolf) as well!

## Rules

The rules for Byte Golf are surprisingly simple, just write the least code using the standard library of a language of your choice.

1. The player with the lowest total bytes at the end wins, just like Golf!
1. Total score consists of bytes in a working code solution (whitespace and comments included). Some holes have multiple or hidden test cases, especially dealing with Stdin.
1. All code must be written in one of the supported languages and must only use standard libraries.

ex. _"Print 'Hello, World!' to console"_ could be solved with python3 in just 22 bytes!

```python
print('Hello, World!')
```

## Design

### Origin

A few years ago, my roommates and I decided that it would be cool to have a game where code problems could be solved in interesting and short ways. I wanted to open source the entire solution as a way to show my full stack development skills. The frontend was designed using Adobe XD, and the logos and pages were created in Adobe Illustrator (I am by no means a graphic designer). The entire stack is hosted using serverless technologies on Google Cloud Platform, and will be described in more detail in the next section.

### Tech Stack

The stack consists of a Go backend, Typescript React frontend, while leveraging Google Cloud Platform's Cloud Run serverless runtime for cheap hosting. I am using Google Cloud Firestore as my database in Native mode. This gives me scalability without increasing my bill to the moon. The backend is stateless and containerized so it can run basically anywhere, having secrets injected via Google Cloud Secrets.

### Running Submissions Remotely

Currently, I am using the [Jdoodle API](https://docs.jdoodle.com/compiler-api/compiler-api) to compile code remotely, but there is a custom remote compiler in the works ;)

### Site Design

Earlier I mentioned that all of the logos, loading icon, and pages were all designed and created by me as well. Here are some examples of those!

#### Logos!

##### Default Bytegolf Logo

![Bytegolf Logo](https://raw.githubusercontent.com/Squwid/bytegolf/master/assets/logo/bytegolf_logo-half.png)

##### Not Found Logo

![Bytegolf Not Found Logo](https://raw.githubusercontent.com/Squwid/bytegolf/master/assets/logo/bytegolf-logo-not-found-half.png)

##### Loading Icon

![Bytegolf Loading Icon](https://raw.githubusercontent.com/Squwid/bytegolf/master/assets/golf-ball/loading-icon.gif)

This was animated in Adobe After Effects!


### TODOs:

- Profiles
- <s>Output of some test cases, also some hidden test cases</s>
    - <s>It would be good to see exactly what test cases failed and which ones passed</s>
- Pagination of Submissions
- Logout
- Footer
- Automated build pipelines
- <s>Submitted Time</s>
- Colored text editor in submissions modal (based on language)

## FAQ

- How were the logos and icons made?
    - They were made in Photoshop and After Effects