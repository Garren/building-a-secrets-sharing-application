package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func setupParseArgs(w io.Writer, args []string) (Config, error) {
	c := Config{}

	fs := flag.NewFlagSet("secret-client", flag.ContinueOnError)
	fs.SetOutput(w)

	c.ApiURL = fs.String("url", "", "API Endpoint")
	c.Action = fs.String("action", "", "action to perform (create or view a secret)")
	c.PlainText = fs.String("text", "", "secret text (create)")
	c.SecretId = fs.String("id", "", "secret id (get)")

	showHelp := fs.Bool("h", false, "show help")

	err := fs.Parse(args)
	if err != nil {
		return c, err
	}

	if *showHelp {
		fmt.Println("Usage: ...")
		fs.PrintDefaults()
		os.Exit(0)
	}

	return c, err
}

func validateConfig(c Config) []error {
	var validationErrors []error

	if len(*c.ApiURL) == 0 {
		validationErrors = append(validationErrors, errors.New("API URL not defined"))
	}

	switch *c.Action {
	case "create":
		if len(*c.PlainText) == 0 {
			validationErrors = append(validationErrors, errors.New("create action requires a text argument"))
		}
		if len(*c.SecretId) > 0 {
			validationErrors = append(validationErrors, errors.New("create action does not accept a secret id"))
		}
	case "view":
		if len(*c.SecretId) == 0 {
			validationErrors = append(validationErrors, errors.New("view action requires a secret id"))
		}
		if len(*c.PlainText) > 0 {
			validationErrors = append(validationErrors, errors.New("view action does not accept a text argument"))
		}
	default:
		if len(*c.Action) == 0 {
			validationErrors = append(validationErrors, errors.New("no action specified"))
		} else {
			validationErrors = append(validationErrors, errors.New("unknown action"))
		}
	}

	return validationErrors
}
