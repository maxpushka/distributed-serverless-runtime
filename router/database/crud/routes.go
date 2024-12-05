package crud

import (
	"database/sql"
	"log"

	"serverless/router/schema"
)

func SaveRoute(db *sql.DB, user schema.User, route schema.RouteName) (*schema.Route, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO routes (user_id, name) VALUES ($1, $2) RETURNING id;",
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
	rows, err := db.Query(
		"SELECT id, name, executable_exists FROM routes WHERE user_id = $1",
		user.UserId,
	)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

	routes := []schema.Route{}
	for rows.Next() {
		var id int
		var name string
		var executableExists bool
		err := rows.Scan(&id, &name, &executableExists)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		routes = append(routes, schema.Route{
			Id:               id,
			Name:             name,
			ExecutableExists: executableExists,
		})
	}
	return routes, nil
}

func GetRoute(db *sql.DB, user schema.User, routeId int) (*schema.Route, error) {
	var id int
	var name string
	var executableExists bool
	err := db.QueryRow(
		"SELECT id, name, executable_exists FROM routes WHERE id = $1 AND user_id = $2",
		routeId,
		user.UserId,
	).Scan(&id, &name, &executableExists)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &schema.Route{
		Id:               id,
		Name:             name,
		ExecutableExists: executableExists,
	}, nil
}

func UpdateRoute(db *sql.DB, user schema.User, routeId int, route schema.RouteName) error {
	_, err := db.Exec(
		"UPDATE routes SET name = $1 WHERE id = $2 AND user_id = $3",
		route.Name,
		routeId,
		user.UserId,
	)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func DeleteRoute(db *sql.DB, user schema.User, routeId int) error {
	_, err := db.Exec(
		"DELETE FROM routes WHERE id = $1 AND user_id = $2",
		routeId,
		user.UserId,
	)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func SetExecutable(db *sql.DB, user schema.User, routeId int) error {
	_, err := db.Exec(
		"UPDATE routes SET executable_exists = true WHERE id = $1 AND user_id = $2",
		routeId,
		user.UserId,
	)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
