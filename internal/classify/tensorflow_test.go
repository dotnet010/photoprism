package classify

import (
	"io/ioutil"
	"testing"

	tensorflow "github.com/tensorflow/tensorflow/tensorflow/go"

	"github.com/stretchr/testify/assert"
)

var resourcesPath = "../../assets/resources"
var modelPath = resourcesPath + "/nasnet"
var examplesPath = resourcesPath + "/examples"

func TestTensorFlow_Init(t *testing.T) {
	t.Run("disabled true", func(t *testing.T) {
		tensorFlow := New(resourcesPath, true)

		result := tensorFlow.Init()
		assert.Nil(t, result)
	})
	t.Run("disabled false", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		result := tensorFlow.Init()
		assert.Nil(t, result)
	})
}

func TestTensorFlow_LabelsFromFile(t *testing.T) {
	t.Run("chameleon_lime.jpg", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		result, err := tensorFlow.File(examplesPath + "/chameleon_lime.jpg")

		assert.Nil(t, err)

		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.NotNil(t, result)
		assert.IsType(t, Labels{}, result)
		assert.Equal(t, 1, len(result))

		t.Log(result)

		assert.Equal(t, "chameleon", result[0].Name)

		assert.Equal(t, 7, result[0].Uncertainty)
	})
	t.Run("not existing file", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		result, err := tensorFlow.File(examplesPath + "/notexisting.jpg")
		assert.Contains(t, err.Error(), "no such file or directory")
		assert.Empty(t, result)
	})
	t.Run("disabled true", func(t *testing.T) {
		tensorFlow := New(resourcesPath, true)

		result, err := tensorFlow.File(examplesPath + "/chameleon_lime.jpg")
		assert.Nil(t, err)

		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, result)
		assert.IsType(t, Labels{}, result)
		assert.Equal(t, 0, len(result))

		t.Log(result)
	})
}

func TestTensorFlow_Labels(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Run("chameleon_lime.jpg", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		if imageBuffer, err := ioutil.ReadFile(examplesPath + "/chameleon_lime.jpg"); err != nil {
			t.Error(err)
		} else {
			result, err := tensorFlow.Labels(imageBuffer)

			t.Log(result)

			assert.NotNil(t, result)

			assert.Nil(t, err)
			assert.IsType(t, Labels{}, result)
			assert.Equal(t, 1, len(result))

			assert.Equal(t, "chameleon", result[0].Name)

			assert.Equal(t, 100-93, result[0].Uncertainty)
		}
	})
	t.Run("dog_orange.jpg", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		if imageBuffer, err := ioutil.ReadFile(examplesPath + "/dog_orange.jpg"); err != nil {
			t.Error(err)
		} else {
			result, err := tensorFlow.Labels(imageBuffer)

			t.Log(result)

			assert.NotNil(t, result)

			assert.Nil(t, err)
			assert.IsType(t, Labels{}, result)
			assert.Equal(t, 1, len(result))

			assert.Equal(t, "dog", result[0].Name)

			assert.Equal(t, 34, result[0].Uncertainty)
		}
	})
	t.Run("Random.docx", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		if imageBuffer, err := ioutil.ReadFile(examplesPath + "/Random.docx"); err != nil {
			t.Error(err)
		} else {
			result, err := tensorFlow.Labels(imageBuffer)
			assert.Empty(t, result)
			assert.Contains(t, err.Error(), "invalid image")
		}
	})
	t.Run("6720px_white.jpg", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		if imageBuffer, err := ioutil.ReadFile(examplesPath + "/6720px_white.jpg"); err != nil {
			t.Error(err)
		} else {
			result, err := tensorFlow.Labels(imageBuffer)
			assert.Empty(t, result)
			assert.Nil(t, err)
		}
	})
	t.Run("disabled true", func(t *testing.T) {
		tensorFlow := New(resourcesPath, true)

		if imageBuffer, err := ioutil.ReadFile(examplesPath + "/dog_orange.jpg"); err != nil {
			t.Error(err)
		} else {
			result, err := tensorFlow.Labels(imageBuffer)

			t.Log(result)

			assert.Nil(t, result)

			assert.Nil(t, err)
			assert.IsType(t, Labels{}, result)
			assert.Equal(t, 0, len(result))
		}
	})
}

func TestTensorFlow_LoadModel(t *testing.T) {
	t.Run("model path exists", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		result := tensorFlow.loadModel()
		assert.Nil(t, result)
	})
	t.Run("model path does not exist", func(t *testing.T) {
		tensorFlow := New(resourcesPath+"foo", false)

		err := tensorFlow.loadModel()

		if err == nil {
			t.FailNow()
		}

		assert.Contains(t, err.Error(), "Could not find SavedModel")
	})
}

func TestTensorFlow_BestLabels(t *testing.T) {
	t.Run("labels not loaded", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		p := make([]float32, 1000)

		p[666] = 0.5

		result := tensorFlow.bestLabels(p)
		assert.Empty(t, result)
	})
	t.Run("labels loaded", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)
		tensorFlow.loadLabels(modelPath)

		p := make([]float32, 1000)

		p[8] = 0.7
		p[1] = 0.5

		result := tensorFlow.bestLabels(p)
		assert.Equal(t, "chicken", result[0].Name)
		assert.Equal(t, "bird", result[0].Categories[0])
		assert.Equal(t, "image", result[0].Source)
		t.Log(result)
	})
}

func TestTensorFlow_MakeTensor(t *testing.T) {
	t.Run("cat_brown.jpg", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		imageBuffer, err := ioutil.ReadFile(examplesPath + "/cat_brown.jpg")
		assert.Nil(t, err)
		result, err := tensorFlow.makeTensor(imageBuffer, "jpeg")
		assert.Equal(t, tensorflow.DataType(0x1), result.DataType())
		assert.Equal(t, int64(1), result.Shape()[0])
		assert.Equal(t, int64(224), result.Shape()[2])
	})
	t.Run("Random.docx", func(t *testing.T) {
		tensorFlow := New(resourcesPath, false)

		imageBuffer, err := ioutil.ReadFile(examplesPath + "/Random.docx")
		assert.Nil(t, err)
		result, err := tensorFlow.makeTensor(imageBuffer, "jpeg")
		assert.Empty(t, result)
		assert.Equal(t, "image: unknown format", err.Error())
	})
}

func Test_ConvertTF(t *testing.T) {
	result := convertTF(uint32(98765432))
	assert.Equal(t, float32(3024.898), result)
}
