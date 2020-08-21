package dynamic

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"

	yamlutil "github.com/criticalstack/crit/pkg/kubernetes/yaml"
	"github.com/criticalstack/crit/pkg/log"
)

type client struct {
	dynClient  dynamic.Interface
	clientset  *kubernetes.Clientset
	restMapper meta.RESTMapper
}

func newClient(config *rest.Config) (*client, error) {
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		return nil, err
	}
	restMapper := restmapper.NewDiscoveryRESTMapper(groupResources)
	return &client{
		dynClient:  dynClient,
		clientset:  clientset,
		restMapper: restMapper,
	}, nil
}

func (c *client) getMapping(obj *unstructured.Unstructured) (*meta.RESTMapping, error) {
	gvk := obj.GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	return c.restMapper.RESTMapping(gk, gvk.Version)
}

func (c *client) getDynamicResource(ns string, scope meta.RESTScope, resource schema.GroupVersionResource) dynamic.ResourceInterface {
	if scope.Name() == meta.RESTScopeNameRoot {
		return c.dynClient.Resource(resource)
	}
	if ns == "" {
		ns = metav1.NamespaceSystem
	}
	return c.dynClient.Resource(resource).Namespace(ns)
}

func (c *client) Create(ctx context.Context, v interface{}) (*unstructured.Unstructured, error) {
	obj, err := convertUnstructuredObject(v)
	if err != nil {
		return nil, err
	}
	mapping, err := c.getMapping(obj)
	if err != nil {
		return nil, err
	}
	return c.getDynamicResource(obj.GetNamespace(), mapping.Scope, mapping.Resource).Create(ctx, obj, metav1.CreateOptions{})
}

func (c *client) Update(ctx context.Context, v interface{}) (*unstructured.Unstructured, error) {
	obj, err := convertUnstructuredObject(v)
	if err != nil {
		return nil, err
	}
	mapping, err := c.getMapping(obj)
	if err != nil {
		return nil, err
	}
	return c.getDynamicResource(obj.GetNamespace(), mapping.Scope, mapping.Resource).Update(ctx, obj, metav1.UpdateOptions{})
}

func convertUnstructuredObject(v interface{}) (*unstructured.Unstructured, error) {
	switch t := v.(type) {
	case runtime.Object:
		obj := &unstructured.Unstructured{}
		if err := clientsetscheme.Scheme.Convert(t, obj, nil); err != nil {
			return nil, err
		}
		return obj, nil
	case unstructured.Unstructured:
		return &t, nil
	default:
		return nil, errors.Errorf("invalid type for dynamic client: %T", v)
	}
}

// TODO(chrism): add strategic merge?
func Apply(ctx context.Context, config *rest.Config, data []byte) error {
	client, err := newClient(config)
	if err != nil {
		return err
	}
	objs, err := yamlutil.UnmarshalFromYamlUnstructured(data)
	if err != nil {
		return err
	}
	for _, obj := range objs {
		if _, err := client.Create(ctx, obj); err != nil {
			if apierrors.IsAlreadyExists(err) || apierrors.IsInvalid(err) {
				log.Debug("kube.Apply", zap.Error(err))
				continue
			}
			return err
		}
	}
	return nil
}
