python3 -m http.server 5678 > server-5678.log 2>&1 &
python3 -m http.server 5679 > server-5679.log 2>&1 &
su root
iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner $(id admin -u) -m multiport --dport 5678 -j REDIRECT --to-port 15001
iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner $(id admin -u) -m multiport --dport 5679 -j REDIRECT --to-port 15006
curl 0.0.0.0:5678 
curl 0.0.0.0:5679
