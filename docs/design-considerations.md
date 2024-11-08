# Design considerations

## used technologies
the application is a 2 layer(service/controller) application that uses Gorm for database access and Gin as a web framework. 

For ease of development the backing database is SQLite which saves to a file and uploaded files are stored locally in a directory which is exposed as a static filesystem through Gin. a real application dealing with files would likely use S3 buckets or similiar solutions to store the files.

the api Input/output models are spec-first, and the relevant go models and routes are generated using the OpenAPI generator from the openAPI specification to ensure that the spec can be used as the main source of truth.

gomock is used to mock services for controller testing.
## Improvements given time
* first and foremost would be using a real database instead of file storage SQLite, as well as a more sophisticated file storage solution such as S3 buckets or similiar.
* the major functionality of the app is unit tested, but unit tests for smaller utilities and more edge cases could be improved
* integration testing 
* many hardcoded values such as stored file location, listening port, graceful shutdown window, etc... would be either removed due to changes outlined above or otherwise extracted into environment variables