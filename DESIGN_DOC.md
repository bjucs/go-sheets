# go-sheets
go-sheets will be a CLI tool written Go for managing assignments related to courses. It will use the Google Sheets API as a backend. Students can organize their assignments using a simple command-line interface, and a selection of support commands available at their disposal. (May be added to in the future)

## Supported commands
- `create-course <course_name>`
    - User can optionally include a `<class_description>` parameter 
- `create-assignment <assignment_name>` 
    - User will be prompted for other info, such as `due_date` and `info` (optional)
- `list-courses` 
- `list-assignments <course_name>`
- `remove-course <course_name>`
- `remove-assignment <course_name>`
    - User will be prompted for the `assignment_name` and a confirmation 

## Implementation plan
1. Create initial functionality for parsing in user commands 
2. Verify parsing functionality via unit tests
3. Implement backend logic for writing/retrieving data via the Google Sheets API 
4. Verify full functional logic via more comprehensive unit tests
5. Verify go-sheets works as intended across machines & Google Cloud credentials 