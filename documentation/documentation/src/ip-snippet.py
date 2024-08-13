from ipaddress import ip_address, ip_network
from websocket import create_connection
...
ws.send("ip a\n")
ip_address = # ...

network = ip_network(f"{ip_address}/24", strict=False)
ip_list = [str(ip) for ip in list(network.hosts())[:10]]
ip_list = [ip for ip in ip_list if ip != ip_address and not ip.endswith(".1")]
