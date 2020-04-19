
hup ./consul agent -data-dir /tmp/consul/ -bind=192.168.163.93 -client 0.0.0.0 -ui -config-dir=/etc/consul.d/ -join 192.168.163.184 > consul.log 2>&1 &
