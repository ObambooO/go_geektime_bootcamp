# Web核心-Server

从特性上看，至少要提供三部分功能

1、生命周期控制：即启动、关闭

2、路由注册接口：提供路由注册功能

3、作为http包到web框架的桥梁

实现的路由树不是线程安全的。要求用户必须要注册完路由才能启动HTTPServer，而正常用法都是在启动之依次注册路由，不存在并发的场景。
至于运行期间动态注册路由，没必要支持。这是典型的为了解决1%的问题，引入99%的代码

路由树查找的性能受限于路由树的高度（深度），其次是路由树的宽度


作业：
加强路由树功能：
通配符匹配：修改为，如果通配符出现在路由的末尾，例如:/a/b/*，那么它能够匹配到后面多段路由，例如/a/b/c/d/e/f,而且前面支持/a/b/c

正则匹配：如，开发者注册一个正则匹配路由/user/:id(^[0-9]+$)，那么这种情况下，用户输入/user/123，可以匹配这个路由，
而/user/xionger则无法匹配