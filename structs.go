package core

// Input structs

type SettingInput struct {
	IsSafeDeal *bool `json:"is_safe_deal" binding:"required"`
	IsBudget   *bool `json:"is_budget" binding:"required"`
	IsTerm     *bool `json:"is_term" binding:"required"`
	Categories []int `json:"categories" binding:"required"`
}

type ChannelInput struct {
	ApiId   int          `json:"api_id" binding:"required"`
	ApiHash string       `json:"api_hash" binding:"required"`
	Name    string       `json:"name" binding:"required"`
	Setting SettingInput `json:"setting" binding:"required"`
}

type UserInput struct {
	TgId     int          `json:"tg_id" binding:"required"`
	Username string       `json:"username" binding:"required"`
	Setting  SettingInput `json:"setting" binding:"required"`
}

type ApiIdInput struct {
	ApiId int `json:"api_id" binding:"required"`
}

type TgIdInput struct {
	TgId int `json:"tg_id" binding:"required"`
}

type TaskDataInput struct {
	FLName          string `json:"fl_name" binding:"required"`
	FLUrl           string `json:"fl_url" binding:"required"`
	TaskUrl         string `json:"task_url" binding:"required"`
	Category        string `json:"category" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Body            string `json:"body" binding:"required"`
	Budget          int    `json:"budget"`
	IsBudgetPerHour bool   `json:"is_budget_per_hour"`
	Term            string `json:"term"`
	IsSafeDeal      bool   `json:"is_safe_deal" binding:"required"`
	DateTime        string `json:"datetime" binding:"required"`
}

type TasksInput struct {
	Tasks []TaskDataInput `json:"tasks" binding:"required"`
}

// Response structs

type SettingResponse struct {
	IsSafeDeal bool  `json:"is_safe_deal"`
	IsBudget   bool  `json:"is_budget"`
	IsTerm     bool  `json:"is_term"`
	Categories []int `json:"categories"`
}

type ChannelResponse struct {
	Id      int             `json:"id"`
	ApiId   int             `json:"api_id"`
	ApiHash string          `json:"api_hash"`
	Name    string          `json:"name"`
	Setting SettingResponse `json:"setting"`
}

type UserResponse struct {
	Id       int             `json:"id"`
	TgId     int             `json:"tg_id"`
	Username string          `json:"username"`
	Setting  SettingResponse `json:"setting"`
}

type ChannelTaskResponse struct {
	ApiId   int    `json:"api_id"`
	ApiHash string `json:"api_hash"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Url     string `json:"url"`
}

type ChannelTasksResponse struct {
	Tasks []ChannelTaskResponse `json:"tasks"`
}

type UserTaskResponse struct {
	TgId  int    `json:"tg_id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Url   string `json:"url"`
}

type UserTasksResponse struct {
	Tasks []UserTaskResponse `json:"tasks"`
}
