package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
    app := &cli.App{
        Name:  "greet",
        Usage: "fight the loneliness!",
        Action: func(*cli.Context) error {
            fmt.Println("Hello friend!")
            return nil
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}


func config(address string){
}

func upload() string{
    return ""
}

func list() string{
    return ""
}

func delete(){

}