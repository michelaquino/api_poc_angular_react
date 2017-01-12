# api_poc_angular_react

- Create user
`curl -i -X POST -d  "{\"name\":\"Michel\",\"email\":\"michel.has@gmail.com\", \"gender\": \"m\"}" http://localhost:8080/user -H "Content-Type: application/json"`

- Read all users
`curl http://localhost:8080/users`