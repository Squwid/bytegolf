package receiver

import (
	"context"
	"io"
	"net/http"

	"github.com/Squwid/bytegolf/compiler/processor"
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
	logrus.Infof("Initializing HTTP receiver at port %s", h.port)
	h.client = http.DefaultClient
	return nil
}

func (h *Http) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/compile", httpHandler).Methods("POST")
	logrus.Fatalln(http.ListenAndServe(":"+h.port, r))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("Action", "Http Handler")

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Error reading body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger = logger.WithField("SubmissionID", string(bs))
	logger.Infof("Recieved message")

	processor.ProcessMessage(context.Background(), string(bs))
}
