// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"go-api/internal/pkg/models"
)

func newOrderProduct(db *gorm.DB, opts ...gen.DOOption) orderProduct {
	_orderProduct := orderProduct{}

	_orderProduct.orderProductDo.UseDB(db, opts...)
	_orderProduct.orderProductDo.UseModel(&models.OrderProduct{})

	tableName := _orderProduct.orderProductDo.TableName()
	_orderProduct.ALL = field.NewAsterisk(tableName)
	_orderProduct.ID = field.NewString(tableName, "id")
	_orderProduct.OrderID = field.NewString(tableName, "order_id")
	_orderProduct.Discount = field.NewFloat64(tableName, "discount")
	_orderProduct.ProductID = field.NewString(tableName, "product_id")
	_orderProduct.CreatedAt = field.NewTime(tableName, "created_at")
	_orderProduct.UpdatedAt = field.NewTime(tableName, "updated_at")
	_orderProduct.TotalAmounts = field.NewInt32(tableName, "total_amounts")
	_orderProduct.TotalProductPrice = field.NewFloat64(tableName, "total_product_price")

	_orderProduct.fillFieldMap()

	return _orderProduct
}

type orderProduct struct {
	orderProductDo

	ALL               field.Asterisk
	ID                field.String
	OrderID           field.String
	Discount          field.Float64
	ProductID         field.String
	CreatedAt         field.Time
	UpdatedAt         field.Time
	TotalAmounts      field.Int32
	TotalProductPrice field.Float64

	fieldMap map[string]field.Expr
}

func (o orderProduct) Table(newTableName string) *orderProduct {
	o.orderProductDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o orderProduct) As(alias string) *orderProduct {
	o.orderProductDo.DO = *(o.orderProductDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *orderProduct) updateTableName(table string) *orderProduct {
	o.ALL = field.NewAsterisk(table)
	o.ID = field.NewString(table, "id")
	o.OrderID = field.NewString(table, "order_id")
	o.Discount = field.NewFloat64(table, "discount")
	o.ProductID = field.NewString(table, "product_id")
	o.CreatedAt = field.NewTime(table, "created_at")
	o.UpdatedAt = field.NewTime(table, "updated_at")
	o.TotalAmounts = field.NewInt32(table, "total_amounts")
	o.TotalProductPrice = field.NewFloat64(table, "total_product_price")

	o.fillFieldMap()

	return o
}

func (o *orderProduct) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *orderProduct) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 8)
	o.fieldMap["id"] = o.ID
	o.fieldMap["order_id"] = o.OrderID
	o.fieldMap["discount"] = o.Discount
	o.fieldMap["product_id"] = o.ProductID
	o.fieldMap["created_at"] = o.CreatedAt
	o.fieldMap["updated_at"] = o.UpdatedAt
	o.fieldMap["total_amounts"] = o.TotalAmounts
	o.fieldMap["total_product_price"] = o.TotalProductPrice
}

func (o orderProduct) clone(db *gorm.DB) orderProduct {
	o.orderProductDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o orderProduct) replaceDB(db *gorm.DB) orderProduct {
	o.orderProductDo.ReplaceDB(db)
	return o
}

type orderProductDo struct{ gen.DO }

type IOrderProductDo interface {
	gen.SubQuery
	Debug() IOrderProductDo
	WithContext(ctx context.Context) IOrderProductDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOrderProductDo
	WriteDB() IOrderProductDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOrderProductDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOrderProductDo
	Not(conds ...gen.Condition) IOrderProductDo
	Or(conds ...gen.Condition) IOrderProductDo
	Select(conds ...field.Expr) IOrderProductDo
	Where(conds ...gen.Condition) IOrderProductDo
	Order(conds ...field.Expr) IOrderProductDo
	Distinct(cols ...field.Expr) IOrderProductDo
	Omit(cols ...field.Expr) IOrderProductDo
	Join(table schema.Tabler, on ...field.Expr) IOrderProductDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IOrderProductDo
	RightJoin(table schema.Tabler, on ...field.Expr) IOrderProductDo
	Group(cols ...field.Expr) IOrderProductDo
	Having(conds ...gen.Condition) IOrderProductDo
	Limit(limit int) IOrderProductDo
	Offset(offset int) IOrderProductDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOrderProductDo
	Unscoped() IOrderProductDo
	Create(values ...*models.OrderProduct) error
	CreateInBatches(values []*models.OrderProduct, batchSize int) error
	Save(values ...*models.OrderProduct) error
	First() (*models.OrderProduct, error)
	Take() (*models.OrderProduct, error)
	Last() (*models.OrderProduct, error)
	Find() ([]*models.OrderProduct, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.OrderProduct, err error)
	FindInBatches(result *[]*models.OrderProduct, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.OrderProduct) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOrderProductDo
	Assign(attrs ...field.AssignExpr) IOrderProductDo
	Joins(fields ...field.RelationField) IOrderProductDo
	Preload(fields ...field.RelationField) IOrderProductDo
	FirstOrInit() (*models.OrderProduct, error)
	FirstOrCreate() (*models.OrderProduct, error)
	FindByPage(offset int, limit int) (result []*models.OrderProduct, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Rows() (*sql.Rows, error)
	Row() *sql.Row
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOrderProductDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (o orderProductDo) Debug() IOrderProductDo {
	return o.withDO(o.DO.Debug())
}

func (o orderProductDo) WithContext(ctx context.Context) IOrderProductDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o orderProductDo) ReadDB() IOrderProductDo {
	return o.Clauses(dbresolver.Read)
}

func (o orderProductDo) WriteDB() IOrderProductDo {
	return o.Clauses(dbresolver.Write)
}

func (o orderProductDo) Session(config *gorm.Session) IOrderProductDo {
	return o.withDO(o.DO.Session(config))
}

func (o orderProductDo) Clauses(conds ...clause.Expression) IOrderProductDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o orderProductDo) Returning(value interface{}, columns ...string) IOrderProductDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o orderProductDo) Not(conds ...gen.Condition) IOrderProductDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o orderProductDo) Or(conds ...gen.Condition) IOrderProductDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o orderProductDo) Select(conds ...field.Expr) IOrderProductDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o orderProductDo) Where(conds ...gen.Condition) IOrderProductDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o orderProductDo) Order(conds ...field.Expr) IOrderProductDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o orderProductDo) Distinct(cols ...field.Expr) IOrderProductDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o orderProductDo) Omit(cols ...field.Expr) IOrderProductDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o orderProductDo) Join(table schema.Tabler, on ...field.Expr) IOrderProductDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o orderProductDo) LeftJoin(table schema.Tabler, on ...field.Expr) IOrderProductDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o orderProductDo) RightJoin(table schema.Tabler, on ...field.Expr) IOrderProductDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o orderProductDo) Group(cols ...field.Expr) IOrderProductDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o orderProductDo) Having(conds ...gen.Condition) IOrderProductDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o orderProductDo) Limit(limit int) IOrderProductDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o orderProductDo) Offset(offset int) IOrderProductDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o orderProductDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IOrderProductDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o orderProductDo) Unscoped() IOrderProductDo {
	return o.withDO(o.DO.Unscoped())
}

func (o orderProductDo) Create(values ...*models.OrderProduct) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o orderProductDo) CreateInBatches(values []*models.OrderProduct, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o orderProductDo) Save(values ...*models.OrderProduct) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o orderProductDo) First() (*models.OrderProduct, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.OrderProduct), nil
	}
}

func (o orderProductDo) Take() (*models.OrderProduct, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.OrderProduct), nil
	}
}

func (o orderProductDo) Last() (*models.OrderProduct, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.OrderProduct), nil
	}
}

func (o orderProductDo) Find() ([]*models.OrderProduct, error) {
	result, err := o.DO.Find()
	return result.([]*models.OrderProduct), err
}

func (o orderProductDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.OrderProduct, err error) {
	buf := make([]*models.OrderProduct, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o orderProductDo) FindInBatches(result *[]*models.OrderProduct, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o orderProductDo) Attrs(attrs ...field.AssignExpr) IOrderProductDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o orderProductDo) Assign(attrs ...field.AssignExpr) IOrderProductDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o orderProductDo) Joins(fields ...field.RelationField) IOrderProductDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o orderProductDo) Preload(fields ...field.RelationField) IOrderProductDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o orderProductDo) FirstOrInit() (*models.OrderProduct, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.OrderProduct), nil
	}
}

func (o orderProductDo) FirstOrCreate() (*models.OrderProduct, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.OrderProduct), nil
	}
}

func (o orderProductDo) FindByPage(offset int, limit int) (result []*models.OrderProduct, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o orderProductDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o orderProductDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o orderProductDo) Delete(models ...*models.OrderProduct) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *orderProductDo) withDO(do gen.Dao) *orderProductDo {
	o.DO = *do.(*gen.DO)
	return o
}
