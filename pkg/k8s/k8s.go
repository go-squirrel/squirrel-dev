package k8s

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type Client struct {
	ClientSet     *kubernetes.Clientset
	DynamicClient *dynamic.DynamicClient
}

func New(config *rest.Config) *Client {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		zap.S().Error("new clientSet err: ", err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		zap.S().Error("new clientSet err: ", err)
	}
	return &Client{
		ClientSet:     clientSet,
		DynamicClient: dynamicClient,
	}
}

func (c *Client) CreateByYaml(yamlString string) (err error) {

	docs := strings.Split(string(yamlString), "---")
	for _, doc := range docs {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		obj := &unstructured.Unstructured{}
		decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		_, gvk, err := decoder.Decode([]byte(doc), nil, obj)
		if err != nil {
			panic(err.Error())
		}

		mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(c.ClientSet))

		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return fmt.Errorf("get GVR: %v", err)
		}

		var dr dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			namespace := obj.GetNamespace()
			if namespace == "" {
				namespace = "default"
			}
			dr = c.DynamicClient.Resource(mapping.Resource).Namespace(namespace)
		} else {
			dr = c.DynamicClient.Resource(mapping.Resource)
		}

		_, err = dr.Apply(context.TODO(), obj.GetName(), obj, metav1.ApplyOptions{
			FieldManager: obj.GetName(),
		})
		if err != nil {
			zap.S().Error("create err: ", doc, err)
			return fmt.Errorf("create err: %v", err)
		}
		zap.S().Error("create success: %s/%s", obj.GetKind(), obj.GetName())
	}
	return nil
}
