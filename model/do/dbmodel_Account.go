package do

// dbmodel orm封装文件

import (
	orm "demo_backend/model/orm"			//这个包名从go.mod获取
	dbmodel "demo_backend/model/dbmodel" 	//这个包名从go.mod获取
)

func init() {
	orm.Conn().AutoMigrate(&dbmodel.Account{})
}

// 增 -------------------------------------------------------------------
func Account_CreateOne(data *dbmodel.Account) *dbmodel.Account {
	// 创建一个
	result := orm.Conn().Create(data)
	orm.PanicGormResultError("do.Account_CreateOne", result)
	return data
}

func Account_CreateMulti(data *[]dbmodel.Account) *[]dbmodel.Account {
	// 创建多个
	result := orm.Conn().Create(data)
	orm.PanicGormResultError("do.Account_CreateMulti", result)
	return data
}

// 删 -------------------------------------------------------------------
func Account_DeleteOneById(id uint) int64 {
	result := orm.Conn().Delete(&dbmodel.Account{}, id)
	orm.PanicGormResultError("do.Account_DeleteOneById", result)
	return result.RowsAffected
}

func Account_DeleteMultiByFields(withFields *dbmodel.Account) int64 {
	// 根据withFields的非零字段，删除符合条件的多条数据。where里的struct只有非零生效
	// 对于主键，可以指定ID，也可以指定其他复合主键
	result := orm.Conn().Where(withFields).Delete(withFields)
	orm.PanicGormResultError("do.Account_DeleteMultiByFields", result)
	return result.RowsAffected
}

func Account_DeleteMultiByIds(ids []uint) int64 {
	// 根据指定的多个id，删除多条数据
	result := orm.Conn().Delete(&dbmodel.Account{}, ids)
	orm.PanicGormResultError("do.Account_DeleteMultiByIds", result)
	return result.RowsAffected
}

// 查 -------------------------------------------------------------------
/*
	如果where有更复杂的条件，需要handler自己实现，或者补充本模板
*/

func Account_QueryOneById(id uint) *dbmodel.Account {
	// 根据id查询
	var queryResult dbmodel.Account
	result := orm.Conn().Find(&queryResult, id)
	orm.PanicGormResultError("do.Account_QueryOneById", result)
	return &queryResult
}

func Account_QueryMultiByIds(ids []uint) *[]dbmodel.Account {
	//根据多个id，查询多条数据
	var queryResult []dbmodel.Account
	result := orm.Conn().Find(&queryResult, ids)
	orm.PanicGormResultError("do.Account_QueryMultiByIds", result)
	return &queryResult
}

func Account_QueryMultiByFields(withFields *dbmodel.Account) *[]dbmodel.Account {
	//根据非零字段指定，查询多条数据。where里的struct只有非零字段生效
	var queryResult []dbmodel.Account
	result := orm.Conn().Where(withFields).Find(&queryResult)
	orm.PanicGormResultError("do.Account_QueryMultiByFields", result)
	return &queryResult
}

func Account_QueryOneByFields(withFields *dbmodel.Account) *dbmodel.Account {
	// 根据非零字段指定，查询1条数据。where里的struct只有非零字段生效
	// 这两个指针不相同，以返回数据为准，未找到则返回值的ID=0
	var queryResult dbmodel.Account
	result := orm.Conn().Limit(1).Where(withFields).Find(&queryResult)
	orm.PanicGormResultError("do.Account_QueryMultiByFields", result)
	return &queryResult
}

func Account_QueryExistsByFields(withFields *dbmodel.Account) bool {
	// 根据非零字段查找，看该数据是否存在
	// 若存在: 返回true，withFields的ID也会被设置
	one:=Account_QueryOneByFields(withFields)
	if 0!=one.ID{
		withFields.ID = one.ID
		return true
	}else{
		return false
	}
}

func Account_QueryAll() *[]dbmodel.Account {
	// 返回所有表中数据
	var queryResult []dbmodel.Account
	result := orm.Conn().Find(&queryResult)
	orm.PanicGormResultError("do.Account_QueryAll", result)
	return &queryResult
}

func Account_QueryAllCount() int64 {
	// 返回所有表中数据个数
	var count int64
	result:= orm.Conn().Model(&dbmodel.Account{}).Where("1=1").Count(&count)
	orm.PanicGormResultError("do.Account_QueryAllCount", result)
	return count
}

func Account_QueryMultiByLike(field, reg string) *[]dbmodel.Account {
	// where like, ("Name","pyx%") → where name like pyx%
	// 自动转换field命名格式
	var queryResult []dbmodel.Account
	result := orm.Conn().Where(orm.ToOrmFieldName(field)+" LIKE ?", reg).Find(&queryResult)
	orm.PanicGormResultError("do.Account_QueryLike", result)
	return &queryResult
}

func Account_QueryMultiByIn(field string, vals []string) *[]dbmodel.Account {
	// where in, ("ID", []string{"1","3"}) → where id in ["1","4"]
	// 自动转换field命名格式
	var queryResult []dbmodel.Account
	result := orm.Conn().Where(orm.ToOrmFieldName(field)+" IN ?", vals).Find(&queryResult)
	orm.PanicGormResultError("do.Account_QueryMultiByIn", result)
	return &queryResult
}

func Account_QueryMultiByBetween(field string, a, b interface{}) *[]dbmodel.Account {
	// where between, ("ID",2,4) → where id between '2' and '4'
	// 自动转换field命名格式
	var queryResult []dbmodel.Account
	result := orm.Conn().Where(orm.ToOrmFieldName(field)+" between ? and ?", a, b).Find(&queryResult)
	orm.PanicGormResultError("do.Account_QueryMultiByBetween", result)
	return &queryResult
}

/*
	select取出指定字段，而非所有字段。取出指定字段，本身就是一个模式，可以和其他模式配合。

func Account_QueryMulti_Columns1_ByFields(withFields *dbmodel.Account) *[]iomodel.Columns1 {
	// 这是个模板，这个模式不必自动生成，手动改改Columns1类型就能用
	// 使用gorm高级查询的智能选择字段，只要修改下类型，就能取出指定字段
	var queryResult []iomodel.Columns1
	result := orm.Conn().Model(withFields).Where(withFields).Find(&queryResult)
	orm.PanicGormResultError("do.Account_QueryColumns1ByFields", result)
	return &queryResult
}
*/

// 改 -------------------------------------------------------------------
func Account_UpdateOneById(withIdAndFields *dbmodel.Account) int64 {
	// 根据id，更新非零字段。Model指定主键。Updates里的struct只有非零字段生效。
	result := orm.Conn().Model(withIdAndFields).Updates(withIdAndFields)
	orm.PanicGormResultError("do.Account_UpdateOneById", result)
	return result.RowsAffected
}

func Account_UpdateMultiByIds(ids *[]uint, updas *dbmodel.Account) int64 {
	// 根据id，更新非零字段。Model指定主键。Updates里的struct只有非零字段生效。
	result := orm.Conn().Model(updas).Where("id IN ?", *ids).Updates(updas)
	orm.PanicGormResultError("do.Account_UpdateMultiByIds", result)
	return result.RowsAffected
}

func Account_UpdateMultiByFields(selec, updas *dbmodel.Account) int64 {
	// 根据selec指定字段选择，选中后用updas里非零字段更新。Where/Updates都只有非零字段生效
	// 返回更新的行数
	result := orm.Conn().Where(selec).Updates(updas)
	orm.PanicGormResultError("do.Account_UpdateMultiByFields", result)
	return result.RowsAffected
}
