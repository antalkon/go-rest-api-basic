package res

// UserData example
// @Description User profile data
type AuthTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
