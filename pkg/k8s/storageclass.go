package k8s

import (
	"context"
	"fmt"

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StorageClass 结构体定义
type StorageClass struct {
	APIVersion        string   `yaml:"apiVersion"`
	Kind              string   `yaml:"kind"`
	Metadata          Metadata `yaml:"metadata"`
	Provisioner       string   `yaml:"provisioner"`
	VolumeBindingMode string   `yaml:"volumeBindingMode,omitempty"`
}

// 创建StorageClass的函数
func (c *Client) GenerateSC(
	name string,
	provisioner string,
	volumeBindingMode string,
	labelKey, labelValue string) (sc *StorageClass, err error) {

	sc = &StorageClass{
		APIVersion: "storage.k8s.io/v1",
		Kind:       "StorageClass",
		Metadata: Metadata{
			Name: name,
			Labels: map[string]string{
				labelKey: labelValue,
			},
		},
		Provisioner:       provisioner,
		VolumeBindingMode: volumeBindingMode,
	}
	return sc, nil
}

func (c *Client) CreateStorageClass(name string) {
	// 定义StorageClass对象
	storageClass := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Provisioner: "kubernetes.io/no-provisioner",
		VolumeBindingMode: func() *storagev1.VolumeBindingMode {
			mode := storagev1.VolumeBindingWaitForFirstConsumer
			return &mode
		}(),
	}

	// 创建StorageClass
	result, err := c.ClientSet.StorageV1().StorageClasses().Create(context.TODO(), storageClass, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("创建StorageClass失败: %v\n", err)
		return
	}

	fmt.Printf("成功创建StorageClass: %s\n", result.Name)
}
