package processor

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/stretchr/testify/assert"
)

const boilerplate = `
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	const input = "input.txt"

	file, err := os.Open(input)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	fifthWord, i := mostFrequentWords(scanner, 5)
	fmt.Printf("%v:%v\n", fifthWord, i)
}

{{.UserSolution}}
`

const script = `
func mostFrequentWords(scanner *bufio.Scanner, top int) (string, int) {
	// Create a map to store the frequency of each word
	wordFrequency := make(map[string]int)

	// Scan through each word in the file
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		word = strings.Trim(word, ".,;:!?\"'()[]{}")
		wordFrequency[word]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Create a slice to store the top n most frequent words
	topWords := make([]string, top)

	// Find the top n most frequently occurring words
	for i := 0; i < top; i++ {
		mostFrequentWord := ""
		mostFrequentWordCount := 0

		for word, frequency := range wordFrequency {
			if frequency > mostFrequentWordCount && !contains(topWords, word) {
				mostFrequentWord = word
				mostFrequentWordCount = frequency
			}
		}
		topWords[i] = mostFrequentWord
	}

	// for i := 0; i < len(topWords); i++ {
	// 	fmt.Printf("%v: %v:%v\n", i+1, topWords[i], wordFrequency[topWords[i]])
	// }

	return topWords[top-1], wordFrequency[topWords[top-1]]
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

`

const expectedOutput = `
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	const input = "input.txt"

	file, err := os.Open(input)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	fifthWord, i := mostFrequentWords(scanner, 5)
	fmt.Printf("%v:%v\n", fifthWord, i)
}


func mostFrequentWords(scanner *bufio.Scanner, top int) (string, int) {
	// Create a map to store the frequency of each word
	wordFrequency := make(map[string]int)

	// Scan through each word in the file
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		word = strings.Trim(word, ".,;:!?\"'()[]{}")
		wordFrequency[word]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Create a slice to store the top n most frequent words
	topWords := make([]string, top)

	// Find the top n most frequently occurring words
	for i := 0; i < top; i++ {
		mostFrequentWord := ""
		mostFrequentWordCount := 0

		for word, frequency := range wordFrequency {
			if frequency > mostFrequentWordCount && !contains(topWords, word) {
				mostFrequentWord = word
				mostFrequentWordCount = frequency
			}
		}
		topWords[i] = mostFrequentWord
	}

	// for i := 0; i < len(topWords); i++ {
	// 	fmt.Printf("%v: %v:%v\n", i+1, topWords[i], wordFrequency[topWords[i]])
	// }

	return topWords[top-1], wordFrequency[topWords[top-1]]
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


`

func TestSubmissionInit_CreatesFolderAndFiles(t *testing.T) {
	sub := &Submission{
		ID: "100",
		Hole: &api.HoleDB{
			TestsDB: []*api.TestDB{
				{
					ID:          math.MaxInt64 - 1,
					Boilerplate: boilerplate,
				},
			},
			LanguageDB: &api.LanguageDB{
				Extension: "go",
			},
		},
		Submission: &api.SubmissionDB{
			Script: script,
		},
	}

	assert.NoError(t, sub.Init())
	assert.Contains(t, sub.tempDir, "/tmp/bg-100_")
	assert.DirExists(t, sub.tempDir)
	assert.FileExists(t, filepath.Join(sub.tempDir, "main-9223372036854775806.go"))

	bs, err := os.ReadFile(filepath.Join(sub.tempDir, "main-9223372036854775806.go"))
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, string(bs))

	sub.Clean()
}

func TestSubmissionClean_DeletesFolderAndFiles(t *testing.T) {
	sub := &Submission{
		ID: "999",
		Hole: &api.HoleDB{
			TestsDB: []*api.TestDB{
				{
					ID:          9,
					Boilerplate: boilerplate,
				},
			},
			LanguageDB: &api.LanguageDB{
				Extension: "go",
			},
		},
		Submission: &api.SubmissionDB{
			Script: script,
		},
	}

	assert.NoError(t, sub.Init())
	assert.FileExists(t, filepath.Join(sub.tempDir, "main-9.go"))
	assert.DirExists(t, sub.tempDir)

	assert.NoError(t, sub.Clean())
	assert.NoFileExists(t, filepath.Join(sub.tempDir, "main-9.go"))
	assert.NoDirExists(t, sub.tempDir)
}
