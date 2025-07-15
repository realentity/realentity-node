@echo off
echo Testing /api/services/execute endpoint
echo =====================================

echo.
echo 1. Testing text.process service with uppercase operation:
curl -X POST http://localhost:8080/api/services/execute ^
  -H "Content-Type: application/json" ^
  -d "{\"service\":\"text.process\",\"payload\":{\"text\":\"hello world\",\"operation\":\"uppercase\"}}"

echo.
echo.
echo 2. Testing text.process service with lowercase operation:
curl -X POST http://localhost:8080/api/services/execute ^
  -H "Content-Type: application/json" ^
  -d "{\"service\":\"text.process\",\"payload\":{\"text\":\"HELLO WORLD\",\"operation\":\"lowercase\"}}"

echo.
echo.
echo 3. Testing text.process service with reverse operation:
curl -X POST http://localhost:8080/api/services/execute ^
  -H "Content-Type: application/json" ^
  -d "{\"service\":\"text.process\",\"payload\":{\"text\":\"hello\",\"operation\":\"reverse\"}}"

echo.
echo.
echo 4. Testing echo service:
curl -X POST http://localhost:8080/api/services/execute ^
  -H "Content-Type: application/json" ^
  -d "{\"service\":\"echo\",\"payload\":{\"message\":\"Hello from API!\"}}"

echo.
echo.
echo 5. Testing invalid service (should fail):
curl -X POST http://localhost:8080/api/services/execute ^
  -H "Content-Type: application/json" ^
  -d "{\"service\":\"nonexistent\",\"payload\":{}}"

echo.
echo.
echo Testing complete!
