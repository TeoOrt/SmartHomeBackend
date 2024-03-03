# SmartHomeBackend

SmartHomeBackend is a backend system designed to facilitate the management and control of smart home devices. It provides a RESTful API for interacting with various smart home devices such as lights, thermostats, and more.

## Clone the repository:

```
git clone https://github.com/TeoOrt/SmartHomeBackend.git
```

Navigate into the project directory:

```
cd SmartHomeBackend
```

Install dependencies:
```
go mod tiddy
```
Start the server:

```
go build .
./main
```

### Usage
Once the server is running, you can start interacting with the API using tools like Postman or curl. Ensure you have appropriate authentication tokens for accessing protected endpoints.

API Endpoints

- POST /upload_video: Upload video to database.
- GET /get_counter: Retrieves all of the video gestures recorded.
- GET /get_expert/video/: Retrieves Expert videos for demos.


Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

License
This project is licensed under the MIT License.

Feel free to adjust the content as needed to accurately reflect the specifics of your project and its functionalities.