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
)

// APIURI is the uri that execututes the api
const APIURI = "https://api.jdoodle.com/v1/execute"
const subBucket = "bytegolf-submissions"

// Send sends a CodeSubmission to the server compiler to check against the output
func (s *CodeSubmission) Send(storeLocal bool) (*CodeResponse, error) {
	var r CodeResponse

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

	// pass the information from one object to the other
	r.Info = s.Info
	r.UUID = s.UUID
	r.awsSess = s.awsSess // pass the aws session through

	// only store if the question is correct
	if storeLocal {
		go StoreLocal(s, &r)
	}
	return &r, nil
}

// StoreLocal creates a new CodeFile and stores it in the correct spot in the file system
// since this is run on its own go routine (even when it returns an error) it will log the error
func StoreLocal(s *CodeSubmission, r *CodeResponse) error {
	// log.SetPrefix("[debug] ")
	if !r.Check() {
		// the response was not correct
		log.Printf("%s was not correct\n", r.Output)
		return nil
	}
	// question was correct, so this holds the new score
	newLength := int(s.Score())

	var p = path.Join("localfiles", "codefiles", s.Info.QuestionID, s.Info.User)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.MkdirAll(p, os.ModePerm)
	} else {
		// The folder already exists, see if the file beats the old one and add it
		// check to make sure the previous submission does not exist
		filelist, err := ioutil.ReadDir(p)
		if err != nil {
			log.Printf("error reading files: %v\n", err)
			return err
		}
		if len(filelist) > 1 {
			// there should only be a single file per person per question
			log.Printf("expected less than 1 file but got %v\n", len(filelist))
			return fmt.Errorf("expected less than 1 file but got %v", len(filelist))
		}

		// there is either 1 file or 0 files, so just iterate to ignore the logic
		for _, fileinfo := range filelist {
			if fileinfo.Mode().IsRegular() {
				contents, err := ioutil.ReadFile(path.Join(p, fileinfo.Name()))
				if err != nil {
					log.Printf("error reading file %s : %v\n", fileinfo.Name(), err)
					return err
				}
				var prev CodeFile
				err = json.Unmarshal(contents, &prev)
				if err != nil {
					log.Printf("error unmarshalling file : %v\n", err)
					return err
				}
				// only update the file if the old is longer than the new
				if newLength >= prev.Length {
					// nothing needs to happen if the previous submission is better
					return nil
				}
				os.Remove(path.Join(p, fileinfo.Name()))
				log.Printf("removing file %s\n", path.Join(p, fileinfo.Name()))
			}
		}
	}
	// all of the folders are created and ready to insert a file into
	var cf = CodeFile{
		Submission: *s,
		Response:   *r,
		Correct:    true,
		Length:     newLength,
	}
	bs, err := json.Marshal(cf)
	if err != nil {
		log.Printf("err marshalling %v\n", err)
		return err
	}

	fileName := path.Join(p, cf.Submission.ID+".json")
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("err creating file: %v\n", err)
		return err
	}

	_, err = f.Write(bs)
	if err != nil {
		log.Printf("err writing to file: %v\n", err)
		return err
	}
	log.Printf("successfully wrote file %s", fileName)
	return nil
}

/*
REENABLE IN FUTURE VERSIONS FOR DB STORAGE
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
*/
