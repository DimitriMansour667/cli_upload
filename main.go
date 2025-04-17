package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name: "cmdim",
        HelpName: "cmdim",
        Usage: "CMDIM is a CLI tool to upload files to a PocketBase Instance and get a link to the file.",
        Version: "0.0.1",
        Authors: []*cli.Author{
            {
                Name:  "Dimitri Mansour",
                Email: "dimitri@dimitrimansour.com",
            },
        },
        EnableBashCompletion: true,
        Commands: []*cli.Command{
            {
                Name:  "config",
                Usage: "Configure the url of PocketBase Instance",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name:  "set",
                        Aliases: []string{"s"},
                        Usage: "Set the url of the PocketBase Instance",
                    },
                    &cli.BoolFlag{
                        Name:  "get",
                        Aliases: []string{"g"},
                        Usage: "Get the url of the PocketBase Instance",
                    },
                    &cli.BoolFlag{
                        Name:  "path",
                        Aliases: []string{"p"},
                        Usage: "Get the path of the PocketBase Instance Config",
                    },
                },
                Action: func(ctx *cli.Context) error {
                    if ctx.String("set") != "" {
                        saveConfig(ctx.String("set"))
                        fmt.Println("PocketBase Instance URL set to:", ctx.String("set"))
                    }
                    if ctx.Bool("get") {
                        url, err := loadConfig()
                        if err != nil {
                            fmt.Println("Error loading config:", err)
                        } else {
                            fmt.Println(url)
                        }
                    }
                    if ctx.Bool("path") {
                        fmt.Println(getConfigPath())
                    }
                    return nil
                },
            },
            {
                Name:  "check",
                Usage: "Check if the PocketBase Instance is running",
                Action: func(*cli.Context) error {
                    url, err := loadConfig()
                    if err != nil {
                        log.Fatal("Error loading config:", err)
                    }
                    if checkInstance() {
                        fmt.Println("PocketBase Instance is running at:", url)
                    }
                    return nil
                },
            },
            {
                Name:  "upload",
                Usage: "Upload a file to the PocketBase Instance",
                Action: func(*cli.Context) error {
                    fmt.Println(`lol
                    `)
                    return nil
                },
            },
            {
                Name:  "list",
                Usage: "List all files in the PocketBase Instance",
                Action: func(*cli.Context) error {
                    fmt.Println(`lol
                    `)
                    return nil
                },
            },
            {
                Name:  "delete",
                Usage: "Delete a file from the PocketBase Instance",
                Action: func(*cli.Context) error {
                    fmt.Println(`lol
                    `)
                    return nil
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

// Config Functions

func getConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal("Could not find home directory:", err)
    }
    return filepath.Join(homeDir, ".cmdim")
}

func saveConfig(url string) error {
    return os.WriteFile(getConfigPath(), []byte(url), 0600)
}

func loadConfig() (string, error) {
    data, err := os.ReadFile(getConfigPath())
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(data)), nil
}

// PocketBase Functions

func checkInstance() bool{
    url, err := loadConfig()
    if err != nil {
        log.Fatal("Error loading config:", err)
        return false
    }
    if url == "" {
        log.Fatal("No PocketBase Instance URL found")
        return false
    }


    resp, err := http.Get(url)
    if err != nil {
        log.Fatal("Error checking PocketBase Instance:", err)
        return false
    }
    defer resp.Body.Close()
    if resp.StatusCode == 200 {
        return true
    }
    log.Fatal("PocketBase Instance is not running at:", url)
    return false
}

func checkConfig() bool{
    return false
}

func upload() string{
    return ""
}

func list() string{
    return ""
}

func delete(){

}