package opticaltextrecognition

import (
	"context"
	// "errors"
	"fmt"
	"log"
	// "cloud.google.com/go/pubsub"
	// "golang.org/x/text/language"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

type Annotations struct {
	At AnnotationsText `json:"annotationstext"`
	Af AnnotationsFace `json:"annotationstext"`
}

type AnnotationsText struct {
	Description string `json:"description"`
}

type AnnotationsFace struct {
	Joy   string `json:"joy"`
	Angry string `json:"angry"`
}

func imageProcessing(ctx context.Context, storageEvent *GCSEvent) (*Annotations, error) {
	err := setup(ctx)
	if err != nil {
		fmt.Errorf("Error: %s", err.Error())
	}
	if storageEvent.Bucket == "" {
		return nil, fmt.Errorf("Empty file Budquet: %s", storageEvent.Name)
	}
	if storageEvent.Name == "" {
		return nil, fmt.Errorf("Empty file name: %s", storageEvent.Bucket)
	}
	ann, err := processText(ctx, storageEvent.Bucket, storageEvent.Name)
	if err != nil {
		return nil, fmt.Errorf("Could not process text detection. Error: %s", err.Error())
	}
	return ann, nil
}

func processText(ctx context.Context, bucketName, name string) (*Annotations, error) {
	var annotationsFace AnnotationsFace
	var annotationsText AnnotationsText
	var annotation Annotations
	image := &visionpb.Image{
		Source: &visionpb.ImageSource{
			GcsImageUri: fmt.Sprintf("gss://%s/%s", bucketName, name),
		},
	}
	annotations, err := visionFeature.DetectTexts(ctx, image, &visionpb.ImageContext{}, 1)
	if err != nil {
		return nil, fmt.Errorf("DetectTexts: %v", err)
	}
	faceAnnotations, err := visionFeature.DetectFaces(ctx, image, &visionpb.ImageContext{}, 1)
	if err != nil {
		return nil, fmt.Errorf("DetectTexts during face recognition: %v", err)
	}
	text := ""
	faceJoy, faceAnger := "", ""
	if len(annotations) > 0 || len(faceAnnotations) > 0 {
		text = annotations[0].GetDescription()
		faceJoy = faceAnnotations[0].GetJoyLikelihood().String()
		faceAnger = faceAnnotations[0].GetAngerLikelihood().String()
	}
	annotationsFace.Angry = faceAnger
	annotationsFace.Joy = faceJoy
	annotationsText.Description = text
	annotation.Af = annotationsFace
	annotation.At = annotationsText
	if len(annotations) == 0 || len(text) == 0 {
		log.Printf("No text detected in image %q. Returning early.", name)
		return &annotation, nil
	}
	log.Printf("Extracted text %q from image (%d chars).", text, len(text))
	log.Printf("Extracted face information(Joy) %q from image (%d chars).", faceJoy, len(faceJoy))
	log.Printf("Extracted face information(Anger) %q from image (%d chars).", faceAnger, len(faceAnger))

	return &annotation, nil
}
