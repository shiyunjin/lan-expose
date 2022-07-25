# Lan Expose

[![GitHub Workflow Status](https://github.com/shiyunjin/lan-expose/actions/workflows/go.yml/badge.svg)](https://github.com/shiyunjin/lan-expose/actions/workflows/go.yml)
[![GitHub release](https://img.shields.io/github/tag/shiyunjin/lan-expose.svg?label=release)](https://github.com/shiyunjin/lan-expose/releases)

[README](README.md)

Lan Expose 是一个可以优雅的在被封禁 `443,80` 端口的情况下，使你和往常一样使用浏览器 (Chrome, Firefox, Edge) 
访问暴露到公网的网站， **无需指定端口号**。 仅支持 HTTPS。

## 为什么使用 Lan Expose

通过 Lan Expose 公开网站可以在有公网IP情况下，但未开放 `443,80` 端口的情况下，你可以获得:
 * 提供和正常网站一样的访问体验
 * 不使用云服务器中转，节省了流量消耗 (几乎不产生流量)
 * 使用 `QUIC` 协议, 提升访问速度
 * 避免了 `SNI` 泄漏的问题 (可能用于判断IP是否绑定了域名)
 * 兼容了 Websocket 协议 (proxy, 302 模式)
 * 极低的内存使用 ≈ 4MB
 * 提供 Docker 和多种部署方式和平台支持

## 原理
Lan Expose 原理上依托于 `HTTP Alternative Services` **RFC7838** ，并提供一整套开箱即用的以 `Upgrade` 和 `Porxy` 服务为基础的解决方案。

`Upgrade` 将访问用户接收和完成SSL握手，并向其发送 `Alt-Svc` 请求，将其重定向到对应的 `Proxy` 服务。 `Proxy` 服务接受到用户请求后，使用配置文件路由将协议降级转换到目标服务器。

## 快速使用
目前可以在 Github 的 [Release](https://github.com/shiyunjin/lan-expose/releases) 页面中下载到最新版本的客户端和服务端二进制文件。你也可以在 [Actions](https://github.com/shiyunjin/lan-expose/actions) 下载到每个合并到主线版本的 `Commit` 版本。

我们也提供 Docker 镜像以方便部署 （仅当发布 Release 版本后，Docker镜像才会进行推送）
```
# Ghcr
#  - Proxy 
docker pull ghcr.io/shiyunjin/lan-expose-proxy:v0.1.0
#  - Upgrade
docker pull ghcr.io/shiyunjin/lan-expose-upgrade:v0.1.0
```

## 部署 Example
 > `Proxy` 需要部署在  未开放 `443,80` 端口 的局域网内，并确保外网端口映射配置正确（默认端口：690）
 > 
 > `Upgrade` 需要部署在 开放了 `80,443` 端口的服务器上，以提供正确的握手服务（可经过Nginx等服务中转）。


Proxy Docker 部署
```
#!/usr/bin/env bash

# setup sysctl max udp
echo "net.core.rmem_max=2500000" >> /etc/sysctl.conf
sysctl -p


## docker 镜像内不 含有 `conf` 配置文件，需要 `-v` 映射到镜像里启动 
## 并且请用 `-c <dir>/<file>.ini` 指定配置文件目录
# run
docker run -d --restart=always --name="lan-expose-proxy" \
  -v /volume1/docker/lan-expose-proxy/proxy.ini:/config/proxy.ini \
  -v /volume1/homes/syj/syno-acme/acme.sh/ **#YOUR DOMAIN#** /fullchain.cer:/config/ssl.crt \
  -v /volume1/homes/syj/syno-acme/acme.sh/ **#YOUR DOMAIN#** / **#YOUR DOMAIN#** .key:/config/ssl.key \
  -p 690:690/tcp \
  -p 690:690/udp \
  ghcr.io/shiyunjin/lan-expose-proxy:v0.1.0 \
  -c /config/proxy.ini

```

## 其他部署方案
更多部署方案可以查看[ 【run】 文件夹](run)内的脚本。

## 特性
 > README 文档正在撰写中 ***Draft**

### 配置文件

你可以在这里查看完整的配置文件和注解，来查看未在这里描述的所有功能。

[完整配置 - 在服务器搭建 (Upgrade)](./conf/upgrade.ini)

[完整配置 - 在局域网搭建 (Proxy)](./conf/proxy.ini)

### Websocket 兼容

由于 QUIC 推行 `WebTransport`，导致支持 `Websocket over HTTP/3`**RFC9220** 至今没有通过。
所以协议和主流客户端并没有支持直接进行连接，会自动降级到 `HTTP/2` 导致无法连接。

我提供了几种方式可以选择如何处理 `Websocket` 流量，以达到兼容的目的。

 * `Block` 阻止 `Websocket` 的访问请求 ***默认***
 * `Proxy` 通过服务器代理 `Websocket` 流量 (兼容性最好, 且不泄露 `SNI`，须消耗服务器流量)
 * `302`   重定向 `Websocket` 流量到直连地址，需要客户端支持 (存在泄漏 `SNI` 风险)

#### Block Mode

这是默认值，如无需使用 `Websocket`，请保持为这个值。 
这将会阻止所有尝试通过 `Websocket` 方式进行的请求。

#### Proxy Mode

 > 如果你需要使用 `Websocket`，推荐使用这个模式。

通过服务器转发 `Websocket` 流量，可以达到完美的兼容性。但是缺点也很明显，会消耗服务器的流量。
我自用的应用 `Websocket` 流量较小，其实不会消耗特别多。

#### 302 Mode

重定向 `Websocket` 流量到直连地址，无需消耗流量。看起来是这么的美好，并且符合 `RFC6455` 规范， 但其中有一句话。

 >   1.  If the status code received from the server is not 101, the
         client handles the response per HTTP [RFC2616] procedures.  In
         particular, the client might perform authentication if it
         receives a 401 status code; the server might redirect the client
         using a 3xx status code **(but clients are not required to follow
         them)**, etc.  Otherwise, proceed as follows.

***但客户端不需要遵循它们***，所以经过测试，绝大多数客户端并没有做兼容（包括 Chrome）。

事情不是绝对的，你可以很容易的自己完成对其兼容的适配。比如说：
 * 使用这个 [3p3r/websocket-redirect-shim](https://github.com/3p3r/websocket-redirect-shim) 包
 * 使用 `ws` 包，并将其中的 [followRedirects](https://github.com/websockets/ws/blob/d2c935a477fa6999c8fa85b89dfae27b85b807e7/doc/ws.md?plain=1#L272) 设为 `true`
 * [阿里云应用高可用服务 AHAS - WebSocket多活实践](https://help.aliyun.com/document_detail/188595.html) 中提到：
   
``` nodejs
您可以重点关注Client，以下示例采用NodeJS的WebSocket Library：

const WebSocket = require('ws');

let host = "http://websocket.msha.tech/";

let routerId = 1111;
// routerId = 6249;
routerId = 8330;

let options = {
    'headers': {
        routerId: routerId,
        unitType: "unit_type",
    }
};

let ws = handleWs();

function handleWs(){
    let ws = new WebSocket(host,[],options);
    ws.on('upgrade', function open(resp) {
        // console.log('upgrade ',resp);
    });

    ws.on('open', function open() {
        console.log('connected:'+routerId);
        ws.send(Date.now());
    });

    ws.on('error', function error(e) {
        console.log('err',e);
    });

    ws.on('close', function close() {
        console.log('disconnected');
    //断连后重连。
        let retryTime = 1500;
        setTimeout(()=>{
            console.log('!!! reconnecting in ...'+retryTime+' ms');
            ws = handleWs();
        }, retryTime);
    });

    ws.on('message', function incoming(data) {
        console.log(`msg: ${data} `);
    });

    ws.on('unexpected-response', function handleerr(req,resp) {
        //处理重定向。
    if ((resp.statusCode+'').startsWith("30")){
            console.log("!!! redirecting... from ", host," to",resp.headers.location);
            host = resp.headers.location;
            ws = handleWs();
        }

    });
    return ws;
}
```

### Check Page
> README 文档正在撰写中 ***Draft**
