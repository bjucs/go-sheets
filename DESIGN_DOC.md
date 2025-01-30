# go-sheets
go-sheets will be a CLI tool written Go for managing assignments related to courses. It will use the Google Sheets API as a backend. Students can organize their assignments using a simple command-line interface, and a selection of support commands available at their disposal. (May be added to in the future)

## Supported commands
- `create-course <course_name>`
    - User can optionally include an additional `<class_description>` parameter 
- `create-assignment <course_name> <assignment_name>` 
    - User will be prompted for other info, such as `due_date` (required) and `info` (optional notes)
- `list-courses` 
- `list-assignments <course_name>`
- `remove-course <course_name>`
- `remove-assignment <assignment_number>`
    - `assignment_number` is a 1-based index viewable using the `list-assignments <coursename>` command

## Implementation plan
1. Create initial functionality for parsing in user commands 
2. Verify parsing functionality via unit tests
3. Implement initial backend logic for writing/retrieving data via local persistent read/write. 
4. Verify locally-written data logic works via unit tests
4. Implement final backend logic for writing/retrieving data via the Google Sheets API 
5. Verify full functional logic via more comprehensive unit tests
6. Verify go-sheets works as intended across machines & Google Cloud credentials 