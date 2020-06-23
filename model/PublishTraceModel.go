package model

import sq "github.com/Masterminds/squirrel"

const publishTraceTable = "`publish_trace`"

// PublishTrace mysql table for rsync trace
type PublishTrace struct {
	ID            int64  `json:"id"`
	Token         string `json:"token"`
	ProjectID     int64  `json:"projectId"`
	ProjectName   string `json:"projectName"`
	Detail        string `json:"detail"`
	State         int    `json:"state"`
	PublisherID   int64  `json:"publisherId"`
	PublisherName string `json:"publisherName"`
	Type          int    `json:"type"`
	Ext           string `json:"ext"`
	PublishState  int    `json:"publishState"`
	InsertTime    string `json:"insertTime"`
	UpdateTime    string `json:"updateTime"`
}

// PublishTraces PublishTrace slice
type PublishTraces []PublishTrace

// publish trace type
const (
	BeforePull = 1

	Pull = 2

	AfterPull = 3

	BeforeDeploy = 4

	Deploy = 5

	AfterDeploy = 6
)

// AddRow add one row to table deploy and add id to deploy.ID
func (pt PublishTrace) AddRow() (int64, error) {
	result, err := sq.
		Insert(publishTraceTable).
		Columns("token", "project_id", "project_name", "detail", "state", "publisher_id", "publisher_name", "type", "ext").
		Values(pt.Token, pt.ProjectID, pt.ProjectName, pt.Detail, pt.State, pt.PublisherID, pt.PublisherName, pt.Type, pt.Ext).
		RunWith(DB).
		Exec()

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

// GetListByToken PublishTrace row
func (pt PublishTrace) GetListByToken() (PublishTraces, error) {
	rows, err := sq.
		Select("id, token, project_id, project_name, detail, state, publisher_id, publisher_name, type, ext, insert_time, update_time").
		From(publishTraceTable).
		Where(sq.Eq{"token": pt.Token}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	publishTraces := PublishTraces{}
	for rows.Next() {
		var publishTrace PublishTrace

		if err := rows.Scan(
			&publishTrace.ID,
			&publishTrace.Token,
			&publishTrace.ProjectID,
			&publishTrace.ProjectName,
			&publishTrace.Detail,
			&publishTrace.State,
			&publishTrace.PublisherID,
			&publishTrace.PublisherName,
			&publishTrace.Type,
			&publishTrace.Ext,
			&publishTrace.InsertTime,
			&publishTrace.UpdateTime); err != nil {
			return nil, err
		}
		publishTraces = append(publishTraces, publishTrace)
	}
	return publishTraces, nil
}

// GetPreview PublishTrace row
func (pt PublishTrace) GetPreview(pagination Pagination) (PublishTraces, Pagination, error) {
	builder := sq.
		Select("id, token, project_id, project_name, detail, state, publisher_id, publisher_name, type, ext, insert_time, update_time").
		Column("!EXISTS (SELECT id FROM " + publishTraceTable + " AS pt where pt.state = 0 AND pt.token = publish_trace.token) as publish_state").
		From(publishTraceTable).
		Where(sq.Eq{"type": Pull})
	if pt.ProjectID != 0 {
		builder = builder.Where(sq.Eq{"project_id": pt.ProjectID})
	}
	if pt.PublisherID != 0 {
		builder = builder.Where(sq.Eq{"publisher_id": pt.PublisherID})
	}
	if pt.PublishState != -1 {
		builder = builder.Having(sq.Eq{"publish_state": pt.PublishState})
	}
	rows, err := builder.RunWith(DB).
		OrderBy("update_time DESC").
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		Query()
	if err != nil {
		return nil, pagination, err
	}
	publishTraces := PublishTraces{}
	for rows.Next() {
		var publishTrace PublishTrace

		if err := rows.Scan(
			&publishTrace.ID,
			&publishTrace.Token,
			&publishTrace.ProjectID,
			&publishTrace.ProjectName,
			&publishTrace.Detail,
			&publishTrace.State,
			&publishTrace.PublisherID,
			&publishTrace.PublisherName,
			&publishTrace.Type,
			&publishTrace.Ext,
			&publishTrace.InsertTime,
			&publishTrace.UpdateTime,
			&publishTrace.PublishState); err != nil {
			return nil, pagination, err
		}
		publishTraces = append(publishTraces, publishTrace)
	}

	builder = sq.
		Select("COUNT(*) AS count").
		From(publishTraceTable).
		Where(sq.Eq{"type": Pull})

	if pt.ProjectID != 0 {
		builder = builder.Where(sq.Eq{"project_id": pt.ProjectID})
	}
	if pt.PublisherID != 0 {
		builder = builder.Where(sq.Eq{"publisher_id": pt.PublisherID})
	}
	if pt.PublishState == 0 {
		builder = builder.Where("EXISTS (SELECT id FROM " + publishTraceTable + " AS pt where pt.state = 0 AND pt.token = publish_trace.token)")
	} else if pt.PublishState == 1 {
		builder = builder.Where("! EXISTS (SELECT id FROM " + publishTraceTable + " AS pt where pt.state = 0 AND pt.token = publish_trace.token)")
	}

	err = builder.RunWith(DB).
		QueryRow().
		Scan(&pagination.Total)
	if err != nil {
		return nil, pagination, err
	}
	return publishTraces, pagination, nil
}
