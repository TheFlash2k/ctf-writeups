package models

import (
	"database/sql"
	"strings"
)

type User struct {
	Profile         Profile
	Username        string
	ActivatedPlugin []string
	Directory       string
}

type Profile struct {
	Nickname    string
	Nationality string
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `
		SELECT 
		    p.nickname, 
		    p.nationality, 
		    u.username, 
		    u.activated_plugin, 
		    u.directory 
		FROM 
		    users u 
		JOIN 
		    profiles p ON u.profile_id = p.id 
		WHERE 
		    u.username = ?
	`

	row := db.QueryRow(query, username)

	var (
		nickname         string
		nationality      string
		uName            string
		activated        string
		directory        string
		activatedPlugins []string
	)

	if err := row.Scan(&nickname, &nationality, &uName, &activated, &directory); err != nil {
		return nil, err
	}

	if activated != "" {
		activatedPlugins = strings.Split(activated, ",")
	}

	return &User{
		Profile: Profile{
			Nickname:    nickname,
			Nationality: nationality,
		},
		Username:        uName,
		ActivatedPlugin: activatedPlugins,
		Directory:       directory,
	}, nil
}
