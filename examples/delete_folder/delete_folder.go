// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tiagomelo/go-google-drive-client/googledrive"
)

func main() {
	ctx := context.Background()
	const credsFilePath = "creds.json"
	client, err := googledrive.New(ctx, credsFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	folderId := "folderId"
	err = client.DeleteFile(folderId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("folder deleted successfully")
}
