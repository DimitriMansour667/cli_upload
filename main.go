package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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
                Email: "dimi.mansour03@gmail.com",
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
                        fmt.Println(`
                        Note:
                        - Make sure that you provided the link to PocketBase valid record
                        - Example: https://your-pocketbase-instance.com/api/collections/{record-name}/records
                        - Make sure that your record is valid and has the fields "name", "file", "link", "created"
                        - Use "cmdim check" to check if the PocketBase Instance is running
                        `)
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
                Action: func(ctx *cli.Context) error {
                    upload(ctx.Args().First())
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


func upload(filePath string) error {
    fmt.Println("Uploading file:", filePath)
    url, err := loadConfig()
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    part, err := writer.CreateFormFile("file", filepath.Base(filePath))
    if err != nil {
        return fmt.Errorf("failed to create form file: %w", err)
    }

    _, err = io.Copy(part, file)
    if err != nil {
        return fmt.Errorf("failed to copy file content: %w", err)
    }

    _ = writer.WriteField("name", filepath.Base(filePath))

    err = writer.Close()
    if err != nil {
        return fmt.Errorf("failed to close writer: %w", err)
    }

    req, err := http.NewRequest("POST", url, body)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        respBody, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
    }

    fmt.Printf("Successfully uploaded %s\n", filepath.Base(filePath))
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }
    var respBodyMap map[string]interface{}
    err = json.Unmarshal(respBody, &respBodyMap)
    if err != nil {
        return fmt.Errorf("failed to unmarshal response body: %w", err)
    }

    // http://127.0.0.1:8090/api/files/COLLECTION_ID_OR_NAME/RECORD_ID/FILENAME
    splitUrl := strings.Split(url, "/")
    link := splitUrl[0] + "//" + splitUrl[2] + "/api/files/" + respBodyMap["collectionName"].(string) + "/" + respBodyMap["id"].(string) + "/" + respBodyMap["file"].(string)
    // Update link in db
    // /api/collections/files/records/:id
    updateUrl := url + "/" + respBodyMap["id"].(string)
    updateReq, err := http.NewRequest("PATCH", updateUrl, bytes.NewBuffer([]byte(`{"link": "` + link + `"}`)))
    if err != nil {
        return fmt.Errorf("failed to create update request to update link in db: %w", err)
    }
    updateReq.Header.Set("Content-Type", "application/json")
    updateClient := &http.Client{}
    updateResp, err := updateClient.Do(updateReq)
    if err != nil {
        return fmt.Errorf("failed to send update request to update link in db: %w", err)
    }
    defer updateResp.Body.Close()
    if updateResp.StatusCode != http.StatusOK && updateResp.StatusCode != http.StatusCreated {
        return fmt.Errorf("update failed with status to update link in db: %d", updateResp.StatusCode)
    }
    fmt.Println("Link:", link)
    return nil
}