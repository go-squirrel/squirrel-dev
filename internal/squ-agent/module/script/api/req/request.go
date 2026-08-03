package req

type ExecuteScript struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TaskID  uint   `json:"task_id"`
}
