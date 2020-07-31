package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/models"
)

// GFUserRepository struct
type GFUserRepository struct {
	db *database.DB
}

// Init method
func (repo *GFUserRepository) Init(db *database.DB) {
	repo.db = db
}

// AddUser makes changes to the user
func (repo *GFUserRepository) AddUser(user models.GFUser) (models.GFUser, error) {
	statement := `
    insert into gfuser (relatedId, email, key, createdBy, updatedBy)
    values ($1, $2, $3, $4, $5)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, user.RelatedID, user.Email, user.Key, user.CreatedBy, user.UpdatedBy).Scan(&id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return user, err
	}
	user.ID = id
	return user, nil
}

// UpdateUser makes changes to the user
func (repo *GFUserRepository) UpdateUser(user models.GFUser) error {
	query := `
    update gfuser
	set 
		key=$2,
		updatedBy=$3,
		enabled=$4		
    where id=$1 and deleted=false
  	`
	res, err := repo.db.Conn.Exec(query, user.ID, user.Key, user.UpdatedBy, user.Enabled)
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

// GetUsers method returns users from the database
func (repo *GFUserRepository) GetUsers(page int64) ([]models.GFUser, error) {
	query := `
	select id, relatedId, email, key, createdBy, createdOn, updatedBy, updatedOn, enabled, deleted 
	from gfuser 
	where deleted = false
	order by createdOn asc limit 500 OFFSET $1-1
	`
	rows, err := repo.db.Conn.Query(query, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return userMapper(rows)
}

// DeleteUser makes changes to the user delete status
func (repo *GFUserRepository) DeleteUser(id int64, userID int64) error {
	query := `
    update gfuser
	set 
		email=CONCAT(email, '-', 'DELETED'),
		updatedBy=$2,
		updatedOn=CURRENT_TIMESTAMP,
		enabled=false,
		deleted=true
    where id=$1
  	`
	res, err := repo.db.Conn.Exec(query, id, userID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated user for %d", id)
	}

	return nil
}

// prepare users row
func userMapper(rows *sql.Rows) ([]models.GFUser, error) {
	users := []models.GFUser{}
	for rows.Next() {
		var user models.GFUser
		err := rows.Scan(&user.ID, &user.RelatedID, &user.Email, &user.Key, &user.CreatedBy, &user.CreatedOn, &user.UpdatedBy, &user.UpdatedOn, &user.Enabled, &user.Deleted)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
