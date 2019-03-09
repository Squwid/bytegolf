package runner

import (
	"log"
	"strings"

	"github.com/Squwid/bytegolf/questions"
)

// Check checks the response against the correct answer and sees if it is correct
// only returns a bool, logs any errors that have occurred
func (resp *CodeResponse) Check() bool {
	q := questions.GetQuestion(resp.Info.QuestionID)
	if q.ID == "" {
		log.Printf("question %s is not live\n", resp.Info.QuestionID)
		return false
	}

	// todo: Explore more jdoodle status codes
	// if resp.StatusCode != http.StatusOK {
	// 	// make sure status is 200
	// 	// log.Printf("expected status 200 but got %v\n", resp.StatusCode)
	// 	return false
	// }
	if q.Answer == strings.TrimSpace(resp.Output) {
		// Output is the same as the answer so it is correct
		return true
	}

	return false
}

/*
	Scoring System:
	Each character that is not white space or newline is 1 point
	Max score per problem is 300 (Maybe this needs to change)
*/

type strPair struct {
	begin string
	end   string
}

func gatherCustomCommentTags(code string, tags []strPair, sizeTags uint, tagFinder strPair) ([]strPair, uint) {
	var strFind = -3

	for true {
		code = code[strFind+3:]
		strFind = strings.Index(code, tagFinder.begin)
		//fmt.Printf("strFind: %d\n", strFind)

		if strFind == -1 {
			return tags, sizeTags
		}

		var tagStart = strFind + len(tagFinder.begin)
		var tagEnd = tagStart
		if len(tagFinder.end) == 0 {
			tagEnd = tagEnd + 1
		} else {
			tagEnd = tagEnd + strings.Index(code[tagStart:], tagFinder.end)
		}
		var pair = strPair{code[strFind:tagEnd] + "\n", "\n" + code[tagStart:tagEnd]}
		//fmt.Printf("tags.append(%s, %s)\n", pair.begin, pair.end)
		tags = append(tags, pair)
		sizeTags = sizeTags + 1
	}

	return tags, sizeTags
}

func genericScore(code string, comments []strPair, commentsSize uint, strings []strPair, stringsSize uint) uint {
	var isInComment, isInString, isEscaped bool

	var score uint
	var codeLen = len(code)
	var endCommentTag, endStringTag string
	var endTagSize int

	for i, c := range code {
		//fmt.Printf("loop @ %d %c score %d string:%t comment:%t escaped:%t|%t\n", i, c, score, isInString, isInComment, isEscaped, c == '\\')

		if !isEscaped {
			if isInComment {
				isInComment = false

				for j := range endCommentTag {
					//fmt.Printf("i: %d j: %d endTagSize: %d codeLen: %d\n", i, j, endTagSize, codeLen)
					if code[i+j-endTagSize] != endCommentTag[j] {
						isInComment = true
						break
					}
				}

				//if !isInComment {
				//fmt.Printf("Comment ended by '%s' tag!\n", endCommentTag)
				//}

				continue
			} else if isInString {
				if i+1 < codeLen && i-endTagSize > -2 {
					isInString = false

					for j := range endStringTag {
						//fmt.Printf("i: %d j: %d endTagSize: %d codeLen: %d\n", i, j, endTagSize, codeLen)
						if code[i+j-endTagSize] != endStringTag[j] {
							isInString = true
							break
						}
					}
				}

				//if !isInString {
				//fmt.Printf("String ended by '%s' tag!\n", endStringTag)
				//}

				if c == '\\' {
					isEscaped = true
				}

				score = score + 1
				continue
			} else if i < codeLen {
				var j uint

				for j = 0; j < commentsSize; j++ {
					var match = true
					var k int
					for k = 0; k < len(comments[j].begin); k++ {
						if i+k >= codeLen {
							//fmt.Printf("Not a comment because were at the end of the code!\n")
							match = false
							break
						}

						if code[i+k] != comments[j].begin[k] {
							//fmt.Printf("%d,%d Not a comment because %c != %c (%d)\n", i, k, code[i + k], comments[j].begin[k], j)
							match = false
							break
						}
					}
					if match {
						//fmt.Printf("Comment started by '%s' tag!\n", comments[j].begin)
						isInComment = true
						endCommentTag = comments[j].end
						endTagSize = len(endCommentTag) + 1
						break
					}
				}

				for j = 0; j < stringsSize; j++ {
					var match = true
					var k int
					for k = 0; k < len(strings[j].begin); k++ {
						if i+k >= codeLen {
							//fmt.Printf("Not a string because were at the end of the code!\n")
							match = false
							break
						}

						if code[i+k] != strings[j].begin[k] {
							//fmt.Printf("%d,%d Not a string because %c != %c (%d)\n", i, k, code[i + k], strings[j].begin[k], j)
							match = false
							break
						}
					}
					if match {
						//fmt.Printf("String started by '%s' tag!\n", strings[j].begin)
						isInString = true
						endStringTag = strings[j].end
						endTagSize = len(endStringTag) + 1
						break
					}
				}
			}
		}

		if !isInComment && (isInString || (c != ' ' && c != '\t' && c != '\n')) {
			score = score + 1
		}

		isEscaped = false
	}

	return score
}

// Score checks the score of a submission, however does not check if the submission is correct so that needs to
// be done ahead of time.
// Your score is determined by any characters except for spaces
// func (q *Question) Score(sub *CodeSubmission) uint {
func (sub *CodeSubmission) Score() uint {
	if sub.Language == LangJava || sub.Language == LangCPP || sub.Language == LangCPP14 || sub.Language == LangC {
		return genericScore(sub.Script, []strPair{{"//", "\n"}, {"/*", "*/"}}, 2, []strPair{{"\"", "\""}, {"'", "'"}}, 2)
	}

	if sub.Language == LangPHP {
		// handle <<<*\n<anything>\n*; style comments
		var tags, sizeTags = gatherCustomCommentTags(sub.Script, []strPair{{"\"", "\""}, {"'", "'"}}, 2, strPair{"<<<", "\n"})
		return genericScore(sub.Script, []strPair{{"//", "\n"}, {"#", "\n"}, {"/*", "*/"}}, 3, tags, sizeTags)
	}

	if sub.Language == LangPy2 || sub.Language == LangPy3 {
		return genericScore(sub.Script, []strPair{{"#", "\n"}, {"\"\"\"", "\"\"\""}}, 2, []strPair{{"\"\"\"", "\"\"\""}, {"\"", "\""}, {"'", "'"}}, 3)
	}

	if sub.Language == LangRuby {
		// need to catch heredoc <<*\nDATA\n*\n
		var tags, sizeTags = gatherCustomCommentTags(sub.Script, []strPair{{"\"", "\""}, {"'", "'"}, {"%(", ")"}, {"%[", "]"}, {"%{", "}"}}, 5, strPair{"<<", "\n"})
		// need to catch %*-* comments
		tags, sizeTags = gatherCustomCommentTags(sub.Script, []strPair{{"\"", "\""}, {"'", "'"}}, 2, strPair{"%", ""})
		return genericScore(sub.Script, []strPair{{"#", "\n"}, {"=begin", "=end"}}, 2, tags, sizeTags)
	}

	if sub.Language == LangGo {
		return genericScore(sub.Script, []strPair{{"//", "\n"}, {"/*", "*/"}}, 2, []strPair{{"\"", "\""}, {"'", "'"}, {"`", "`"}}, 3)
	}

	if sub.Language == LangBash {
		// the first condition should realistically catch all bash comments, the second is rare, but there are other very rare comment tecniques i dont know
		return genericScore(sub.Script, []strPair{{"#", "\n"}, {"${IFS#", "\n"}}, 2, []strPair{{"\"", "\""}, {"'", "'"}}, 2)
	}

	if sub.Language == LangSwift {
		return genericScore(sub.Script, []strPair{{"//", "\n"}, {"/*", "*/"}}, 2, []strPair{{"\"\"\"", "\"\"\""}, {"\"", "\""}, {"'", "'"}}, 3)
	}

	if sub.Language == LangR {
		return genericScore(sub.Script, []strPair{{"#", "\n"}}, 1, []strPair{{"\"", "\""}, {"'", "'"}}, 2)
	}

	if sub.Language == LangNode {
		return genericScore(sub.Script, []strPair{{"//", "\n"}, {"/*", "*/"}}, 2, []strPair{{"\"", "\""}, {"'", "'"}}, 2)
	}

	if sub.Language == LangFS {
		return genericScore(sub.Script, []strPair{{"//", "\n"}, {"(*", "*)"}}, 2, []strPair{{"\"\"\"", "\"\"\""}, {"@\"", "\""}, {"\"", "\""}, {"'", "'"}}, 4)
	}

	return count(sub.Script)
}

// TODO: remove newlines and spaces inside of strings
func count(s string) uint {
	var c uint
	for _, l := range s {
		if len(strings.TrimSpace(string(l))) == 0 {
			continue
		}
		c++
	}
	return c
}
