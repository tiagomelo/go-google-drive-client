// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package googledrive

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Constant definitions for various roles and grantee types.
const (
	// roles
	OwnerRole         = "owner"
	OrganizerRole     = "organizer"
	FileOrganizerRole = "fileOrganizer"
	WriterRole        = "writer"
	CommenterRole     = "commenter"
	ReaderRole        = "reader"

	// grantee types
	UserGranteeType   = "user"
	GroupGranteeType  = "group"
	DomainGranteeType = "domain"
	AnyoneGranteeType = "anyone"
)

// For ease of unit testing.
var (
	newDriveService = drive.NewService
	osCreate        = os.Create
	ioCopy          = io.Copy
)

// Client defines the interface for interacting with Google Drive.
type Client interface {
	// CreateFolder creates a new folder in Google Drive with the specified name.
	// The new folder can optionally be created within specified parent folders.
	// Returns the ID of the created folder or an error if the operation fails.
	CreateFolder(folderName string, parentFolders ...string) (string, error)

	// UploadFile uploads a file to Google Drive.
	// The file to upload is specified as an os.File pointer.
	// The file can optionally be uploaded within specified parent folders.
	// Returns the ID of the uploaded file or an error if the operation fails.
	UploadFile(file *os.File, parentFolders ...string) (string, error)

	// UpdateFile updates the content of an existing file in Google Drive.
	// The file to update is specified by its file ID, and the new content is
	// provided as an os.File pointer.
	// Returns the ID of the updated file or an error if the operation fails.
	UpdateFile(fileId string, newContent *os.File) (string, error)

	// GetFileById retrieves the details of a file from Google Drive using its file ID.
	// Returns a drive.File object containing file details or an error if retrieval fails.
	GetFileById(fileId string) (*drive.File, error)

	// DeleteFile deletes a file from Google Drive using its file ID.
	// Returns an error if the deletion operation fails.
	DeleteFile(fileId string) error

	// AssignRoleToUserOnFile assigns a specified role to a user for a file in Google Drive.
	// The role, user's email address, and file ID are specified.
	// Returns an error if the role assignment operation fails.
	AssignRoleToUserOnFile(role, emailAddress, fileId string) error

	// AssignRoleToGroupOnFile assigns a specified role to a group for a file in Google Drive.
	// The role, group's email address, and file ID are specified.
	// Returns an error if the role assignment operation fails.
	AssignRoleToGroupOnFile(role, emailAddress, fileId string) error

	// AssignRoleToDomainOnFile assigns a specified role to a domain for a file in Google Drive.
	// The role, domain, and file ID are specified.
	// Returns an error if the role assignment operation fails.
	AssignRoleToDomainOnFile(role, domain, fileId string) error

	// AssignRoleToAnyoneOnFile assigns a specified role to anyone for a file in Google Drive.
	// The role and file ID are specified.
	// Returns an error if the role assignment operation fails.
	AssignRoleToAnyoneOnFile(role, fileId string) error

	// DownloadFile downloads a file from Google Drive using its file ID.
	// The downloaded file is saved to the specified outputFile path.
	// Returns the path of the downloaded file or an error if the download operation fails.
	DownloadFile(fileId, outputFile string) (string, error)
}

// client is the concrete implementation of the Client interface, providing
// methods to interact with Google Drive's file and permission services.
type client struct {
	srv  driveService
	pSrv permissionsService
}

// New creates a new Client instance for interacting with Google Drive.
// It initializes the drive service with the provided context and credentials file.
func New(ctx context.Context, credsFilePath string) (Client, error) {
	srv, err := newDriveService(ctx, option.WithCredentialsFile(credsFilePath))
	if err != nil {
		return nil, errors.Wrap(err, "creating drive service")
	}
	dsw := &driveServiceWrapper{
		fsw: &fileServiceWrapper{srv.Files},
		psw: &permissionsServiceWrapper{srv.Permissions},
	}
	pSrv := dsw.Permissions()
	return &client{
		srv:  dsw,
		pSrv: pSrv,
	}, nil
}

func (c *client) CreateFolder(folderName string, parentFolders ...string) (string, error) {
	const mimeType = "application/vnd.google-apps.folder"
	createdFolder, err := c.srv.Files().Create(&drive.File{
		Name:     folderName,
		MimeType: mimeType,
		Parents:  parentFolders,
	}).Do()
	if err != nil {
		return "", errors.Wrapf(err, "creating folder %s under parent folders %v", folderName, parentFolders)
	}
	return createdFolder.Id, nil
}

func (c *client) GetFileById(fileId string) (*drive.File, error) {
	driveFile, err := c.srv.Files().Get(fileId).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "getting file with id %s", fileId)
	}
	return driveFile, nil
}

func (c *client) UploadFile(file *os.File, parentFolders ...string) (string, error) {
	driveFile, err := c.srv.Files().Create(&drive.File{
		Name:    filepath.Base(file.Name()),
		Parents: parentFolders,
	}).Media(file).Do()
	if err != nil {
		return "", errors.Wrapf(err, "creating file %s under parent folders %v", file.Name(), parentFolders)
	}
	return driveFile.Id, nil
}

func (c *client) DeleteFile(fileId string) error {
	if err := c.srv.Files().Delete(fileId).Do(); err != nil {
		return errors.Wrapf(err, "deleting file with id %s", fileId)
	}
	return nil
}

func (c *client) assignPermissionOnFile(permission *drive.Permission, fileId string) error {
	_, err := c.pSrv.Create(fileId, permission).Do()
	return err
}

func (c *client) AssignRoleToUserOnFile(role, emailAddress, fileId string) error {
	if err := c.assignPermissionOnFile(&drive.Permission{
		EmailAddress: emailAddress,
		Type:         UserGranteeType,
		Role:         role,
	}, fileId); err != nil {
		return errors.Wrapf(err, "assigning role %s on file with id %s to email address %s", role, fileId, emailAddress)
	}
	return nil
}

func (c *client) AssignRoleToGroupOnFile(role, emailAddress, fileId string) error {
	if err := c.assignPermissionOnFile(&drive.Permission{
		EmailAddress: emailAddress,
		Type:         GroupGranteeType,
		Role:         role,
	}, fileId); err != nil {
		return errors.Wrapf(err, "assigning role %s on file with id %s to email address %s", role, fileId, emailAddress)
	}
	return nil
}

func (c *client) AssignRoleToDomainOnFile(role, domain, fileId string) error {
	if err := c.assignPermissionOnFile(&drive.Permission{
		Domain: domain,
		Type:   DomainGranteeType,
		Role:   role,
	}, fileId); err != nil {
		return errors.Wrapf(err, "assigning role %s on file with id %s to domain %s", role, fileId, domain)
	}
	return nil
}

func (c *client) AssignRoleToAnyoneOnFile(role, fileId string) error {
	if err := c.assignPermissionOnFile(&drive.Permission{
		Type: AnyoneGranteeType,
		Role: role,
	}, fileId); err != nil {
		return errors.Wrapf(err, "assigning role %s on file with id %s to anyone", role, fileId)
	}
	return nil
}

func (c *client) DownloadFile(fileId, outputFile string) (string, error) {
	resp, err := c.srv.Files().Get(fileId).Download()
	if err != nil {
		return "", errors.Wrapf(err, "downloading file with id %s", fileId)
	}
	defer resp.Body.Close()
	f, err := osCreate(outputFile)
	if err != nil {
		return "", errors.Wrapf(err, "creating output file %s", outputFile)
	}
	defer f.Close()
	if _, err = ioCopy(f, resp.Body); err != nil {
		return "", errors.Wrapf(err, "writing output file %s", outputFile)
	}
	return f.Name(), nil
}

func (c *client) UpdateFile(fileId string, newContent *os.File) (string, error) {
	updatedFile, err := c.srv.Files().Update(fileId, nil).Media(newContent).Do()
	if err != nil {
		return "", errors.Wrapf(err, "updating file %s", fileId)
	}
	return updatedFile.Id, nil
}
