package main

import "errors"

func performAction(c Config) (string, error) {
	switch *c.Action {
	case "view":
		resp, err := getSecret(*c.ApiURL, *c.SecretId)
		if err != nil {
			return "", err
		}
		return resp.Data, nil
	case "create":
		resp, err := createSecret(*c.ApiURL, *c.PlainText)
		if err != nil {
			return "", err
		}
		return resp.Id, nil
	default:
		return "", errors.New("unknown action")
	}
}
