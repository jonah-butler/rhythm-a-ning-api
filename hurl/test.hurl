# TEST: get rhythm by ID
# TAG: positive
GET http://127.0.0.1:8080/api/v1/rhythm/1

HTTP 200
[Asserts]
status == 200

# TEST: get subdivision types
# TAG: positive
GET http://127.0.0.1:8080/api/v1/rhythm/subdivisions

HTTP 200
[Asserts]
status == 200