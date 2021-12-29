client
dev tun
proto udp
remote ${endpoint} 443
remote-random-hostname
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
cipher AES-256-GCM
verb 3

dhcp-option DNS "${vpc_dns_resolver_ip}"
dhcp-option DOMAIN "${dns_domain}"
register-dns
block-outside-dns

ca ${environment}-ca.crt
cert ${environment}-client.crt
key ${environment}-client.key

reneg-sec 0
