package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
)

// 通过以下代码配置，再使用 gorm.io/gen 生成 dao 和 model 的代码
// 在本目录下执行 go run main.go 即可
// 成功生成代码后，将根目录下的 do 和 entity 文件夹迁移到 app/user/models 目录下

var (
	// 根据本地环境，配置 username，password，dbName，host，port
	dsn = "username:password@tcp(host:port)/dbName?charset=utf8mb4&parseTime=True&loc=Local"
	// curd 代码的输出路径
	outPath = "./do"
	// entity 代码的输出路径
	modelPkgPath = "entity"
)

func main() {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	// 构造生成器实例
	g := gen.NewGenerator(gen.Config{

		OutPath:      outPath,
		ModelPkgPath: modelPkgPath,

		// WithDefaultQuery 生成默认查询结构体(作为全局变量使用), 即 `Q`结构体和其字段(各表模型)
		// WithoutContext 生成没有 context 调用限制的代码供查询
		// WithQueryInterface 生成 interface 形式的查询代码(可导出), 如 `Where()` 方法返回的就是一个可导出的接口类型
		Mode: gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,

		// 表字段可为 null 值时, 对应结体字段使用指针类型
		FieldNullable: false, // generate pointer when field is nullable

		// 表字段默认值与模型结构体字段零值不一致的字段, 在插入数据时需要赋值该字段值为零值的, 结构体字段须是指针类型才能成功, 即`FieldCoverable:true`配置下生成的结构体字段.
		// 因为在插入时遇到字段为零值的会被GORM赋予默认值. 如字段`age`表默认值为10, 即使你显式设置为0最后也会被GORM设为10提交.
		// 如果该字段没有上面提到的插入时赋零值的特殊需要, 则字段为非指针类型使用起来会比较方便.
		FieldCoverable: false, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values

		// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
		FieldSignable: false, // detect integer field's unsigned type, adjust generated data type
		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: false, // generate with gorm index tag
		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: true, // generate with gorm column type tag
	})

	// 设置目标 db
	g.UseDB(db)

	// 自定义字段的数据类型
	// 统一数字类型为 int64, 兼容 protobuf
	dataMap := map[string]func(detailType gorm.ColumnType) (dataType string){
		"tinyint":   func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint":    func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"int":       func(detailType gorm.ColumnType) (dataType string) { return "int64" },
	}
	g.WithDataTypeMap(dataMap)

	// 自定义模型结体字段的标签
	// 将特定字段名的 json 标签加上`string`属性,即 MarshalJSON 时该字段由数字类型转成字符串类型
	jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
		toStringField := `id, `
		if strings.Contains(toStringField, columnName) {
			return columnName + ",string"
		}
		return columnName
	})

	// 将非默认字段名的字段定义为自动时间戳和软删除字段;
	// 自动时间戳默认字段名为: `updated_at`、`created_at, 表字段数据类型为: INT 或 DATETIME
	// 软删除默认字段名为: `deleted_at`, 表字段数据类型为: DATETIME
	//autoUpdateTimeField := gen.FieldGORMTag("update_time", "column:update_time;type:int unsigned;autoUpdateTime")
	//autoCreateTimeField := gen.FieldGORMTag("create_time", "column:create_time;type:int unsigned;autoCreateTime")
	//softDeleteField := gen.FieldType("delete_time", "gorm.DeletedAt")
	// 模型自定义选项组
	//fieldOpts := []gen.ModelOpt{jsonField, autoCreateTimeField, autoUpdateTimeField, softDeleteField}
	fieldOpts := []gen.ModelOpt{jsonField}

	// 创建模型的结构体,生成文件在 model 目录; 先创建的结果会被后面创建的覆盖
	// 这里创建个别模型仅仅是为了拿到`*generate.QueryStructMeta`类型对象用于后面的模型关联操作中
	//User := g.GenerateModel("user")
	// 创建全部模型文件, 并覆盖前面创建的同名模型
	allModel := g.GenerateAllTable(fieldOpts...)

	// 创建有关联关系的模型文件
	// 可以用于指定外键
	//Score := g.GenerateModel("score",
	// append(
	//    fieldOpts,
	//    // user 一对多 address 关联, 外键`uid`在 address 表中
	//    gen.FieldRelate(field.HasMany, "user", User, &field.RelateConfig{GORMTag: "foreignKey:UID"}),
	// )...,
	//)

	// 创建模型的方法,生成文件在 query 目录; 先创建结果不会被后创建的覆盖
	//g.ApplyBasic(User)
	g.ApplyBasic(allModel...)

	g.Execute()
}
