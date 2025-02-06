# go-sheets
go-sheets is a CLI tool written in Go for managing assignments related to courses. It uses a simple Google Sheet (created upon first run / initial setup) via the Google Sheets API as a backend. 

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
- `exit` 

## Setup and usage
In order to setup `go-sheets`, you will need a Google Cloud project (and service account credentials), which you can set up by following [this](https://developers.google.com/sheets/api/quickstart/go) Google Cloud Go tutorial.

Instead of using an OAuth 2.0 client ID, you're going to setup service account credentials, which you will download as `service-account.json` and place into the `cmd/go-sheets/cli` directory (where the `main.go` is located). This will give go-sheets credentials to freely create and access a sheet used for persistent storage. If not already present, create a `.env` file for environment variables in this folder as well. 

*Make sure you run the main.go executable from within this directory* (`cmd/go-sheets/cli`) so it has access to your environment variables and service account credentials. You should also notice a `gosheets-cli.log` file that will re-write itself upon each run of go-sheets, notably logging all API usage-related errors and irregularities that occur. 

Finally, once you have all credential and Google Cloud setup done - you can simply run go-sheets via `go run main.go` from within the directory (`cmd/go-sheets-cli`) that it resides in.

## Sheets API 
If you're curious, you can look at `courseapi/course_api.go` to view our various types, but essentially we use two columns in the sheet where the first (A) is a string of `course_name`, and the second (B) is a serialized JSON of a `CourseItem`. 

Since CourseItems include all AssignmentItems (i.e. assignments are tied to a course), we're able to retrieve all course and assignment data by going down each row and un-serializing these JSONs back into CourseItem objects. This is then used to populate a central CourseMap which is the highest-level struct used for course and assignment lookup. 


