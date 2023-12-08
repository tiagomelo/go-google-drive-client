// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package googledrive

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name                string
		mockJsonMarshal     func(v any) ([]byte, error)
		mockNewDriveService func(ctx context.Context, opts ...option.ClientOption) (*drive.Service, error)
		expectedError       error
	}{
		{
			name: "happy path",
			mockJsonMarshal: func(v any) ([]byte, error) {
				return []byte("{}"), nil
			},
			mockNewDriveService: func(ctx context.Context, opts ...option.ClientOption) (*drive.Service, error) {
				return new(drive.Service), nil
			},
		},
		{
			name: "error when creating service",
			mockJsonMarshal: func(v any) ([]byte, error) {
				return []byte("{}"), nil
			},
			mockNewDriveService: func(ctx context.Context, opts ...option.ClientOption) (*drive.Service, error) {
				return nil, errors.New("create service error")
			},
			expectedError: errors.New("creating drive service: create service error"),
		},
	}
	for _, tc := range testCases {
		newDriveService = tc.mockNewDriveService
		t.Run(tc.name, func(t *testing.T) {
			c, err := New(context.TODO(), "path/to/credsfile")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
				require.NotNil(t, c)
			}
		})
	}
}

func TestCreateFolder(t *testing.T) {
	testCases := []struct {
		name           string
		mockClosure    func(mfcc *mockFilesCreateCall)
		expectedOutput string
		expectedError  error
	}{
		{
			name: "happy path",
			mockClosure: func(mfcc *mockFilesCreateCall) {
				mfcc.file = &drive.File{
					Id: "someFolderId",
				}
			},
			expectedOutput: "someFolderId",
		},
		{
			name: "error",
			mockClosure: func(mfcc *mockFilesCreateCall) {
				mfcc.doErr = errors.New("create error")
			},
			expectedError: errors.New("creating folder someFolder under parent folders [parentFolder]: create error"),
		},
	}
	for _, tc := range testCases {
		mfcc := new(mockFilesCreateCall)
		mfs := new(mockFileService)
		mds := new(mockDriveService)
		mds.fs = mfs
		mfs.fcc = mfcc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mfcc)
			client := &client{
				srv: mds,
			}
			output, err := client.CreateFolder("someFolder", "parentFolder")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestUploadFile(t *testing.T) {
	testCases := []struct {
		name           string
		file           *os.File
		mockClosure    func(mfcc *mockFilesCreateCall)
		expectedOutput string
		expectedError  error
	}{
		{
			name: "happy path",
			file: os.NewFile(555, "newFile"),
			mockClosure: func(mfcc *mockFilesCreateCall) {
				mfcc.file = &drive.File{
					Id: "someFileId",
				}
			},
			expectedOutput: "someFileId",
		},
		{
			name: "error",
			file: os.NewFile(555, "newFile"),
			mockClosure: func(mfcc *mockFilesCreateCall) {
				mfcc.doErr = errors.New("Do error")
			},
			expectedError: errors.New("creating file newFile under parent folders [parentFolder]: Do error"),
		},
	}
	for _, tc := range testCases {
		mfcc := new(mockFilesCreateCall)
		mfs := new(mockFileService)
		mds := new(mockDriveService)
		mds.fs = mfs
		mfs.fcc = mfcc
		mfcc.fcc = mfcc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mfcc)
			client := &client{
				srv: mds,
			}
			output, err := client.UploadFile(tc.file, "parentFolder")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestGetFileById(t *testing.T) {
	testCases := []struct {
		name           string
		mockClosure    func(mfgc *mockFilesGetCall)
		expectedOutput *drive.File
		expectedError  error
	}{
		{
			name: "happy path",
			mockClosure: func(mfgc *mockFilesGetCall) {
				mfgc.File = &drive.File{
					Id: "fileId",
				}
			},
			expectedOutput: &drive.File{
				Id: "fileId",
			},
		},
		{
			name: "error",
			mockClosure: func(mfgc *mockFilesGetCall) {
				mfgc.DoErr = errors.New("get error")
			},
			expectedError: errors.New("getting file with id fileId: get error"),
		},
	}
	for _, tc := range testCases {
		mfgc := new(mockFilesGetCall)
		mfs := new(mockFileService)
		mds := new(mockDriveService)
		mds.fs = mfs
		mfs.fgc = mfgc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mfgc)
			client := &client{
				srv: mds,
			}
			output, err := client.GetFileById("fileId")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	testCases := []struct {
		name          string
		mockClosure   func(mfdc *mockFilesDeleteCall)
		expectedError error
	}{
		{
			name:        "happy path",
			mockClosure: func(mfdc *mockFilesDeleteCall) {},
		},
		{
			name: "error",
			mockClosure: func(mfdc *mockFilesDeleteCall) {
				mfdc.Err = errors.New("delete error")
			},
			expectedError: errors.New("deleting file with id fileId: delete error"),
		},
	}
	for _, tc := range testCases {
		mfdc := new(mockFilesDeleteCall)
		mfs := new(mockFileService)
		mds := new(mockDriveService)
		mds.fs = mfs
		mfs.fdc = mfdc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mfdc)
			client := &client{
				srv: mds,
			}
			err := client.DeleteFile("fileId")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
			}
		})
	}
}

func TestAssignRoleToUserOnFile(t *testing.T) {
	testCases := []struct {
		name          string
		mockClosure   func(m *mockPermissionsCreateCall)
		expectedError error
	}{
		{
			name:        "happy path",
			mockClosure: func(m *mockPermissionsCreateCall) {},
		},
		{
			name: "error",
			mockClosure: func(m *mockPermissionsCreateCall) {
				m.Err = errors.New("create error")
			},
			expectedError: errors.New("assigning role some role on file with id fileId to email address email: create error"),
		},
	}
	for _, tc := range testCases {
		mps := new(mockPermissionsService)
		mpcc := new(mockPermissionsCreateCall)
		mps.pcc = mpcc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mpcc)
			client := &client{
				pSrv: mps,
			}
			err := client.AssignRoleToUserOnFile("some role", "email", "fileId")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
			}
		})
	}
}

func TestAssignRoleToGroupOnFile(t *testing.T) {
	testCases := []struct {
		name          string
		mockClosure   func(m *mockPermissionsCreateCall)
		expectedError error
	}{
		{
			name:        "happy path",
			mockClosure: func(m *mockPermissionsCreateCall) {},
		},
		{
			name: "error",
			mockClosure: func(m *mockPermissionsCreateCall) {
				m.Err = errors.New("create error")
			},
			expectedError: errors.New("assigning role some role on file with id fileId to email address email: create error"),
		},
	}
	for _, tc := range testCases {
		mps := new(mockPermissionsService)
		mpcc := new(mockPermissionsCreateCall)
		mps.pcc = mpcc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mpcc)
			client := &client{
				pSrv: mps,
			}
			err := client.AssignRoleToGroupOnFile("some role", "email", "fileId")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
			}
		})
	}
}

func TestAssignRoleToDomainOnFile(t *testing.T) {
	testCases := []struct {
		name          string
		mockClosure   func(m *mockPermissionsCreateCall)
		expectedError error
	}{
		{
			name:        "happy path",
			mockClosure: func(m *mockPermissionsCreateCall) {},
		},
		{
			name: "error",
			mockClosure: func(m *mockPermissionsCreateCall) {
				m.Err = errors.New("create error")
			},
			expectedError: errors.New("assigning role some role on file with id fileId to domain domain: create error"),
		},
	}
	for _, tc := range testCases {
		mps := new(mockPermissionsService)
		mpcc := new(mockPermissionsCreateCall)
		mps.pcc = mpcc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mpcc)
			client := &client{
				pSrv: mps,
			}
			err := client.AssignRoleToDomainOnFile("some role", "domain", "fileId")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
			}
		})
	}
}

func TestAssignRoleToAnyoneOnFile(t *testing.T) {
	testCases := []struct {
		name          string
		mockClosure   func(m *mockPermissionsCreateCall)
		expectedError error
	}{
		{
			name:        "happy path",
			mockClosure: func(m *mockPermissionsCreateCall) {},
		},
		{
			name: "error",
			mockClosure: func(m *mockPermissionsCreateCall) {
				m.Err = errors.New("create error")
			},
			expectedError: errors.New("assigning role some role on file with id fileId to anyone: create error"),
		},
	}
	for _, tc := range testCases {
		mps := new(mockPermissionsService)
		mpcc := new(mockPermissionsCreateCall)
		mps.pcc = mpcc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mpcc)
			client := &client{
				pSrv: mps,
			}
			err := client.AssignRoleToAnyoneOnFile("some role", "fileId")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	var tmpFile *os.File
	testCases := []struct {
		name          string
		mockClosure   func(m *mockFilesGetCall)
		mockOsCreate  func(name string) (*os.File, error)
		mockIoCopy    func(dst io.Writer, src io.Reader) (written int64, err error)
		expectedError error
	}{
		{
			name: "happy path",
			mockClosure: func(m *mockFilesGetCall) {
				m.HttpResponse = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}
			},
			mockOsCreate: func(name string) (*os.File, error) {
				var err error
				tmpFile, err = os.CreateTemp("", "test")
				require.Nil(t, err)
				return tmpFile, nil
			},
			mockIoCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 10, nil
			},
		},
		{
			name: "error when downloading file",
			mockClosure: func(m *mockFilesGetCall) {
				m.DownloadErr = errors.New("download error")
			},
			mockOsCreate: func(name string) (*os.File, error) {
				return nil, nil
			},
			mockIoCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, nil
			},
			expectedError: errors.New("downloading file with id filedId: download error"),
		},
		{
			name: "error when creating output file",
			mockClosure: func(m *mockFilesGetCall) {
				m.HttpResponse = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}
			},
			mockOsCreate: func(name string) (*os.File, error) {
				return nil, errors.New("os open error")
			},
			mockIoCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, nil
			},
			expectedError: errors.New("creating output file /path/to/file/test: os open error"),
		},
		{
			name: "error when writing output file",
			mockClosure: func(m *mockFilesGetCall) {
				m.HttpResponse = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}
			},
			mockOsCreate: func(name string) (*os.File, error) {
				var err error
				tmpFile, err = os.CreateTemp("", "test")
				require.Nil(t, err)
				return tmpFile, nil
			},
			mockIoCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, errors.New("copy error")
			},
			expectedError: errors.New("writing output file /path/to/file/test: copy error"),
		},
	}
	for _, tc := range testCases {
		mfgc := new(mockFilesGetCall)
		mfs := new(mockFileService)
		mds := new(mockDriveService)
		mds.fs = mfs
		mfs.fgc = mfgc
		osCreate = tc.mockOsCreate
		ioCopy = tc.mockIoCopy
		defer func() {
			os.Remove(tmpFile.Name())
		}()
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mfgc)
			client := &client{
				srv: mds,
			}
			output, err := client.DownloadFile("filedId", "/path/to/file/test")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
				require.NotEmpty(t, output)
			}
		})
	}
}

func TestUpdateFile(t *testing.T) {
	testCases := []struct {
		name           string
		mockClosure    func(mfcc *mockFilesUpdateCall)
		expectedOutput string
		expectedError  error
	}{
		{
			name: "happy path",
			mockClosure: func(mfuc *mockFilesUpdateCall) {
				mfuc.file = &drive.File{
					Id: "someFileId",
				}
			},
			expectedOutput: "someFileId",
		},
		{
			name: "happy path",
			mockClosure: func(mfuc *mockFilesUpdateCall) {
				mfuc.doErr = errors.New("update error")
			},
			expectedError: errors.New("updating file fileId: update error"),
		},
	}
	for _, tc := range testCases {
		mfuc := new(mockFilesUpdateCall)
		mfs := new(mockFileService)
		mds := new(mockDriveService)
		mds.fs = mfs
		mfs.fuc = mfuc
		mfuc.fuc = mfuc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mfuc)
			client := &client{
				srv: mds,
			}
			output, err := client.UpdateFile("fileId", &os.File{})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

type mockFilesCreateCall struct {
	file  *drive.File
	doErr error
	fcc   filesCreateCall
}

func (m *mockFilesCreateCall) Do(opts ...googleapi.CallOption) (*drive.File, error) {
	return m.file, m.doErr
}

func (m *mockFilesCreateCall) Media(r io.Reader, options ...googleapi.MediaOption) filesCreateCall {
	return m.fcc
}

type mockFilesGetCall struct {
	File         *drive.File
	DoErr        error
	HttpResponse *http.Response
	DownloadErr  error
}

func (m *mockFilesGetCall) Do(opts ...googleapi.CallOption) (*drive.File, error) {
	return m.File, m.DoErr
}

func (m *mockFilesGetCall) Download(opts ...googleapi.CallOption) (*http.Response, error) {
	return m.HttpResponse, m.DownloadErr
}

type mockFilesDeleteCall struct {
	Err error
}

func (m *mockFilesDeleteCall) Do(opts ...googleapi.CallOption) error {
	return m.Err
}

type mockFilesUpdateCall struct {
	file  *drive.File
	doErr error
	fuc   filesUpdateCall
}

func (m *mockFilesUpdateCall) Do(opts ...googleapi.CallOption) (*drive.File, error) {
	return m.file, m.doErr
}

func (m *mockFilesUpdateCall) Media(r io.Reader, options ...googleapi.MediaOption) filesUpdateCall {
	return m.fuc
}

type mockFileService struct {
	fcc filesCreateCall
	fdc filesDeleteCall
	fgc filesGetCall
	fuc filesUpdateCall
}

func (m *mockFileService) Create(file *drive.File) filesCreateCall {
	return m.fcc
}

func (m *mockFileService) Get(fileId string) filesGetCall {
	return m.fgc
}

func (m *mockFileService) Delete(fileId string) filesDeleteCall {
	return m.fdc
}

func (m *mockFileService) Update(fileId string, file *drive.File) filesUpdateCall {
	return m.fuc
}

type mockDriveService struct {
	fs fileService
	ps permissionsService
}

func (m *mockDriveService) Files() fileService {
	return m.fs
}

func (m *mockDriveService) Permissions() permissionsService {
	return m.ps
}

type mockPermissionsCreateCall struct {
	Err error
}

func (m *mockPermissionsCreateCall) Do(opts ...googleapi.CallOption) (*drive.Permission, error) {
	return new(drive.Permission), m.Err
}

type mockPermissionsService struct {
	pcc permissionsCreateCall
}

func (m *mockPermissionsService) Create(fileId string, permission *drive.Permission) permissionsCreateCall {
	return m.pcc
}
