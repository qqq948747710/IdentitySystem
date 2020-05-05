一个征信录入系统

fabric images v1.0.0

Gopkg.toml里面有fabric-sdk-go的版本，按照版本下载即可.

fixtures是超级帐本的配置文件：
     进入cd ./fixtures
     docker-compose up来启动超级帐本镜像。
     cd ../IdentitySystem

![1](https://raw.githubusercontent.com/qqq948747710/IdentitySystem/master/images/1.JPG)

config.yaml是fabric-sdk-go配置文件:
     一个组织两个节点：orgid:OrgIdentityMSP
     通道：identitychannel

通道策略：使用fabric1.0默认的策略

sdk背书的策略：OrgIdentityMSP.member

关于添加信息：只有管理员才能添加信息，具体修改方法可以在mysql的user把isadmin字段修改为0

![2](https://raw.githubusercontent.com/qqq948747710/IdentitySystem/master/images/2.JPG)

![3](https://raw.githubusercontent.com/qqq948747710/IdentitySystem/master/images/3.JPG)

![4](https://raw.githubusercontent.com/qqq948747710/IdentitySystem/master/images/4.JPG)

![5](https://raw.githubusercontent.com/qqq948747710/IdentitySystem/master/images/5.JPG)

![6](https://raw.githubusercontent.com/qqq948747710/IdentitySystem/master/images/6.JPG)

![7](https://raw.githubusercontent.com/qqq948747710/IdentitySystem/master/images/7.JPG)
