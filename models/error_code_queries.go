// Code generated by SQLBoiler 4.13.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// ErrorCodeQuery is an object representing the database table.
type ErrorCodeQuery struct {
	ID            string            `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserDeviceID  string            `boil:"user_device_id" json:"user_device_id" toml:"user_device_id" yaml:"user_device_id"`
	ErrorCodes    types.StringArray `boil:"error_codes" json:"error_codes" toml:"error_codes" yaml:"error_codes"`
	QueryResponse string            `boil:"query_response" json:"query_response" toml:"query_response" yaml:"query_response"`
	CreatedAt     time.Time         `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt     time.Time         `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *errorCodeQueryR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L errorCodeQueryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ErrorCodeQueryColumns = struct {
	ID            string
	UserDeviceID  string
	ErrorCodes    string
	QueryResponse string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "id",
	UserDeviceID:  "user_device_id",
	ErrorCodes:    "error_codes",
	QueryResponse: "query_response",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

var ErrorCodeQueryTableColumns = struct {
	ID            string
	UserDeviceID  string
	ErrorCodes    string
	QueryResponse string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "error_code_queries.id",
	UserDeviceID:  "error_code_queries.user_device_id",
	ErrorCodes:    "error_code_queries.error_codes",
	QueryResponse: "error_code_queries.query_response",
	CreatedAt:     "error_code_queries.created_at",
	UpdatedAt:     "error_code_queries.updated_at",
}

// Generated where

type whereHelpertypes_StringArray struct{ field string }

func (w whereHelpertypes_StringArray) EQ(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertypes_StringArray) NEQ(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertypes_StringArray) LT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_StringArray) LTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_StringArray) GT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_StringArray) GTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var ErrorCodeQueryWhere = struct {
	ID            whereHelperstring
	UserDeviceID  whereHelperstring
	ErrorCodes    whereHelpertypes_StringArray
	QueryResponse whereHelperstring
	CreatedAt     whereHelpertime_Time
	UpdatedAt     whereHelpertime_Time
}{
	ID:            whereHelperstring{field: "\"devices_api\".\"error_code_queries\".\"id\""},
	UserDeviceID:  whereHelperstring{field: "\"devices_api\".\"error_code_queries\".\"user_device_id\""},
	ErrorCodes:    whereHelpertypes_StringArray{field: "\"devices_api\".\"error_code_queries\".\"error_codes\""},
	QueryResponse: whereHelperstring{field: "\"devices_api\".\"error_code_queries\".\"query_response\""},
	CreatedAt:     whereHelpertime_Time{field: "\"devices_api\".\"error_code_queries\".\"created_at\""},
	UpdatedAt:     whereHelpertime_Time{field: "\"devices_api\".\"error_code_queries\".\"updated_at\""},
}

// ErrorCodeQueryRels is where relationship names are stored.
var ErrorCodeQueryRels = struct {
	UserDevice string
}{
	UserDevice: "UserDevice",
}

// errorCodeQueryR is where relationships are stored.
type errorCodeQueryR struct {
	UserDevice *UserDevice `boil:"UserDevice" json:"UserDevice" toml:"UserDevice" yaml:"UserDevice"`
}

// NewStruct creates a new relationship struct
func (*errorCodeQueryR) NewStruct() *errorCodeQueryR {
	return &errorCodeQueryR{}
}

func (r *errorCodeQueryR) GetUserDevice() *UserDevice {
	if r == nil {
		return nil
	}
	return r.UserDevice
}

// errorCodeQueryL is where Load methods for each relationship are stored.
type errorCodeQueryL struct{}

var (
	errorCodeQueryAllColumns            = []string{"id", "user_device_id", "error_codes", "query_response", "created_at", "updated_at"}
	errorCodeQueryColumnsWithoutDefault = []string{"id", "user_device_id", "error_codes", "query_response"}
	errorCodeQueryColumnsWithDefault    = []string{"created_at", "updated_at"}
	errorCodeQueryPrimaryKeyColumns     = []string{"id"}
	errorCodeQueryGeneratedColumns      = []string{}
)

type (
	// ErrorCodeQuerySlice is an alias for a slice of pointers to ErrorCodeQuery.
	// This should almost always be used instead of []ErrorCodeQuery.
	ErrorCodeQuerySlice []*ErrorCodeQuery
	// ErrorCodeQueryHook is the signature for custom ErrorCodeQuery hook methods
	ErrorCodeQueryHook func(context.Context, boil.ContextExecutor, *ErrorCodeQuery) error

	errorCodeQueryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	errorCodeQueryType                 = reflect.TypeOf(&ErrorCodeQuery{})
	errorCodeQueryMapping              = queries.MakeStructMapping(errorCodeQueryType)
	errorCodeQueryPrimaryKeyMapping, _ = queries.BindMapping(errorCodeQueryType, errorCodeQueryMapping, errorCodeQueryPrimaryKeyColumns)
	errorCodeQueryInsertCacheMut       sync.RWMutex
	errorCodeQueryInsertCache          = make(map[string]insertCache)
	errorCodeQueryUpdateCacheMut       sync.RWMutex
	errorCodeQueryUpdateCache          = make(map[string]updateCache)
	errorCodeQueryUpsertCacheMut       sync.RWMutex
	errorCodeQueryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var errorCodeQueryAfterSelectHooks []ErrorCodeQueryHook

var errorCodeQueryBeforeInsertHooks []ErrorCodeQueryHook
var errorCodeQueryAfterInsertHooks []ErrorCodeQueryHook

var errorCodeQueryBeforeUpdateHooks []ErrorCodeQueryHook
var errorCodeQueryAfterUpdateHooks []ErrorCodeQueryHook

var errorCodeQueryBeforeDeleteHooks []ErrorCodeQueryHook
var errorCodeQueryAfterDeleteHooks []ErrorCodeQueryHook

var errorCodeQueryBeforeUpsertHooks []ErrorCodeQueryHook
var errorCodeQueryAfterUpsertHooks []ErrorCodeQueryHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *ErrorCodeQuery) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range errorCodeQueryAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *ErrorCodeQuery) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range errorCodeQueryBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *ErrorCodeQuery) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range errorCodeQueryAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *ErrorCodeQuery) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range errorCodeQueryBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *ErrorCodeQuery) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range errorCodeQueryAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *ErrorCodeQuery) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range errorCodeQueryBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *ErrorCodeQuery) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range errorCodeQueryAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *ErrorCodeQuery) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range errorCodeQueryBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *ErrorCodeQuery) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range errorCodeQueryAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddErrorCodeQueryHook registers your hook function for all future operations.
func AddErrorCodeQueryHook(hookPoint boil.HookPoint, errorCodeQueryHook ErrorCodeQueryHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		errorCodeQueryAfterSelectHooks = append(errorCodeQueryAfterSelectHooks, errorCodeQueryHook)
	case boil.BeforeInsertHook:
		errorCodeQueryBeforeInsertHooks = append(errorCodeQueryBeforeInsertHooks, errorCodeQueryHook)
	case boil.AfterInsertHook:
		errorCodeQueryAfterInsertHooks = append(errorCodeQueryAfterInsertHooks, errorCodeQueryHook)
	case boil.BeforeUpdateHook:
		errorCodeQueryBeforeUpdateHooks = append(errorCodeQueryBeforeUpdateHooks, errorCodeQueryHook)
	case boil.AfterUpdateHook:
		errorCodeQueryAfterUpdateHooks = append(errorCodeQueryAfterUpdateHooks, errorCodeQueryHook)
	case boil.BeforeDeleteHook:
		errorCodeQueryBeforeDeleteHooks = append(errorCodeQueryBeforeDeleteHooks, errorCodeQueryHook)
	case boil.AfterDeleteHook:
		errorCodeQueryAfterDeleteHooks = append(errorCodeQueryAfterDeleteHooks, errorCodeQueryHook)
	case boil.BeforeUpsertHook:
		errorCodeQueryBeforeUpsertHooks = append(errorCodeQueryBeforeUpsertHooks, errorCodeQueryHook)
	case boil.AfterUpsertHook:
		errorCodeQueryAfterUpsertHooks = append(errorCodeQueryAfterUpsertHooks, errorCodeQueryHook)
	}
}

// One returns a single errorCodeQuery record from the query.
func (q errorCodeQueryQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ErrorCodeQuery, error) {
	o := &ErrorCodeQuery{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for error_code_queries")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all ErrorCodeQuery records from the query.
func (q errorCodeQueryQuery) All(ctx context.Context, exec boil.ContextExecutor) (ErrorCodeQuerySlice, error) {
	var o []*ErrorCodeQuery

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ErrorCodeQuery slice")
	}

	if len(errorCodeQueryAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all ErrorCodeQuery records in the query.
func (q errorCodeQueryQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count error_code_queries rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q errorCodeQueryQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if error_code_queries exists")
	}

	return count > 0, nil
}

// UserDevice pointed to by the foreign key.
func (o *ErrorCodeQuery) UserDevice(mods ...qm.QueryMod) userDeviceQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserDeviceID),
	}

	queryMods = append(queryMods, mods...)

	return UserDevices(queryMods...)
}

// LoadUserDevice allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (errorCodeQueryL) LoadUserDevice(ctx context.Context, e boil.ContextExecutor, singular bool, maybeErrorCodeQuery interface{}, mods queries.Applicator) error {
	var slice []*ErrorCodeQuery
	var object *ErrorCodeQuery

	if singular {
		var ok bool
		object, ok = maybeErrorCodeQuery.(*ErrorCodeQuery)
		if !ok {
			object = new(ErrorCodeQuery)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeErrorCodeQuery)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeErrorCodeQuery))
			}
		}
	} else {
		s, ok := maybeErrorCodeQuery.(*[]*ErrorCodeQuery)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeErrorCodeQuery)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeErrorCodeQuery))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &errorCodeQueryR{}
		}
		args = append(args, object.UserDeviceID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &errorCodeQueryR{}
			}

			for _, a := range args {
				if a == obj.UserDeviceID {
					continue Outer
				}
			}

			args = append(args, obj.UserDeviceID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`devices_api.user_devices`),
		qm.WhereIn(`devices_api.user_devices.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load UserDevice")
	}

	var resultSlice []*UserDevice
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice UserDevice")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for user_devices")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for user_devices")
	}

	if len(errorCodeQueryAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.UserDevice = foreign
		if foreign.R == nil {
			foreign.R = &userDeviceR{}
		}
		foreign.R.ErrorCodeQueries = append(foreign.R.ErrorCodeQueries, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserDeviceID == foreign.ID {
				local.R.UserDevice = foreign
				if foreign.R == nil {
					foreign.R = &userDeviceR{}
				}
				foreign.R.ErrorCodeQueries = append(foreign.R.ErrorCodeQueries, local)
				break
			}
		}
	}

	return nil
}

// SetUserDevice of the errorCodeQuery to the related item.
// Sets o.R.UserDevice to related.
// Adds o to related.R.ErrorCodeQueries.
func (o *ErrorCodeQuery) SetUserDevice(ctx context.Context, exec boil.ContextExecutor, insert bool, related *UserDevice) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"devices_api\".\"error_code_queries\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_device_id"}),
		strmangle.WhereClause("\"", "\"", 2, errorCodeQueryPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserDeviceID = related.ID
	if o.R == nil {
		o.R = &errorCodeQueryR{
			UserDevice: related,
		}
	} else {
		o.R.UserDevice = related
	}

	if related.R == nil {
		related.R = &userDeviceR{
			ErrorCodeQueries: ErrorCodeQuerySlice{o},
		}
	} else {
		related.R.ErrorCodeQueries = append(related.R.ErrorCodeQueries, o)
	}

	return nil
}

// ErrorCodeQueries retrieves all the records using an executor.
func ErrorCodeQueries(mods ...qm.QueryMod) errorCodeQueryQuery {
	mods = append(mods, qm.From("\"devices_api\".\"error_code_queries\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"devices_api\".\"error_code_queries\".*"})
	}

	return errorCodeQueryQuery{q}
}

// FindErrorCodeQuery retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindErrorCodeQuery(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*ErrorCodeQuery, error) {
	errorCodeQueryObj := &ErrorCodeQuery{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"devices_api\".\"error_code_queries\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, errorCodeQueryObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from error_code_queries")
	}

	if err = errorCodeQueryObj.doAfterSelectHooks(ctx, exec); err != nil {
		return errorCodeQueryObj, err
	}

	return errorCodeQueryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ErrorCodeQuery) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no error_code_queries provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(errorCodeQueryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	errorCodeQueryInsertCacheMut.RLock()
	cache, cached := errorCodeQueryInsertCache[key]
	errorCodeQueryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			errorCodeQueryAllColumns,
			errorCodeQueryColumnsWithDefault,
			errorCodeQueryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(errorCodeQueryType, errorCodeQueryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(errorCodeQueryType, errorCodeQueryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"devices_api\".\"error_code_queries\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"devices_api\".\"error_code_queries\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into error_code_queries")
	}

	if !cached {
		errorCodeQueryInsertCacheMut.Lock()
		errorCodeQueryInsertCache[key] = cache
		errorCodeQueryInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the ErrorCodeQuery.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ErrorCodeQuery) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	errorCodeQueryUpdateCacheMut.RLock()
	cache, cached := errorCodeQueryUpdateCache[key]
	errorCodeQueryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			errorCodeQueryAllColumns,
			errorCodeQueryPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update error_code_queries, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"devices_api\".\"error_code_queries\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, errorCodeQueryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(errorCodeQueryType, errorCodeQueryMapping, append(wl, errorCodeQueryPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update error_code_queries row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for error_code_queries")
	}

	if !cached {
		errorCodeQueryUpdateCacheMut.Lock()
		errorCodeQueryUpdateCache[key] = cache
		errorCodeQueryUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q errorCodeQueryQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for error_code_queries")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for error_code_queries")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ErrorCodeQuerySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), errorCodeQueryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"devices_api\".\"error_code_queries\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, errorCodeQueryPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in errorCodeQuery slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all errorCodeQuery")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ErrorCodeQuery) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no error_code_queries provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(errorCodeQueryColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	errorCodeQueryUpsertCacheMut.RLock()
	cache, cached := errorCodeQueryUpsertCache[key]
	errorCodeQueryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			errorCodeQueryAllColumns,
			errorCodeQueryColumnsWithDefault,
			errorCodeQueryColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			errorCodeQueryAllColumns,
			errorCodeQueryPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert error_code_queries, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(errorCodeQueryPrimaryKeyColumns))
			copy(conflict, errorCodeQueryPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"devices_api\".\"error_code_queries\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(errorCodeQueryType, errorCodeQueryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(errorCodeQueryType, errorCodeQueryMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert error_code_queries")
	}

	if !cached {
		errorCodeQueryUpsertCacheMut.Lock()
		errorCodeQueryUpsertCache[key] = cache
		errorCodeQueryUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single ErrorCodeQuery record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ErrorCodeQuery) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ErrorCodeQuery provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), errorCodeQueryPrimaryKeyMapping)
	sql := "DELETE FROM \"devices_api\".\"error_code_queries\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from error_code_queries")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for error_code_queries")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q errorCodeQueryQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no errorCodeQueryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from error_code_queries")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for error_code_queries")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ErrorCodeQuerySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(errorCodeQueryBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), errorCodeQueryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"devices_api\".\"error_code_queries\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, errorCodeQueryPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from errorCodeQuery slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for error_code_queries")
	}

	if len(errorCodeQueryAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ErrorCodeQuery) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindErrorCodeQuery(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ErrorCodeQuerySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ErrorCodeQuerySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), errorCodeQueryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"devices_api\".\"error_code_queries\".* FROM \"devices_api\".\"error_code_queries\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, errorCodeQueryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ErrorCodeQuerySlice")
	}

	*o = slice

	return nil
}

// ErrorCodeQueryExists checks if the ErrorCodeQuery row exists.
func ErrorCodeQueryExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"devices_api\".\"error_code_queries\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if error_code_queries exists")
	}

	return exists, nil
}