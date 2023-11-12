package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
)

var vue int = 0

type person struct {
	Nom      string
	Prenom   string
	Birthday string
	Sexe     string
}

var dataForm person = person{}

func main() {
	temp, err := template.ParseGlob("./temp/*.html")
	if err != nil {
		fmt.Println(fmt.Sprintf("Erreur %s", err.Error()))
	}
	type PageVariables struct {
		Nom     string
		Filiere string
		Niv     int
		Nbr     int
	}

	type Eleve struct {
		Prenom string
		Nom    string
		Age    int
		Sexe   bool
	}

	type data struct {
		PV  PageVariables
		Elv []Eleve
	}

	type even struct {
		Value int
		Check bool
	}

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		vue++
		var page even
		if vue%2 == 0 {
			page = even{vue, true}
		} else {
			page = even{vue, false}
		}
		temp.ExecuteTemplate(w, "change", page)
	})

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		lstelv := []Eleve{{"Cyril", "RODRIGUES", 22, true}, {"Kheir-eddine", "MEDERREG", 22, false}, {"Alan", "PHILIPIERT", 26, true}}
		page := PageVariables{"Mentor'ac", "Informatique", 5, len(lstelv)}
		d := data{page, lstelv}
		temp.ExecuteTemplate(w, "promo", d)
	})

	http.HandleFunc("/user/init", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "init", nil)
	})

	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {
		dataForm = person{
			r.FormValue("Name"),
			r.FormValue("Firstname"),
			r.FormValue("Date"),
			r.FormValue("Sexe")}
		checkValue, _ := regexp.MatchString("^[a-zA-Z-]{1,64}$", dataForm.Nom)
		if !checkValue {
			dataForm.Nom = "Invalide"
		}
		checkValue, _ = regexp.MatchString("^[a-zA-Z-]{1,64}$", dataForm.Prenom)
		if !checkValue {
			dataForm.Prenom = "Invalide"
		}
		fmt.Println(dataForm)
		http.Redirect(w, r, "/user/display", http.StatusMovedPermanently)
	})

	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "display", dataForm)

	})
	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/asset"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.ListenAndServe("localhost:8080", nil)
}
