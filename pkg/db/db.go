package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"goweb/pkg/config"
	gowebLogger "goweb/pkg/logger"
)

// ==================== 数据库管理器 ====================

// Manager 数据库管理器
type Manager struct {
	config *config.DatabaseConfig
	logger gowebLogger.Logger
	db     *gorm.DB
}

// NewManager 创建数据库管理器
func NewManager(cfg *config.Config, log gowebLogger.Logger) *Manager {
	dbConfig := cfg.GetDatabaseConfig()
	return &Manager{
		config: &dbConfig,
		logger: log,
	}
}

// Connect 连接数据库
func (m *Manager) Connect() error {
	var dialector gorm.Dialector
	var dsn string

	switch m.config.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			m.config.Username,
			m.config.Password,
			m.config.Host,
			m.config.Port,
			m.config.Database,
			m.config.Charset,
		)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
			m.config.Host,
			m.config.Username,
			m.config.Password,
			m.config.Database,
			m.config.Port,
			m.config.Timezone,
		)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dsn = m.config.Database
		dialector = sqlite.Open(dsn)
	default:
		return fmt.Errorf("unsupported database driver: %s", m.config.Driver)
	}

	// 配置 GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 根据日志级别设置 GORM 日志
	switch m.config.LogLevel {
	case "debug":
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	case "info":
		gormConfig.Logger = logger.Default.LogMode(logger.Warn)
	case "warn", "error":
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	case "silent":
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(m.config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(m.config.ConnMaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	m.db = db
	m.logger.Info("Database connected successfully",
		"driver", m.config.Driver,
		"database", m.config.Database,
	)

	return nil
}

// GetDB 获取 GORM 实例
func (m *Manager) GetDB() *gorm.DB {
	return m.db
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	if m.db != nil {
		sqlDB, err := m.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate 自动迁移模型
func (m *Manager) AutoMigrate(models ...interface{}) error {
	if m.db == nil {
		return fmt.Errorf("database not connected")
	}

	err := m.db.AutoMigrate(models...)
	if err != nil {
		m.logger.Error("AutoMigrate failed", "error", err)
		return err
	}

	m.logger.Info("AutoMigrate completed", "models", len(models))
	return nil
}

// Transaction 执行事务
func (m *Manager) Transaction(fn func(*gorm.DB) error) error {
	if m.db == nil {
		return fmt.Errorf("database not connected")
	}

	return m.db.Transaction(fn)
}

// Health 健康检查
func (m *Manager) Health() error {
	if m.db == nil {
		return fmt.Errorf("database not connected")
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// Stats 获取连接池统计信息
func (m *Manager) Stats() map[string]interface{} {
	if m.db == nil {
		return map[string]interface{}{
			"status": "disconnected",
		}
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"status":               "connected",
		"driver":               m.config.Driver,
		"database":             m.config.Database,
		"max_open_conns":       stats.MaxOpenConnections,
		"open_conns":           stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// ==================== 通用仓库 ====================

// Repository 通用仓库接口
type Repository[T any] interface {
	// 基础CRUD操作
	Create(entity *T) error
	GetByID(id uint) (*T, error)
	Update(entity *T) error
	Delete(id uint) error
	DeleteByCondition(condition interface{}) error

	// 查询操作
	Find(condition interface{}) ([]*T, error)
	FindOne(condition interface{}) (*T, error)
	FindByPage(page, pageSize int, condition interface{}) ([]*T, int64, error)
	Count(condition interface{}) (int64, error)
	Exists(condition interface{}) (bool, error)

	// 原生查询
	RawQuery(sql string, args ...interface{}) ([]*T, error)
	RawExec(sql string, args ...interface{}) error

	// 事务操作
	WithTransaction(fn func(Repository[T]) error) error
}

// BaseRepository 基础仓库实现
type BaseRepository[T any] struct {
	db    *gorm.DB
	model T
}

// NewRepository 创建新的仓库实例
func NewRepository[T any](db *gorm.DB) Repository[T] {
	var model T
	return &BaseRepository[T]{
		db:    db,
		model: model,
	}
}

// Create 创建记录
func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// GetByID 根据ID获取记录
func (r *BaseRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update 更新记录
func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete 根据ID删除记录
func (r *BaseRepository[T]) Delete(id uint) error {
	return r.db.Delete(&r.model, id).Error
}

// DeleteByCondition 根据条件删除记录
func (r *BaseRepository[T]) DeleteByCondition(condition interface{}) error {
	return r.db.Where(condition).Delete(&r.model).Error
}

// Find 根据条件查找记录
func (r *BaseRepository[T]) Find(condition interface{}) ([]*T, error) {
	var entities []*T
	query := r.db
	if condition != nil {
		query = query.Where(condition)
	}
	err := query.Find(&entities).Error
	return entities, err
}

// FindOne 根据条件查找单条记录
func (r *BaseRepository[T]) FindOne(condition interface{}) (*T, error) {
	var entity T
	query := r.db
	if condition != nil {
		query = query.Where(condition)
	}
	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindByPage 分页查询
func (r *BaseRepository[T]) FindByPage(page, pageSize int, condition interface{}) ([]*T, int64, error) {
	var entities []*T
	var total int64

	query := r.db.Model(&r.model)
	if condition != nil {
		query = query.Where(condition)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// Count 统计记录数
func (r *BaseRepository[T]) Count(condition interface{}) (int64, error) {
	var count int64
	query := r.db.Model(&r.model)
	if condition != nil {
		query = query.Where(condition)
	}
	err := query.Count(&count).Error
	return count, err
}

// Exists 检查记录是否存在
func (r *BaseRepository[T]) Exists(condition interface{}) (bool, error) {
	count, err := r.Count(condition)
	return count > 0, err
}

// RawQuery 原生查询
func (r *BaseRepository[T]) RawQuery(sql string, args ...interface{}) ([]*T, error) {
	var entities []*T
	err := r.db.Raw(sql, args...).Scan(&entities).Error
	return entities, err
}

// RawExec 原生执行
func (r *BaseRepository[T]) RawExec(sql string, args ...interface{}) error {
	return r.db.Exec(sql, args...).Error
}

// WithTransaction 事务操作
func (r *BaseRepository[T]) WithTransaction(fn func(Repository[T]) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &BaseRepository[T]{
			db:    tx,
			model: r.model,
		}
		return fn(txRepo)
	})
}

// ==================== 查询构建器 ====================

// QueryBuilder 查询构建器
type QueryBuilder struct {
	db    *gorm.DB
	model interface{}
	query *gorm.DB
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder(db *gorm.DB, model interface{}) *QueryBuilder {
	return &QueryBuilder{
		db:    db,
		model: model,
		query: db.Model(model),
	}
}

// Where 添加WHERE条件
func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
	qb.query = qb.query.Where(query, args...)
	return qb
}

// Order 排序
func (qb *QueryBuilder) Order(value interface{}) *QueryBuilder {
	qb.query = qb.query.Order(value)
	return qb
}

// Limit 限制数量
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.query = qb.query.Limit(limit)
	return qb
}

// Offset 偏移量
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.query = qb.query.Offset(offset)
	return qb
}

// Preload 预加载关联
func (qb *QueryBuilder) Preload(query string, args ...interface{}) *QueryBuilder {
	qb.query = qb.query.Preload(query, args...)
	return qb
}

// Find 查找记录
func (qb *QueryBuilder) Find(dest interface{}) error {
	return qb.query.Find(dest).Error
}

// First 查找第一条记录
func (qb *QueryBuilder) First(dest interface{}) error {
	return qb.query.First(dest).Error
}

// Count 统计数量
func (qb *QueryBuilder) Count(count *int64) error {
	return qb.query.Count(count).Error
}

// Delete 删除记录
func (qb *QueryBuilder) Delete() error {
	return qb.query.Delete(qb.model).Error
}

// Update 更新记录
func (qb *QueryBuilder) Update(column string, value interface{}) error {
	return qb.query.Update(column, value).Error
}

// Updates 批量更新
func (qb *QueryBuilder) Updates(values interface{}) error {
	return qb.query.Updates(values).Error
}

// ==================== 自定义类型 ====================

// JSON 自定义JSON类型
type JSON json.RawMessage

// Scan 实现sql.Scanner接口
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = JSON(v)
	case string:
		*j = JSON(v)
	default:
		return fmt.Errorf("cannot scan %T into JSON", value)
	}

	return nil
}

// Value 实现driver.Valuer接口
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return []byte(j), nil
}

// MarshalJSON 实现json.Marshaler接口
func (j JSON) MarshalJSON() ([]byte, error) {
	return []byte(j), nil
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (j *JSON) UnmarshalJSON(data []byte) error {
	*j = JSON(data)
	return nil
}

// String 实现fmt.Stringer接口
func (j JSON) String() string {
	return string(j)
}

// Unmarshal 解析JSON到指定结构体
func (j JSON) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(j), v)
}

// Marshal 将结构体序列化为JSON
func (j *JSON) Marshal(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	*j = JSON(data)
	return nil
}

// IsNull 检查是否为null
func (j JSON) IsNull() bool {
	return j == nil || string(j) == "null"
}

// IsEmpty 检查是否为空
func (j JSON) IsEmpty() bool {
	return j == nil || len(j) == 0 || string(j) == "null" || string(j) == "{}" || string(j) == "[]"
}

// BaseModel 基础模型
type BaseModel struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// GetID 获取ID
func (bm BaseModel) GetID() uint {
	return bm.ID
}

// SetID 设置ID
func (bm *BaseModel) SetID(id uint) {
	bm.ID = id
}

// GetCreatedAt 获取创建时间
func (bm BaseModel) GetCreatedAt() time.Time {
	return bm.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (bm BaseModel) GetUpdatedAt() time.Time {
	return bm.UpdatedAt
}

// GetDeletedAt 获取删除时间
func (bm BaseModel) GetDeletedAt() *time.Time {
	return bm.DeletedAt
}

// IsDeleted 检查是否已删除
func (bm BaseModel) IsDeleted() bool {
	return bm.DeletedAt != nil
}

// Delete 软删除
func (bm *BaseModel) Delete() {
	now := time.Now()
	bm.DeletedAt = &now
}

// Restore 恢复删除
func (bm *BaseModel) Restore() {
	bm.DeletedAt = nil
}

// Status 状态类型
type Status int

const (
	StatusInactive Status = 0
	StatusActive   Status = 1
	StatusPending  Status = 2
	StatusDeleted  Status = 3
)

// Scan 实现sql.Scanner接口
func (s *Status) Scan(value interface{}) error {
	if value == nil {
		*s = StatusInactive
		return nil
	}

	switch v := value.(type) {
	case int64:
		*s = Status(v)
	case int:
		*s = Status(v)
	case string:
		// 这里可以添加字符串解析逻辑
		*s = StatusInactive
	default:
		return fmt.Errorf("cannot scan %T into Status", value)
	}

	return nil
}

// Value 实现driver.Valuer接口
func (s Status) Value() (driver.Value, error) {
	return int(s), nil
}

// String 实现fmt.Stringer接口
func (s Status) String() string {
	switch s {
	case StatusInactive:
		return "inactive"
	case StatusActive:
		return "active"
	case StatusPending:
		return "pending"
	case StatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// IsActive 检查是否激活
func (s Status) IsActive() bool {
	return s == StatusActive
}

// IsInactive 检查是否未激活
func (s Status) IsInactive() bool {
	return s == StatusInactive
}

// IsPending 检查是否待处理
func (s Status) IsPending() bool {
	return s == StatusPending
}

// IsDeleted 检查是否已删除
func (s Status) IsDeleted() bool {
	return s == StatusDeleted
}

// SetActive 设置为激活状态
func (s *Status) SetActive() {
	*s = StatusActive
}

// SetInactive 设置为未激活状态
func (s *Status) SetInactive() {
	*s = StatusInactive
}

// SetPending 设置为待处理状态
func (s *Status) SetPending() {
	*s = StatusPending
}

// SetDeleted 设置为已删除状态
func (s *Status) SetDeleted() {
	*s = StatusDeleted
}

// ==================== 工具函数 ====================

// Paginate 分页查询工具函数
func Paginate[T any](db *gorm.DB, page, pageSize int, condition interface{}) ([]*T, int64, error) {
	var entities []*T
	var total int64

	query := db
	if condition != nil {
		query = query.Where(condition)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// Transaction 事务执行工具函数
func Transaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	return db.Transaction(fn)
}

// BatchCreate 批量创建工具函数
func BatchCreate[T any](db *gorm.DB, entities []*T, batchSize int) error {
	return db.CreateInBatches(entities, batchSize).Error
}

// BatchUpdate 批量更新工具函数
func BatchUpdate[T any](db *gorm.DB, condition interface{}, updates map[string]interface{}) error {
	var model T
	return db.Model(&model).Where(condition).Updates(updates).Error
}

// BatchDelete 批量删除工具函数
func BatchDelete[T any](db *gorm.DB, ids []uint) error {
	var model T
	return db.Delete(&model, ids).Error
}

// FirstOrCreate 查找或创建工具函数
func FirstOrCreate[T any](db *gorm.DB, condition interface{}, entity *T) error {
	return db.Where(condition).FirstOrCreate(entity).Error
}

// UpdateOrCreate 更新或创建工具函数
func UpdateOrCreate[T any](db *gorm.DB, condition interface{}, entity *T) error {
	return db.Where(condition).Assign(entity).FirstOrCreate(entity).Error
}

// Pluck 提取字段值工具函数
func Pluck[T any](db *gorm.DB, column string, dest interface{}) error {
	var model T
	return db.Model(&model).Pluck(column, dest).Error
}

// Exists 检查记录是否存在工具函数
func Exists[T any](db *gorm.DB, condition interface{}) (bool, error) {
	var model T
	var count int64
	query := db.Model(&model)
	if condition != nil {
		query = query.Where(condition)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// Count 统计记录数工具函数
func Count[T any](db *gorm.DB, condition interface{}) (int64, error) {
	var model T
	var count int64
	query := db.Model(&model)
	if condition != nil {
		query = query.Where(condition)
	}
	err := query.Count(&count).Error
	return count, err
}

// Find 查找记录工具函数
func Find[T any](db *gorm.DB, condition interface{}) ([]*T, error) {
	var entities []*T
	query := db
	if condition != nil {
		query = query.Where(condition)
	}
	err := query.Find(&entities).Error
	return entities, err
}

// FindOne 查找单条记录工具函数
func FindOne[T any](db *gorm.DB, condition interface{}) (*T, error) {
	var entity T
	query := db
	if condition != nil {
		query = query.Where(condition)
	}
	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetByID 根据ID获取记录工具函数
func GetByID[T any](db *gorm.DB, id uint) (*T, error) {
	var entity T
	err := db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Create 创建记录工具函数
func Create[T any](db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

// Update 更新记录工具函数
func Update[T any](db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

// Delete 删除记录工具函数
func Delete[T any](db *gorm.DB, id uint) error {
	var model T
	return db.Delete(&model, id).Error
}

// DeleteByCondition 根据条件删除记录工具函数
func DeleteByCondition[T any](db *gorm.DB, condition interface{}) error {
	var model T
	return db.Where(condition).Delete(&model).Error
}

// RawQuery 原生查询工具函数
func RawQuery[T any](db *gorm.DB, sql string, args ...interface{}) ([]*T, error) {
	var entities []*T
	err := db.Raw(sql, args...).Scan(&entities).Error
	return entities, err
}

// RawExec 原生执行工具函数
func RawExec(db *gorm.DB, sql string, args ...interface{}) error {
	return db.Exec(sql, args...).Error
}

// ==================== 查询条件构建器 ====================

// Condition 查询条件构建器
type Condition struct {
	query string
	args  []interface{}
}

// NewCondition 创建查询条件
func NewCondition() *Condition {
	return &Condition{
		args: make([]interface{}, 0),
	}
}

// Where 添加WHERE条件
func (c *Condition) Where(query string, args ...interface{}) *Condition {
	if c.query != "" {
		c.query += " AND "
	}
	c.query += query
	c.args = append(c.args, args...)
	return c
}

// Or 添加OR条件
func (c *Condition) Or(query string, args ...interface{}) *Condition {
	if c.query != "" {
		c.query += " OR "
	}
	c.query += query
	c.args = append(c.args, args...)
	return c
}

// WhereIn 添加IN条件
func (c *Condition) WhereIn(column string, values interface{}) *Condition {
	return c.Where(fmt.Sprintf("%s IN ?", column), values)
}

// WhereNotIn 添加NOT IN条件
func (c *Condition) WhereNotIn(column string, values interface{}) *Condition {
	return c.Where(fmt.Sprintf("%s NOT IN ?", column), values)
}

// WhereLike 添加LIKE条件
func (c *Condition) WhereLike(column string, value string) *Condition {
	return c.Where(fmt.Sprintf("%s LIKE ?", column), "%"+value+"%")
}

// WhereBetween 添加BETWEEN条件
func (c *Condition) WhereBetween(column string, start, end interface{}) *Condition {
	return c.Where(fmt.Sprintf("%s BETWEEN ? AND ?", column), start, end)
}

// WhereNull 添加IS NULL条件
func (c *Condition) WhereNull(column string) *Condition {
	return c.Where(fmt.Sprintf("%s IS NULL", column))
}

// WhereNotNull 添加IS NOT NULL条件
func (c *Condition) WhereNotNull(column string) *Condition {
	return c.Where(fmt.Sprintf("%s IS NOT NULL", column))
}

// WhereDate 添加日期条件
func (c *Condition) WhereDate(column string, date string) *Condition {
	return c.Where(fmt.Sprintf("DATE(%s) = ?", column), date)
}

// WhereTime 添加时间条件
func (c *Condition) WhereTime(column string, time string) *Condition {
	return c.Where(fmt.Sprintf("TIME(%s) = ?", column), time)
}

// WhereYear 添加年份条件
func (c *Condition) WhereYear(column string, year int) *Condition {
	return c.Where(fmt.Sprintf("YEAR(%s) = ?", column), year)
}

// WhereMonth 添加月份条件
func (c *Condition) WhereMonth(column string, month int) *Condition {
	return c.Where(fmt.Sprintf("MONTH(%s) = ?", column), month)
}

// WhereDay 添加日期条件
func (c *Condition) WhereDay(column string, day int) *Condition {
	return c.Where(fmt.Sprintf("DAY(%s) = ?", column), day)
}

// GetQuery 获取查询字符串
func (c *Condition) GetQuery() string {
	return c.query
}

// GetArgs 获取参数
func (c *Condition) GetArgs() []interface{} {
	return c.args
}

// IsEmpty 检查是否为空
func (c *Condition) IsEmpty() bool {
	return c.query == ""
}

// ==================== 排序构建器 ====================

// Order 排序构建器
type Order struct {
	orders []string
}

// NewOrder 创建排序构建器
func NewOrder() *Order {
	return &Order{
		orders: make([]string, 0),
	}
}

// Asc 添加升序排序
func (o *Order) Asc(column string) *Order {
	o.orders = append(o.orders, fmt.Sprintf("%s ASC", column))
	return o
}

// Desc 添加降序排序
func (o *Order) Desc(column string) *Order {
	o.orders = append(o.orders, fmt.Sprintf("%s DESC", column))
	return o
}

// Custom 添加自定义排序
func (o *Order) Custom(order string) *Order {
	o.orders = append(o.orders, order)
	return o
}

// GetOrder 获取排序字符串
func (o *Order) GetOrder() string {
	return strings.Join(o.orders, ", ")
}

// IsEmpty 检查是否为空
func (o *Order) IsEmpty() bool {
	return len(o.orders) == 0
}

// ==================== 分页信息 ====================

// PageInfo 分页信息
type PageInfo struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}

// NewPageInfo 创建分页信息
func NewPageInfo(page, pageSize int, total int64) *PageInfo {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	return &PageInfo{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Pages:    pages,
	}
}

// HasNext 是否有下一页
func (p *PageInfo) HasNext() bool {
	return p.Page < p.Pages
}

// HasPrev 是否有上一页
func (p *PageInfo) HasPrev() bool {
	return p.Page > 1
}

// GetOffset 获取偏移量
func (p *PageInfo) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// ==================== 查询结果 ====================

// QueryResult 查询结果
type QueryResult[T any] struct {
	Data     []*T      `json:"data"`
	PageInfo *PageInfo `json:"page_info"`
}

// NewQueryResult 创建查询结果
func NewQueryResult[T any](data []*T, page, pageSize int, total int64) *QueryResult[T] {
	return &QueryResult[T]{
		Data:     data,
		PageInfo: NewPageInfo(page, pageSize, total),
	}
}

// ==================== 错误处理 ====================

// DBError 数据库错误
type DBError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error 实现error接口
func (e *DBError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap 实现错误包装
func (e *DBError) Unwrap() error {
	return e.Err
}

// NewDBError 创建数据库错误
func NewDBError(code, message string, err error) *DBError {
	return &DBError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 常见错误代码
const (
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeDuplicate     = "DUPLICATE"
	ErrCodeInvalidInput  = "INVALID_INPUT"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeForbidden     = "FORBIDDEN"
	ErrCodeInternalError = "INTERNAL_ERROR"
)

// 常见错误
var (
	ErrNotFound      = NewDBError(ErrCodeNotFound, "记录不存在", nil)
	ErrDuplicate     = NewDBError(ErrCodeDuplicate, "记录已存在", nil)
	ErrInvalidInput  = NewDBError(ErrCodeInvalidInput, "输入参数无效", nil)
	ErrUnauthorized  = NewDBError(ErrCodeUnauthorized, "未授权", nil)
	ErrForbidden     = NewDBError(ErrCodeForbidden, "禁止访问", nil)
	ErrInternalError = NewDBError(ErrCodeInternalError, "内部错误", nil)
)

// ==================== 工具函数 ====================

// IsNotFound 检查是否为未找到错误
func IsNotFound(err error) bool {
	if dbErr, ok := err.(*DBError); ok {
		return dbErr.Code == ErrCodeNotFound
	}
	return err == gorm.ErrRecordNotFound
}

// IsDuplicate 检查是否为重复错误
func IsDuplicate(err error) bool {
	if dbErr, ok := err.(*DBError); ok {
		return dbErr.Code == ErrCodeDuplicate
	}
	// 检查GORM的重复键错误
	return strings.Contains(err.Error(), "Duplicate entry")
}

// WrapError 包装错误
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return NewDBError(ErrCodeInternalError, message, err)
}

// WrapNotFound 包装未找到错误
func WrapNotFound(err error, message string) error {
	if err == nil {
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		return NewDBError(ErrCodeNotFound, message, err)
	}
	return NewDBError(ErrCodeNotFound, message, err)
}

// WrapDuplicate 包装重复错误
func WrapDuplicate(err error, message string) error {
	if err == nil {
		return nil
	}
	return NewDBError(ErrCodeDuplicate, message, err)
}

// WrapInvalidInput 包装无效输入错误
func WrapInvalidInput(err error, message string) error {
	if err == nil {
		return nil
	}
	return NewDBError(ErrCodeInvalidInput, message, err)
}
