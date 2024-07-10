POST http://{{host}}/v1/insert
Content-Type: application/json
{
  "value": 2, "index": 0
}
HTTP 201
[Asserts]
body jsonpath "$.message" == "Insert successful"

POST http://{{host}}/v1/insert
Content-Type: application/json
{
  "value": 4, "index": 1
}
HTTP 201
[Asserts]
body jsonpath "$.message" == "Insert successful"

POST http://{{host}}/v1/insert
Content-Type: application/json
{
  "value": 3, "index": 3
}
HTTP 404
[Asserts]
body == "Index out of range\n"

GET http://{{host}}/v1/get/0
HTTP 200
[Asserts]
body jsonpath "$.value" == 2

GET http://{{host}}/v1/get/2
HTTP 404
[Asserts]
body == "Index out of range\n"

GET http://{{host}}/v1/find/4
HTTP 200
[Asserts]
body jsonpath "$.index" == 1

GET http://{{host}}/v1/find/6
HTTP 404
[Asserts]
body == "Value not found\n"

DELETE http://{{host}}/v1/remove/1
HTTP 200

GET http://{{host}}/v1/get/1
HTTP 404
[Asserts]
body == "Index out of range\n"

POST http://{{host}}/v2/numbers/0/2
Content-Type: application/json
HTTP 201
[Asserts]
jsonpath "$.index" == 0
jsonpath "$.value" == 2

POST http://{{host}}/v2/numbers/1/4
Content-Type: application/json
HTTP 201
[Asserts]
jsonpath "$.index" == 1
jsonpath "$.value" == 4

POST http://{{host}}/v2/numbers/3/3
Content-Type: application/json
HTTP 400
[Asserts]
jsonpath "$.message" == "Invalid index"

GET http://{{host}}/v2/numbers/index/0
HTTP 200
[Asserts]
jsonpath "$.index" == 0
jsonpath "$.value" == 2

GET http://{{host}}/v2/numbers/index/2
HTTP 404
[Asserts]
jsonpath "$.message" == "Index not found"

GET http://{{host}}/v2/numbers/value/4
HTTP 200
[Asserts]
jsonpath "$.index" == 1
jsonpath "$.value" == 4

GET http://{{host}}/v2/numbers/value/6
HTTP 404
[Asserts]
jsonpath "$.message" == "Value not found"

DELETE http://{{host}}/v2/numbers/1
HTTP 200

GET http://{{host}}/v2/numbers/index/1
HTTP 404
[Asserts]
jsonpath "$.message" == "Index not found"
