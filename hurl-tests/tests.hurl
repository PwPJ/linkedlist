POST http://{{host}}/insert
Content-Type: application/json
{
  "value": 2, "index": 0
}
HTTP 201
[Asserts]
body jsonpath "$.message" == "Insert successful"

POST http://{{host}}/insert
Content-Type: application/json
{
  "value": 4, "index": 1
}
HTTP 201
[Asserts]
body jsonpath "$.message" == "Insert successful"

POST http://{{host}}/insert
Content-Type: application/json
{
  "value": 3, "index": 3
}
HTTP 404
[Asserts]
body == "Index out of range\n"

GET http://{{host}}/get/0
HTTP 200
[Asserts]
body jsonpath "$.value" == 2

GET http://{{host}}/get/2
HTTP 404
[Asserts]
body == "Index out of range\n"

GET http://{{host}}/find/4
HTTP 200
[Asserts]
body jsonpath "$.index" == 1

GET http://{{host}}/find/6
HTTP 404
[Asserts]
body == "Value not found\n"

DELETE http://{{host}}/remove/1
HTTP 200

GET http://{{host}}/get/1
HTTP 404
[Asserts]
body == "Index out of range\n"
