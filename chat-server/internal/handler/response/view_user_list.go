package response

type ViewUserListResponse struct {
	Response string   `json:"response"`
	Messages []string `json:"users"`
}
