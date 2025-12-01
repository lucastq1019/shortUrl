grpcurl -plaintext -d '{"long_url":"www.google.com"}' localhost:9090 shortener.ShortenerService/CreateShortLink

grpcurl -plaintext localhost:9090 list|xargs -I {} grpcurl -plaintext localhost:9090 describe {}  > all.txt

# 用于查询当前grpc server支持的接口
grpcurl -plaintext ip:port list

grpcurl -plaintext ip:port describe {interface_name}

# 用于查询当前grpc server支持的接口的详细信息
grpcurl -plaintext ip:port list|xargs -I {} grpcurl -plaintext ip:port describe {}  > all.txt


grpcurl -plaintext localhost:9090 list|xargs -I {} grpcurl -plaintext localhost:9090 describe {}  > all.txt

