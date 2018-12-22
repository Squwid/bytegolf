package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// APIURI is the uri that execututes the api
const APIURI = "https://api.jdoodle.com/execute"
const subBucket = "bytegolf-submissions"

// Send sends a CodeSubmission to the server compiler to check against the output
func (s *CodeSubmission) Send(storeLocal bool) (*CodeResponse, error) {
	var r CodeResponse

	if storeLocal {
		go s.storeLocal()
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

	// todo: storing concurrently does not check for an error
	if storeLocal {
		go r.storeLocal()
	}

	return &r, nil
}

func (s *CodeSubmission) storeLocal() error {
	var path = fmt.Sprintf("./subs/%s/", s.Info.User)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	f, err := os.Create(path + s.UUID)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(s.Script))
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (s *CodeResponse) storeLocal() error {
	var path = fmt.Sprintf("./resp/%s/", s.Info.User)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	f, err := os.Create(path + s.UUID)
	if err != nil {
		return err
	}
	defer f.Close()

	store, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = f.Write(store)
	if err != nil {
		return err
	}
	return nil
}

// Store stores a submission to an S3 Bucket
func (s *CodeSubmission) store() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("us-east-1")},
	}))

	uploader := s3manager.NewUploader(sess)

	//
	key := fmt.Sprintf("%s/%s/sub_%s_%s", s.Info.Hole, s.Info.User, s.Info.Name, s.UUID)

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
	key := fmt.Sprintf("%s/%s/resp_%s_%s", s.Info.Hole, s.Info.User, s.Info.Name, s.UUID)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(subBucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(fmt.Sprintf("Status: %v\nMemory: %s\nCPU Time: %s\nOutput:\n%s\n", s.StatusCode, s.Memory, s.CPUTime, s.Output)),
	})
	if err != nil {
		log.Fatalf("an error occurred storing a code response: %s\n", err.Error())
		return
	}
}
