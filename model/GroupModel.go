package model

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

const groupTable = "`group`"

// Group mysql table group
type Group struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	InsertTime string  `json:"insertTime"`
	UpdateTime string  `json:"updateTime"`
}

// Groups many Group
type Groups []Group

// AddRow add one row to table Group
func (g Group) AddRow() (int64, error) {
	result, err := sq.
		Insert(groupTable).
		Columns("name").
		Values(g.Name).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow edit one row to table group
func (g Group) EditRow() error {
	_, err := sq.
		Update(groupTable).
		Set("name", g.Name).
		Where(sq.Eq{"id": g.ID}).
		RunWith(DB).
		Exec()
	return err
}

// Remove Server
func (g Group) Remove() error {
	tx, err := DB.Begin()
	if err != nil {
		return errors.New("开启事务失败")
	}

	_, err = sq.
		Delete(groupTable).
		Where(sq.Eq{"id": g.ID}).
		RunWith(tx).
		Exec()

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sq.
		Update(projectTable).
		Set("group_id", 0).
		Where(sq.Eq{"group_id": g.ID}).
		RunWith(tx).
		Exec()

	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.New("事务提交失败")
	}
	return nil
}

// GetList Group row
func (g Group) GetList(pagination Pagination, groupIDs []string) (Groups, error) {
	builder := sq.
		Select("id, name, insert_time, update_time").
		From(groupTable)
	if len(groupIDs) > 0 {
		builder = builder.Where(sq.Eq{"id": groupIDs})
	}
	rows, err := builder.
		OrderBy("id DESC").
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}

	groups := Groups{}
	for rows.Next() {
		var group Group

		if err := rows.Scan(&group.ID, &group.Name, &group.InsertTime, &group.UpdateTime); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// GetList group total
func (g Group) GetTotal(groupIDs []string) (int64, error) {
	var total int64
	builder := sq.
		Select("COUNT(*) AS count").
		From(groupTable)
	if len(groupIDs) > 0 {
		builder = builder.Where(sq.Eq{"id": groupIDs})
	}
	err := builder.RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetAll Group row
func (g Group) GetAllInGroupIDs(groupIDs []string) (Groups, error) {
	rows, err := sq.
		Select("id, name, insert_time, update_time").
		From(groupTable).
		Where(sq.Eq{"id": groupIDs}).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	projectGroups := Groups{}
	for rows.Next() {
		var projectGroup Group

		if err := rows.Scan(&projectGroup.ID, &projectGroup.Name, &projectGroup.InsertTime, &projectGroup.UpdateTime); err != nil {
			return nil, err
		}
		projectGroups = append(projectGroups, projectGroup)
	}
	return projectGroups, nil
}

// GetAll Group row
func (g Group) GetAll() (Groups, error) {
	rows, err := sq.
		Select("id, name, insert_time, update_time").
		From(groupTable).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	projectGroups := Groups{}
	for rows.Next() {
		var projectGroup Group

		if err := rows.Scan(&projectGroup.ID, &projectGroup.Name, &projectGroup.InsertTime, &projectGroup.UpdateTime); err != nil {
			return nil, err
		}
		projectGroups = append(projectGroups, projectGroup)
	}
	return projectGroups, nil
}

// GetData to Group
func (g Group) GetData() (Group, error) {
	var group Group
	err := sq.
		Select("name, insert_time, update_time").
		From(groupTable).
		Where(sq.Eq{"id": g.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&group.Name, &group.InsertTime, &group.UpdateTime)
	if err != nil {
		return group, errors.New("数据查询失败")
	}
	return group, nil
}
