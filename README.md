# ResNet50v2 scoring (inference) in Go

Tool to classify an image with the ResNet50v2 model from the ONNX Model Zoo. The tool expects a running container based on the gbaeke/onnxresnet50v2 image.

The container exposes a scoring URI at http://localhost:5001/score.
