package opticaltextrecognition

import(
	"context"
	"encoding/json"
	"os"
	"fmt"
	"time"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/text/language"
)

//Global variables
const(
	projectID = "swamphacks2020-266920"
)

var(
	confi *ConfigurationFile
	visionFeature *vision.ImageAnnotatorClient
	translateFeature *translate.Client
	storageFeature *storage.Client
	pubSubCommunication *pubsub.Client
)
type ocrMessage struct {
	Text     string       `json:"text"`
	FileName string       `json:"fileName"`
	Lang     language.Tag `json:"lang"`
	SrcLang  language.Tag `json:"srcLang"`
}

// GCSEvent is the payload of a Google cloud storage event.
type GCSEvent struct {
	Bucket         string    `json:"bucket"`
	Name           string    `json:"name"`
	Metageneration string    `json:"metageneration"`
	ResourceState  string    `json:"resourceState"`
	TimeCreated    time.Time `json:"timeCreated"`
	Updated        time.Time `json:"updated"`
}


// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}


func setup(ctx context.Context) error {
	var err error
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	d := json.NewDecoder(file)
	err = d.Decode(&d)
	if err != nil {
		fmt.Errorf("Error: %s", err.Error())
	}
	
	return initFunction(ctx)
}


func initFunction(ctx context.Context) (err error) {
	if visionFeature == nil {
		visionFeature, err = vision.NewImageAnnotatorClient(ctx)
		if err != nil {
			return err
		}
	}
	if translateFeature == nil {
		translateFeature, err = translate.NewClient(ctx)
		if err != nil {
			return err
		}
	}
	if storageFeature == nil {
		storageFeature, err = storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}
	}
	if pubSubCommunication == nil {
		pubSubCommunication, err = pubsub.NewClient(ctx, projectID)
		if err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}
	}
	return 
}
