#### easyquery 拥有通用的CRUD接口；特别是强大的查询功能，前端只要遵守查询约定，就能解放后端双手（基于GORM和GIN） 。



#### 安装:

```textile
go get github.com/JohnnyTing/easyquery
```



#### 前后端查询约定:

    `字段[操作符]=值`  (字段对应数据库表字段的lowerCamel)

1. 等于(eq)
   
   ```textile
   url：name[eq]=你好
   
   sql: where name = "你好"
   ```

2. 大于(gt)
   
   ```textile
   url: age[gt]=10
   
   sql: where age > 10
   ```

3. 大于等于(gteq)
   
   ```textile
   url: age[gteq]=10
   
   sql: where age >= 10
   ```

4. 小于(lt)
   
   ```textile
   url: age[lt]=10
   
   sql: where age < 10
   ```

5. 小于等于(lteq)
   
   ```textile
   url: age[lteq]=10
   
   sql: where age <= 10
   ```

6. 模糊查询(like)
   
   ```textile
   url: name[like]="你"
   
   sql: where name like %你%
   ```

7. 查询多个值(in)
   
   ```textile
   url: id[in]=1&id[in]=2
   
   sql: where id in (1,2)
   ```

8. 非(not)
   
   ```textile
   url: id[not]=1
   
   sql: where id <> 1
   ```

9. Null值(is_null)
   
   ```textile
   url: name[is_null]=true
   
   sql: where name is null
   ```

10. 空值(is_empty)
    
    ```textile
    url: name[is_empty]=true
    
    sql: name is null or trim(name) = ''
    ```

11. 非空(not_null)
    
    ```textile
    url: name[not_null]=true
    
    sql: where name is not null and trim(name) != ''
    ```

12. 特殊判断非空(s_not_null)
    
    ```textile
    url: name[s_not_null]=true
    
    sql: where name is not null and trim(name) != '' and name != '无' and name != '不涉及'
    ```

13. 日期date查询(date_gt、date_gteq、date_lt、date_lteq)
    
    ```textile
    url: created_at[date_gt]=2021-04-11
    
    sql: Date(created_at) >= 2021-04-11
    ```

14. 内部or查询(or_in_eq)
    
    ```textile
    url: name[or_in_eq]=dd&name[or_in_eq]=ss&age[eq]=10
    
    sql: where (name = dd or name = ss) and age = 10
    ```

15. 外部or查询(or_out_eq)
    
    ```textile
    url: name[or_out_eq]=dd&name[or_out_eq]=ss&age[eq]=10
    
    sql: where (name = dd or name = ss) or age = 10
    ```

16. 不在集合值里面(not_in)
    
    ```textile
    url: id[not_in]=1&id[not_in]=2
    
    sql: where id not in (1,2)
    ```

17. 排序(order)
    
    ```textile
    url: id[order]=desc&mobile[order]=asc
    
    sql: order by id desc, mobile asc
    ```

18. 分组(group)
    
    ```textile
    url: /users/group?userName[group]=true
    
    sql: SELECT user_name as label, count(1) as value FROM "users" WHERE "users"."deleted_at" IS NULL GROUP BY "user_name" ORDER BY value desc
    
    
    # join字段分组
    
    url: /users/group?company_j_name[group]=true
    
    sql:  SELECT "Company"."name" as label, count(1) as value FROM "users" LEFT JOIN "companies" "Company" ON "users"."company_id" = "Company"."id" WHERE "users"."deleted_at" IS NULL GROUP BY "Company"."name" ORDER BY value desc
    ```

19. Join连接查询(join)，仅适用于一对一的关系(has one, belongs to)
    
    ```textile
    url: /users/?company_j_name[eq]=qylz
    
    sql: SELECT * FROM "users" LEFT JOIN "companies" "Company" ON "users"."company_id" = "Company"."id" WHERE "Company"."name" = 'qylz' AND "users"."deleted_at" IS NULL
    
    
    url: /users?company_j_id[in]=2&company_j_id[in]=1
    
    sql: SELECT * FROM "users" LEFT JOIN "companies" "Company" ON "users"."company_id" = "Company"."id" WHERE "Company"."id" IN ('2','1') AND "users"."deleted_at" IS NULL
    ```
    
    说明：
    
    ```textile
    company_j_name 以_j_为分隔符，company为关联对象，name为关联对象字段。
    user属于一个company，需要在User结构体添加CompanyID外键及Company对象。
    并在User结构体添加Joins方法，如有多个关联则添加多个关联对象，具体看examples/pkg/user/user.go：
    
    func (user *User) Joins() []interface{} {
     return []interface{}{Company{}}
    }
    
    ```

17. 支持预加载(Preload)
    
    ```textile
    在User结构体添加Preload方法：
    
    func (user *User) Preload() []string {
        return []string{clause.Associations}
    }
    
    clause.Associations 代表预加载全部，如果想自定义要加载的对象，修改即可，比如:
    
    func (user *User) Preload() []string {
        []string{"Company"}
    }
    ```

18. 自定义查询操作符，仅支持单值，不支持数组(in)、无值(is_null, not_null，order等)
    
    ```go
    # 定义的操作符为lowerCamel
    
    easyquery.Clause["Dayu"] = func(field string) string {
    	return fmt.Sprintf("%s > ?", field)
    }
    
    url: id[dayu]=1
    
    sql: WHERE id > '1' 
    ```

19. 组合查询
    
    ```textile
    url: /users/?id[gteq]=1&loginName[eq]=dingxu&company_j_name[like]=qy
    
    
    sql: SELECT * FROM "users" LEFT JOIN "companies" "Company" ON "users"."company_id" = "Company"."id" WHERE "users"."id" >= '1' AND "users"."login_name" = 'dingxu' AND "Company"."name" like '%qy%' AND "users"."deleted_at" IS NULL ORDER BY "users"."id" desc LIMIT 10
    
    ```



#### 用法 [完整例子](https://github.com/JohnnyTing/easyquery/tree/master/examples)

1. 定义模型
   
   ```go
   package user
   
   import (
   	"gorm.io/gorm/clause"
   
   	"gorm.io/gorm"
   )
   
   // 模型
   var (
   	Model  = &User{}
   	Models = &[]User{}
   )
   
   type User struct {
   	gorm.Model
   	LoginName *string
   	UserName  *string
   	Mobile    *int32
   	Password  *string
   	Gender    *string
   	CompanyID uint
   	Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
   	Roles     []Role  `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
   }
   
   // 设置预加载对象
   func (user *User) Preload() []string {
   	return []string{clause.Associations}
   }
   
   // 设置join查询对象，join查询时需要
   func (user *User) Joins() []interface{} {
   	return []interface{}{Company{}}
   }
   
   ```

2. 定义service
   
   ```go
   package user
   
   import (
   	"easyquery"
   	"easyquery/examples/pkg/db"
   
   	"github.com/gin-gonic/gin"
   	"gorm.io/gorm"
   
   	"golang.org/x/crypto/bcrypt"
   )
   
   var UserCrudService *UserService
   
   // 继承easyquery Crud接口
   type UserService struct {
   	easyquery.Crud
   }
   
   func init() {
   	UserCrudService = NewDefaultUserService()
   }
   
   func NewUserService(crud easyquery.Crud) *UserService {
   	return &UserService{crud}
   }
   
   // 需要传crud模型与数据库连接gorm.DB
   func NewDefaultUserService() *UserService {
   	return NewUserService(easyquery.NewCrudService(Model, Models, retreiveGormDB))
   }
   
   // 数据库连接，用户自定义postgresql、mysql等
   func retreiveGormDB() *gorm.DB {
   	return db.Postgres
   }
   
   ```

3. 定义handler
   
   ```go
   package user
   
   import (
   	"easyquery"
   	"easyquery/tools/stringutil"
   
   	"github.com/gin-gonic/gin"
   )
   
   // 继承easyquery handler
   type UserHandler struct {
   	easyquery.BaseHandler
   }
   
   func (handler *UserHandler) List(c *gin.Context) {
   	var models []User
   	err := UserCrudService.List(&models, handler.Transform(c, Model))
   	handler.HandleList(c, &models, err)
   }
   
   func (handler *UserHandler) Group(c *gin.Context) {
   	var models []easyquery.GroupVO
   	err := UserCrudService.Group(&models, handler.Transform(c, Model))
   	handler.Handle(c, &models, err)
   }
   
   func (handler *UserHandler) Create(c *gin.Context) {
   	var model User
   	c.ShouldBind(&model)
   	err := UserCrudService.Create(&model)
   	handler.Handle(c, &model, err)
   }
   
   func (handler *UserHandler) Update(c *gin.Context) {
   	var model User
   	c.ShouldBind(&model)
   	model.ID = stringutil.Str2Uint(c.Param("id"))
   	err := UserCrudService.Update(&model)
   	handler.Handle(c, &model, err)
   }
   
   func (handler *UserHandler) Find(c *gin.Context) {
   	model := User{}
   	err := UserCrudService.Find(&model, handler.Transform(c, Model))
   	handler.Handle(c, &model, err)
   }
   
   func (handler *UserHandler) Delete(c *gin.Context) {
   	model := User{}
   	model.ID = stringutil.Str2Uint(c.Param("id"))
   	err := UserCrudService.Delete(&model)
   	handler.Handle(c, &model, err)
   }
   ```



注意：如果实际项目没有用gin，自己实现QueryParamer 接口即可，具体数据格式参考handler.go、queries.go、pagination.go

```go
type QueryParamer interface {
    GetFields() []*QueryField
    GetPagination() Paginater
    GetJoin() bool
}
```
