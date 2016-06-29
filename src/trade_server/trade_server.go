package trade_server

import (
	"github.com/apache/thrift/lib/go/thrift"
	"logger"
	"os"
	"time"
	"trade_service"
)

type TraderServer struct {
	host             string
	handler          *TradeServiceHandler
	processor        *trade_service.TradeServiceProcessor
	transport        *thrift.TServerSocket
	transportFactory thrift.TTransportFactory
	protocolFactory  *thrift.TBinaryProtocolFactory
	server           *thrift.TSimpleServer
}

func NewTraderServer(host string) *TraderServer {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory := thrift.NewTCompactProtocolFactory()

	transport, err := thrift.NewTServerSocketTimeout(host, 30*time.Second)
	if err != nil {
		logger.Infoln("NewTServerSocketTimeout Error!", err)
		os.Exit(1)
	}

	handler := &TradeServiceHandler{}
	processor := trade_service.NewTradeServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	return &TraderServer{
		host:             host,
		handler:          handler,
		processor:        processor,
		transport:        transport,
		transportFactory: transportFactory,
		protocolFactory:  protocolFactory,
		server:           server,
	}
}

func (ts *TraderServer) Run() {
	logger.Infoln("Thrift server listening on", ts.host)
	ts.server.Serve()
}

func (ts *TraderServer) Stop() {
	logger.Infoln("Thrift stopping server...")
	ts.server.Stop()
}
