package response

type ShowPublicMessageResponse struct {
	Response string   `json:"response"`
	Messages []string `json:"messages"`
}
