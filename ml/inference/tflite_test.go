package inference

import (
	"path/filepath"
	"runtime"
	"testing"

	tflite "github.com/mattn/go-tflite"
	"go.viam.com/test"
)

const badPath string = "bad path"

var (
	// used to get the path from the root to current directory.
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

type fakeInterpreter struct{}

func (fI *fakeInterpreter) AllocateTensors() tflite.Status {
	return tflite.OK
}

func (fI *fakeInterpreter) Invoke() tflite.Status {
	return tflite.OK
}

func (fI *fakeInterpreter) GetOutputTensorCount() int {
	return 1
}

func (fI *fakeInterpreter) GetInputTensorCount() int {
	return 1
}

func (fI *fakeInterpreter) GetOutputTensor(i int) *tflite.Tensor {
	return &tflite.Tensor{}
}

func (fI *fakeInterpreter) GetInputTensor(i int) *tflite.Tensor {
	return &tflite.Tensor{}
}

func (fI *fakeInterpreter) Delete() {}

var goodOptions *tflite.InterpreterOptions = &tflite.InterpreterOptions{}

func goodGetInfo(i Interpreter) *TFLiteInfo {
	return &TFLiteInfo{}
}

// TestLoadModel uses a real tflite model to test loading.
func TestLoadModel(t *testing.T) {
	tfliteModelPath := basePath + "/testing_files/model_with_metadata.tflite"
	loader, err := NewDefaultTFLiteModelLoader()
	test.That(t, err, test.ShouldBeNil)
	tfliteStruct, err := loader.Load(tfliteModelPath)
	test.That(t, tfliteStruct, test.ShouldNotBeNil)
	test.That(t, err, test.ShouldBeNil)

	structInfo := tfliteStruct.Info
	test.That(t, structInfo, test.ShouldNotBeNil)

	h := structInfo.InputHeight
	w := structInfo.InputWidth
	c := structInfo.InputChannels
	test.That(t, h, test.ShouldEqual, 640)
	test.That(t, w, test.ShouldEqual, 640)
	test.That(t, c, test.ShouldEqual, 3)
	test.That(t, structInfo.InputTensorType, test.ShouldEqual, "Float32")
	test.That(t, structInfo.InputTensorCount, test.ShouldEqual, 1)
	test.That(t, structInfo.OutputTensorCount, test.ShouldEqual, 4)
	test.That(t, structInfo.OutputTensorTypes, test.ShouldResemble, []string{"Float32", "Float32", "Float32", "Float32"})

	buf := make([]float32, c*h*w)
	outTensors, err := tfliteStruct.Infer(buf)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, outTensors, test.ShouldNotBeNil)
	test.That(t, len(outTensors), test.ShouldEqual, 4)

	tfliteStruct.Close()
}

func TestLoadRealBadPath(t *testing.T) {
	tfliteModelPath := basePath + "/testing_files/does_not_exist.tflite"
	loader, err := NewDefaultTFLiteModelLoader()
	test.That(t, err, test.ShouldBeNil)
	tfliteStruct, err := loader.Load(tfliteModelPath)
	test.That(t, tfliteStruct, test.ShouldBeNil)
	test.That(t, err, test.ShouldBeError, FailedToLoadError("model"))
}

func TestLoadTFLiteStruct(t *testing.T) {
	goodInterpreterLoader := func(model *tflite.Model, options *tflite.InterpreterOptions) (Interpreter, error) {
		return &fakeInterpreter{}, nil
	}

	loader := &TFLiteModelLoader{
		newModelFromFile:   modelLoader,
		newInterpreter:     goodInterpreterLoader,
		interpreterOptions: goodOptions,
		getInfo:            goodGetInfo,
	}

	tfStruct, err := loader.Load("random path")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, tfStruct, test.ShouldNotBeNil)
	test.That(t, tfStruct.model, test.ShouldNotBeNil)
	test.That(t, tfStruct.interpreter, test.ShouldNotBeNil)
	test.That(t, tfStruct.interpreterOptions, test.ShouldNotBeNil)

	tfStruct, err = loader.Load(badPath)
	test.That(t, err, test.ShouldBeError, FailedToLoadError("model"))
	test.That(t, tfStruct, test.ShouldBeNil)
}

func TestMetadataReader(t *testing.T) {
	val, err := getTFLiteMetadataBytes(badPath)
	test.That(t, err, test.ShouldBeError)
	test.That(t, val, test.ShouldBeNil)
}

func TestBadInterpreter(t *testing.T) {
	badInterpreter := func(model *tflite.Model, options *tflite.InterpreterOptions) (Interpreter, error) {
		return nil, FailedToLoadError("interpreter")
	}

	loader := &TFLiteModelLoader{
		newModelFromFile:   modelLoader,
		newInterpreter:     badInterpreter,
		interpreterOptions: goodOptions,
		getInfo:            goodGetInfo,
	}

	tfStruct, err := loader.Load("ok path")
	test.That(t, err, test.ShouldBeError, FailedToLoadError("interpreter"))
	test.That(t, tfStruct, test.ShouldBeNil)
}

func TestHasMetadata(t *testing.T) {
	tfliteModelPath := basePath + "/testing_files/model_with_metadata.tflite"
	loader, err := NewDefaultTFLiteModelLoader()
	test.That(t, err, test.ShouldBeNil)
	tfliteStruct, err := loader.Load(tfliteModelPath)
	test.That(t, tfliteStruct, test.ShouldNotBeNil)
	test.That(t, err, test.ShouldBeNil)

	meta, err := tfliteStruct.GetMetadata()
	test.That(t, err, test.ShouldBeNil)
	test.That(t, meta, test.ShouldNotBeNil)
	test.That(t, meta.SubgraphMetadata[0].OutputTensorGroups[0].TensorNames, test.ShouldResemble, []string{"location", "category", "score"})

	tfliteStruct.Close()
}

func TestNoMetadata(t *testing.T) {
	tfliteModelPath := basePath + "/testing_files/fizzbuzz_model.tflite"
	loader, err := NewDefaultTFLiteModelLoader()
	test.That(t, err, test.ShouldBeNil)
	tfliteStruct, err := loader.Load(tfliteModelPath)
	test.That(t, tfliteStruct, test.ShouldNotBeNil)
	test.That(t, err, test.ShouldBeNil)

	fizzMeta, err := tfliteStruct.GetMetadata()
	test.That(t, err, test.ShouldBeError, MetadataDoesNotExistError())
	test.That(t, fizzMeta, test.ShouldBeNil)

	tfliteStruct.Close()
}

func modelLoader(path string) *tflite.Model {
	if path == badPath {
		return nil
	}

	return &tflite.Model{}
}
