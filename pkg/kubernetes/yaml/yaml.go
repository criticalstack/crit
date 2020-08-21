package yaml

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
)

// MarshalToYaml marshals an object into yaml.
func MarshalToYaml(obj runtime.Object, gv runtime.GroupVersioner) ([]byte, error) {
	return MarshalToYamlForCodecs(obj, gv, clientsetscheme.Codecs)
}

// MarshalToYamlForCodecs marshals an object into yaml using the specified codec
func MarshalToYamlForCodecs(obj runtime.Object, gv runtime.GroupVersioner, codecs serializer.CodecFactory) ([]byte, error) {
	const mediaType = runtime.ContentTypeYAML
	info, ok := runtime.SerializerInfoForMediaType(codecs.SupportedMediaTypes(), mediaType)
	if !ok {
		return []byte{}, errors.Errorf("unsupported media type %q", mediaType)
	}

	encoder := codecs.EncoderForVersion(info.Serializer, gv)
	return runtime.Encode(encoder, obj)
}

// UnmarshalFromYaml unmarshals yaml into an object.
func UnmarshalFromYaml(buffer []byte, gv runtime.GroupVersioner) (runtime.Object, error) {
	return UnmarshalFromYamlForCodecs(buffer, gv, clientsetscheme.Codecs)
}

// UnmarshalFromYamlForCodecs unmarshals yaml into an object using the specified codec
func UnmarshalFromYamlForCodecs(buffer []byte, gv runtime.GroupVersioner, codecs serializer.CodecFactory) (runtime.Object, error) {
	const mediaType = runtime.ContentTypeYAML
	info, ok := runtime.SerializerInfoForMediaType(codecs.SupportedMediaTypes(), mediaType)
	if !ok {
		return nil, errors.Errorf("unsupported media type %q", mediaType)
	}

	decoder := codecs.DecoderToVersion(info.Serializer, gv)
	return runtime.Decode(decoder, buffer)
}

func UnmarshalFromYamlUnstructured(data []byte) ([]*unstructured.Unstructured, error) {
	objs := make([]*unstructured.Unstructured, 0)
	reader := utilyaml.NewYAMLReader(bufio.NewReader(bytes.NewBuffer(data)))
	for {
		b, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(b) == 0 {
			break
		}
		var m map[string]interface{}
		if err := yaml.Unmarshal(b, &m); err != nil {
			return nil, err
		}
		if len(m) == 0 {
			continue
		}
		obj := &unstructured.Unstructured{
			Object: m,
		}
		objs = append(objs, obj)
	}
	return objs, nil
}
