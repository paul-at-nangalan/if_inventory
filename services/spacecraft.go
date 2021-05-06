package services

import (
	"database/sql"
	"encoding/json"
	"if_inventory/errorhandler"
	"if_inventory/models/spacecraft"
	"net/http"
	"strconv"
)

type Spacecrafts struct{
	spacecraftmod *spacecraft.Model
}

func NewSpacesrafts(db *sql.DB)*Spacecrafts{

	spacecraftmod := spacecraft.NewModel(db)
	spacecrafts := &Spacecrafts{
		spacecraftmod: spacecraftmod,
	}

	http.HandleFunc("/find", spacecrafts.Find)
	http.HandleFunc("/get", spacecrafts.Find)

	return spacecrafts
}

func errorJson(reason string)string{
	resp := map[string]interface{}{
		"success": false,
		"reason": reason,
	}
	b := make([]byte, 0)
	err := json.Unmarshal(b, resp)
	errorhandler.PanicOnErr(err)

	return string(b)
}

func (p *Spacecrafts)Find(w http.ResponseWriter, req *http.Request){
	search := req.URL.Query().Get("search")
	if search == ""{
		http.Error(w, errorJson("No search term"), http.StatusBadRequest)
		return
	}
	ships := p.spacecraftmod.Find(search)
	
	data := make(map[string]interface{})
	data["data"] = ships
	jsondata, err := json.Marshal(data)
	errorhandler.PanicOnErr(err)

	w.Write(jsondata)
}

func (p *Spacecrafts)GetById(w http.ResponseWriter, req *http.Request){
	shipidstr := req.URL.Query().Get("shipid")
	if shipidstr == ""{
		http.Error(w, errorJson("No ship ID"), http.StatusBadRequest)
		return
	}
	shipid, err := strconv.ParseInt(shipidstr, 10, 64)
	if err != nil{
		http.Error(w, errorJson("Invalid ID"), http.StatusBadRequest)
		return
	}
	ships, found := p.spacecraftmod.GetById(shipid)
	if !found{
		http.Error(w, errorJson("ID not found in the database"), http.StatusGone)
		return
	}
	jsondata, err := json.Marshal(ships)
	errorhandler.PanicOnErr(err)

	w.Write(jsondata)
}