package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// APIURI is the uri that execututes the api
const APIURI = "https://api.jdoodle.com/v1/execute"
const subBucket = "bytegolf-submissions"

// Send sends a CodeSubmission to the server compiler to check against the output
func (s *CodeSubmission) Send(storeLocal bool) (*CodeResponse, error) {
	var r CodeResponse
	if storeLocal {
		go s.storeLocal()
	}

	reqBody, err := json.Marshal(*s)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, APIURI, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	r.Info = s.Info
	r.UUID = s.UUID
	r.awsSess = s.awsSess // pass the aws session through
	if storeLocal {
		go r.storeLocal()
	}
	return &r, nil
}

/* FUNCTIONS RELATING TO STORING THE SUBMISSIONS AND RESPONSES LOCALLY */
func (s *CodeSubmission) storeLocal() error {
	var p = path.Join("./subs", s.Info.User)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.MkdirAll(p, os.ModePerm)
	}

	f, err := os.Create(path.Join(p, s.UUID))
	if err != nil {
		return err
	}
	defer f.Close()

	bs, err := json.Marshal(*s)
	if err != nil {
		return err
	}
	_, err = f.Write(bs)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}

// Store local stores a code response to the local file system rather than an s3 bucket
func (s *CodeResponse) storeLocal() error {
	var p = path.Join("./resp", s.Info.User)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.MkdirAll(p, os.ModePerm)
	}

	f, err := os.Create(path.Join(p, s.UUID))
	if err != nil {
		return err
	}
	defer f.Close()

	store, err := json.Marshal(*s)
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
	key := path.Join(s.Info.Hole, s.Info.User, fmt.Sprintf("sub_%s_%s", s.Info.Name, s.UUID))
	// key := fmt.Sprintf("%s/%s/sub_%s_%s", s.Info.Hole, s.Info.User, s.Info.Name, s.UUID)
	uploader := s3manager.NewUploader(s.awsSess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(subBucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(s.Script),
	})
	if err != nil {
		log.Printf("error storing %s submission : %v\n", s.UUID, err)
		return
	}
}

// Store a code submission to an S3 bucket
func (s *CodeResponse) store() {
	key := path.Join(s.Info.Hole, s.Info.User, fmt.Sprintf("resp_%s_%s", s.Info.Name, s.UUID))
	uploader := s3manager.NewUploader(s.awsSess)
	// key := fmt.Sprintf("%s/%s/resp_%s_%s", s.Info.Hole, s.Info.User, s.Info.Name, s.UUID)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(subBucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(fmt.Sprintf("Status: %v\nMemory: %s\nCPU Time: %s\nOutput:\n%s\n", s.StatusCode, s.Memory, s.CPUTime, s.Output)),
	})
	if err != nil {
		log.Printf("error storing %s response : %v\n", s.UUID, err)
		return
	}
}
