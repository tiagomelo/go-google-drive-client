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

// filesCreateCallWrapper wraps the Google Drive API's FilesCreateCall.
// This wrapper allows for method chaining and easier testing/mocking.
type filesCreateCallWrapper struct {
	fcc *drive.FilesCreateCall
}

func (fccw *filesCreateCallWrapper) Do(opts ...googleapi.CallOption) (*drive.File, error) {
	return fccw.fcc.Do(opts...)
}

func (fccw *filesCreateCallWrapper) Media(r io.Reader, options ...googleapi.MediaOption) filesCreateCall {
	fccw.fcc.Media(r, options...)
	return fccw
}

// filesGetCallWrapper wraps the Google Drive API's FilesGetCall.
// It provides a simplified interface for file retrieval operations.
type filesGetCallWrapper struct {
	fgc *drive.FilesGetCall
}

func (fgcw *filesGetCallWrapper) Do(opts ...googleapi.CallOption) (*drive.File, error) {
	return fgcw.fgc.Do(opts...)
}

func (fgcw *filesGetCallWrapper) Download(opts ...googleapi.CallOption) (*http.Response, error) {
	return fgcw.fgc.Download(opts...)
}

// filesDeleteCallWrapper wraps the Google Drive API's FilesDeleteCall.
// It provides a simplified interface for file deletion operations.
type filesDeleteCallWrapper struct {
	fdc *drive.FilesDeleteCall
}

func (fdcw *filesDeleteCallWrapper) Do(opts ...googleapi.CallOption) error {
	return fdcw.fdc.Do(opts...)
}

type filesUpdateCallWrapper struct {
	fuc *drive.FilesUpdateCall
}

func (fucw *filesUpdateCallWrapper) Do(opts ...googleapi.CallOption) (*drive.File, error) {
	return fucw.fuc.Do(opts...)
}

func (fucw *filesUpdateCallWrapper) Media(r io.Reader, options ...googleapi.MediaOption) filesUpdateCall {
	fucw.fuc.Media(r, options...)
	return fucw
}

// fileServiceWrapper wraps Google Drive's FilesService.
// It provides methods for file-related operations like creation, retrieval, and deletion.
type fileServiceWrapper struct {
	srv *drive.FilesService
}

func (fsw *fileServiceWrapper) Create(file *drive.File) filesCreateCall {
	return &filesCreateCallWrapper{fsw.srv.Create(file)}
}

func (fsw *fileServiceWrapper) Get(fileId string) filesGetCall {
	return &filesGetCallWrapper{fsw.srv.Get(fileId)}
}

func (fsw *fileServiceWrapper) Delete(fileId string) filesDeleteCall {
	return &filesDeleteCallWrapper{fsw.srv.Delete(fileId)}
}

func (fsw *fileServiceWrapper) Update(fileId string, file *drive.File) filesUpdateCall {
	return &filesUpdateCallWrapper{fsw.srv.Update(fileId, file)}
}

// driveServiceWrapper wraps Google Drive's Service.
// It provides a high-level interface to access file and permission operations.
type driveServiceWrapper struct {
	fsw *fileServiceWrapper
	psw *permissionsServiceWrapper
}

func (dsw *driveServiceWrapper) Files() fileService {
	return dsw.fsw
}

func (dsw *driveServiceWrapper) Permissions() permissionsService {
	return dsw.psw
}

// permissionsCreateCallWrapper wraps the Google Drive API's PermissionsCreateCall.
// It provides a simplified interface for setting permissions on files and folders.
type permissionsCreateCallWrapper struct {
	pcc *drive.PermissionsCreateCall
}

func (pccw *permissionsCreateCallWrapper) Do(opts ...googleapi.CallOption) (*drive.Permission, error) {
	return pccw.pcc.Do(opts...)
}

// permissionsServiceWrapper wraps Google Drive's PermissionsService.
// It provides methods for managing permissions on files and folders.
type permissionsServiceWrapper struct {
	pSrv *drive.PermissionsService
}

func (psw *permissionsServiceWrapper) Create(fileId string, permission *drive.Permission) permissionsCreateCall {
	return &permissionsCreateCallWrapper{psw.pSrv.Create(fileId, permission)}
}
