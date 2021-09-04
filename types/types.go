package types

type CreateSecretPayload struct {
	PlainText string `json:"plain_text"`
}

type GetSecretResponse struct {
	Data string `json:"data"`
}

type CreateSecretResponse struct {
	Id string `json:"id"`
}

type SecretData struct {
	Id     string
	Secret string
}
