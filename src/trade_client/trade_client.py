import sys

sys.path.append('../lib')

from trade_service import TradeService
from trade_service.ttypes import *

from thrift import Thrift
from thrift.transport import TSocket
from thrift.transport import TTransport
from thrift.protocol import TBinaryProtocol


OKCoin_api_key     = ""
OKCoin_secret_key  = ""
Huobi_access_key   = ""
Huobi_secret_key   = ""
Chbtc_access_key   = ""
Chbtc_secret_key   = ""


try:

  # Make socket
  transport = TSocket.TSocket("127.0.0.1", 19090)

  # Buffering is critical. Raw sockets are very slow
  transport = TTransport.TFramedTransport(transport)

  # Wrap in a protocol
  protocol = TBinaryProtocol.TBinaryProtocol(transport)

  # Create a client to use the protocol encoder
  client = TradeService.Client(protocol)

  # Connect!
  transport.open()

  if True:
    configs = []
    hb_config = ExchangeConfig("huobi", Huobi_access_key, Huobi_secret_key)
    configs.append(hb_config)

    ok_config = ExchangeConfig()
    ok_config.exchange = "okcoin"
    ok_config.access_key = OKCoin_api_key
    ok_config.secret_key = OKCoin_secret_key
    configs.append(ok_config)

    ch_config = ExchangeConfig()
    ch_config.exchange = "chbtc"
    ch_config.access_key = Chbtc_access_key
    ch_config.secret_key = Chbtc_secret_key
    configs.append(ch_config)

    print("config->", configs)
    client.config(configs)

    accounts = client.get_account()
    print("accounts->", accounts)

    ticker= client.get_ticker()
    print("ticker->", ticker)

    for i in range(1, 2):
      ticker= client.get_ticker()
      print("ticker->", i, ticker)

      if True: 
        buyOrder = Trade("20", 5900)
        order= client.buy(buyOrder)
        print("buy->",order)

        order= client.get_order(order.id)
        print("buy order->",  order)

      if False: 
        sellOrder = Trade("2", 0.015)
        order= client.sell(sellOrder)
        print("Sell->",  order)

        order = client.get_order(order.id)
        print("Sell order->", order)

    alertOrders= client.get_alert_orders()
    print("errorOrders->", len(alertOrders))

  # Close!
  transport.close()

except Thrift.TException, tx:
  print '%s' % (tx.message)