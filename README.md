# README

感谢[阿里云的容器镜像服务](https://help.aliyun.com/document_detail/64340.html?spm=a2c4g.11186623.6.550.704d33deS6pChu)为每个租户提供的免费额度.
可以绕过一些极端的网络问题.

## Install and Usage
``go get github.com/j2gg0s/gproxy/cmd/gproxy``

使用默认的 registry.cn-huhehaote.aliyuncs.com/gproxy
``gproxy --source alpine:3``

指定具体的 acr repo
``gproxy --source alpine:3 --dest registry.cn-hangzhou.aliyuncs.com/xxx/xxx --username xxx --password xxx``
