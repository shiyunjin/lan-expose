# Lan Expose

[![GitHub Workflow Status](https://github.com/shiyunjin/lan-expose/actions/workflows/go.yml/badge.svg)](https://github.com/shiyunjin/lan-expose/actions/workflows/go.yml)
[![GitHub release](https://img.shields.io/github/tag/shiyunjin/lan-expose.svg?label=release)](https://github.com/shiyunjin/lan-expose/releases)

[README](README.md)

Lan Expose 是一个可以优雅的在被封禁 `443,80` 端口的情况下，使你和往常一样使用浏览器 (Chrome, Firefox, Edge) 
访问暴露到公网的网站， **无需指定端口号**. 仅支持 HTTPS。

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
 > README 文档正在撰写中 ***Draft**

## 快速使用
 > README 文档正在撰写中 ***Draft**

## 特性
 > README 文档正在撰写中 ***Draft**
### 配置文件

你可以在这里查看完整的配置文件和注解，来查看未在这里描述的所有功能。

[完整配置 - 在服务器搭建 (Upgrade)](./conf/upgrade.ini)

[完整配置 - 在局域网搭建 (Proxy)](./conf/proxy.ini)

### Websocket

#### Mode
 * Block
 * Proxy
 * 302

#### Block

#### Proxy

#### 302

### Check Page
