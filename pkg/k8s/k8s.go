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
		// 不应该panic，因为k8s挂了，不能影程序运行
		zap.S().Error("new clientSet err: ", err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		// 不应该panic，因为k8s挂了，不能影程序运行
		zap.S().Error("new clientSet err: ", err)
	}
	return &Client{
		ClientSet:     clientSet,
		DynamicClient: dynamicClient,
	}
}

func (c *Client) CreateByYaml(yamlString string) (err error) {

	// 3. 处理多文档YAML
	docs := strings.Split(string(yamlString), "---")
	for _, doc := range docs {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		// 4. 解析为Unstructured对象
		obj := &unstructured.Unstructured{}
		decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		_, gvk, err := decoder.Decode([]byte(doc), nil, obj)
		if err != nil {
			panic(err.Error())
		}

		// 5. 准备RESTMapper以发现API资源
		mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(c.ClientSet))

		// 6. 获取GVR
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return fmt.Errorf("获取GVR: %v", err)
		}

		// 8. 确定资源接口（命名空间作用域或集群作用域）
		var dr dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			namespace := obj.GetNamespace()
			if namespace == "" {
				namespace = "default" // 默认命名空间
			}
			dr = c.DynamicClient.Resource(mapping.Resource).Namespace(namespace)
		} else {
			dr = c.DynamicClient.Resource(mapping.Resource)
		}

		// 9. 创建资源
		_, err = dr.Apply(context.TODO(), obj.GetName(), obj, metav1.ApplyOptions{
			FieldManager: obj.GetName(),
		})
		if err != nil {
			zap.S().Error("错误资源： ", doc, err)
			return fmt.Errorf("创建资源失败: %v", err)
		}
		zap.S().Error("创建成功: %s/%s", obj.GetKind(), obj.GetName())
		// fmt.Printf("创建成功: %s/%s\n", obj.GetKind(), obj.GetName())
	}
	return nil
}
