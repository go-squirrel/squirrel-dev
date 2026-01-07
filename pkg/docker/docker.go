package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/moby/moby/client"
	"go.uber.org/zap"
)

type Docker struct {
	Client  *client.Client
	AuthStr string `json:"authStr"`
}

type PullResponse struct {
	Status      string         `json:"status"`
	ErrorDetail map[string]any `json:"errorDetail"`
}

func New(regUser, regPassword, regUrl string) *Docker {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		zap.S().Error("Failed connect docker")
		return nil
	}
	authConfig := client.RegistryLoginOptions{
		Username:      regUser,     // 替换为你的用户名
		Password:      regPassword, // 替换为你的密码
		ServerAddress: regUrl,      // 替换为你的镜像仓库地址，包括端口号，例如 "registry.example.com:80" 或 "registry.example.com:443"
	}
	_, err = cli.RegistryLogin(context.TODO(), authConfig)
	if err != nil {
		zap.S().Error("Failed loggin docker", err)
		return nil
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		zap.S().Errorf("Error encoding auth config: %v", err)
		return nil
	}
	return &Docker{
		Client:  cli,
		AuthStr: base64.URLEncoding.EncodeToString(encodedJSON),
	}
}

// Pull pulls an image from a specified repository using docker command.

func (d *Docker) Pull(imageName string) (err error) {
	out, err := d.Client.ImagePull(context.Background(), imageName, client.ImagePullOptions{
		RegistryAuth: d.AuthStr,
	})
	if err != nil {
		return err
	}
	defer out.Close()
	decoder := json.NewDecoder(out)
	for {
		var response PullResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				// 读取结束
				break
			}
			zap.S().Errorf("读取输出时出错: %v", err)
			return err
		}

		// 打印拉取镜像的名称（status 字段包含了镜像名称）
	}

	return err
}

// Tag tags an existing image with a new tag using docker command.
func (d *Docker) Tag(oldImage string, newImage string) {
	opt := client.ImageTagOptions{

		Source: oldImage,

		Target: newImage,
	}
	if _, err := d.Client.ImageTag(context.Background(), opt); err != nil {
		log.Fatalf("err: %v", err)
	}
}

// Push pushes a tagged image to a specified repository using docker command.
func (d *Docker) Push(imageName string) (err error) {
	out, err := d.Client.ImagePush(context.Background(), imageName, client.ImagePushOptions{
		RegistryAuth: d.AuthStr,
	})
	if err != nil {
		log.Fatalf("Failed push image %v,err: %v", imageName, err)
		return err
	}
	defer out.Close()
	decoder := json.NewDecoder(out)
	for {
		var response PullResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				// 读取结束
				break
			}
			log.Printf("读取输出时出错: %v", err)
			return err
		}

		// 打印拉取镜像的名称（status 字段包含了镜像名称）
		if response.ErrorDetail != nil {
			log.Printf("push输出时出错: %v", response.ErrorDetail)
			return fmt.Errorf("push输出时出错: %v", response.ErrorDetail)
		}
	}
	return err
}

func (d *Docker) SaveImagesToTar(ctx context.Context, images []string, outputFile string) error {
	// 打开目标 tar 文件进行写入
	tarFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("无法创建文件 %s: %v", outputFile, err)
	}
	defer tarFile.Close()

	// 使用 docker save API 保存镜像
	reader, err := d.Client.ImageSave(ctx, images)
	if err != nil {
		return fmt.Errorf("无法保存镜像: %v", err)
	}
	defer reader.Close()

	buffer := make([]byte, 32*1024) // 32KB 缓冲区，可以根据需要调整大小
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("读取镜像数据流时出错: %v", err)
		}
		if n == 0 {
			// 到达流的结尾
			break
		}

		// 将读取到的数据块写入到 tar 文件
		if _, err := tarFile.Write(buffer[:n]); err != nil {
			return fmt.Errorf("写入 tar 文件时出错: %v", err)
		}
	}

	return nil
}
