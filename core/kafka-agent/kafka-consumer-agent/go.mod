module github.com/mycreditchain/kafka-agent/kafka-consumer-agent

go 1.14

require (
	github.com/Shopify/sarama v1.27.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/labstack/echo v3.3.10+incompatible // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mycreditchain/common/msg v0.0.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0

)

replace github.com/mycreditchain/common/msg v0.0.0 => ./../../common/msg
