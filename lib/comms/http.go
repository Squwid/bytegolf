package comms

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Http is the receiver type that skips using GCP and is able to receive messages
// directly from the orchestrator for easier local running.
type Http struct {
	client *http.Client
	port   string
}

func (h *Http) Init() error {
	h.port = "8081"
	logrus.Infof("Initializing HTTP at port %s", h.port)
	h.client = http.DefaultClient
	return nil
}

func (h *Http) Publish(ctx context.Context, message []byte) error {
	url := fmt.Sprintf("http://%s:%s/compile", os.Getenv("BG_COMPILER_URL"), h.port)

	req, err := http.NewRequestWithContext(ctx, "POST",
		url+"/compile", bytes.NewReader(message))
	if err != nil {
		return err
	}
	_, err = h.client.Do(req)
	return err
}

func (h *Http) Listen(processor func(context.Context, string)) {
	r := mux.NewRouter()
	r.HandleFunc("/compile", httpHandler(processor)).Methods("POST")
	logrus.Fatalln(http.ListenAndServe(":"+h.port, r))
}

func httpHandler(processor func(context.Context, string)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithField("Action", "Http Handler")

		bs, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Errorf("Error reading body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger = logger.WithField("SubmissionID", string(bs))
		logger.Infof("Recieved message")

		go processor(context.Background(), string(bs))
	})
}
