package execute

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

func Command(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(out) == 0 {
		return "", nil
	}
	s := string(out)
	return strings.ReplaceAll(s, "\n", ""), nil
}

// 引入template模版，渲染字段
// command: //go:embed 引入的文本字符串
// data: 需要被渲染的数值
func TemplateCommandBash(command string, data any) (string, error) {
	// 解析模板内容
	tmpl, err := template.New("example").Parse(command)
	if err != nil {
		return "", err
	}

	// 渲染模板到缓冲区
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", err
	}
	cmd := exec.Command("bash", "-c", rendered.String())
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(out) == 0 {
		return "", nil

	}
	s := string(out)
	return strings.ReplaceAll(s, "\n", ""), nil
}

func CommandBash(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(out) == 0 {
		return "", nil
	}
	s := string(out)
	return strings.ReplaceAll(s, "\n", ""), nil
}

// 执行，且反馈错误详情
func CommandError(command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)
	// 创建一个 bytes.Buffer 来捕获标准输出和标准错误
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	// 运行命令并捕获错误
	err := cmd.Run()
	// 返回标准输出、标准错误和错误
	return outBuf.String(), errBuf.String(), err
}

// 执行，且反馈错误详情
func CommandToLog(toPath string, command string, args ...string) error {

	cmd := exec.Command(command, args...)

	file, err := os.OpenFile(toPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()
	// 创建一个 bytes.Buffer 来捕获标准输出和标准错误
	cmd.Stdout = file
	cmd.Stderr = file
	// 运行命令并捕获错误
	err = cmd.Run()
	// 返回标准输出、标准错误和错误
	return err
}

// 不能使用> 这个操作打印日志
func Shell(command string) (string, error) {
	hasPipe := strings.Contains(command, "|")
	if hasPipe {
		parts := strings.Split(command, "|")
		var commands []*exec.Cmd
		var lastOut bytes.Buffer

		// 创建管道中的每个子命令
		for _, part := range parts {
			args := strings.Fields(part)
			cmd := exec.Command(args[0], args[1:]...)
			if len(commands) > 0 {
				cmd.Stdin = &lastOut
			}
			commands = append(commands, cmd)
			lastOut.Reset()
			cmd.Stdout = &lastOut
			if err := cmd.Start(); err != nil {
				return "", err
			}
		}

		// 等待所有命令完成
		for _, cmd := range commands {
			if err := cmd.Wait(); err != nil {
				return "", err
			}
		}

		// 返回最后一个命令的输出作为结果
		return lastOut.String(), nil
	}

	// 如果没有管道，直接执行命令
	commandSlice := strings.Fields(command)
	cmd := exec.Command(commandSlice[0], commandSlice[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func PkillByPID(processName string) error {
	cmd := exec.Command("pgrep", "-f", processName)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	processIDs := strings.Split(strings.TrimSpace(string(output)), "\n")

	// 遍历每个进程 ID 并杀死相应的进程
	for _, pidStr := range processIDs {
		pid := strings.TrimSpace(pidStr)
		killCmd := exec.Command("kill", pid)
		if err := killCmd.Run(); err != nil {
			zap.S().Errorf("Failed to kill process with PID %s: %v\n", pid, err)
		}
	}

	return nil
}
