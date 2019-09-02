# sshAutoRenew_txDomain
certbot-auto renew自动续期SSL证书


---
# 使用方法
## 1. 前期设置工作
  - 在sslAuto.sh设置好腾讯API的秘钥对
  - 安装好Golang环境
  - 编译腾讯api程序 ` go build -o txDomain ` (这个程序可以跨平台编译，具体架构和系统可以看我博客)
## 2. 准备开始
  - 把sslAuto.sh和编译好的txDomain，拷贝进certbot-auto路径(我的certbot-auto路径是/ssl/letsencrypt改成你自己的)

## 3. 执行以下命令即可
```
/ssl/letsencrypt/certbot-auto renew  --manual --preferred-challenges dns --manual-auth-hook "/ssl/letsencrypt/sslAuto.sh add" --manual-cleanup-hook "/ssl/letsencrypt/sslAuto.sh clean"
```
---
## 后续添加阿里云的自动脚本，我其他仓库里有阿里云的API调用代码，有能力的朋友稍微改一下也行


