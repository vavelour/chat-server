package response

type ShowPrivateMessageResponse struct {
	Response string   `json:"response"`
	Messages []string `json:"messages"`
}
