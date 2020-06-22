package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/go-frame/database"
)

// NotifyRepository struct
type NotifyRepository struct {
	db *database.DB
}

// Init method
func (repo *NotifyRepository) Init(db *database.DB) {
	repo.db = db
}

// CreateUser method
func (repo *NotifyRepository) CreateUser(user models.User) (models.User, error) {
	statement := `
    insert into users (type, firstname, middlename, lastname, mobilenumber, email, password, expireddate, status)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, user.Type, user.FirstName, user.MiddleName, user.LastName,
		user.MobileNumber, user.Email, user.Password, user.ExpiredDate, user.Status).Scan(&id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return user, err
	}
	createdUser := user
	createdUser.ID = id
	return createdUser, nil
}

// GetUserByEmail method
func (repo *NotifyRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `
	select id, firstname, middlename, lastname, mobilenumber, email
	from users 
	where email = $1
    `
	row := repo.db.Conn.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber, &user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetPasswordByEmail method
func (repo *NotifyRepository) GetPasswordByEmail(email string) (models.User, error) {
	var user models.User
	query := `
	select id, firstname, middlename, lastname, mobilenumber, email, password, lastattempt, failedattempts, status, enabled
	from users 
	where email = $1 and deleted=$2 and enabled=$3
    `
	row := repo.db.Conn.QueryRow(query, email, false, true)
	err := row.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber, &user.Email,
		&user.Password, &user.LastAttempt, &user.FailedAttempts, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByEmailOrMobileNumber method
func (repo *NotifyRepository) GetUserByEmailOrMobileNumber(email string, mobilenumber string) ([]models.User, error) {
	query := `
	select mobilenumber, email
	from users 
	where email = $1 or mobilenumber = $2
	`
	rows, err := repo.db.Conn.Query(query, email, mobilenumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getUsersFromRows(rows)
}

// UpdateUser method
func (repo *NotifyRepository) UpdateUser(user models.User) error {
	query := `
    update users
	set 
		status=$3, 
		enabled=$4,
		failedattempts=$5,
		expireddate=$6,
		updatedat=CURRENT_TIMESTAMP
    where id=$1 and deleted=$2
  	`

	res, err := repo.db.Conn.Exec(query, user.ID, false, user.Status, user.Enabled, user.FailedAttempts, user.ExpiredDate)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated User for %d", user.ID)
	}

	return nil
}

// UpdateUserLoginAttempt method
func (repo *NotifyRepository) UpdateUserLoginAttempt(user models.User) error {
	query := `
    update users
	set 
		lastattempt=$2, 
		failedattempts=$3,
		lastchange=$4,
		status=$5		
    where id=$1
  	`

	res, err := repo.db.Conn.Exec(query, user.ID, user.LastAttempt, user.FailedAttempts, user.LastChange, user.Status)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated User for %d", user.ID)
	}

	return nil
}

// GetUsers method
func (repo *NotifyRepository) GetUsers(page int64) ([]models.User, error) {
	query := `
	select id, type, firstname, middlename, lastname, mobilenumber, email, failedattempts, lastattempt, lastchange, expireddate, createdat, updatedat, status, enabled
	from users 
	where deleted=$1
	order BY createdat limit 50 OFFSET $2-1
    `
	rows, err := repo.db.Conn.Query(query, false, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Type, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber,
			&user.Email, &user.FailedAttempts, &user.LastAttempt, &user.LastChange, &user.ExpiredDate, &user.CreatedAt, &user.UpdatedAt, &user.Status, &user.Enabled)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUser method
func (repo *NotifyRepository) GetUser(id int64) (models.User, error) {
	var user models.User
	query := `
	select id, type, firstname, middlename, lastname, mobilenumber, email, failedattempts, lastattempt, lastchange, expireddate, createdat, updatedat, status, enabled
	from users 
	where id=$1 and deleted=$2 and enabled=$3
	`
	row := repo.db.Conn.QueryRow(query, id, false, true)
	err := row.Scan(&user.ID, &user.Type, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber,
		&user.Email, &user.FailedAttempts, &user.LastAttempt, &user.LastChange, &user.ExpiredDate, &user.CreatedAt, &user.UpdatedAt, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// prepare users row
func getUsersFromRows(rows *sql.Rows) ([]models.User, error) {
	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.MobileNumber)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
