package types

type SecretPost struct {
	PlainText string `json:"plain_text"`
}

type SecretGetResponse struct {
	Data string `json:"data"`
}

type SecretPostResponse struct {
	Id string `json:"id"`
}

