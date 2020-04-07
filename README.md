一个征信录入系统

fabric images v1.0.0

Gopkg.toml里面有fabric-sdk-go的版本，按照版本下载即可.

fixtures是超级帐本的配置文件：
     进入cd ./fixtures
     docker-compose up来启动超级帐本镜像。
     cd ../IdentitySystem

config.yaml是fabric-sdk-go配置文件:
     一个组织两个节点：orgid:OrgIdentityMSP
     通道：identitychannel

通道策略：使用fabric1.0默认的策略

sdk背书的策略：OrgIdentityMSP.member
