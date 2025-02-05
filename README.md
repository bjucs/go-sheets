# go-sheets
go-sheets is a CLI tool written Go for managing assignments related to courses. It uses a simple Google Sheet (created upon first run / initial setup) via the Google Sheets API as a backend. 

Students can organize their assignments using a simple command-line interface, and a selection of support commands available at their disposal. 

## Supported commands
- `info`
    - Displays the various commands that a user has available to use
- `create-course <course_name> [<class_description>]`
    - User can optionally include an additional `<class_description>` parameter 
- `create-assignment <course_name> <assignment_name>` 
    - User will be prompted for other info, such as `due_date` (required) and `info` (optional notes)
- `list-courses` 
- `list-assignments <course_name>`
- `remove-course <course_name>`
- `remove-assignment <course_name> <assignment_number>`
    - `assignment_number` is a 1-based index viewable using the `list-assignments <course_name>` command

## Setup and usage
In order to setup `go-sheets`, you will need a Google Cloud project (and service account credentials), which you can set up by following [this](https://developers.google.com/sheets/api/quickstart/go) Google Cloud Go tutorial.

Instead of using an OAuth 2.0 client ID, you're going to setup service account credentials, which you will download as `service-account.json` and place into the `cmd/go-sheets/cli` directory (where the `main.go` is located). This will give go