### Consul部署

1. 完成Consul安装后，必须运行agent。agent可以运行为server或client模式。每个数据中心必须拥有一台server。建议在一个集群中有3或5个server。部署单一的server，在出现失败时会不可避免的造成数据丢失。
2. 其他agent运行为client模式。一个client是一个非常轻量级的进程，用于注册服务，运行健康监测和转发对server的查询。
3. agent必须在集群中的每个主机上运行。

#### Consul目录需求

1. linux: /etc/consul/data  用户存储数据；/etc/consul/conf  用户存储服务配置文件
2. win: C:/etc/consul/data  用户存储数据；C:/etc/consul/conf  用户存储服务配置文件；由于consul不识别其他磁盘，只能放C盘下

#### 启动Consul Server

实际应用中，部署3个server。

数据中心统一用 dc1

##### beta server1(10.62.62.25)

```
!#consul config file: /etc/consul/master.json
{
  "server":true,
  "ui":true,
  "datacenter":"dc1",
  "data_dir":"/etc/consul/data",
  "bootstrap_expect":3,
  "node_name":"s1_25",
  "bind_addr":"10.62.62.25",
  "client_addr":"0.0.0.0"
}
//command
consul agent -config-file=/etc/consul/master.json -config-dir=/etc/consul/conf -rejoin
```

- -server：定义agent运行server模式

- -bootstrap-expect：在一个datacenter中期望提供的server节点数，当该值提供的时候，consul一直等到达到指定server数量的时候才会引导整个集群（不能和bootstrap公用）

- -bind (bind_addr)：改地址用来集群内部的通讯，集群内的所有节点到地址都必须的可达的，默认0.0.0.0

- -node (node_name)：节点在集群中的名称，在一个集群中必须是唯一的，默认的该节点的主机名

- -ui：启用UI界面

- -rejoin：使consul忽略先前的离开，在再次启动后仍旧尝试加入集群中

- -client (client_addr)：consul服务侦听地址，这个地址提供http,DNS,RPC等服务，默认是127.0.0.1不对外，如果要对外需改成0.0.0.0

  ​

##### beta server2(10.62.62.26)

```
!#consul config file: /etc/consul/master.json
{
  "server":true,
  "ui":true,
  "datacenter":"dc1",
  "data_dir":"/etc/consul/data",
  "bootstrap_expect":3,
  "node_name":"s2_26",
  "bind_addr":"10.62.62.26",
  "client_addr":"0.0.0.0"
}
consul agent -config-file=/etc/consul/master.json -config-dir=/etc/consul/conf -rejoin -join 10.62.62.25
```

- -join：加入集群

##### beta server3(10.62.62.27)

```
!#consul config file: /etc/consul/master.json
{
  "server":true,
  "ui":true,
  "datacenter":"dc1",
  "data_dir":"/etc/consul/data",
  "bootstrap_expect":3,
  "node_name":"s3_27",
  "bind_addr":"10.62.62.27",
  "client_addr":"0.0.0.0"
}
consul agent -config-file=/etc/consul/master.json -config-dir=/etc/consul/conf -rejoin -join 10.62.62.25
```

##### beta client(10.62.62.28)

```
!#consul config file: /etc/consul/client.json
{
  "datacenter":"dc1",
  "data_dir":"/etc/consul/data",
  "node_name":"c1_28",
  "bind_addr":"10.62.62.28",
  "client_addr":"0.0.0.0"
}
consul agent -config-file=/etc/consul/acl_client.json -config-dir=/etc/consul/conf -join 10.62.62.25
```

