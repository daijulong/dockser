# Dockser

```
     ________                  ___
    /  ___   \                /  / ___ 
   /  /   |  /_____  ______  /  /_/  /_____________  ______
  /  /   /  /  __  \/  ____\/     __/ ______/ ___  \/  ___/
 /  /___/  /  /__/ /  /__ _   /\  \/_____  \  _____/  /
/_________/\______/\______/__/  \_/________/\_____/__/
```

[升级说明](upgrade.md)

## 更方便灵活地管理 docker-compose.yml

你不必再从冗长的 `docker-compose.yml` 中苦苦追寻，`dockser` 已经为你将 `services` 拆散，让你可以灵活组装。

你可以将每个 `service` 都放到单独的 `yml` 文件，做为最基本的“零件”，用时按需组装即可：

- 想增加一个 `service`，你要做的仅仅是在 `group` 中增加一个名字，而已。移除亦然；
- 你可将一些常用的组合放到一个 `yml` 文件中，比如你只需要在名单中加入 `elk`，就能把 `Elasticsearch` + `Logstash` + `Kibana` 这 3 个服务的黄金组合一起加入；
- 想要只管理一套可用于多个项目的 `docker` 开发环境，却因 `yaml` 中的 `key` 无法使用环境变量而束手，`dockser` 为你解忧。

## 安装前准备

dockser 是为了辅助 `docker-compose` 而存在，所以请先确保系统中已经安装了 `docker-compose`。

## 安装

下载地址：https://github.com/daijulong/dockser/releases

先下载对应平台的二进制程序。

### Linux & MacOS

移动到 path 目录

```
sudo mv dockser /usr/local/bin/dockser
```

加上可执行权限 ：

```sudo chmod +x /usr/local/bin/dockser```

可以在任意位置执行：

```dockser -v```

### Windows

访问 https://github.com/daijulong/dockser/releases 并下载 dockser.exe 文件，放到 path 文件夹即可。

## 使用说明

### 帮助手册

```
dockser
```

或在每个命令后加上 `-h` 就会显示当前命令的帮助信息。例如：

```
dockser make -h
```

建议必要时看一下帮助信息，会有本文档中无法兼顾的细节信息。

### 初始化

执行命令：

```
dockser init
```

你将在当前目录中得到以下目录和文件：

```
compose/              存放 docker-compose 相关配置文件 
  |- services/               存放 service 文件，一般每个 service 一个文件
    |- nginx.yml             示例 service 文件
  |- templates/              存放 docker-compose.yml 模板
    |- docker-compose.yml    示例 docker-compose.yml 模板文件
  |- groups.yml              分组配置文件，默认会有一个 default 分组
.env 环境变量
.env.example  环境变量示例
```

如果加上 `--with-demo`，你将得到额外得到 `compose/templates/docker-compose-demo.yml` 模板文件，并且 `compose/groups.yml` 中也会多一个 `demo` 的分组配置。

### 生成 docker-compose.yml 文件

执行命令：

```
dockser make
```

此命令会根据 groups.yml 中的指定分组配置，生成对应的 docker-compose.yml 文件。默认分组为 `default`，如果需要使用其他分组配置，则加上分组名即可：

```
dockser make demo
```

### 其他命令

当子命令不存在时，将会直接调用系统的 `docker-compose` 来执行子命令，如 `dockser ps -a`，因为 `ps` 不是 `dockser` 的子命令，所以会直接执行 `docker-compose ps -a` 命令。

> 目前仅支持同步输出的命令，交互式命令（如 `dockser exec xxx sh`）暂不支持。

## 配置详解

### service

service 等同于常规 docker-compose.yml 中 services 下的配置。

你可以考虑将你平时在 docker-compose.yml 中的每个 service 都拆分成独立的文件并放到此处，不再需要到冗长的文件中查找或修改，管理起来会十分方便。

对于一些黄金组合，例如 `ELK`，也可以放到同一个 yml 中。总之，service 可以任意拆分组合到各文件中。

### 附加指令

各服务之间可能会有耦合，例如使用 nginx 时，其他服务可能需要由其进行反代，例如 blog 服务如果被使用，则应向 nginx 中新增一个类似于 blog.xxx.com.conf 的配置文件。为了解决这个问题，在 2.0 版本开始，使用附加指令进行处理。

附加指令依托于 service 的定义文件，需要定义一个 dockser 的数据，以下为 blog/blog.yml 文件的内容：

```
blog:
  ...... //此处省略具体内容
dockser:
  add:
    copy:
      - ./components/blog/nginx.conf:./components/nginx/etc/nginx/templates/blog.xxx.com.conf.template:override
    #remove:
    #  - ./components/nginx/etc/nginx/templates/blog.xxx.com.conf.template
  remove:
    remove:
      - ./components/nginx/etc/nginx/templates/blog.xxx.com.conf.template
```

在此 service 文件中，除了声明了 blog 这个实际 service 内容外，还有一个 dockser 的声明，专门用于声明要执行的附加指令，dockser 作为保留名称，且不会作为 service 被编译到 docker-compose.yml 文件中。

工具预设了 add 和 remove 两个时机，分别对应服务被添加和移除时需要执行的指令，目前仅实现了 add ，remove 将在后续版本中实现。

当前支持的附加指令有：
- copy：复制文件，每个文件声明包括 3 段，由”:“分隔，格式为：”源文件名:目标文件名:是否覆盖“
  - 源文件名：复制哪个文件
  - 目标文件名：复制到哪里，具体的文件名
  - 是否覆盖：可以省略，默认为不覆盖，如果每次 make 时都需要覆盖一次文件，则需要指定为固定值”override“
- remove：删除文件

### 模板

模板是生成 docker-compose.yml 的依据。在执行 `make` 命令时，即是根据某个模板，替换掉特定内容（占位符或环境变量）后，生成最终的 docker-compose.yml 文件。

### 分组配置

分组配置文件 compose/groups.yml 是本工具的核心配置，以下示例内容定义了非常常见的 php 开发环境：

```
lnmp:
  services:
    - nginx
    - php
    - mysql
    - redis
  template: "docker-compose-lnmp.yml"
  output: "docker-compose.yml"
  override: "auto"
```

- 分组名：lnmp
- services: 组合 4 个 service，对应 compose/services 下应该有 nginx.yml,php.yml,mysql.yml,redis.yml 这 4 个文件。
- template: 使用的模板文件是 docker-compose/templates/docker-compose-lnmp.yml
- output：生成文件名，将在当前目录下生成 docker-compose.yml 文件
- override：输出文件覆盖机制。设置为 auto 时，如果不存在 output 文件，直接生成，如果存在，会在文件名后面增加一个日期时间，不会覆盖原有文件内容。设置为 force 时，将会强行覆盖已存在的文件内容。默认为 auto。

其中 `template` 和 `output` 可以在 `make` 命令中重新指定，例如：

```
dockser make lnmp -tpl=docker-compose.yml -out=docker-compose-php.yml
```

### 变量替换

下面是一个极为简单的模板内容示例：

```
version: "3"
networks:
  #@_NETWORK_@#:
      driver: bridge
```

`#@_NETWORK_@#` 在生成文件时，会被替换为 .env 文件中定义的环境变量“NETWORK”。如果要使用环境变量，只需要在 .env 文件中定义，并在文件中以 `#@_` + `环境变量名` + `_@#` 的格式占位即可。在模板文件和服务文件中都可使用。这也就解决了 yaml 文件中 key 不能使用变量来灵活设置的问题。

模板中也可以声明有效的 service，且会与 `group` 中声明的 service 内容合并，但一般不建议这样做。

## 其他

`dockser` 以后将会提供更多的特性和功能，欢迎大家提出宝贵的意见和建议。

联系方式：

- 邮箱：daijulong@qq.com
- 微信：julongdai
- QQ：88622090