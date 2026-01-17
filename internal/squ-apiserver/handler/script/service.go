package script

import (
	"strings"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/script/req"
	"squirrel-dev/internal/squ-apiserver/handler/script/res"
	"squirrel-dev/internal/squ-apiserver/model"

	scriptRepository "squirrel-dev/internal/squ-apiserver/repository/script"

	"go.uber.org/zap"
)

type Script struct {
	Config     *config.Config
	Repository scriptRepository.ScriptRepository
}

func (s *Script) List() response.Response {
	var scripts []res.Script
	daoScripts, err := s.Repository.List()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	for _, daoS := range daoScripts {
		scripts = append(scripts, res.Script{
			ID:      daoS.ID,
			Name:    daoS.Name,
			Content: daoS.Content,
		})
	}
	return response.Success(scripts)
}

func (s *Script) Get(id uint) response.Response {
	var scriptRes res.Script
	daoS, err := s.Repository.Get(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	scriptRes = res.Script{
		ID:      daoS.ID,
		Name:    daoS.Name,
		Content: daoS.Content,
	}

	return response.Success(scriptRes)
}

func (s *Script) Delete(id uint) response.Response {
	err := s.Repository.Delete(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Script) Add(request req.Script) response.Response {
	// 验证脚本名称和内容
	if request.Name == "" {
		zap.S().Error("script name is empty")
		return response.Error(res.ErrInvalidScriptContent)
	}

	if request.Content == "" {
		zap.S().Error("script content is empty")
		return response.Error(res.ErrInvalidScriptContent)
	}

	// 验证脚本内容是否以 shebang 开头
	if !strings.HasPrefix(request.Content, "#!") {
		zap.S().Error("script must start with shebang (#!)")
		return response.Error(res.ErrInvalidScriptContent)
	}

	// 清理脚本内容（去除首尾空白）
	request.Content = strings.TrimSpace(request.Content)

	modelReq := model.Script{
		Name:    request.Name,
		Content: request.Content,
	}

	err := s.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (s *Script) Update(request req.Script) response.Response {
	// 验证脚本名称和内容
	if request.Name == "" {
		zap.S().Error("script name is empty")
		return response.Error(res.ErrInvalidScriptContent)
	}

	if request.Content == "" {
		zap.S().Error("script content is empty")
		return response.Error(res.ErrInvalidScriptContent)
	}

	// 验证脚本内容是否以 shebang 开头
	if !strings.HasPrefix(request.Content, "#!") {
		zap.S().Error("script must start with shebang (#!)")
		return response.Error(res.ErrInvalidScriptContent)
	}

	// 清理脚本内容（去除首尾空白）
	request.Content = strings.TrimSpace(request.Content)

	modelReq := model.Script{
		Name:    request.Name,
		Content: request.Content,
	}
	modelReq.ID = request.ID
	err := s.Repository.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}
