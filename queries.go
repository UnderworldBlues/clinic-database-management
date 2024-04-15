package main

import (
	"database/sql"
	"log"
	"time"
)

func insertClient(c *client, db *sql.DB) int {
	res, err := db.Exec("INSERT INTO client (clientName, phone, birthdate) VALUES (?, ?, ?)", c.name, c.phone, c.birthdate)
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	// check if the auto increment worked
	if err != nil {
		log.Fatal(err)
	}
	// return the id
	return int(id)
}

func insertDoctor(doc *employee, db *sql.DB) int {
	currentTime := time.Now() // gets the date the doctor was hired
	res, err := db.Exec("INSERT INTO employee (empName, phone, birthdate, hiringdate, salary) VALUES (?, ?, ?, ?, ?)", doc.name, doc.phone, doc.birthdate, currentTime.Format("2006-01-02"), doc.salary)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO doctor (empID, CRM, specialty) VALUES (?, ?, ?)", id, doc.CRM, doc.specialty)

	if err != nil {
		log.Fatal(err)
	}

	// return the id
	return int(id)

}

func insertNurse(nurse *employee, db *sql.DB) int {
	currentTime := time.Now() // gets the date the nurse was hired
	res, err := db.Exec("INSERT INTO employee (empName, phone, birthdate, hiringdate, salary) VALUES (?, ?, ?, ?, ?)", nurse.name, nurse.phone, nurse.birthdate, currentTime.Format("2006-01-02"), nurse.salary)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO nurse (empID, RN) VALUES (?, ?)", id, nurse.RN)
	if err != nil {
		log.Fatal(err)
	}

	return int(id)
}

// creates a new tuple in the sector table
func createSector(db *sql.DB, sectorName string) int {
	// manager is the doctor's CRM
	res, err := db.Exec("INSERT INTO sector ( sector_name ) VALUES ( ? )", sectorName)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return int(id)

}

// inserts/alters the manager of a sector
func alterSector(db *sql.DB, sectorID int, managerCRM int) {
	currentTime := time.Now()
	_, err := db.Exec("UPDATE sector SET manager = ?, manager_start_date = ? WHERE id = ?", managerCRM, currentTime.Format("2006-01-02"), sectorID)
	if err != nil {
		log.Fatal(err)
	}
}

// assigns an employee to a sector
func assignToSector(employee int, sectorID int, db *sql.DB) {
	_, err := db.Exec("INSERT INTO works_in (emp_ID, sector_ID) VALUES (?, ?)", employee, sectorID)
	if err != nil {
		log.Fatal(err)
	}
}

// inserts a new entry into the screening table
func fileScreening(db *sql.DB, data screening) {
	// nurse id actually points to the nurse's RN that was a mistake on my part, my bad >w<
	_, err := db.Exec("INSERT INTO screening (patient, nurse_id, date, diagnosis) VALUES (?, ?, ?, ?)", data.patient, data.nurse_id, data.date, data.diagnosis)
	if err != nil {
		log.Fatal(err)
	}
}

// alters a screening tuple to add a doctor which the patient is fowarded to
func screeningAlter(db *sql.DB, doctor int, data screening) {
	_, err := db.Exec("UPDATE screening SET fowards_to = ? WHERE nurse_id = ? AND patient = ? and date = ?", doctor, data.nurse_id, data.patient, data.date)
	if err != nil {
		log.Fatal(err)
	}
}

// inserts a new entry into the appointment table
func fileAppointment(db *sql.DB, data appointment) {
	_, err := db.Exec("INSERT INTO appointment (patient, doctor, dateApt) VALUES (?, ?, ?)", data.patient, data.doctor, data.dateApt)
	if err != nil {
		log.Fatal(err)
	}
}
