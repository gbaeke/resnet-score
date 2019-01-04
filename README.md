# ResNet50v2 scoring (inference) in Go

Tool to classify an image with the ResNet50v2 model from the ONNX Model Zoo. The tool expects a running container based on the gbaeke/onnxresnet50v2 image.

Use the following command to start the container:

docker run -d -p 5001:5001 gbaeke/onnxresnet50v2

Use tool as follows:

./class -image image-path

e.g. ./class -image images/cat.jpg


The container exposes a scoring URI at http://localhost:5001/score.

Download tool for Linux from releases: [Download](https://github.com/gbaeke/resnet-score/releases/download/v0.1/class)
