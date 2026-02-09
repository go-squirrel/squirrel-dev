package script

import (
	"squirrel-dev/internal/squ-apiserver/handler/script/req"
	"squirrel-dev/internal/squ-apiserver/handler/script/res"
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

func (s *Script) modelToResponse(daoS model.Script) res.Script {
	return res.Script{
		ID:      daoS.ID,
		Name:    daoS.Name,
		Content: daoS.Content,
	}
}

func (s *Script) modelToRequest(daoS model.Script) req.Script {
	return req.Script{
		ID:      daoS.ID,
		Name:    daoS.Name,
		Content: daoS.Content,
	}
}

func (s *Script) requestToModel(request req.Script) model.Script {
	return model.Script{
		Name:    request.Name,
		Content: request.Content,
	}
}

func (s *Script) scriptResultToResponse(r model.ScriptResult) res.ScriptResult {
	return res.ScriptResult{
		ID:           r.ID,
		TaskID:       r.TaskID,
		ScriptID:     r.ScriptID,
		ServerID:     r.ServerID,
		ServerIP:     r.ServerIP,
		AgentPort:    r.AgentPort,
		Output:       r.Output,
		Status:       r.Status,
		ErrorMessage: r.ErrorMessage,
		CreatedAt:    r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// returnScriptErrCode 根据错误类型返回精确的脚本错误码
func returnScriptErrCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return res.ErrScriptNotFound
	case gorm.ErrDuplicatedKey:
		return res.ErrDuplicateScript
	}
	return res.ErrInvalidScriptContent
}
