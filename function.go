package cryptocurrencies

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/fauna/faunadb-go/v4/faunadb"
	"github.com/sirupsen/logrus"
)

func HelloHTTP(w http.ResponseWriter, r *http.Request) {
	// initialize the storage service
	store, err := New()
	if err != nil {
		Log.Fatal("failed to initialize database ", err)
	}

	basicError := BasicError{
		Code:    "SERVER_ERROR",
		Message: "failed to get cryptocurrencies",
	}

	l, err := store.List(r.Context())
	if err != nil {
		Log.Info(err)
		b, err := json.Marshal(basicError)
		if err != nil {
			Log.Info(err)
			return
		}

		w.Write(b)
		return
	}

	b, err := json.Marshal(l)
	if err != nil {
		Log.Info(err)
		b, err := json.Marshal(basicError)
		if err != nil {
			Log.Info(err)
			return
		}

		w.Write(b)
		return
	}

	w.Write(b)
}

// Cryptocurrency ...
type Cryptocurrency struct {
	Name   string  `fauna:"name" json:"name"`
	Price  float32 `fauna:"price" json:"price"`
	Symbol string  `fauna:"symbol" json:"symbol"`
}

type ListRes struct {
	Data Cryptocurrency `fauna:"data" json:"data"`
}

// BasicError is an error structure that is returned to the API caller indicating the error code and the message
type BasicError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Log is the main logger
var Log *logrus.Logger

func init() {
	Log = logrus.New()

	env := os.Getenv("CRYPTOCURRENCIES")

	Log.SetReportCaller(true)

	if env == "dev" {
		Log.SetLevel(logrus.DebugLevel)
	}

}

// Fauna database
type Fauna struct {
	Client *faunadb.FaunaClient
}

// New creates a new client connection to the Fauna database
func New() (*Fauna, error) {

	client := faunadb.NewFaunaClient(os.Getenv("FAUNA_SECRET"))

	return &Fauna{Client: client}, nil
}

// List cryptocurrencies
func (f Fauna) List(c context.Context) ([]*ListRes, error) {

	res, err := f.Client.Query(faunadb.Map(faunadb.Paginate(faunadb.Documents(faunadb.Collection("cryptocurrencies"))), faunadb.Lambda("ref", faunadb.Get(faunadb.Var("ref")))))

	if err != nil {
		Log.Info(err)
		return nil, err
	}

	var cryptocurrencies []*ListRes

	Log.Info(res)

	if err := res.At(faunadb.ObjKey("data")).Get(&cryptocurrencies); err != nil {
		Log.Info(err)
		return nil, err
	}

	return cryptocurrencies, nil
}
