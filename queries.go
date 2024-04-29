package main

import (
	"database/sql"
	"log"
	"time"
)

// insertions

func insertClient(c client, db *sql.DB) int {
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

func insertDoctor(doc employee, db *sql.DB) int {
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

func insertNurse(nurse employee, db *sql.DB) int {
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

// assigns an employee to a sector
func insertIntoSector(employee int, sectorID int, db *sql.DB) {
	_, err := db.Exec("INSERT INTO works_in (emp_ID, sector_ID) VALUES (?, ?)", employee, sectorID)
	if err != nil {
		log.Fatal(err)
	}
}

// inserts a new entry into the screening table
func insertScreening(db *sql.DB, data screening) {
	// nurse id actually points to the nurse's RN that was a mistake on my part, my bad >w<
	_, err := db.Exec("INSERT INTO screening (patient, nurse_id, date, diagnosis) VALUES (?, ?, ?, ?)", data.patient, data.nurse_id, data.date, data.diagnosis)
	if err != nil {
		log.Fatal(err)
	}
}

// inserts a new entry into the appointment table
func insertAppointment(db *sql.DB, data appointment) {
	_, err := db.Exec("INSERT INTO appointment (patient, doctor, dateApt) VALUES (?, ?, ?)", data.patient, data.doctor, data.dateApt)
	if err != nil {
		log.Fatal(err)
	}
}

// updates

// alters a screening tuple to add a doctor which the patient is fowarded to
func updateScreening(db *sql.DB, doctor int, data screening) {
	_, err := db.Exec("UPDATE screening SET fowards_to = ? WHERE nurse_id = ? AND patient = ? and date = ?", doctor, data.nurse_id, data.patient, data.date)
	if err != nil {
		log.Fatal(err)
	}
}

// inserts/alters the manager of a sector
func updateSectorManager(db *sql.DB, sectorID int, managerCRM int) {
	currentTime := time.Now()
	_, err := db.Exec("UPDATE sector SET manager = ?, manager_start_date = ? WHERE id = ?", managerCRM, currentTime.Format("2006-01-02"), sectorID)
	if err != nil {
		log.Fatal(err)
	}
}

func updateClientName(db *sql.DB, c client, id int) {
	_, err := db.Exec("UPDATE client SET clientName = ? WHERE id = ?", c.name, id)
	if err != nil {
		log.Fatal(err)
	}
}

func updateClientPhone(db *sql.DB, c client, id int) {
	_, err := db.Exec("UPDATE client SET phone = ? WHERE clientID = ?", c.phone, id)
	if err != nil {
		log.Fatal(err)
	}
}

func updateEmployeeName(db *sql.DB, e employee, id int) {
	_, err := db.Exec("UPDATE employee SET empName = ? WHERE id = ?", e.name, id)
	if err != nil {
		log.Fatal(err)
	}
}

func updateEmployeeSalary(db *sql.DB, e employee, id int) {
	_, err := db.Exec("UPDATE employee SET salary = ? WHERE id = ?", e.salary, id)
	if err != nil {
		log.Fatal(err)
	}
}

func updateEmployeePhone(db *sql.DB, e employee, id int) {
	_, err := db.Exec("UPDATE employee SET phone = ? WHERE id = ?", e.phone, id)
	if err != nil {
		log.Fatal(err)
	}
}

func updateDoctorSpecialty(db *sql.DB, specialty string, CRM int) {
	_, err := db.Exec("UPDATE doctor SET specialty = ? and CRM = ?", specialty, CRM)
	if err != nil {
		log.Fatal(err)
	}
}

// deletions

func dropEmployee(db *sql.DB, id int) {
	_, err := db.Exec("DELETE FROM employee WHERE id =?", id)
	if err != nil {
		log.Fatal(err)
	}
}

func dropClient(db *sql.DB, id int) {
	_, err := db.Exec("DELETE FROM client WHERE clientID =?", id)
	if err != nil {
		log.Fatal(err)
	}
}

func dropAppointment(db *sql.DB, apt appointment) {
	_, err := db.Exec("DELETE FROM appointment WHERE patient = ? and doctor = ? and dateApt = ?", apt.patient, apt.doctor, apt.dateApt)
	if err != nil {
		log.Fatal(err)
	}
}

func dropScreening(db *sql.DB, data screening) {
	_, err := db.Exec("DELETE FROM screening WHERE nurse_id = ? and patient = ? and date = ?", data.nurse_id, data.patient, data.date)
	if err != nil {
		log.Fatal(err)
	}
}

func dropSectorEmployee(db *sql.DB, id int, sectorID int) {
	_, err := db.Exec("DELETE FROM works_in WHERE emp_ID = ? AND sector_ID = ?", id, sectorID)
	if err != nil {
		log.Fatal(err)
	}
}

//queries

// query to check basic info regarding a given client
func queryClient(db *sql.DB, id int) (client, error) {
	var c client
	err := db.QueryRow("SELECT clientName, phone, birthdate FROM client WHERE clientID = ?", id).Scan(&c.name, &c.phone, &c.birthdate)
	if err != nil {
		return client{}, err
	}
	return c, nil
}

// query to check an appointment of a given patient with a given appointment date
func queryAppointment(db *sql.DB, id int, date string) (appointment, error) {
	var apt appointment
	err := db.QueryRow("SELECT patient, doctor, dateApt FROM appointment WHERE patient = ? and dateApt = ?", id, date).Scan(&apt.patient, &apt.doctor, &apt.dateApt)
	if err != nil {
		return appointment{}, err
	}
	return apt, nil
}

// query to check a screening of a given patient with a given date
func queryScreening(db *sql.DB, id int, date string) (screening, error) {
	var scr screening
	err := db.QueryRow("SELECT nurse_id, patient, date, diagnosis, fowards_to FROM screening WHERE patient = ? and date = ?", id, date).Scan(&scr.nurse_id, &scr.patient, &scr.date, &scr.diagnosis, &scr.fowards_to)
	if err != nil {
		return screening{}, err
	}
	return scr, nil
}

// queries all appointments of a given patient
func queryAllAppointments(db *sql.DB, id int) ([]appointment, error) {
	// query for rows
	rows, err := db.Query("SELECT c.clientName, e.empName, a.dateApt FROM appointment a, employee e, doctor d, client c WHERE a.patient = ? AND a.patient = c.clientID AND a.doctor = d.CRM AND d.empID = e.id", id)
	// check for errors
	if err != nil {
		return nil, err
	}
	// defer so the rows close when the function finishes
	defer rows.Close()

	// make a slice
	apts := make([]appointment, 0, 50)
	// insert the other rows into a temporary var
	for rows.Next() {
		var apt appointment
		err := rows.Scan(&apt.patient, &apt.doctor, &apt.dateApt)
		// check for errors
		if err != nil {
			return nil, err
		}
		// append to the main slice
		apts = append(apts, apt)
	}

	return apts, nil
}

// queries all screenings of a given patient
func queryAllScreenings(db *sql.DB, id int) ([]screening, error) {
	// query for rows
	rows, err := db.Query("SELECT c.clientName, e.empName, s.date, s.diagnosis, s.fowards_to FROM screening s, client c, employee e, nurse n WHERE patient = ? AND s.patient = c.clientID AND s.nurse_id = n.RN AND n.empID = e.id", id)
	// check for errors
	if err != nil {
		return nil, err
	}
	// defer so the rows close when the function finishes
	defer rows.Close()

	// make a slice
	scrs := make([]screening, 0, 50)
	// insert each row into a temporary var
	for rows.Next() {
		var scr screening
		err := rows.Scan(&scr.nurse_id, &scr.patient, &scr.date, &scr.diagnosis, &scr.fowards_to)
		// check for errors
		if err != nil {
			return nil, err
		}
		// append to the main slice
		scrs = append(scrs, scr)
	}

	return scrs, nil

}

// queries employee info given their ID
func queryEmployeeByID(db *sql.DB, id int) (struct {
	id         int
	name       string
	phone      string
	birthdate  string
	hiringdate string
	salary     string
	sector     string
}, error) {
	// declare variables
	var e employee
	var idGet int
	var sector string
	// query for row
	err := db.QueryRow("SELECT e.id, empName, e.phone, e.birthdate, e.hiringdate, e.salary, s.sector_name FROM employee e, sector s, works_in w WHERE e.id = ? AND e.id = w.emp_ID AND w.sector_ID = s.id", id).Scan(&idGet, &e.name, &e.phone, &e.birthdate, &e.hiringdate, &e.salary, &sector)
	if err != nil {
		// return an empty struct in case of error
		return struct {
			id         int
			name       string
			phone      string
			birthdate  string
			hiringdate string
			salary     string
			sector     string
		}{}, err
	}

	// else return queried values
	return struct {
		id         int
		name       string
		phone      string
		birthdate  string
		hiringdate string
		salary     string
		sector     string
	}{idGet, e.name, e.phone, e.birthdate, e.hiringdate, e.salary, sector}, err
}

// queries doctor info given their CRM
func queryDoctorByCRM(db *sql.DB, CRM int) (struct {
	id         int
	name       string
	phone      string
	birthdate  string
	hiringdate string
	salary     string
	sector     string
	specialty  string
}, error) {
	// declare variables
	var e employee
	var idGet int
	var sector string
	var specialty string
	// query for row
	err := db.QueryRow("SELECT e.id, empName, e.phone, e.birthdate, e.hiringdate, e.salary, s.sector_name, d.specialty FROM employee e, sector s, works_in w, doctor d WHERE d.CRM = ? AND e.id = d.empID AND e.id = w.emp_ID AND w.sector_ID = s.id", CRM).Scan(&idGet, &e.name, &e.phone, &e.birthdate, &e.hiringdate, &e.salary, &sector, &specialty)
	if err != nil {
		// return an empty struct in case of error
		return struct {
			id         int
			name       string
			phone      string
			birthdate  string
			hiringdate string
			salary     string
			sector     string
			specialty  string
		}{}, err
	}

	// else return queried values
	return struct {
		id         int
		name       string
		phone      string
		birthdate  string
		hiringdate string
		salary     string
		sector     string
		specialty  string
	}{idGet, e.name, e.phone, e.birthdate, e.hiringdate, e.salary, sector, specialty}, err
}

func queryNurseByRN(db *sql.DB, RN int) (struct {
	id         int
	name       string
	phone      string
	birthdate  string
	hiringdate string
	salary     string
	sector     string
}, error) {
	// declare variables
	var e employee
	var idGet int
	var sector string
	// query for row
	err := db.QueryRow("SELECT e.id, empName, e.phone, e.birthdate, e.hiringdate, e.salary, s.sector_name FROM employee e, sector s, works_in w, nurse n WHERE n.RN = ? AND e.id = n.empID AND e.id = w.emp_ID AND w.sector_ID = s.id", RN).Scan(&idGet, &e.name, &e.phone, &e.birthdate, &e.hiringdate, &e.salary, &sector)
	if err != nil {
		// return an empty struct in case of error
		return struct {
			id         int
			name       string
			phone      string
			birthdate  string
			hiringdate string
			salary     string
			sector     string
		}{}, err
	}

	// else return queried values
	return struct {
		id         int
		name       string
		phone      string
		birthdate  string
		hiringdate string
		salary     string
		sector     string
	}{idGet, e.name, e.phone, e.birthdate, e.hiringdate, e.salary, sector}, err
}

// returns the CRM of a doctor given his employee ID
func queryCRMbyID(db *sql.DB, id int) (int, error) {
	var CRM int
	err := db.QueryRow("SELECT d.CRM FROM doctor d WHERE d.empID = ?", id).Scan(&CRM)
	if err != nil {
		return 0, err
	}
	return CRM, err
}

func queryRNbyID(db *sql.DB, id int) (int, error) {
	var RN int
	err := db.QueryRow("SELECT n.RN FROM nurse n WHERE n.empID = ?", id).Scan(&RN)
	if err != nil {
		return 0, err
	}
	return RN, err
}

// queries a sector's info given its ID
func querySectorByID(db *sql.DB, sectID int) (struct {
	sectName           string
	manager            string
	managerID          int
	manager_start_date string
}, error) {
	// declare variables
	var sectorName string
	var manager string
	var managerID int
	var manager_start_date string
	// query for row
	err := db.QueryRow("SELECT s.sector_name, e.empName, e.id, s.manager_start_date FROM sector s, works_in w, employee e WHERE s.id = ? AND s.id = w.sector_ID AND w.emp_ID = e.id", sectID).Scan(&sectorName, &manager, &managerID, &manager_start_date)
	if err != nil {
		// return an empty struct in case of error
		return struct {
			sectName           string
			manager            string
			managerID          int
			manager_start_date string
		}{}, err
	}

	// else return queried values
	return struct {
		sectName           string
		manager            string
		managerID          int
		manager_start_date string
	}{sectorName, manager, managerID, manager_start_date}, err

}

// queries a sector's info given its ID
func querySectorByName(db *sql.DB, sectName string) (struct {
	manager_start_date string
	manager            string
	managerID          int
	sectID             int
}, error) {
	// declare variables
	var sectID int
	var manager string
	var managerID int
	var manager_start_date string
	// query for row
	err := db.QueryRow("SELECT s.id, e.empName, e.id, s.manager_start_date FROM sector s, works_in w, employee e WHERE s.sector_name = ? AND s.id = w.sector_ID AND w.emp_ID = e.id", sectName).Scan(&sectID, &manager, &managerID, &manager_start_date)
	if err != nil {
		// return an empty struct in case of error
		return struct {
			manager_start_date string
			manager            string
			managerID          int
			sectID             int
		}{}, err
	}

	// else return queried values
	return struct {
		manager_start_date string
		manager            string
		managerID          int
		sectID             int
	}{manager_start_date, manager, managerID, sectID}, err

}
