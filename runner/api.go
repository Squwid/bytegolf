package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// APIURI is the uri that execututes the api
const APIURI = "https://api.jdoodle.com/execute"
const subBucket = "bytegolf-submissions"

// Send todo
func (s *CodeSubmission) Send() (*CodeResponse, error) {
	var r CodeResponse

	if s.Config.SaveSubmissions {
		go s.store() // store concurrently during the send
	}

	reqBody, err := json.Marshal(*s)
	if err != nil {
		return &CodeResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, APIURI, bytes.NewBuffer(reqBody))
	if err != nil {
		return &CodeResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &CodeResponse{}, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return &CodeResponse{}, err
	}

	r.Info = s.Info
	r.UUID = s.UUID
	if s.Config.SaveSubmissions {
		go r.store() // store response concurrently
	}

	return &r, nil
}

// Store stores a submission to an S3 Bucket
func (s *CodeSubmission) store() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("us-east-1")},
	}))

	uploader := s3manager.NewUploader(sess)

	//
	key := fmt.Sprintf("%s_%s/%s/sub_%s_%s", s.Info.GameName, s.Info.Game, s.Info.User, s.Info.FileName, s.UUID)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(subBucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(s.Script),
	})
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

// Store a code submission to an S3 bucket
func (s *CodeResponse) store() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("us-east-1")},
	}))

	uploader := s3manager.NewUploader(sess)

	//
	key := fmt.Sprintf("%s_%s/%s/res_%s_%s", s.Info.GameName, s.Info.Game, s.Info.User, s.Info.FileName, s.UUID)

	// TODO: deal with this error in the future
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(subBucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(fmt.Sprintf("Status: %v\nMemory: %s\nCPU Time: %s\nOutput:\n%s\n", s.StatusCode, s.Memory, s.CPUTime, s.Output)),
	})
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
