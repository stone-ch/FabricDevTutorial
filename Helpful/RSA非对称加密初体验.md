## 概述

区块链离不开非对称加密，本文章将简单介绍一下非对称加密的一些基本概念，体验密钥的生成和加密解密的过程。

## 基本概念

- PKCS: The Public-Key Cryptography Standards (PKCS)是由美国RSA数据安全公司及其合作伙伴制定的一组公钥密码学标准。

- x.509: 是密码学里公钥证书的格式标准。

- RSA: 一种非对称加密算法。1977年由罗纳德·李维斯特（Ron Rivest）、阿迪·萨莫尔（Adi Shamir）和伦纳德·阿德曼（Leonard Adleman）一起提出。RSA就是他们三人姓氏开头字母拼在一起组成的。

- openssl: 是目前最流行的 SSL 密码库工具，其提供了一个通用、健壮、功能完备的工具套件，用以支持SSL/TLS 协议的实现。

## 加解密初体验

生成一个密钥：
```
$ openssl genrsa -out test.key 1024
```

根据密钥生成公钥:
```
$ openssl rsa -in test.key -pubout -out test_pub.key
```

随意创建一个hello的文本文件，写入一些字符。然后加密该文件：
```
$ echo "Hello RSA！" > hello.txt
$ openssl rsautl -encrypt -in hello.txt -inkey test_pub.key -pubin -out hello.en
```

解密文件并查看：
```
$ openssl rsautl -decrypt -in hello.en -inkey test.key -out hello.de
$ cat hello.de
```

正常情况下，解密结果应该和原文件一致。

参考文献：
1. https://blog.csdn.net/ly0303521/article/details/53391741
2. https://www.cnblogs.com/littleatp/p/5878763.html
