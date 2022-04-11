package core

type ChannelTask struct {
	APIID int `json:"api_id"`
	APIHash string `json:"api_hash"`
	Title string `json:"title"`
	Body string `json:"body"`
	Url	string `json:"url"`
}

type ChannelTaskResponse struct {
	Tasks []ChannelTask `json:"tasks"`
}

type UserTask struct {
	TGID int `json:"tg_id"`
	Title string `json:"title"`
	Body string `json:"body"`
	Url	string `json:"url"`
}

type UserTaskResponse struct {
	Tasks []UserTask `json:"tasks"`
}

type Setting struct {
	ID int `json:"id"`
	IsSafeDeal bool `json:"is_safe_deal"`
	IsBudget bool `json:"is_budget"`
	IsTerm bool `json:"is_term"`
}

type Channel struct {
	ID int `json:"id"`
	APIID int `json:"api_id"`
	APIHash string `json:"api_hash"`
	Name string `json:"name"`
}

type User struct {
	ID int `json:"id"`
	TGID int `json:"tg_id"`
	Username string `json:"username"`
}

type SettingResponse struct {
	IsSafeDeal bool `json:"is_safe_deal"`
	IsBudget bool `json:"is_budget"`
	IsTerm bool `json:"is_term"`
	Categories []int `json:"categories"`
}

type ChannelResponse struct {
	ID int `json:"id"`
	APIID int `json:"api_id"`
	APIHash string `json:"api_hash"`
	Name string `json:"name"`
	Setting SettingResponse `json:"setting"`
}

type UserResponse struct {
	ID int `json:"id"`
	TGID int `json:"tg_id"`
	Username string `json:"username"`
	Setting SettingResponse `json:"setting"`
}

type SettingInput struct {
	IsSafeDeal *bool `json:"is_safe_deal" binding:"required"`
	IsBudget *bool `json:"is_budget" binding:"required"`
	IsTerm *bool `json:"is_term" binding:"required"`
	Categories []int `json:"categories" binding:"required"`
}

type ChannelInput struct {
	APIID int `json:"api_id" binding:"required"`
	APIHash string `json:"api_hash" binding:"required"`
	Name string `json:"name" binding:"required"`
	Setting SettingInput `json:"setting" binding:"required"`
}

type UserInput struct {
	TGID int `json:"tg_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Setting SettingInput `json:"setting" binding:"required"`
}

type ChannelAPIIDInput struct {
	APIID int `json:"api_id" binding:"required"`
}

type UserAPIIDInput struct {
	TGID int `json:"tg_id" binding:"required"`
}
