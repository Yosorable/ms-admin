# ms-admin
用户服务模块，主要功能:
1. 用户登录、注册等基本功能
2. 其它服务与用户相关的表的创建、表内记录的crud。表名根据本服务名与其它服务名和自定义的功能名称拼接。模型定义如下。

    * UserRecord
    
        浏览记录、收藏、点赞等功能，表的主键为用户id与其它服务的某个实体id。字段固定，分别为：用户id、某个实体id、创建时间、更新时间。

    * UserMultipleRecord (todo)
    
        评论、用户操作记录等功能，表的主键为自增int。字段不固定，但在创建表时确定，包含：主键id、用户id、某个实体id、创建时间、更新时间。
        