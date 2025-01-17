package kafka

const (
	fieldBrokers           = "brokers"
	fieldClient            = "client"
	fieldGroup             = "group"
	fieldTopics            = "topics"
	fieldError             = "error"
	fieldErrors            = "errors"
	msgConsumeInit         = "mq consumer initializing..."
	msgConsumeInitFail     = "mq consumer init failed"
	msgConsumeStart        = "mq consumer started"
	msgConsumePollFail     = "mq consumer poll/fetch failed"
	msgConsumeProcessFail  = "mq consumer processing failed"
	msgConsumeCommitFail   = "mq consumer commit failed"
	msgConsumeStop         = "mq consumer stopped"
	msgProduceInit         = "mq producer initializing..."
	msgProduceInitFail     = "mq producer init failed"
	msgProduceStart        = "mq producer started"
	msgProduceSendFail     = "mq producer send failed"
	msgProducePingFail     = "mq producer ping failed"
	msgProduceStop         = "mq producer stopped"
	msgTransactInitFail    = "mq transacter init failed"
	msgTransactBeginFail   = "mq transaction begin failed"
	msgTransactProcessFail = "mq transaction processing failed"
	msgTransactEndFail     = "mg transaction end failed"
)
