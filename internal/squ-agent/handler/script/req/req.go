package req

type Script struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TaskID  uint   `json:"task_id"` // APIServer 分配的任务ID（ScriptResult.ID）
}
