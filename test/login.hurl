POST http://localhost:8080/api/public/users
HTTP 201
[Captures]
token: body

POST http://localhost:8080/api/public/login
[Form]
token: {{token}}
HTTP 303

POST http://localhost:8080/api/public/login
[Form]
token: abc
HTTP 400

POST http://localhost:8080/api/public/login
HTTP 400
