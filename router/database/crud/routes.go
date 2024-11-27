package crud

import (
	"database/sql"
	"log"
	"serverless/router/schema"
)

func SaveRoute(db *sql.DB, user schema.User, route schema.RouteName) (*schema.Route, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO routes (user_id, route) VALUES ($1, $2) RETURNING id;",
		user.UserId,
		route.Name,
	).Scan(&id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &schema.Route{Id: id, Name: route.Name}, nil
}

func GetRoutes(db *sql.DB, user schema.User) ([]schema.Route, error) {
	rows, err := db.Query("SELECT id, name FROM routes WHERE user_id = $1", user.UserId)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

	routes := []schema.Route{}
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		routes = append(routes, schema.Route{Id: id, Name: name})
	}
	return routes, nil
}
