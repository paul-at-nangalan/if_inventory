package armourments

import (
	"database/sql"
	"fmt"
	"if_inventory/errorhandler"
	"if_inventory/models"
)

type Model struct{
	getbyship *sql.Stmt
	insbyship *sql.Stmt
	findbyTitleShip *sql.Stmt
	del *sql.Stmt
	upd *sql.Stmt
}

func NewModel(db *sql.DB)*Model {

	cols := map[string]string{
		"id":           "BIGINT AUTO_INCREMENT",
		"title":        "VARCHAR (600)", ///may want to be able to index this, hence use varchar
		"qty":          "BIGINT",
		"spaceship_id": "BIGINT",
		"deleted":      "BOOL DEFAULT FALSE",
		"created":      "TIMESTAMP DEFAULT NOW()",
	}
	indx := []string{"id", "title"}
	primes := []string{"id"}
	migrator := models.NewMigrator(db)
	migrator.CreateTable("armourments", "create-table-armourments",
		cols, indx, primes)
	get := `SELECT id, title, qty FROM armourments WHERE spacecraft_id=? AND deleted=false`
	getstmt, err := db.Prepare(get)
	errorhandler.PanicOnErr(err)

	ins := `INSERT INTO armourments (title, qty, spacecraft_id)
					VALUES(?, ?, ?)`
	insstmt, err := db.Prepare(ins)

	upd := `UPDATE armourments SET title=?, qty=? WHERE id=?`
	updstmt, err := db.Prepare(upd)
	errorhandler.PanicOnErr(err)

	del := `UPDATE armourments SET deleted=true WHERE spacecraft_id=?`
	delstmt, err := db.Prepare(del)
	errorhandler.PanicOnErr(err)

	find := `SELECT id FROM armourments WHERE title=? AND spacecraft_id=? AND deleted=false`
	findstmt, err := db.Prepare(find)
	errorhandler.PanicOnErr(err)

	return &Model{
		getbyship: getstmt,
		insbyship: insstmt,
		upd: updstmt,
		findbyTitleShip: findstmt,
		del: delstmt,
	}
}

func (p *Model)GetByShip(shipid int64)[]models.Armourments{
	res, err := p.getbyship.Query(shipid)
	errorhandler.PanicOnErr(err)
	defer res.Close()

	armourments := make([]models.Armourments, 0)
	for res.Next(){
		id := int64(0)
		title := sql.NullString{}
		qty := int64(0)

		err := res.Scan(&id, title, qty)
		errorhandler.PanicOnErr(err)

		armourment := models.Armourments{}
		armourment.Id = id
		if title.Valid {
			armourment.Title = title.String
		}
		armourment.Qty = qty
		armourments = append(armourments, armourment)
	}
	return armourments
}

func (p *Model)findByTitleAndShip(title string, spacecraftid int64)bool{
	res, err := p.findbyTitleShip.Query(title, spacecraftid)
	errorhandler.PanicOnErr(err)
	defer res.Close()

	if res.Next(){
		return true
	}
	return false
}

func (p *Model)Insert(armourment models.Armourments, shipid int64)error{
	if p.findByTitleAndShip(armourment.Title, shipid){
		return models.NewErrorAlreadyExists(armourment.Title)
	}

	_, err := p.insbyship.Exec(armourment.Title, armourment.Qty, shipid)
	errorhandler.PanicOnErr(err)
	return nil
}

func (p *Model)Update(armament models.Armourments)error  {
	res, err := p.upd.Exec(armament.Title, armament.Qty, armament.Id)
	errorhandler.PanicOnErr(err)

	rows, err := res.RowsAffected()
	errorhandler.PanicOnErr(err)
	if rows == 0{
		return models.NewErrorNotExist(fmt.Sprint(armament.Id))
	}
	return nil
}

func (p *Model)Delete(shipid int64){
	_, err := p.del.Exec(shipid)
	errorhandler.PanicOnErr(err)
}