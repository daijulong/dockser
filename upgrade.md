### 升级说明

#### v2.0.0

日期：2023-02-08

内容：

- [新] service 文件中支持 dockser 指令，用于支持 service 相关附加操作，如新增一个服务 xxx 时，需要将 xxx.conf 向 nginx 的配置文件目录中，使服务（xxx）与服务（nginx）之间解耦
- [改] 生成的最终的 docker-compose.yml 文件内容不再保持与 template 和 service 内容中相同的顺序，由字符串替换进阶为数据结构，以使以后有更多可能
- [删] 模板中不再支持 `@@_SERVICES_@@` 替换，service 用新的方式进行支持

#### v1.1.0

日期：2021-05-04

内容：

- [新] 支持通过 `.env` 文件中的 `DEFAULT_GROUP` 来指定默认分组，便于直接管理多个项目
- [新] 当子命令不存在时，将转交给 `docker-compose` 命令进行执行
- [改] 原 `docker-compose` 目录变更为 `compose`，以便于操作系统的自动提示
