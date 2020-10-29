package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	bonusly "github.com/kimchelly/go-bonusly"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func main() {
	if err := app().Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func app() *cli.App {
	app := cli.NewApp()
	app.Name = "bonusly"
	app.Usage = "Bonusly CLI"

	app.Commands = []*cli.Command{
		bonus(),
		userInfo(),
	}

	return app
}

const (
	idFlagName     = "id"
	reasonFlagName = "reason"
)

func bonus() *cli.Command {
	return &cli.Command{
		Name: "bonus",
		Subcommands: []*cli.Command{
			createBonus(),
			getBonus(),
			updateBonus(),
			deleteBonus(),
		},
	}
}

func createBonus() *cli.Command {
	const (
		parentIDFlagName = "parent_id"
	)

	return &cli.Command{
		Name:  "create",
		Usage: "create a new bonus",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     reasonFlagName,
				Usage:    "the reason for the bonus",
				Required: true,
			},
			&cli.StringFlag{
				Name:  parentIDFlagName,
				Usage: "the ID of the parent bonus",
			},
		},
		Action: func(c *cli.Context) error {
			return withClient(func(ctx context.Context, client bonusly.Client) error {
				req := bonusly.CreateBonusRequest{
					Reason:        c.String(reasonFlagName),
					ParentBonusID: c.String(parentIDFlagName),
				}
				resp, err := client.CreateBonus(ctx, req)
				if err != nil {
					return err
				}
				output, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				_, err = fmt.Fprintf(os.Stdout, string(output))
				return err
			})
		},
	}
}

func getBonus() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "get an existing bonus",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     idFlagName,
				Usage:    "the bonus ID",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			return withClient(func(ctx context.Context, client bonusly.Client) error {
				resp, err := client.GetBonus(ctx, c.String(idFlagName))
				if err != nil {
					return err
				}
				output, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				_, err = fmt.Fprintf(os.Stdout, string(output))
				return err
			})
		},
	}
}

func updateBonus() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "update an existing bonus",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     idFlagName,
				Usage:    "the bonus ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     reasonFlagName,
				Usage:    "the new reason message for the bonus",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			return withClient(func(ctx context.Context, client bonusly.Client) error {
				resp, err := client.UpdateBonus(ctx, c.String(idFlagName), c.String(reasonFlagName))
				if err != nil {
					return err
				}
				output, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				_, err = fmt.Fprintf(os.Stdout, string(output))
				return err
			})
		},
	}
}

func deleteBonus() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "delete an existing bonus",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     idFlagName,
				Usage:    "the bonus ID",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			return withClient(func(ctx context.Context, client bonusly.Client) error {
				if err := client.DeleteBonus(ctx, c.String(idFlagName)); err != nil {
					return err
				}
				_, err := fmt.Fprintf(os.Stdout, "Successfully deleted bonus.")
				return err
			})
		},
	}
}

func userInfo() *cli.Command {
	return &cli.Command{
		Name: "user",
		Subcommands: []*cli.Command{
			myUserInfo(),
		},
	}
}

func myUserInfo() *cli.Command {
	return &cli.Command{
		Name: "me",
		Action: func(c *cli.Context) error {
			return withClient(func(ctx context.Context, client bonusly.Client) error {
				info, err := client.MyUserInfo(ctx)
				if err != nil {
					return err
				}
				output, err := json.Marshal(info)
				if err != nil {
					return err
				}
				_, err = fmt.Fprintf(os.Stdout, string(output))
				return err
			})
		},
	}
}

func withClient(clientOp func(ctx context.Context, client bonusly.Client) error) error {
	token, err := getBonuslyToken()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, err := bonusly.NewClient(bonusly.ClientOptions{
		AccessToken: token,
	})
	if err != nil {
		return err
	}
	defer c.Close(ctx)

	return clientOp(ctx, c)
}

func getBonuslyToken() (string, error) {
	token := os.Getenv("BONUSLY_TOKEN")
	if token == "" {
		return "", errors.Errorf("BONUSLY_TOKEN environment variable must be set to your Bonusly access token")
	}
	return token, nil
}
