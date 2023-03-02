# 已实现的功能
根据数据库和模版自动生成go文件

根据数据库生成graphqls文件

可以通过参数选择表格、数据库、模板、模版和输出路径

# 存在的问题
gqlgen配置文件固定了路径，是否需要随项目更改

gqlgen版本兼容性问题


# 安装说明
go1.16
```
go get github.com/akari-yd/codegen
```

go1.17及以上 
```
go get -b github.com/akari-yd/codegen

go install github.com/akari-yd/codegen
```

# 使用说明
```makefile
codegen --help 查看帮助
```

样例
```makefile
codegen --sqltype=mysql \
-c "username:password@(ip:port)/dbname" \
-d dbName \
-t "tabName1,tabName2" \
-T graph
```
