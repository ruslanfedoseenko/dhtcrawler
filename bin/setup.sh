#!/usr/bin/env bash
iptables -t raw -I PREROUTING -p tcp --dport 9999 -j NOTRACK
iptables -t raw -I PREROUTING -p tcp --sport 9999 -j NOTRACK
iptables -t raw -I PREROUTING -p tcp --sport 5432 -j NOTRACK
iptables -t raw -I PREROUTING -p tcp --dport 5432 -j NOTRACK
iptables -t raw -I PREROUTING -p tcp --dport 9200 -j NOTRACK
iptables -t raw -I PREROUTING -p tcp --sport 9200 -j NOTRACK
iptables -t raw -I PREROUTING -p udp --dport 9999 -j NOTRACK
iptables -t raw -I PREROUTING -p udp --sport 9999 -j NOTRACK
iptables -t raw -I PREROUTING -p udp --sport 5432 -j NOTRACK
iptables -t raw -I PREROUTING -p udp --dport 5432 -j NOTRACK
iptables -F