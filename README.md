core_filter
===========

control coredump in distribution system, avoid multi core dump make system unvaliable.

# 背景
一个服务下多台机器, 由于程序异常等原因可能导致程序批量coredump, 特别是对于php 这种大量进程或者是在运行期间需要开辟大内存的程序, 会产生大量的coredump. 大量的IO 会导致系统大量进程处于 D state, 系统假死. 


# 解决方案
1: 修改linux `/proc/sys/kernel/core_pattern`  截取 coredump 控制权.  从 /home/coresave/core.%e.%p.%t 修改为   

			  |/home/foo/go/src/core_filter/core_filter -p %p -e %e -t %t  
			 
2: 开发core_filter 模块来进行core 文件处理, 实现一台机器在一小时内只允许一个进程core dump 一次.  每次core dump 文件大小限制小于 1G, 对于core 大于 1G 只写入开始文件的1G 部分.
