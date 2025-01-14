package scratchgonnect

// Structs

type Comment struct {
	Id           int    `json:"id"`
	ParentId     int    `json:"parent_id"`
	ComenteeId   int    `json:"commentee_id"`
	Content      string `json:"content"`
	DateCreated  string `json:"datetime_created"`
	DateModified string `json:"datetime_modified"`
	Visiblity    string `json:"visibility"`
	ReplyCount   int    `json:"reply_count"`
	Author       User   `json:"author"`
}

// Struct functions
