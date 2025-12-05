curl -X POST http://localhost:8080/shortener/v1/c \
     -H "Content-Type: application/json" \
     -d '{"long_url":"www.baidu.com"}'


curl http://localhost:8080/shortener/v1/all