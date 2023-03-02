# 代码提交前代码检查使用手册
## 介绍
该部分用于提交前的代码检查，检查工具包括golangci-lint静态代码检查工具和基于语法树检查的go程序。
操作指令集成在[makefile](makefile)中，可包含多个部分，如需要对部分模块进行精细化检查，可新增指令集，并在all中添加即可在提交前进行自动检查。

## golangci-lint
该静态代码检查工具集成多种检查器，使用golangci-lint run运行代码检查，详情可见[makefile](makefile)。
检查器可通过配置文件.golangci.yml进行配置，具体配置方式参考[配置文件](tools/.golangci.yml)，也可参考[官方文档](https://golangci-lint.run/usage/linters/).

## 语法树检查
见[go程序](ast.go)，程序使用官方抽象语法树包"go/ast",该包提供了抽象语法树的相关结构，具体使用可参考[中文开发手册](https://www.php.cn/manual/view/35199.html).
基于语法树，我们可以自定义代码检查工具。

### checkFunc(sourceDir,sourceFun, targetDir,targetFun string)
该接口用于检查sourceDir路径下函数名匹配sourceFun正则表达式的函数是否被targetDir路径下函数名匹配targetFun正则表达式的函数匹配，
一般用于检查特定类型的函数是否编写相应Test函数进行自测。 可在go程序中的main函数中修改相应的参数进行更改。