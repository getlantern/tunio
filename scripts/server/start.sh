#!/bin/bash

badvpn-udpgw --listen-addr 0.0.0.0:5353 &

/usr/bin/lantern -headless -addr 0.0.0.0:2099
