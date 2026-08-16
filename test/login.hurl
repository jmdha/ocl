POST http://localhost:8080/api/public/users
HTTP 201
[Captures]
token: body

POST http://localhost:8080/api/public/login
[Form]
token: {{token}}
HTTP 303
