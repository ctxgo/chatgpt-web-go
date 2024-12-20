// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"chatgpt-web-new-go/model"
)

func newInviteRecord(db *gorm.DB, opts ...gen.DOOption) inviteRecord {
	_inviteRecord := inviteRecord{}

	_inviteRecord.inviteRecordDo.UseDB(db, opts...)
	_inviteRecord.inviteRecordDo.UseModel(&model.InviteRecord{})

	tableName := _inviteRecord.inviteRecordDo.TableName()
	_inviteRecord.ALL = field.NewAsterisk(tableName)
	_inviteRecord.ID = field.NewInt64(tableName, "id")
	_inviteRecord.UserID = field.NewInt64(tableName, "user_id")
	_inviteRecord.InviteCode = field.NewString(tableName, "invite_code")
	_inviteRecord.SuperiorID = field.NewInt64(tableName, "superior_id")
	_inviteRecord.Reward = field.NewString(tableName, "reward")
	_inviteRecord.RewardType = field.NewString(tableName, "reward_type")
	_inviteRecord.Status = field.NewInt32(tableName, "status")
	_inviteRecord.Remarks = field.NewString(tableName, "remarks")
	_inviteRecord.IP = field.NewString(tableName, "ip")
	_inviteRecord.UserAgent = field.NewString(tableName, "user_agent")
	_inviteRecord.CreateTime = field.NewTime(tableName, "create_time")
	_inviteRecord.UpdateTime = field.NewTime(tableName, "update_time")
	_inviteRecord.IsDelete = field.NewInt32(tableName, "is_delete")

	_inviteRecord.fillFieldMap()

	return _inviteRecord
}

type inviteRecord struct {
	inviteRecordDo

	ALL        field.Asterisk
	ID         field.Int64
	UserID     field.Int64
	InviteCode field.String // 邀请码
	SuperiorID field.Int64  // 上级ID（一旦确定将不可修改）
	Reward     field.String // 奖励
	RewardType field.String // 奖励类型
	Status     field.Int32  // 0-异常｜1-正常发放｜3-审核中
	Remarks    field.String // 评论
	IP         field.String
	UserAgent  field.String // ua
	CreateTime field.Time
	UpdateTime field.Time
	IsDelete   field.Int32

	fieldMap map[string]field.Expr
}

func (i inviteRecord) Table(newTableName string) *inviteRecord {
	i.inviteRecordDo.UseTable(newTableName)
	return i.updateTableName(newTableName)
}

func (i inviteRecord) As(alias string) *inviteRecord {
	i.inviteRecordDo.DO = *(i.inviteRecordDo.As(alias).(*gen.DO))
	return i.updateTableName(alias)
}

func (i *inviteRecord) updateTableName(table string) *inviteRecord {
	i.ALL = field.NewAsterisk(table)
	i.ID = field.NewInt64(table, "id")
	i.UserID = field.NewInt64(table, "user_id")
	i.InviteCode = field.NewString(table, "invite_code")
	i.SuperiorID = field.NewInt64(table, "superior_id")
	i.Reward = field.NewString(table, "reward")
	i.RewardType = field.NewString(table, "reward_type")
	i.Status = field.NewInt32(table, "status")
	i.Remarks = field.NewString(table, "remarks")
	i.IP = field.NewString(table, "ip")
	i.UserAgent = field.NewString(table, "user_agent")
	i.CreateTime = field.NewTime(table, "create_time")
	i.UpdateTime = field.NewTime(table, "update_time")
	i.IsDelete = field.NewInt32(table, "is_delete")

	i.fillFieldMap()

	return i
}

func (i *inviteRecord) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := i.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (i *inviteRecord) fillFieldMap() {
	i.fieldMap = make(map[string]field.Expr, 13)
	i.fieldMap["id"] = i.ID
	i.fieldMap["user_id"] = i.UserID
	i.fieldMap["invite_code"] = i.InviteCode
	i.fieldMap["superior_id"] = i.SuperiorID
	i.fieldMap["reward"] = i.Reward
	i.fieldMap["reward_type"] = i.RewardType
	i.fieldMap["status"] = i.Status
	i.fieldMap["remarks"] = i.Remarks
	i.fieldMap["ip"] = i.IP
	i.fieldMap["user_agent"] = i.UserAgent
	i.fieldMap["create_time"] = i.CreateTime
	i.fieldMap["update_time"] = i.UpdateTime
	i.fieldMap["is_delete"] = i.IsDelete
}

func (i inviteRecord) clone(db *gorm.DB) inviteRecord {
	i.inviteRecordDo.ReplaceConnPool(db.Statement.ConnPool)
	return i
}

func (i inviteRecord) replaceDB(db *gorm.DB) inviteRecord {
	i.inviteRecordDo.ReplaceDB(db)
	return i
}

type inviteRecordDo struct{ gen.DO }

type IInviteRecordDo interface {
	gen.SubQuery
	Debug() IInviteRecordDo
	WithContext(ctx context.Context) IInviteRecordDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IInviteRecordDo
	WriteDB() IInviteRecordDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IInviteRecordDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IInviteRecordDo
	Not(conds ...gen.Condition) IInviteRecordDo
	Or(conds ...gen.Condition) IInviteRecordDo
	Select(conds ...field.Expr) IInviteRecordDo
	Where(conds ...gen.Condition) IInviteRecordDo
	Order(conds ...field.Expr) IInviteRecordDo
	Distinct(cols ...field.Expr) IInviteRecordDo
	Omit(cols ...field.Expr) IInviteRecordDo
	Join(table schema.Tabler, on ...field.Expr) IInviteRecordDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IInviteRecordDo
	RightJoin(table schema.Tabler, on ...field.Expr) IInviteRecordDo
	Group(cols ...field.Expr) IInviteRecordDo
	Having(conds ...gen.Condition) IInviteRecordDo
	Limit(limit int) IInviteRecordDo
	Offset(offset int) IInviteRecordDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IInviteRecordDo
	Unscoped() IInviteRecordDo
	Create(values ...*model.InviteRecord) error
	CreateInBatches(values []*model.InviteRecord, batchSize int) error
	Save(values ...*model.InviteRecord) error
	First() (*model.InviteRecord, error)
	Take() (*model.InviteRecord, error)
	Last() (*model.InviteRecord, error)
	Find() ([]*model.InviteRecord, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.InviteRecord, err error)
	FindInBatches(result *[]*model.InviteRecord, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.InviteRecord) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IInviteRecordDo
	Assign(attrs ...field.AssignExpr) IInviteRecordDo
	Joins(fields ...field.RelationField) IInviteRecordDo
	Preload(fields ...field.RelationField) IInviteRecordDo
	FirstOrInit() (*model.InviteRecord, error)
	FirstOrCreate() (*model.InviteRecord, error)
	FindByPage(offset int, limit int) (result []*model.InviteRecord, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IInviteRecordDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	FilterWithNameAndRole(name string, role string) (result []model.InviteRecord, err error)
}

// FilterWithNameAndRole SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
func (i inviteRecordDo) FilterWithNameAndRole(name string, role string) (result []model.InviteRecord, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	generateSQL.WriteString("SELECT * FROM invite_record WHERE name = ? ")
	if role != "" {
		params = append(params, role)
		generateSQL.WriteString("AND role = ? ")
	}

	var executeSQL *gorm.DB
	executeSQL = i.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (i inviteRecordDo) Debug() IInviteRecordDo {
	return i.withDO(i.DO.Debug())
}

func (i inviteRecordDo) WithContext(ctx context.Context) IInviteRecordDo {
	return i.withDO(i.DO.WithContext(ctx))
}

func (i inviteRecordDo) ReadDB() IInviteRecordDo {
	return i.Clauses(dbresolver.Read)
}

func (i inviteRecordDo) WriteDB() IInviteRecordDo {
	return i.Clauses(dbresolver.Write)
}

func (i inviteRecordDo) Session(config *gorm.Session) IInviteRecordDo {
	return i.withDO(i.DO.Session(config))
}

func (i inviteRecordDo) Clauses(conds ...clause.Expression) IInviteRecordDo {
	return i.withDO(i.DO.Clauses(conds...))
}

func (i inviteRecordDo) Returning(value interface{}, columns ...string) IInviteRecordDo {
	return i.withDO(i.DO.Returning(value, columns...))
}

func (i inviteRecordDo) Not(conds ...gen.Condition) IInviteRecordDo {
	return i.withDO(i.DO.Not(conds...))
}

func (i inviteRecordDo) Or(conds ...gen.Condition) IInviteRecordDo {
	return i.withDO(i.DO.Or(conds...))
}

func (i inviteRecordDo) Select(conds ...field.Expr) IInviteRecordDo {
	return i.withDO(i.DO.Select(conds...))
}

func (i inviteRecordDo) Where(conds ...gen.Condition) IInviteRecordDo {
	return i.withDO(i.DO.Where(conds...))
}

func (i inviteRecordDo) Order(conds ...field.Expr) IInviteRecordDo {
	return i.withDO(i.DO.Order(conds...))
}

func (i inviteRecordDo) Distinct(cols ...field.Expr) IInviteRecordDo {
	return i.withDO(i.DO.Distinct(cols...))
}

func (i inviteRecordDo) Omit(cols ...field.Expr) IInviteRecordDo {
	return i.withDO(i.DO.Omit(cols...))
}

func (i inviteRecordDo) Join(table schema.Tabler, on ...field.Expr) IInviteRecordDo {
	return i.withDO(i.DO.Join(table, on...))
}

func (i inviteRecordDo) LeftJoin(table schema.Tabler, on ...field.Expr) IInviteRecordDo {
	return i.withDO(i.DO.LeftJoin(table, on...))
}

func (i inviteRecordDo) RightJoin(table schema.Tabler, on ...field.Expr) IInviteRecordDo {
	return i.withDO(i.DO.RightJoin(table, on...))
}

func (i inviteRecordDo) Group(cols ...field.Expr) IInviteRecordDo {
	return i.withDO(i.DO.Group(cols...))
}

func (i inviteRecordDo) Having(conds ...gen.Condition) IInviteRecordDo {
	return i.withDO(i.DO.Having(conds...))
}

func (i inviteRecordDo) Limit(limit int) IInviteRecordDo {
	return i.withDO(i.DO.Limit(limit))
}

func (i inviteRecordDo) Offset(offset int) IInviteRecordDo {
	return i.withDO(i.DO.Offset(offset))
}

func (i inviteRecordDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IInviteRecordDo {
	return i.withDO(i.DO.Scopes(funcs...))
}

func (i inviteRecordDo) Unscoped() IInviteRecordDo {
	return i.withDO(i.DO.Unscoped())
}

func (i inviteRecordDo) Create(values ...*model.InviteRecord) error {
	if len(values) == 0 {
		return nil
	}
	return i.DO.Create(values)
}

func (i inviteRecordDo) CreateInBatches(values []*model.InviteRecord, batchSize int) error {
	return i.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (i inviteRecordDo) Save(values ...*model.InviteRecord) error {
	if len(values) == 0 {
		return nil
	}
	return i.DO.Save(values)
}

func (i inviteRecordDo) First() (*model.InviteRecord, error) {
	if result, err := i.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.InviteRecord), nil
	}
}

func (i inviteRecordDo) Take() (*model.InviteRecord, error) {
	if result, err := i.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.InviteRecord), nil
	}
}

func (i inviteRecordDo) Last() (*model.InviteRecord, error) {
	if result, err := i.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.InviteRecord), nil
	}
}

func (i inviteRecordDo) Find() ([]*model.InviteRecord, error) {
	result, err := i.DO.Find()
	return result.([]*model.InviteRecord), err
}

func (i inviteRecordDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.InviteRecord, err error) {
	buf := make([]*model.InviteRecord, 0, batchSize)
	err = i.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (i inviteRecordDo) FindInBatches(result *[]*model.InviteRecord, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return i.DO.FindInBatches(result, batchSize, fc)
}

func (i inviteRecordDo) Attrs(attrs ...field.AssignExpr) IInviteRecordDo {
	return i.withDO(i.DO.Attrs(attrs...))
}

func (i inviteRecordDo) Assign(attrs ...field.AssignExpr) IInviteRecordDo {
	return i.withDO(i.DO.Assign(attrs...))
}

func (i inviteRecordDo) Joins(fields ...field.RelationField) IInviteRecordDo {
	for _, _f := range fields {
		i = *i.withDO(i.DO.Joins(_f))
	}
	return &i
}

func (i inviteRecordDo) Preload(fields ...field.RelationField) IInviteRecordDo {
	for _, _f := range fields {
		i = *i.withDO(i.DO.Preload(_f))
	}
	return &i
}

func (i inviteRecordDo) FirstOrInit() (*model.InviteRecord, error) {
	if result, err := i.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.InviteRecord), nil
	}
}

func (i inviteRecordDo) FirstOrCreate() (*model.InviteRecord, error) {
	if result, err := i.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.InviteRecord), nil
	}
}

func (i inviteRecordDo) FindByPage(offset int, limit int) (result []*model.InviteRecord, count int64, err error) {
	result, err = i.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = i.Offset(-1).Limit(-1).Count()
	return
}

func (i inviteRecordDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = i.Count()
	if err != nil {
		return
	}

	err = i.Offset(offset).Limit(limit).Scan(result)
	return
}

func (i inviteRecordDo) Scan(result interface{}) (err error) {
	return i.DO.Scan(result)
}

func (i inviteRecordDo) Delete(models ...*model.InviteRecord) (result gen.ResultInfo, err error) {
	return i.DO.Delete(models)
}

func (i *inviteRecordDo) withDO(do gen.Dao) *inviteRecordDo {
	i.DO = *do.(*gen.DO)
	return i
}
