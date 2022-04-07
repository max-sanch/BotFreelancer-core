package core

type Task struct {
	ID int `json:"api_id" binding:"required"`
	APIHash string `json:"api_hash" binding:"required"`
	Title string `json:"title" binding:"required"`
	Body string `json:"body" binding:"required"`
	Url	string `json:"url" binding:"required"`
}

type ChannelRequest struct {
	Tasks []Task `json:"tasks" binding:"required"`
}

type Channel struct {
	ID int `json:"id" binding:"required"`
	APIID int `json:"api_id" binding:"required"`
	APIHash string `json:"api_hash" binding:"required"`
	Name string `json:"name" binding:"required"`
}
