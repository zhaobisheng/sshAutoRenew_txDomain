# sshAutoRenew_txDomain
certbot-auto renew自动续期SSL证书



# 使用方法
## 在sslAuto.sh设置好腾讯API的秘钥对
## 编译腾讯api程序` go build -o txDomain ` (这个程序可以也调用腾讯其他的api，具体看代码里面的示例)
## 执行以下命令即可(我的certbot-auto路径是/ssl/letsencrypt改成你自己的)
```
/ssl/letsencrypt/certbot-auto renew  --manual --preferred-challenges dns --manual-auth-hook "/ssl/letsencrypt/sslAuto.sh add" --manual-cleanup-hook "/ssl/letsencrypt/sslAuto.sh clean"
```

## 后续添加阿里云的自动脚本，我其他仓库里有阿里云的API调用代码，有能力的朋友稍微改一下也行


