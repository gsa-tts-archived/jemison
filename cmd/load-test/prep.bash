!/bin/bash

# For this shell only
sudo sysctl -w net.ipv4.ip_local_port_range="1024 65535"
sudo sysctl -w net.ipv4.tcp_tw_reuse=1
sudo sysctl -w net.ipv4.tcp_timestamps=1
# For the whole system, until reboot
sudo ulimit -n 250000
