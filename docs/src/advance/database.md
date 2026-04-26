# 数据库 (Database)

## 开篇故事

想象你在经营一家书店。最开始，你用笔记本记录：

```
Ada - 买了《Go 编程》- 35.9 元
Grace - 买了《Go Web》- 35.9 元，买了《GORM》- 42.9 元
```

生意好了之后，问题出现了：
- 如何快速找到某个顾客的所有订单？
- 如果顾客退货，如何确保订单和库存同时更新？
- 如果收银机在结账时死机了，钱收了但订单没记录怎么办？

**数据库**就是你的"数字化账本"，**ORM**（GORM）就是"自动记账员"，**事务**就是"原子操作保证"——要么全成功，要么全失败，不会出现"收了钱没给货"的情况。

这一章教你用 Go 和 GORM 构建可靠的数据层，从最简单的增删改查到复杂的事务处理。

## 本章适合谁

- ✅ 写过 `db.Query()` 但觉得原始 SQL 繁琐的开发者
- ✅ 想用 ORM 简化数据库操作，但不知道 GORM 如何上手
- ✅ 需要理解"一对多"关系如何建模（用户-订单、文章-评论）
- ✅ 遇到"部分写入成功"导致数据不一致，想了解事务的使用场景

如果你曾经为"如何确保两步数据库操作要么都成功，要么都失败"而困惑，本章必读。

## 你会学到什么

完成本章后，你将能够：

1. **定义 GORM 模型**：用结构体和标签映射数据库表，理解主键、外键、约束
2. **执行自动迁移**：用 `AutoMigrate` 同步模型到数据库表结构
3. **完成 CRUD 操作**：创建、读取、更新、删除记录，理解返回值和错误处理
4. **处理一对多关系**：用 `Preload` 预加载关联数据，避免 N+1 查询问题
5. **使用事务保护一致性**：用 `Transaction` 包装多步操作，确保原子性

## 前置要求

在开始之前，请确保你已掌握：

- Go 结构体（struct）和指针（pointer）语法
- 错误处理模式（`if err != nil`）
- 数据库基础概念（表、字段、主键、外键）
- SQL 基础（SELECT、INSERT、UPDATE、DELETE）

如果不熟悉 SQL，建议先花 30 分钟了解基本概念，但本章不要求手写复杂 SQL。

## 第一个例子

让我们从一个最简单的场景开始：创建一个用户并保存到数据库。

```go
package main

import (
    "fmt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type User struct {
    ID    uint
    Name  string
    Email string
}

func main() {
    // 1. 连接到内存数据库
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    
    // 2. 自动迁移（创建表）
    db.AutoMigrate(&User{})
    
    // 3. 创建用户
    user := User{Name: "Alice", Email: "alice@example.com"}
    db.Create(&user)
    
    // 4. 打印结果（注意 ID 被自动填充）
    fmt.Printf("用户 ID=%d 姓名=%s\n", user.ID, user.Name)
    
    // 5. 查询用户
    var loaded User
    db.First(&loaded, user.ID)
    fmt.Printf("查询结果：%+v\n", loaded)
}
```

**运行结果**：
```
用户 ID=1 姓名=Alice
查询结果：{ID:1 Name:Alice Email:alice@example.com}
```

**关键点**：
- `:memory:` 表示内存数据库，程序退出后数据消失（适合测试）
- `AutoMigrate` 自动创建 `users` 表，包含 `id`、`name`、`email` 字段
- `Create(&user)` 执行后，`user.ID` 被自动填充为数据库生成的主键
- `First(&loaded, id)` 按主键查询，找不到会返回 `ErrRecordNotFound`

## 原理解析

### 1. ORM 的核心思想

**ORM（Object Relational Mapping）** = 对象关系映射。

**通俗理解**：ORM 是一个"翻译官"，它在你的 Go 代码和数据库 SQL 之间翻译：

```
Go 代码              ORM 翻译              SQL 语句
---------            -----------           ---------
db.Create(&user)  →   翻译   →   INSERT INTO users (name, email) VALUES (?, ?)
db.First(&u, 1)   →   翻译   →   SELECT * FROM users WHERE id = ? LIMIT 1
db.Delete(&user)  →   翻译   →   DELETE FROM users WHERE id = ?
```

**好处**：
- 不写 SQL，用 Go 语法操作数据库
- 类型安全（编译时检查字段名）
- 自动处理主键、时间戳等细节

**代价**：
- 复杂查询不如原生 SQL 灵活
- 性能开销（反射、翻译层）
- 需要学习 ORM 的"坑"（如 N+1 查询）

### 2. GORM 模型定义

GORM 通过结构体字段推导数据库列：

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`  // 主键
    Name      string    `gorm:"size:255"`    // VARCHAR(255)
    Email     string    `gorm:"uniqueIndex"` // 唯一索引
    CreatedAt time.Time // 自动管理创建时间
    UpdatedAt time.Time // 自动管理更新时间
}
```

**标签（Tags）的作用**：
- `primaryKey`：标记主键字段
- `size:255`：指定字符串长度
- `uniqueIndex`：创建唯一索引
- `autoIncrement`：自增（uint 类型默认）

**外键和关联**：
```go
type Order struct {
    ID        uint
    UserID    uint           `gorm:"index"` // 外键
    User      *User          // 关联对象
    Item      string
    TotalCents int
}
```

**GORM 自动推断**：
- `UserID` 是外键，关联 `User.ID`
- `User` 字段用于预加载关联数据

### 3. AutoMigrate 的工作原理

`AutoMigrate` 会自动对比模型和数据库表结构，执行必要的 ALTER：

```go
// 第一次运行：创建表
db.AutoMigrate(&User{})
// SQL: CREATE TABLE users (id integer PRIMARY KEY, name text, email text)

// 添加新字段后再次运行：修改表
type User struct {
    ID    uint
    Name  string
    Email string
    Age   int  // 新增字段
}
db.AutoMigrate(&User{})
// SQL: ALTER TABLE users ADD COLUMN age integer
```

**注意事项**：
- AutoMigrate **不会删除字段**（防止数据丢失）
- 生产环境建议生成迁移脚本，而非自动执行
- 不适合复杂 schema 变更（如重命名字段）

### 4. 一对多关系的建模

**场景**：一个用户有多个订单（one-to-many）。

**模型定义**：
```go
type User struct {
    ID     uint
    Name   string
    Orders []Order `gorm:"constraint:OnDelete:CASCADE;"`
}

type Order struct {
    ID        uint
    UserID    uint
    Item      string
    TotalCents int
}
```

**关键点**：
- `Orders []Order`：声明一对多关系
- `constraint:OnDelete:CASCADE`：用户删除时，自动删除其订单（外键约束）
- `UserID`：GORM 自动识别为外键（`User` + `ID`）

**查询关联数据**：
```go
// ❌ 糟糕方式：N+1 查询问题
var users []User
db.Find(&users)
for i := range users {
    db.Model(&users[i]).Association("Orders").Find(&users[i].Orders)
    // 每个用户发一条 SQL → N+1 条查询
}

// ✅ 正确方式：Preload 预加载
var users []User
db.Preload("Orders").Find(&users)
// 只发两条 SQL：查用户 + 查所有订单
```

### 5. 事务的原子性保证

**问题场景**：用户下单时，需要同时写两行数据：
1. 创建订单记录
2. 扣减库存

如果第 1 步成功、第 2 步失败，就会出现"超卖"（库存为负）。

**事务解决**：
```go
db.Transaction(func(tx *gorm.DB) error {
    // 步骤 1：创建订单
    order := Order{UserID: 1, Item: "Book", TotalCents: 3590}
    if err := tx.Create(&order).Error; err != nil {
        return err // 返回错误会触发回滚
    }
    
    // 步骤 2：扣减库存
    result := tx.Exec("UPDATE inventory SET stock = stock - 1 WHERE item = ?", "Book")
    if result.Error != nil || result.RowsAffected == 0 {
        return errors.New("库存不足") // 触发回滚
    }
    
    return nil // 返回 nil 会提交事务
})
```

**工作流程**：
1. 开启事务（BEGIN TRANSACTION）
2. 所有操作在事务内执行（共享锁）
3. 返回 nil → 提交（COMMIT）
4. 返回 error → 回滚（ROLLBACK）

**为什么有效**：事务是原子的，要么全部生效，要么全部撤销。

## 常见错误

### 错误 1：忘记传指针

```go
// ❌ 错误代码
user := User{Name: "Alice"}
db.Create(user) // 编译能通过，但 ID 不会被填充

// 原因：GORM 需要修改结构体，必须传指针
```

**修复**：
```go
// ✅ 修复
db.Create(&user) // 传指针，user.ID 会被填充
```

**规则**：写操作（Create、Update、Delete）需要传指针，读操作可以传指针或值（推荐指针）。

### 错误 2：Preload 名称不匹配

```go
// ❌ 错误代码
type User struct {
    Orders []Order // 字段名是 Orders
}

db.Preload("orders").Find(&users) // 小写 o，无法匹配
```

**修复**：
```go
// ✅ 修复：字段名必须完全匹配（包括大小写）
db.Preload("Orders").Find(&users)
```

**原理**：Preload 通过反射查找结构体字段名，Go 区分大小写。

### 错误 3：事务内不使用 tx

```go
// ❌ 错误代码
db.Transaction(func(tx *gorm.DB) error {
    db.Create(&order)    // ❌ 用了 db，不在事务内！
    tx.Create(&inventory) // ✅ 用了 tx，在事务内
    return nil
})

// 后果：order 创建立即提交，inventory 失败时无法回滚
```

**修复**：
```go
// ✅ 修复：事务内统一用 tx
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&order)
    tx.Create(&inventory)
    return nil
})
```

**规则**：事务闭包内所有数据库操作必须用传入的 `tx`，不要用外层的 `db`。

## 动手练习

### 练习 1：预测输出

阅读以下代码，预测输出（先自己想，再看答案）：

```go
db.AutoMigrate(&User{})

// 创建
user := User{Name: "Bob", Email: "bob@example.com"}
db.Create(&user)
fmt.Println(user.ID) // ?

// 查询
var found User
db.First(&found, user.ID)
fmt.Println(found.Name) // ?

// 删除
db.Delete(&found)

// 再查询
var again User
result := db.First(&again, user.ID)
fmt.Println(result.Error) // ?
```

<details>
<summary>点击查看答案</summary>

**输出**：
```
1
Bob
record not found
```

**解析**：
1. `Create` 后 `user.ID` 自动填充为 1
2. `First` 按 ID 查询，找到记录的 Name
3. `Delete` 后记录不存在，`First` 返回 `ErrRecordNotFound`
</details>

### 练习 2：修复 Preload 错误

以下代码有什么问题？如何修复？

```go
type Author struct {
    ID    uint
    Name  string
    Books []Book
}

type Book struct {
    ID       uint
    AuthorID uint
    Title    string
}

// 查询作者及其书籍
var author Author
db.First(&author, 1)
// TODO: 加载 Books
fmt.Println(len(author.Books)) // 输出 0，但数据库有记录
```

<details>
<summary>点击查看答案</summary>

**问题**：没有用 `Preload`，关联字段不会自动加载。

**修复**：
```go
// 方式 1: 用 Preload
var author Author
db.Preload("Books").First(&author, 1)

// 方式 2: 用 Association
db.First(&author, 1)
db.Model(&author).Association("Books").Find(&author.Books)
```

**推荐**：Preload 更简洁，且能减少查询次数。
</details>

### 练习 3：实现转账事务

实现一个函数 `Transfer(fromID, toID uint, amount int)`，从 A 账户转账到 B 账户。

<details>
<summary>点击查看答案</summary>

```go
type Account struct {
    ID      uint
    Name    string
    Balance int
}

func Transfer(fromID, toID uint, amount int) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 步骤 1: 检查余额
        var from Account
        if err := tx.First(&from, fromID).Error; err != nil {
            return err
        }
        if from.Balance < amount {
            return errors.New("余额不足")
        }
        
        // 步骤 2: 扣款
        if err := tx.Model(&from).Update("balance", from.Balance - amount).Error; err != nil {
            return err
        }
        
        // 步骤 3: 收款
        var to Account
        if err := tx.First(&to, toID).Error; err != nil {
            return err
        }
        if err := tx.Model(&to).Update("balance", to.Balance + amount).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

**关键点**：所有操作在事务内，任何一方失败都会回滚。
</details>

## 故障排查 (FAQ)

### Q1: 为什么 Create 后 ID 还是 0？

**可能原因**：
1. **传了值而非指针**：`db.Create(user)` → `db.Create(&user)`
2. **没有主键字段**：结构体缺少 `ID uint` 或 `gorm:"primaryKey"`
3. **表没迁移**：忘记调用 `AutoMigrate`

**排查步骤**：
```go
user := User{Name: "Alice"}
fmt.Printf("创建前：ID=%d\n", user.ID) // 应该是 0
err := db.Create(&user).Error
fmt.Printf("创建后：ID=%d, err=%v\n", user.ID, err)
```

### Q2: Preload 的查询太慢怎么办？

**优化方法**：

1. **用 Select 指定字段**：
   ```go
   db.Preload("Orders", "item = ?", "Book").Find(&users)
   ```

2. **用 Joins 替代 Preload**（复杂查询）：
   ```go
   db.Joins("LEFT JOIN orders ON orders.user_id = users.id").
      Find(&users)
   ```

3. **避免在循环里查询**：
   ```go
   // ❌ 糟糕
   for _, user := range users {
       db.Model(&user).Association("Orders").Find(&user.Orders)
   }
   
   // ✅ 改进
   db.Preload("Orders").Find(&users)
   ```

### Q3: 如何查看 GORM 生成的 SQL？

**方法 1：开启日志**：
```go
db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

**输出**：
```
[info] CREATE TABLE `users` (`id` integer PRIMARY KEY,`name` text,`email` text)
[info] INSERT INTO `users` (`name`,`email`) VALUES ("Alice","alice@example.com")
```

**方法 2：用 Debug 模式**：
```go
db.Debug().Create(&user)
```

## 知识扩展 (选学)

### GORM 钩子（Hooks）

钩子允许你在数据库操作前后插入自定义逻辑：

```go
type User struct {
    ID       uint
    Name     string
    Password string
}

// 保存前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
    hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
    u.Password = string(hashed)
    return nil
}

// 查询后隐藏密码
func (u *User) AfterFind(tx *gorm.DB) error {
    u.Password = "***"
    return nil
}
```

**钩子类型**：
- `BeforeSave` / `AfterSave`
- `BeforeCreate` / `AfterCreate`
- `BeforeUpdate` / `AfterUpdate`
- `BeforeDelete` / `AfterDelete`
- `AfterFind`

**注意**：钩子会增加耦合，谨慎使用。

### 软删除（Soft Delete）

软删除不是真正删除数据，而是标记为"已删除"：

```go
type User struct {
    gorm.Model // 包含 ID、CreatedAt、UpdatedAt、DeletedAt
    Name       string
}

db.Delete(&user)
// SQL: UPDATE users SET deleted_at = NOW() WHERE id = ?

// 查询时自动过滤已删除
db.Find(&users)
// SQL: SELECT * FROM users WHERE deleted_at IS NULL

// 强制包含已删除
db.Unscoped().Find(&users)
```

**适用场景**：用户注销（保留数据）、订单历史、审计日志。

### 连接池配置

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)      // 最大打开连接数
sqlDB.SetMaxIdleConns(5)       // 空闲连接数
sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接最大存活时间
```

**建议值**：
- 小型服务：`MaxOpenConns = 10-25`
- 高并发服务：`MaxOpenConns = 50-100`
- `MaxIdleConns` 设为 `MaxOpenConns` 的 20-50%

## 工业界应用

### 场景 1：电商订单系统

```go
type Order struct {
    ID          uint
    UserID      uint
    Status      string // pending, paid, shipped
    TotalAmount int
    Items       []OrderItem
    CreatedAt   time.Time
}

type OrderItem struct {
    ID        uint
    OrderID   uint
    ProductID uint
    Quantity  int
    UnitPrice int
}

// 创建订单（带事务）
func CreateOrder(userID uint, items []CartItem) (*Order, error) {
    var order *Order
    err := db.Transaction(func(tx *gorm.DB) error {
        // 创建订单
        order = &Order{UserID: userID, Status: "pending"}
        if err := tx.Create(order).Error; err != nil {
            return err
        }
        
        // 创建订单项（同时扣库存）
        for _, item := range items {
            if err := checkAndDecreaseStock(tx, item.ProductID, item.Quantity); err != nil {
                return err
            }
            tx.Create(&OrderItem{
                OrderID:   order.ID,
                ProductID: item.ProductID,
                Quantity:  item.Quantity,
                UnitPrice: getUnitPrice(tx, item.ProductID),
            })
        }
        
        return nil
    })
    
    return order, err
}
```

**关键点**：订单和库存更新在同一事务中，避免超卖。

### 场景 2：博客文章与评论

```go
type Post struct {
    ID       uint
    Title    string
    Content  string
    AuthorID uint
    Comments []Comment `gorm:"constraint:OnDelete:CASCADE;"`
}

type Comment struct {
    ID      uint
    PostID  uint
    UserID  uint
    Content string
}

// 查询文章及评论
func GetPostWithComments(postID uint) (*Post, error) {
    var post Post
    err := db.Preload("Comments.User").First(&post, postID).Error
    return &post, err
}
```

**优化**：用 `Preload("Comments.User")` 同时预加载评论和评论者信息。

### 场景 3：用户积分系统

```go
type UserPoints struct {
    ID        uint
    UserID    uint `gorm:"uniqueIndex"`
    Points    int
    UpdatedAt time.Time
}

// 增加积分（乐观锁防止并发问题）
func AddPoints(userID uint, points int) error {
    return db.Transaction(func(tx *gorm.DB) error {
        var up UserPoints
        if err := tx.Where("user_id = ?", userID).First(&up).Error; err != nil {
            return err
        }
        
        // 乐观锁：检查UpdatedAt没变化
        result := tx.Model(&up).
            Where("updated_at = ?", up.UpdatedAt).
            Update("points", up.Points + points)
        
        if result.RowsAffected == 0 {
            return errors.New("concurrent update, please retry")
        }
        
        return result.Error
    })
}
```

**为什么需要乐观锁**：多个请求同时修改积分时，防止覆盖。

## 小结

### 核心要点

1. **模型定义**：用结构体和标签映射表结构，GORM 自动推断主键、外键
2. **AutoMigrate**：自动同步模型到数据库，适合开发和测试
3. **CRUD 操作**：Create/First/Update/Delete，注意传指针让 GORM 填充字段
4. **Preload 预加载**：解决 N+1 查询问题，一次性加载关联数据
5. **事务保护一致性**：多步写操作用 Transaction，返回 error 自动回滚

### 关键术语

| 英文 | 中文 | 说明 |
|------|------|------|
| ORM | 对象关系映射 | 用对象操作代替 SQL |
| AutoMigrate | 自动迁移 | 同步模型到数据库表 |
| Preload | 预加载 | 一次性加载关联数据 |
| Transaction | 事务 | 原子操作，要么全成功要么全失败 |
| N+1 query | N+1 查询问题 | 循环内查询导致性能问题 |
| Cascade delete | 级联删除 | 删除主记录时自动删除从记录 |

### 下一步建议

1. 为你的项目定义领域模型（User、Post、Comment 等）
2. 用 AutoMigrate 创建数据库表
3. 实现基础的 CRUD 操作
4. 添加一对多关系，练习 Preload
5. 为关键业务逻辑添加事务保护

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 对象关系映射 | ORM | 用面向对象的方式操作关系数据库的技术 |
| 自动迁移 | AutoMigrate | GORM 根据模型自动创建或更新数据库表的功能 |
| 预加载 | Preload | 使用 JOIN 一次性加载关联数据，避免 N+1 查询 |
| 事务 | Transaction | 数据库操作的原子单元，保证 ACID 特性 |
| 一对多关系 | One-to-Many Relationship | 一个实体关联多个实体的关系（如用户 - 订单） |
| 外键 | Foreign Key | 指向另一张表主键的字段 |
| 级联删除 | Cascade Delete | 删除父记录时自动删除子记录 |
| 主键 | Primary Key | 唯一标识表中记录的字段 |
| 内存数据库 | In-Memory Database | 数据存储在内存中的数据库（如 SQLite :memory:） |
| 软删除 | Soft Delete | 标记删除而非真正删除数据的策略 |
| 乐观锁 | Optimistic Lock | 通过版本号或时间戳检测并发冲突 |

## 源码

完整示例代码位于：[internal/advance/database/database.go](https://github.com/savechina/hello-go/blob/main/internal/advance/database/database.go)
