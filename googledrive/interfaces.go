// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package googledrive

import (
	"io"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// filesCreateCall abstracts the file creation and upload operations of Google Drive.
// It provides functions to set up the file metadata and content for upload.
type filesCreateCall interface {
	// Do executes the file creation/upload call with the specified options.
	// It returns the created file object or an error if the operation fails.
	Do(opts ...googleapi.CallOption) (*drive.File, error)

	// Media sets the content of the file to be uploaded.
	// It accepts an io.Reader for the file content and returns a filesCreateCall
	// for chaining further configurations.
	Media(r io.Reader, options ...googleapi.MediaOption) filesCreateCall
}

// filesGetCall abstracts the operation of retrieving a file from Google Drive.
// It provides a function to execute the retrieval operation.
type filesGetCall interface {
	// Do executes the file retrieval call with the specified options.
	// It returns the retrieved file object or an error if the operation fails.
	Do(opts ...googleapi.CallOption) (*drive.File, error)

	// Download executes the file export and download call with the specified options.
	// It returns an HTTP response containing the file content or an error if the operation fails.
	Download(opts ...googleapi.CallOption) (*http.Response, error)
}

// filesDeleteCall abstracts the operation of deleting a file from Google Drive.
// It provides a function to execute the deletion operation.
type filesDeleteCall interface {
	// Do executes the file deletion call with the specified options.
	// It returns an error if the operation fails.
	Do(opts ...googleapi.CallOption) error
}

// filesUpdateCall abstracts the operation of updating a file in Google Drive.
// It provides functions to execute the update operation and set the file's new content.
type filesUpdateCall interface {
	// Do executes the file update call with the specified options.
	Do(opts ...googleapi.CallOption) (*drive.File, error)

	// Media sets the new content of the file to be updated.
	// It accepts an io.Reader for the new file content and returns a filesUpdateCall
	// for chaining further configurations or executing the update.
	Media(r io.Reader, options ...googleapi.MediaOption) filesUpdateCall
}

// fileService defines functions for file operations in Google Drive.
type fileService interface {
	// Create initializes a new file creation call with the given file metadata.
	// It returns a filesCreateCall to further configure and execute the file creation.
	Create(file *drive.File) filesCreateCall

	// Get initializes a file retrieval call for the specified file ID.
	// It returns a filesGetCall to execute the file retrieval.
	Get(fileId string) filesGetCall

	// Delete initializes a file deletion call for the specified file ID.
	// It returns a filesDeleteCall to execute the file deletion.
	Delete(fileId string) filesDeleteCall

	// Update initializes a file update call for the specified file ID.
	// The file parameter allows updating the file's metadata (such as name or description).
	// It returns a filesUpdateCall, which can be used to set the new content of the file
	// and execute the update operation.
	Update(fileId string, file *drive.File) filesUpdateCall
}

// driveService abstracts the Google Drive service.
// It provides a high-level interface to access various functionalities,
// such as file and permission management.
type driveService interface {
	// Files returns a fileService that provides file-related operations
	// such as creating, retrieving, and deleting files.
	Files() fileService

	// Permissions returns a permissionsService that provides operations
	// related to managing permissions on files and folders.
	Permissions() permissionsService
}

// permissionsCreateCall abstracts the operation of setting permissions
// for a file or folder in Google Drive.
type permissionsCreateCall interface {
	// Do executes the permission creation call with the specified options.
	// It returns the created permission object or an error if the operation fails.
	Do(opts ...googleapi.CallOption) (*drive.Permission, error)
}

// permissionsService defines functions for managing permissions
// on files and folders in Google Drive.
type permissionsService interface {
	// Create initializes a new permission creation call for the specified file ID.
	// It returns a permissionsCreateCall to further configure and execute the permission setting.
	Create(fileId string, permission *drive.Permission) permissionsCreateCall
}
