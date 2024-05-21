package main

import (
	"bufio"
	"database/sql"
	"fmt"

	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type employee struct {
	name       string
	salary     string
	phone      string
	birthdate  string
	hiringdate string
	RN         string // for nurses
	CRM        string // for doctors
	specialty  string // also for doctors
}

type client struct {
	name      string
	birthdate string
	phone     string
}

type screening struct {
	patient    string // the id of the patient, primary key
	nurse_id   string // the RN of the nurse who conducted the screening process, primary key
	diagnosis  string // the preliminary diagnosis
	date       string // the date of the screening, primary key
	fowards_to string // the CRM of the doctor that the patient may be fowarded to
}

type appointment struct {
	patient string // id of the patient
	doctor  string // CRM of the doctor
	dateApt string // date of the appointment
}

func main() {

	var db *sql.DB

	// set connection properties
	cfg := mysql.Config{
		// remember to set the env variables
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "clinicDB",
	}
	// get db handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	// ping db to make sure it's sure it's connected
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	// defer so it closes when the application ends
	defer db.Close()

	fmt.Println("Welcome to the management system.")
	var opt int
	for {
		fmt.Printf("Choose an option:\n1-Querry information\n2-Insert information\n3-Update information\n4-Delete information\n5-Exit\n")
		fmt.Scan(&opt)
		switch opt {
		case 1:
			queryMenu(db)
		case 2:
			insertionMenu(db)
		case 3:
			updateMenu(db)
		case 4:
			deletionMenu(db)
		case 5:
			return
		default:
			fmt.Println("Invalid option")
		}
	}

}

func queryMenu(db *sql.DB) {
	var (
		opt,
		id int
		data string
	)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Choose an option:\n1-Employees\n2-Clients\n3-Screenings\n4-Appointments\n5-Sectors\n6-Back\n")
		fmt.Scan(&opt)
		switch opt {
		case 1:
			fmt.Printf("Choose an option:\n1-Query employee by ID\n2-Query doctor by CRM\n3-Query nurse by RN\n4-Go back\n")
			fmt.Scan(&opt)
			switch opt {
			case 1:
				fmt.Println("Query employee by ID")
				fmt.Scan(&id)
				q, err := queryEmployeeByID(db, id)
				if err != nil {
					fmt.Println("Error while retrieving information.")
					fmt.Println(err)
					break
				}
				fmt.Println(q)
			case 2:
				fmt.Println("Query doctor by CRM")
				fmt.Scan(&id) // read id as CRM
				q, err := queryDoctorByCRM(db, id)
				if err != nil {
					fmt.Println("Error while retrieving information.")
					fmt.Println(err)
					break
				}
				fmt.Println(q)
			case 3:
				fmt.Println("Query nurse by RN")
				fmt.Scan(&id) // read id as RN
				q, err := queryNurseByRN(db, id)
				if err != nil {
					fmt.Println("Error while retrieving information.")
					fmt.Println(err)
					break
				}
				fmt.Println(q)
			case 4:
				break
			default:
				fmt.Println("Invalid option")
			}

		case 2:
			fmt.Println("Please inform the client's ID or enter 0 to go back:")
			fmt.Scan(&id)
			if id == 0 {
				break
			}
			q, err := queryClient(db, id)
			if err != nil {
				fmt.Println("Error while retrieving information.")
				fmt.Println(err)
				break
			}
			fmt.Println(q)

		case 3:
			fmt.Printf("Choose an option:\n1-Query a specific screening\n2-Query all screenings of a specific patient\n3-Back")
			fmt.Scan(&opt)
			switch opt {
			case 1:
				fmt.Println("Please inform the patient's ID:")
				fmt.Scan(&id)
				// call function
				q, err := queryAllScreenings(db, id)
				if err != nil {
					fmt.Println("Error while retrieving information.")
					fmt.Println(err)
					break
				}
				// loop over the resulting slice
				for _, v := range q {
					fmt.Println(v)
				}

			case 2:
				fmt.Println("Please inform the patient's ID:")
				fmt.Scan(&id)
				fmt.Println("Please inform the date of the screening:")
				data, _ = reader.ReadString('\n')
				// call function
				q, err := queryScreening(db, id, data)
				if err != nil {
					fmt.Println("Error while retrieving information.")
					fmt.Println(err)
					break
				}
				fmt.Println(q)
			case 3:
				break
			}

		case 4:
			fmt.Printf("Choose an option:\n1-Query a specific appointment\n2-Query all appointments of a specific patient\n3-Back")
			fmt.Scan(&opt)
			switch opt {
			case 1:
				fmt.Println("Please inform the patient's ID:")
				fmt.Scan(&id)
				// call function
				q, err := queryAllAppointments(db, id)
				if err != nil {
					fmt.Println("Error while retrieving information.")
					fmt.Println(err)
					break
				}
				// loop over the resulting slice
				for _, v := range q {
					fmt.Println(v)
				}

			case 2:
				fmt.Println("Please inform the patient's ID:")
				fmt.Scan(&id)
				fmt.Println("Please inform the date of the appointment:")
				data, _ = reader.ReadString('\n')
				// call function
				q, err := queryAppointment(db, id, data)
				if err != nil {
					fmt.Println("Error while retrieving information.")
					fmt.Println(err)
					break
				}
				fmt.Println(q)
			case 3:
				break
			}
		case 5:
			fmt.Printf("Choose an option:\n1-Query sector by name\n2-Query sector by ID\n3-Back\n")
			fmt.Scan(&opt)
			switch opt {
			case 1:
				fmt.Println("Please inform the sector's name:")
				data, _ = reader.ReadString('\n')
				q, err := querySectorByName(db, data)
				if err != nil {
					fmt.Println("Error while retrieving information.")
					fmt.Println(err)
					break
				}
				fmt.Println(q)
			case 2:
				fmt.Println("Please inform the sector's ID:")
				fmt.Scan(&id)
				q, err := querySectorByID(db, id)
				if err != nil {
					fmt.Println("Error while retrieving information.")
					fmt.Println(err)
					break
				}
				fmt.Println(q)
			case 3:
				break
			}
		case 6:
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func insertionMenu(db *sql.DB) {
	var (
		opt int
		emp employee
		cli client
		scr screening
		apt appointment
	)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Choose an option:\n1-Insert employee\n2-Insert client\n3-Insert screening\n4-Insert appointment\n5-Insert sector\n6-Back\n")
		fmt.Scan(&opt)
		switch opt {
		case 1:
			fmt.Println("[INSERTING AN EMPLOYEE]")
			fmt.Printf("Choose an option:\n1-Insert a doctor\n2-Insert a nurse\n3-Back\n")
			fmt.Scan(&opt)
			switch opt {
			case 1:
				fmt.Println("[INSERTING A NEW DOCTOR]")
				fmt.Println("Enter the employee's name:")
				emp.name, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's birthdate (MM/DD/YYYY):")
				emp.birthdate, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's phone number:")
				emp.phone, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's hiring date (MM/DD/YYYY):")
				emp.hiringdate, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's salary:")
				emp.salary, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's CRM:")
				emp.CRM, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's specialty:")
				emp.specialty, _ = reader.ReadString('\n')
				id, err := insertDoctor(emp, db)
				if err != nil {
					fmt.Println("Error while inserting information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The employee has been inserted with ID:", id)

			case 2:
				fmt.Println("[INSERTING A NEW NURSE]")
				fmt.Println("Enter the employee's name:")
				emp.name, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's birthdate (MM/DD/YYYY):")
				emp.birthdate, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's phone number:")
				emp.phone, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's hiring date (MM/DD/YYYY):")
				emp.hiringdate, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's salary:")
				emp.salary, _ = reader.ReadString('\n')
				fmt.Println("Enter the employee's RN:")
				emp.CRM, _ = reader.ReadString('\n')
				id, err := insertNurse(emp, db)
				if err != nil {
					fmt.Println("Error while inserting information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The employee has been inserted with ID:", id)
			case 3:
				break
			default:
				fmt.Println("Invalid option")
			}
		case 2:
			fmt.Println("[INSERTING A NEW CLIENT]")
			fmt.Println("Enter (1) to continue or (0) to go back:")
			fmt.Scan(&opt)
			if opt == 0 {
				break
			}
			fmt.Println("Enter the client's name:")
			cli.name, _ = reader.ReadString('\n')
			fmt.Println("Enter the client's birthdate (MM/DD/YYYY):")
			cli.birthdate, _ = reader.ReadString('\n')
			fmt.Println("Enter the client's phone number:")
			cli.phone, _ = reader.ReadString('\n')
			id, err := insertClient(cli, db)
			if err != nil {
				fmt.Println("Error while inserting information.")
				fmt.Println(err)
				break
			}
			fmt.Println("The client has been inserted with ID:", id)

		case 3:
			fmt.Println("[INSERTING A NEW SCREENING]")
			fmt.Println("Enter (1) to continue or (0) to go back:")
			fmt.Scan(&opt)
			if opt == 0 {
				break
			}
			fmt.Println("Enter the patient's ID:")
			scr.patient, _ = reader.ReadString('\n')
			fmt.Println("Enter the nurse's ID:")
			scr.nurse_id, _ = reader.ReadString('\n')
			fmt.Println("Enter the screening's date (MM/DD/YYYY):")
			scr.date, _ = reader.ReadString('\n')
			fmt.Println("Is there a preliminary diagnosis? (1)-yes/(0)-no")
			fmt.Scan(&opt)
			if opt == 1 {
				fmt.Println("Enter the preliminary diagnosis:")
				scr.diagnosis, _ = reader.ReadString('\n')
			}
			fmt.Println("Has the patient been fowarded to a doctor? (1)-yes/(0)-no")
			fmt.Scan(&opt)
			if opt == 1 {
				fmt.Println("Enter the doctor's CRM:")
				fmt.Scan(&opt) // this is a bit lazy, might change later
			}
			err := insertScreening(db, scr)
			updateScreening(db, opt, scr)
			if err != nil {
				fmt.Println("Error while inserting information.")
				fmt.Println(err)
				break
			}
			fmt.Println("Successfully inserted screening.")

		case 4:
			fmt.Println("[INSERTING A NEW APPOINTMENT]")
			fmt.Println("Enter (1) to continue or (0) to go back:")
			fmt.Scan(&opt)
			if opt == 0 {
				break
			}
			fmt.Println("Enter the patient's ID:")
			apt.patient, _ = reader.ReadString('\n')
			fmt.Println("Enter the doctor's ID:")
			apt.doctor, _ = reader.ReadString('\n')
			fmt.Println("Enter the appointment's date (MM/DD/YYYY):")
			apt.dateApt, _ = reader.ReadString('\n')
			err := insertAppointment(db, apt)
			if err != nil {
				fmt.Println("Error while inserting information.")
				fmt.Println(err)
				break
			}
			fmt.Println("Successfully inserted appointment.")

		case 5:
			fmt.Println("[INSERTING A NEW SECTOR]")
			fmt.Printf("Choose an option:\n1- Create a new sector\n2- Insert an employee into an existing sector\n3-Back\n")
			fmt.Scan(&opt)
			switch opt {
			case 1:
				var sect string
				fmt.Println("Enter the new sector's name:")
				sect, _ = reader.ReadString('\n')
				id, err := createSector(db, sect)
				if err != nil {
					fmt.Println("Error while inserting information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The sector has been inserted with ID:", id)

			case 2:
				fmt.Println("[INSERTING EMPLOYEE INTO SECTOR]")
				fmt.Println("Enter (1) to continue or (0) to go back:")
				fmt.Scan(&opt)
				if opt == 0 {
					break
				}
				var sect int
				var employee int
				fmt.Println("Enter the sector's ID:")
				fmt.Scan(&sect)
				fmt.Println("Enter the employee's ID:")
				fmt.Scan(&employee)
				err := insertIntoSector(employee, sect, db)
				if err != nil {
					fmt.Println("Error while inserting information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The employee has been inserted into the sector.")
			case 3:
				break
			default:
				fmt.Println("Invalid option")
			}
		case 6:
			return
		default:
			fmt.Println("Invalid option")
		}

	}
}

func updateMenu(db *sql.DB) {
	var (
		opc int
		emp employee
		cli client
		scr screening
	)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Choose an option:\n1- Update employee\n2- Update client\n3- Update screening\n4- Update sector\n5- Back\n")
		fmt.Scan(&opc)
		switch opc {
		case 1:
			fmt.Println("[UPDATING EMPLOYEE INFORMATION]")
			fmt.Printf("Choose an option:\n1- Update employee name\n2- Update employee phone\n3- Update employee salary\n4- Update a doctor's specialty\n5- Back\n")
			fmt.Scan(&opc)
			switch opc {
			case 1:
				fmt.Println("[UPDATING EMPLOYEE NAME]")
				fmt.Println("Enter the employee's ID:")
				var id int
				fmt.Scan(&id)
				fmt.Println("Enter the new name:")
				emp.name, _ = reader.ReadString('\n')
				err := updateEmployeeName(db, emp, id)
				if err != nil {
					fmt.Println("Error while updating information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The employee's name has been updated.")
			case 2:
				fmt.Println("[UPDATING EMPLOYEE PHONE]")
				fmt.Println("Enter the employee's ID:")
				var id int
				fmt.Scan(&id)
				fmt.Println("Enter the new phone:")
				emp.phone, _ = reader.ReadString('\n')
				err := updateEmployeePhone(db, emp, id)
				if err != nil {
					fmt.Println("Error while updating information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The employee's phone has been updated.")
			case 3:
				fmt.Println("[UPDATING EMPLOYEE SALARY]")
				fmt.Println("Enter the employee's ID:")
				var id int
				fmt.Scan(&id)
				fmt.Println("Enter the new salary:")
				emp.salary, _ = reader.ReadString('\n')
				err := updateEmployeeSalary(db, emp, id)
				if err != nil {
					fmt.Println("Error while updating information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The employee's salary has been updated.")
			case 4:
				fmt.Println("[UPDATING DOCTOR'S SPECIALTY]")
				fmt.Println("Enter the employee's ID:")
				var id int
				fmt.Scan(&id)
				fmt.Println("Enter their new specialty:")
				emp.specialty, _ = reader.ReadString('\n')
				err := updateDoctorSpecialty(db, emp.specialty, id)
				if err != nil {
					fmt.Println("Error while updating information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The employee's specialty has been updated.")
			case 5:
				break
			default:
				fmt.Println("Invalid option")
			}

		case 2:
			fmt.Println("[UPDATING CLIENT INFORMATION]")
			fmt.Printf("Choose an option:\n1- Update client name\n2- Update client phone\n3- Back\n")
			fmt.Scan(&opc)
			switch opc {
			case 1:
				fmt.Println("[UPDATING CLIENT NAME]")
				fmt.Println("Enter the client's ID:")
				var id int
				fmt.Scan(&id)
				fmt.Println("Enter the new name:")
				cli.name, _ = reader.ReadString('\n')
				err := updateClientName(db, cli, id)
				if err != nil {
					fmt.Println("Error while updating information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The client's name has been updated.")
			case 2:
				fmt.Println("[UPDATING CLIENT PHONE]")
				fmt.Println("Enter the client's ID:")
				var id int
				fmt.Scan(&id)
				fmt.Println("Enter the new phone:")
				cli.phone, _ = reader.ReadString('\n')
				err := updateClientPhone(db, cli, id)
				if err != nil {
					fmt.Println("Error while updating information.")
					fmt.Println(err)
					break
				}
				fmt.Println("The client's phone has been updated.")
			case 3:
				break
			default:
				fmt.Println("Invalid option")
			}

		case 3:
			fmt.Println("[UPDATING SCREENING INFORMATION]")
			fmt.Println("Enter (1) to alter a screening tuple to add a doctor which the patient will be fowarded to or (0) to go back:")
			fmt.Scan(&opc)
			if opc == 0 {
				break
			}
			fmt.Println("Enter the doctor's CRM:")
			var CRM int
			fmt.Scan(&CRM)
			fmt.Println("Enter the screening's date:")
			scr.date, _ = reader.ReadString('\n')
			fmt.Println("Enter the patient's ID:")
			scr.patient, _ = reader.ReadString('\n')
			fmt.Println("Enter the nurse's RN:")
			scr.nurse_id, _ = reader.ReadString('\n')
			err := updateScreening(db, CRM, scr)
			if err != nil {
				fmt.Println("Error while updating information.")
				fmt.Println(err)
				break
			}
			fmt.Println("The screening has been updated.")
		case 4:
			fmt.Println("[UPDATING SECTOR MANAGER]")
			fmt.Printf("Enter (1) to Update sector manager or (0) to go back")
			fmt.Scan(&opc)
			if opc == 0 {
				break
			}
			fmt.Println("Enter the sector's ID:")
			var id int
			fmt.Scan(&id)
			fmt.Println("Enter the new manager's ID:")
			var manager int
			fmt.Scan(&manager)
			err := updateSectorManager(db, id, manager)
			if err != nil {
				fmt.Println("Error while updating information.")
				fmt.Println(err)
				break
			}
			fmt.Println("The sector's manager has been updated.")

		case 5:
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}
func deletionMenu(db *sql.DB) {
	var (
		opc int
		scr screening
	)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("[DELETING INFORMATION]")
	for {
		fmt.Printf("Choose an option:\n1- Delete employee\n2- Delete client\n3- Delete screening\n4- Delete sector employee\n5- Back\n")
		fmt.Scan(&opc)
		switch opc {
		case 1:
			fmt.Println("[DELETING EMPLOYEE]")
			fmt.Println("Enter the employee's ID:")
			var id int
			fmt.Scan(&id)
			err := dropEmployee(db, id)
			if err != nil {
				fmt.Println("Error while deleting information.")
				fmt.Println(err)
				break
			}
			fmt.Println("The employee has been deleted.")
		case 2:
			fmt.Println("[DELETING CLIENT]")
			fmt.Println("Enter the client's ID:")
			var id int
			fmt.Scan(&id)
			err := dropClient(db, id)
			if err != nil {
				fmt.Println("Error while deleting information.")
				fmt.Println(err)
				break
			}
			fmt.Println("The client has been deleted.")
		case 3:
			fmt.Println("[DELETING SCREENING]")
			fmt.Println("Enter the screening's date:")
			scr.date, _ = reader.ReadString('\n')
			fmt.Println("Enter the patient's ID:")
			scr.patient, _ = reader.ReadString('\n')
			fmt.Println("Enter the nurse's RN:")
			scr.nurse_id, _ = reader.ReadString('\n')
			err := dropScreening(db, scr)
			if err != nil {
				fmt.Println("Error while deleting information.")
				fmt.Println(err)
				break
			}
			fmt.Println("The screening has been deleted.")
		case 4:
			fmt.Println("[DELETING SECTOR EMPLOYEE]")
			fmt.Println("Enter the sector's ID:")
			var id int
			fmt.Scan(&id)
			fmt.Println("Enter the employee's ID:")
			var employee int
			fmt.Scan(&employee)
			err := dropSectorEmployee(db, id, employee)
			if err != nil {
				fmt.Println("Error while deleting information.")
				fmt.Println(err)
				break
			}
			fmt.Println("The sector's employee has been deleted.")
		case 5:
			return
		default:
			fmt.Println("Invalid option")
		}
	}

}
