package runner

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

// GetPlayerSubmissions gets a list of all of the code submissions that they have previously submitted by username
func GetPlayerSubmissions(name string) ([]CodeSubmission, error) {
	var subPath = path.Join("localfiles", "subs", name)
	filelist, err := ioutil.ReadDir(subPath)
	if err != nil {
		return nil, err
	}

	var cs = []CodeSubmission{}
	for _, fileinfo := range filelist {
		if fileinfo.Mode().IsRegular() {
			contents, err := ioutil.ReadFile(path.Join(subPath, fileinfo.Name()))
			if err != nil {
				return nil, err
			}

			var sub CodeSubmission
			err = json.Unmarshal(contents, &sub)
			if err != nil {
				return nil, err
			}
			cs = append(cs, sub)
		}
	}
	return cs, nil
}

// GetPlayerResponses gets a list of all of the code responses that they have previously submitted by username
func GetPlayerResponses(name string) ([]CodeResponse, error) {
	var respPath = path.Join("localfiles", "resp", name)
	filelist, err := ioutil.ReadDir(respPath)
	if err != nil {
		return nil, err
	}

	var cr = []CodeResponse{}
	for _, fileinfo := range filelist {
		if fileinfo.Mode().IsRegular() {
			contents, err := ioutil.ReadFile(path.Join(respPath, fileinfo.Name()))
			if err != nil {
				return nil, err
			}

			var resp CodeResponse
			err = json.Unmarshal(contents, &resp)
			if err != nil {
				return nil, err
			}
			cr = append(cr, resp)
		}
	}
	return cr, nil
}
