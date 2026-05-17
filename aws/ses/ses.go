package ses

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

var ErrDaemonAddrNotFound = errors.New("failed to lookup daemon address")
var ErrRegionNotFound = errors.New("failed to lookup aws region")

func SendEmail(input *sesv2.SendEmailInput) error {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return ErrRegionNotFound
	}

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return err
	}

	client := sesv2.NewFromConfig(cfg)

	_, err = client.SendEmail(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func PrepareAccountVerificationEmail(token string, toAddress string) (*sesv2.SendEmailInput, error) {
	daemon, ok := os.LookupEnv("DAEMON_ADDRESS")
	if !ok {
		return nil, ErrDaemonAddrNotFound
	}

	url := "https://rhythmaning.app/verify-account?token=" + token

	plainText := "Please verify your account to get started.\n\n" +
		"Paste the link below into a new tab to verify your account's email and finish activating your rhythmaning.app account.\n\n" +
		"This link is only valid for 10 minutes" +
		url + "\n\n"

	html := "<div><h3>Please verify your account to get started.</h3></div>" +
		"<div><p>Click the link below to verify your account's email and finish activating your rhythmaning.app account.</p></div>" +
		"<div><strong>This link is only valid for 10 minutes</strong></div>" +
		"<div><a href='" + url + "'><button>Verify Account</button></a>"

	return &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(daemon),
		Destination: &types.Destination{
			ToAddresses: []string{toAddress},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String("Verify Account"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(html),
					},
					Text: &types.Content{
						Data: aws.String(plainText),
					},
				},
			},
		},
	}, nil
}

func PreparePasswordResetEmail(token string, toAddress string) (*sesv2.SendEmailInput, error) {
	daemon, ok := os.LookupEnv("DAEMON_ADDRESS")
	if !ok {
		return nil, ErrDaemonAddrNotFound
	}

	url := "https://rhythmaning.app/password-reset?token=" + token

	plainText := "A request was made to reset your password.\n\n" +
		"If you did not request this, no further action is needed.\n\n" +
		"Paste the link below into a new tab to verify your account's email and finish activating your rhythmaning.app account.\n\n" +
		"This link is only valid for 10 minutes" +
		url + "\n\n"

	html := "<div><h3>A request was made to reset your password.</h3></div>" +
		"<div><p>If you did not request this, no further action is needed.</p></div>" +
		"<div><p>Click the link below to verify your account's email and finish activating your rhythmaning.app account.</p></div>" +
		"<div><strong>This link is only valid for 10 minutes</strong></div>" +
		"<div><a href='" + url + "'><button>Reset Password</button></a>"

	return &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(daemon),
		Destination: &types.Destination{
			ToAddresses: []string{toAddress},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String("Reset Password"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(html),
					},
					Text: &types.Content{
						Data: aws.String(plainText),
					},
				},
			},
		},
	}, nil

}
