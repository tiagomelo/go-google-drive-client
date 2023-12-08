# go-google-drive-client

![logo](go-google-drive-client.png)

A simple Go client for interacting with Google Drive using the Google Drive API (v3).

## available operations

- create folder
- upload file
- update file
- download file
- get file by id
- delete file
- set permissions (for user, group, domain or to anyone)

It is important to mention that for GCP everything is a _file_, be it a regular file or a folder (which is called a _drive_). So whenever you see a _fileId_ param, it can be either a regular file or a folder.

## setup

Before using this client, you need to set up a Google Cloud Project and enable the Google Drive API. Follow these steps:

1. Create a new project in the [Google Cloud Console](https://console.cloud.google.com/).
2. Enable the [Google Drive API](https://console.cloud.google.com/apis/library/drive.googleapis.com) for your project.
3. Create credentials (service account key) for your project:
   - Go to the "Credentials" page.
   - Click "Create Credentials" and select "Service account key".
   - Choose or create a new service account.
   - Set the role to `Editor` or another role that grants the necessary permissions.
   - Select JSON for the key type and download the file.
4. Save the downloaded JSON file in the `creds` directory.

For more detailed instructions, see the [Google Drive API documentation](https://developers.google.com/drive/api/v3/quickstart/go).

## examples

- [create folder](examples/create_folder/create_folder.go)
- [upload file](examples/upload_file/upload_file.go)
- [update file](examples/update_file/update_file.go)
- [download file](examples/download_file/download_file.go)
- [get file by id](examples/get_file_by_id/get_file_by_id.go)
- [delete file](examples/delete_folder/delete_folder.go)
- [set permissions (for user, group, domain or to anyone)](examples/assign_role/assign_role.go)

## unit tests

```
make test
```

## unit tests coverage

```
make coverage
```