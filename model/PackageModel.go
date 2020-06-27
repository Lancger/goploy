package model

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
)

const packageTable = "`package`"

// Package mysql table package
type Package struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	InsertTime string `json:"insertTime"`
	UpdateTime string `json:"updateTime"`
}

// Packages many package
type Packages []Package

// GetList package row
func (p Package) GetList(pagination Pagination) (Packages, error) {
	rows, err := sq.
		Select("id, name, size, insert_time, update_time").
		From(packageTable).
		Limit(pagination.Rows).
		Offset((pagination.Page - 1) * pagination.Rows).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	packages := Packages{}
	for rows.Next() {
		var pkg Package

		if err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Size, &pkg.InsertTime, &pkg.UpdateTime); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

// GetList package total
func (p Package) GetTotal() (int64, error) {
	var total int64
	err := sq.
		Select("COUNT(*) AS count").
		From(packageTable).
		RunWith(DB).
		QueryRow().
		Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetAllInIDStr package row
func (p Package) GetAllInIDStr(IDStr string) (Packages, error) {
	rows, err := sq.
		Select("id, name, size").
		From(packageTable).
		Where(sq.Eq{"id": strings.Split(IDStr, ",")}).
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	packages := Packages{}
	for rows.Next() {
		var pkg Package

		if err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Size); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

// GetAll package row
func (p Package) GetAll() (Packages, error) {
	rows, err := sq.
		Select("id, name, size, insert_time, update_time").
		From(packageTable).
		OrderBy("id DESC").
		RunWith(DB).
		Query()
	if err != nil {
		return nil, err
	}
	packages := Packages{}
	for rows.Next() {
		var pkg Package

		if err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Size, &pkg.InsertTime, &pkg.UpdateTime); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

// GetData get package information
func (p Package) GetData() (Package, error) {
	var pkg Package
	err := sq.
		Select("id, name, size, insert_time, update_time").
		From(packageTable).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		QueryRow().
		Scan(&pkg.ID, &pkg.Name, &pkg.Size, &pkg.InsertTime, &pkg.UpdateTime)
	if err != nil {
		return pkg, err
	}
	return pkg, nil
}

// GetDataByName get package information
func (p Package) GetDataByName() (Package, error) {
	var pkg Package
	err := sq.
		Select("id, name, size, insert_time, update_time").
		From(packageTable).
		Where(sq.Eq{"name": p.Name}).
		RunWith(DB).
		QueryRow().
		Scan(&pkg.ID, &pkg.Name, &pkg.Size, &pkg.InsertTime, &pkg.UpdateTime)
	if err != nil {
		return pkg, err
	}
	return pkg, nil
}

// AddRow add one row to table package and add id to p.ID
func (p Package) AddRow() (int64, error) {
	result, err := sq.
		Insert(packageTable).
		Columns("name", "size").
		Values(p.Name, p.Size).
		RunWith(DB).
		Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// EditRow edit one row to table package
func (p Package) EditRow() error {
	_, err := sq.
		Update(packageTable).
		Set("name", p.Name).
		Set("size", p.Size).
		Where(sq.Eq{"id": p.ID}).
		RunWith(DB).
		Exec()
	return err
}
