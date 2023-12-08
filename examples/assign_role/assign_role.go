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
	FileId        string `short:"f" long:"folder" description:"File Id" required:"true"`
	EmailAddress  string `short:"e" long:"email" description:"Email address" required:"true"`
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

	// Assign role to user
	err = client.AssignRoleToUserOnFile(googledrive.WriterRole, opts.EmailAddress, opts.FileId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("role assigned successfully to user")

	// Assign role to group
	// const groupEmailAddress = "email@email.com"
	// err = client.AssignRoleToGroupOnFile(googledrive.WriterRole, groupEmailAddress, fileId)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println("role assigned successfully to group")

	// Assign role to domain
	// const domain = "somedomain"
	// err = client.AssignRoleToDomainOnFile(googledrive.WriterRole, domain, fileId)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println("role assigned successfully to domain")

	// Assign role to anyone
	// err = client.AssignRoleToAnyoneOnFile(googledrive.WriterRole, domain)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println("role assigned successfully to anyone")
}
