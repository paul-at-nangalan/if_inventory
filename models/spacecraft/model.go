package spacecraft

import (
	"database/sql"
	"fmt"
	"if_inventory/errorhandler"
	"if_inventory/models"
	"if_inventory/models/armourments"
)

type Model struct{
	ins *sql.Stmt
	upd *sql.Stmt
	del *sql.Stmt

	find *sql.Stmt
	getbyid *sql.Stmt
	exists *sql.Stmt

	arms *armourments.Model
}

func NewModel(db *sql.DB)*Model{

	cols := map[string]string{
		"id": "BIGINT AUTO_INCREMENT",
		"name": "VARCHAR (600)", ///want to be able to index this, hence use varchar
		"class": "VARCHAR (600)",
		"status": "VARCHAR (20)",
		"crew": "INT",
		"deleted": "BOOL DEFAULT FALSE",
		"created": "TIMESTAMP DEFAULT NOW()",
	}
	indx := []string{"id", "name", "class", "status"}
	primes := []string{"id"}
	migrator := models.NewMigrator(db)
	migrator.CreateTable("spacecraft", "create-table-space-craft",
		cols, indx, primes)
	
	get := `SELECT id, name, class, status, crew, title, qty FROM spacescraft  s
			LEFT JOIN armourments a ON a.spacecraft_id=s.id
			WHERE s.id=? AND deleted = false`
	getstmt, err := db.Prepare(get)
	errorhandler.PanicOnErr(err)

	find := `SELECT id, name, status FROM spacecraft WHERE 
			name (like ? OR status LIKE ? OR class LIKE ?) AND deleted = false`
	findstmt, err := db.Prepare(find)

	ins := `INSERT INTO spacecraft (name, class, status, crew)VALUES
							(?, ?, ?, ?)`
	insstmt, err := db.Prepare(ins)

	updt := `UPDATE space SET name=?, class=?, status=?, crew=? WHERE id=? AND deleted=false`
	updtstmt, err := db.Prepare(updt)

	exists := `SELECT id FROM spacecraft WHERE name=? AND deleted=false`
	existsstmt, err := db.Prepare(exists)
	errorhandler.PanicOnErr(err)

	arms := armourments.NewModel(db)
	return &Model{
		upd: updtstmt,
		find: findstmt,
		getbyid: getstmt,
		ins: insstmt,
		arms: arms,
		exists: existsstmt,
	}
}

func (p *Model)Find(search string)[]models.Spacecraft{
	res, err := p.find.Query(search + "%", search + "%", search + "%")
	errorhandler.PanicOnErr(err)
	defer res.Close()

	ships := make([]models.Spacecraft, 0)
	for res.Next(){
		id := int64(0)
		name := sql.NullString{}
		status := sql.NullString{}
		err := res.Scan(&id, &name, &status)
		errorhandler.PanicOnErr(err)

		ship := models.Spacecraft{}
		if name.Valid{
			ship.Name = name.String
		}
		if status.Valid{
			ship.Status = status.String
		}
		ship.Id = id

	}
	return ships
}

func ifValid(s sql.NullString)string{
	if s.Valid{
		return s.String
	}
	return ""
}

func (p *Model)GetById(shipid int64)( ship models.Spacecraft, found bool){
	res, err := p.getbyid.Query(shipid)
	errorhandler.PanicOnErr(err)
	defer res.Close()

	if !res.Next(){
		return models.Spacecraft{}, false
	}

	var id int64
	var name sql.NullString
	var class sql.NullString
	var status sql.NullString
	var crew int32
	err = res.Scan(&id, &name, &class, &status, &crew)
	errorhandler.PanicOnErr(err)
	
	arms := p.arms.GetByShip(id)
	
	ship = models.Spacecraft{}
	ship.Id = id
	ship.Name = ifValid(name)
	ship.Class = ifValid(class)
	ship.Status = ifValid(status)
	ship.Crew = crew

	ship.Armourments = arms
	return ship, true
}

func (p *Model)Update(ship models.Spacecraft)[]error{
	res, err := p.upd.Exec(ship.Name, ship.Class, ship.Status, ship.Crew, ship.Id)
	errorhandler.PanicOnErr(err)

	rows, err := res.RowsAffected()
	errorhandler.PanicOnErr(err)
	if rows == 0{
		return []error{models.NewErrorNotExist(fmt.Sprint(ship.Id))}
	}
	errlist := make([]error, 0)
	for _, arm := range ship.Armourments {
		///each armourment should have a unique ID
		/// and already be associated to the ship
		/// therefore we don't need the ship ID
		err = p.arms.Update(arm)
		errlist = append(errlist, err)
	}
	return errlist
}

func (p *Model)shipExists(name string)bool{
	////TODO
	return false ///for now
}

func (p *Model)Insert(ship models.Spacecraft)[]error{
	if p.shipExists(ship.Name){
		return []error{models.NewErrorAlreadyExists(ship.Name)}
	}
	res, err := p.ins.Exec(ship.Name, ship.Class, ship.Status, ship.Crew)
	errorhandler.PanicOnErr(err)

	///Get the ship ID so we can insert the armourments
	shipid, err := res.LastInsertId()
	errorhandler.PanicOnErr(err)
	errlist := make([]error, 0)
	for _, arm := range ship.Armourments {
		err = p.arms.Insert(arm, shipid)
		errlist = append(errlist, err)
	}
	return errlist
}