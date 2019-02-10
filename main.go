package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
)

//InputData is sent to ResNet50 model
type InputData struct {
	Data [1][3][224][224]uint8 `json:"data"`
}

//OutputData is received from ResNet50 model
type OutputData struct {
	Result []float64 `json:"result"`
	Time   float64   `json:"time"`
}

func main() {
	imagePath := flag.String("image", "", "path to image")
	scoreURI := flag.String("uri", "http://localhost:5100/score", "scoring URI")
	flag.Parse()

	if *imagePath == "" {
		log.Fatal("Please specify image path. Use --help for help.")
	}

	//read image and resize
	m, err := readImage(*imagePath, 224, 224)
	if err != nil {
		log.Fatal(err)
	}

	//get bounds of image 0 0 224 224
	bounds := m.Bounds()

	// multidim array as input tensor
	var BCHW [1][3][224][224]uint8

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := m.At(x, y).RGBA()

			// height = y and width = x
			BCHW[0][0][y][x] = uint8(r >> 8)
			BCHW[0][1][y][x] = uint8(g >> 8)
			BCHW[0][2][y][x] = uint8(b >> 8)
		}
	}

	// input is struct with 4D array
	input := InputData{
		Data: BCHW,
	}

	// Create JSON from input struct - inputJSON will be sent to model
	inputJSON, _ := json.Marshal(input)
	body := bytes.NewBuffer(inputJSON)

	// Create the HTTP request - no need for auth with local ResNet50 container
	client := &http.Client{}
	request, err := http.NewRequest("POST", *scoreURI, body)
	request.Header.Add("Content-Type", "application/json")

	fmt.Printf("Scoring %s against %s.\n", *imagePath, *scoreURI)

	// Send the request to the local web service
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("Error calling scoring URI: ", err)
	}

	// read response
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	//Unmarshal returned JSON data
	var modelResult OutputData
	err = json.Unmarshal(respBody, &modelResult)
	if err != nil {
		log.Fatal("Error unmarshalling JSON response ", err)
	}

	// highest result
	maxProb := 0.0
	maxIndex := 0
	for index, prob := range modelResult.Result {
		if prob > maxProb {
			maxProb = prob
			maxIndex = index
		}
	}

	// get the categories
	categories, err := getCategories()
	if err != nil {
		log.Fatal("Error getting categories", err)
	}

	fmt.Println("Highest prob is", maxProb, "at", maxIndex, "(inference time:", modelResult.Time, ")")
	fmt.Println("Probably ", categories[maxIndex])
}

func readImage(imgPath string, width, height int) (image.Image, error) {
	// read the image file
	reader, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// decode the image
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// resize image
	m = imaging.Resize(m, width, height, imaging.Linear)

	return m, nil
}

func getCategories() (map[int][]string, error) {
	// open categories file
	reader, err := os.Open("imagenet_class_index.json")
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// read JSON categories
	catJSON, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// unmarshal into map of int to array of string
	var categories map[int][]string
	err = json.Unmarshal(catJSON, &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
