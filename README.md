# bn-lotto-service
# export "DB_CONNECTION=host=localhost user=root password=11111111 dbname=lotto_db port=5432 sslmode=disable TimeZone=Asia/Bangkok search_path=root@lotto_db"


# payout
# curl --location --request POST 'http://localhost:8080/gateway-api/pay-out'
# "{\"code\":0,\"data\":{\"order_no\":\"120250730f3cd549039cb4b66b4388f93555ec16a\",\"payUrl\":\"https://pay.ghpay. vip/#/pay?order=120250730f3cd549039cb4b66b4388f93555ec16a\"}}"


# payin
# curl --location --request POST 'http://localhost:8080/gateway-api/pay-in'
# "{\"code\":0,\"data\":{\"order_no\":\"12025073038d6600767af4785a64abf4c380b9af9\",\"payUrl\":\"https://pay.ghpay.vip/#/pay?order=12025073038d6600767af4785a64abf4c380b9af9\"}}"
