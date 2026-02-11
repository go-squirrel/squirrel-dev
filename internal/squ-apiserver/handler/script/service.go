package script

import (
	"strings"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/agent"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/script/req"
	"squirrel-dev/internal/squ-apiserver/handler/script/res"

	scriptRepository "squirrel-dev/internal/squ-apiserver/repository/script"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"

	"go.uber.org/zap"
)

type Script struct {
	Config      *config.Config
	Repository  scriptRepository.ScriptRepository
	ServerRepo  serverRepository.Repository
	AgentClient *agent.Client
}

func New(config *config.Config, scriptRepo scriptRepository.ScriptRepository, serverRepo serverRepository.Repository) *Script {
	return &Script{
		Config:      config,
		Repository:  scriptRepo,
		ServerRepo:  serverRepo,
		AgentClient: agent.NewClient(config),
	}
}

func (s *Script) List() response.Response {
	var scripts []res.Script
	daoScripts, err := s.Repository.List()
	if err != nil {
		zap.L().Error("failed to list scripts",
			zap.Error(err),
		)
		return response.Error(returnScriptErrCode(err))
	}
	for _, daoS := range daoScripts {
		scripts = append(scripts, s.modelToResponse(daoS))
	}
	return response.Success(scripts)
}

func (s *Script) Get(id uint) response.Response {
	daoS, err := s.Repository.Get(id)
	if err != nil {
		zap.L().Error("failed to get script",
			zap.Uint("script_id", id),
			zap.Error(err),
		)
		return response.Error(returnScriptErrCode(err))
	}
	scriptRes := s.modelToResponse(daoS)

	return response.Success(scriptRes)
}

func (s *Script) Delete(id uint) response.Response {
	err := s.Repository.Delete(id)
	if err != nil {
		zap.L().Error("failed to delete script",
			zap.Uint("script_id", id),
			zap.Error(err),
		)
		return response.Error(returnScriptErrCode(err))
	}

	return response.Success("success")
}

func (s *Script) Add(request req.Script) response.Response {
	// Validate script name and content
	if request.Name == "" {
		zap.L().Error("script name is empty",
			zap.String("request_id", string(rune(request.ID))),
		)
		return response.Error(res.ErrInvalidScriptContent)
	}

	if request.Content == "" {
		zap.L().Error("script content is empty",
			zap.String("name", request.Name),
		)
		return response.Error(res.ErrInvalidScriptContent)
	}

	// Validate that script content starts with shebang
	if !strings.HasPrefix(request.Content, "#!") {
		zap.L().Error("script must start with shebang (#!)",
			zap.String("name", request.Name),
		)
		return response.Error(res.ErrInvalidScriptContent)
	}

	// Clean up script content (trim leading/trailing whitespace)
	request.Content = strings.TrimSpace(request.Content)

	modelReq := s.requestToModel(request)

	err := s.Repository.Add(&modelReq)
	if err != nil {
		zap.L().Error("failed to add script",
			zap.String("name", request.Name),
			zap.Error(err),
		)
		return response.Error(returnScriptErrCode(err))
	}

	return response.Success("success")
}

func (s *Script) Update(request req.Script) response.Response {
	// Validate script name and content
	if request.Name == "" {
		zap.L().Error("script name is empty",
			zap.Uint("script_id", request.ID),
		)
		return response.Error(res.ErrInvalidScriptContent)
	}

	if request.Content == "" {
		zap.L().Error("script content is empty",
			zap.Uint("script_id", request.ID),
			zap.String("name", request.Name),
		)
		return response.Error(res.ErrInvalidScriptContent)
	}

	// Validate that script content starts with shebang
	if !strings.HasPrefix(request.Content, "#!") {
		zap.L().Error("script must start with shebang (#!)",
			zap.Uint("script_id", request.ID),
			zap.String("name", request.Name),
		)
		return response.Error(res.ErrInvalidScriptContent)
	}

	// Clean up script content (trim leading/trailing whitespace)
	request.Content = strings.TrimSpace(request.Content)

	modelReq := s.requestToModel(request)
	modelReq.ID = request.ID
	err := s.Repository.Update(&modelReq)

	if err != nil {
		zap.L().Error("failed to update script",
			zap.Uint("script_id", request.ID),
			zap.String("name", request.Name),
			zap.Error(err),
		)
		return response.Error(returnScriptErrCode(err))
	}

	return response.Success("success")
}
