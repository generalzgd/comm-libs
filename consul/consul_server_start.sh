
hup 
./consul agent -server -bootstrap-expect=1 -data-dir=/tmp/consul/ -node=service-center -bind=192.168.163.184 -client=0.0.0.0 -ui 
-config-dir=/etc/consul.d/ > consul.log 2>&1 &
