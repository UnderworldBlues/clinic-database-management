CREATE DATABASE clinicDB;
USE clinicDB;

CREATE TABLE employee (
    id INT AUTO_INCREMENT PRIMARY KEY,
    empName VARCHAR(20) NOT NULL,
    salary FLOAT NOT NULL,
    birthdate DATE NOT NULL,
    hiringdate DATE NOT NULL,
    email VARCHAR(320) NOT NULL
);

CREATE TABLE client(
    clientID INT AUTO_INCREMENT PRIMARY KEY,
    clientName VARCHAR(20) NOT NULL,
    birthdate DATE NOT NULL,
    phone VARCHAR(15) NOT NULL
);

CREATE TABLE doctor (
    CRM INT PRIMARY KEY,
    specialty VARCHAR(20) NOT NULL,
    empID INT NOT NULL,
    CONSTRAINT empID FOREIGN KEY (empID) REFERENCES employee(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE nurse (
    RN INT PRIMARY KEY,
    empID INT NOT NULL,
    CONSTRAINT nurseID FOREIGN KEY (empID) REFERENCES employee(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE Sector (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sector_name VARCHAR(20) NOT NULL,
    manager_start_date DATE,
    manager INT,
    CONSTRAINT managerFK FOREIGN KEY(manager) REFERENCES doctor(CRM)
);

CREATE TABLE works_in (
    emp_ID INT,
    sector_ID INT,
    CONSTRAINT empFK FOREIGN KEY (emp_ID) REFERENCES employee(id),
    CONSTRAINT sectorFK FOREIGN KEY (sector_ID) REFERENCES Sector(id),
    CONSTRAINT workPK PRIMARY KEY (emp_ID, sector_ID) 
);

CREATE TABLE shift (
    employee INT PRIMARY KEY,
    shift_start DATETIME NOT NULL,
    shift_end DATETIME NOT NULL,
    CONSTRAINT empShiftFK FOREIGN KEY (employee) REFERENCES employee(id)
);

CREATE TABLE screening (
    patient INT,
    nurse_id INT,
    diagnosis VARCHAR(50),
    dateS DATE NOT NULL,
    fowards_to INT,
    CONSTRAINT clientScreeningFK FOREIGN KEY (patient) REFERENCES client(clientID),
    CONSTRAINT nurseScreeningFK FOREIGN KEY (nurse_id) REFERENCES nurse(RN),
    CONSTRAINT fowardFK FOREIGN KEY (fowards_to) REFERENCES doctor(CRM),
    CONSTRAINT screeningPR PRIMARY KEY (nurse_id,patient)
);

CREATE TABLE appointment (
    patient INT,
    doctor INT,
    dateApt DATE,
    CONSTRAINT patientFK FOREIGN KEY (patient) REFERENCES client(clientID),
    CONSTRAINT doctorAptFK FOREIGN KEY (doctor) REFERENCES doctor(CRM),
    CONSTRAINT aptPR PRIMARY KEY (patient, doctor, dateApt)
);
