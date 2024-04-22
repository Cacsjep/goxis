## Build
To add the models to the final .eap package run for example **goxisbuilder** from repo root directory.

```shell
git clone https://github.com/Cacsjep/goxis
cd goxis
goxisbuilder.exe -appdir="examples/axlarod/yolov5" -files yolov5n.tflite
goxisbuilder.exe -appdir="examples/axlarod/classify" -files converted_model.tflite
goxisbuilder.exe -appdir="examples/axlarod/object_detection" -files ssd_mobilenet_v2_coco_quant_postprocess.tflite
```