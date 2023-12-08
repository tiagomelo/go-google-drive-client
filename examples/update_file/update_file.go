// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tiagomelo/go-google-drive-client/googledrive"
)

var opts struct {
	CredsFilePath string `short:"c" long:"creds" description:"Creds file path" required:"true"`
	FileId        string `long:"fileId" description:"File Id" required:"true"`
	FilePath      string `long:"filePath" description:"File path" required:"true"`
}

func main() {
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx := context.Background()
	client, err := googledrive.New(ctx, opts.CredsFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file, err := os.Open(opts.FilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	updatedFileId, err := client.UpdateFile(opts.FileId, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated file ID: %s\n", updatedFileId)
}
