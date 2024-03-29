# cyprusExercise
Golang REST API microservice to handle Companies

Everything that I wanted to do but did not have time can be viewed in the `todo` file

Company is an entity defined by the following attributes:
- Uuid - uniqueidenfier, generated automatically on creation of object
- Name - string
- Code - string
- Country - string
- Website - string
- Phone - string

!Stub JWT token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"!
Accordance with the exercise to make some requests you should have it in your cookies with name 'token'

The operations are defined by the query method as follows:
- Create - (POST) /api/company.
    - In request body should be an objects with fields of company, we want to insert. Empty fields will be recorded as null.
    - Example of path:
        - POST /api/company
    - Example of body:
        - [{"Name": "AnyNewName", "Country": "AnyNewCountry", "Phone": "AnyNewPhone"}...]
    - DB record:
        - [{"Name": "AnyNewName", "Code": null, "Country": "AnyNewCountry", "Website": null, "Phone": "AnyNewPhone"}...]
    - Returned values:
        - OK:
            - Status - 201;
            - Body - inserted object;
            - Example:
                - [{"Name": "AnyNewName", "Code": null, "Country": "AnyNewCountry", "Website": null, "Phone": "AnyNewPhone"}]
        - Error:
            - Status - according to error;
            - Body - {"error": "Error description"}


- Read - (GET) /api/company.
    - In accordance with the exercise in query parameters should be a fields to filter selected rows.
    - Besides it limit and offset parameters are available.
    - Example of path:
        - PATCH /api/company?Uuid=00000000-0000-0000-0000-000000000000&Name=OldName&Code=OldCode
    - Request body should be empty
    - DB record leaves unchanged
    - Returned values:
        - OK:
            - Status - 200;
            - Body - list of selected rows and their count;
            - Example:
                - {"data": [{"Name": "AnyNewName",...},...], count: 3}
        - Error:
            - Status - according to error;
            - Body - {"error": "Error description"}


- Update(Pure) - PUT /api/company.
    - In accordance with the exercise in query parameters should be a fields to filter updated rows.
    - !Attention! if multiple rows match filters - all of then will be patched.
    - Example of path:
        - PUT /api/company?Uuid=00000000-0000-0000-0000-000000000000&Name=OldName&Code=OldCode
    - In request body should be an object with new values for fields, we want to update. Missing fields will be recorded as null.
    - Example of body:
        - {"Name": "AnyNewName", "Country": "AnyNewCountry", "Phone": "AnyNewPhone"}
    - DB record:
        - {"Name": "AnyNewName", "Code": null, "Country": "AnyNewCountry", "Website": null, "Phone": "AnyNewPhone"}
    - Returned values:
        - OK:
            - Status - 200;
            - Body - list of updated rows and their count;
            - Example:
                - {"data": [{"Name": "AnyNewName",...},...], count: 3}
        - Error:
            - Status - according to error;
            - Body - {"error": "Error description"}


- Update(Merge) - PATCH /api/company.
    - In accordance with the exercise in query parameters should be a fields to filter updated rows.
    - !Attention! if multiple rows match filters - all of then will be updated.
    - Example of path:
        - PATCH /api/company?Uuid=00000000-0000-0000-0000-000000000000&Name=OldName&Code=OldCode
    - In request body should be an object with new values for fields, we want to update. Missing fields will be unchanged.
    - Example of body:
        - {"Name": "AnyNewName", "Country": "AnyNewCountry", "Phone": "AnyNewPhone"}
    - DB record:
        - {"Name": "AnyNewName", "Code": "OldCode", "Country": "AnyNewCountry", "Website": "OldWebsite", "Phone": "AnyNewPhone"}
    - Returned values:
        - OK:
            - Status - 200;
            - Body - list of updated rows and their count;
            - Example:
                - {"data": [{"Name": "AnyNewName",...},...], count: 3}
        - Error:
            - Status - according to error;
            - Body - {"error": "Error description"}


- Delete - DELETE /api/company.
    - In accordance with the exercise in query parameters should be a fields to filter deleted rows.
    - !Attention! if multiple rows match filters - all of then will be deleted.
    - Example of path:
        - DELETE /api/company?Uuid=00000000-0000-0000-0000-000000000000&Name=OldName&Code=OldCode
    - Request body should be empty
    - DB record leaves unchanged
    - DB record:
        - {"Name": "AnyNewName", "Code": "OldCode", "Country": "AnyNewCountry", "Website": "OldWebsite", "Phone": "AnyNewPhone"}
    - Returned values:
        - OK:
            - Status - 200;
            - Body - list of deleted rows and their count;
            - Example:
                - {"data": [{"Name": "AnyNewName",...},...], count: 3}
        - Error:
            - Status - according to error;
            - Body - {"error": "Error description"}
