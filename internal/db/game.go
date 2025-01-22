package db

import (
	"connectfour/internal/model"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type MariaDbGameRepository struct {
	db *sql.DB
}

var _ GameRepository = MariaDbGameRepository{}

func NewMariaDbGameRepository() *MariaDbGameRepository {
	return &MariaDbGameRepository{
		db: connect(),
	}
}

func (r MariaDbGameRepository) Save(g model.Game) bool {
	_, err := r.db.Exec(
		`REPLACE INTO game (
                   game_key, 
                   player1_id, 
                   player2_id, 
                   created_at, 
                   started_at, 
                   finished_at, 
                   player_turn_id, 
                   public, 
                   status,
                   board_json) 
			   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		g.Key, g.Player1.Id, g.Player2.Id, g.CreatedAt, g.StartedAt, g.FinishedAt, g.CurrentPlayer().Id, g.Public, g.Status, g.Board.String())
	if err != nil {
		log.Errorf("Error saving the game into the database: %v\n", err)
		return false
	}
	return true
}

func (r MariaDbGameRepository) Fetch(key string) (model.Game, error) {
	row := r.db.QueryRow(`SELECT 
    game_key, 
    g.board_json, 
    u1.email as player1_email,
    u1.id as player1_id,
    u1.name as player1_name,
    ifnull(u2.email, "") as player2_email,
    ifnull(u2.id, 0) as player2_id,
    ifnull(u2.name, "") as player2_name,
    g.player_turn_id, 
    g.created_at, 
    g.started_at, 
    g.finished_at, 
    g.status, 
    g.public 
	FROM game g
	JOIN user u1 ON u1.id = g.player1_id
	LEFT JOIN user u2 ON u2.id = g.player2_id
	WHERE game_key = ?`, key)

	var g model.Game
	var p1 model.User
	var p2 model.User
	var playerTurnId int64
	var boardJson string
	err := row.Scan(
		&g.Key,
		&boardJson,
		&p1.Email,
		&p1.Id,
		&p1.Name,
		&p2.Email,
		&p2.Id,
		&p2.Name,
		&playerTurnId,
		&g.CreatedAt,
		&g.StartedAt,
		&g.FinishedAt,
		&g.Status,
		&g.Public,
	)

	if err != nil {
		log.Errorf("Error scanning the game row: %v\n", err)
		return model.Game{}, err
	}

	g.Player1 = p1
	g.Player2 = p2
	if playerTurnId == p1.Id {
		g.PlayerTurn = 1
	} else {
		g.PlayerTurn = 2
	}

	g.Board, err = model.BoardFromString(boardJson)
	return g, nil
}

func (r MariaDbGameRepository) List(userId int64, status string) ([]model.Game, error) {
	// select games that are open
	rows, err := r.db.Query(`	SELECT 
    g.game_key, 
    u1.email as player1_email,
    u1.id as player1_id,
    u1.name as player1_name,
    ifnull(u2.email, "") as player2_email,
    ifnull(u2.id, 0) as player2_id,
    ifnull(u2.name, "") as player2_name,
    g.player_turn_id, 
    g.created_at, 
    g.started_at, 
    g.finished_at, 
    g.status, 
    g.public 
	FROM game g
	JOIN user u1 ON u1.id = g.player1_id
	LEFT JOIN user u2 ON u2.id = g.player2_id
	WHERE status = ? `, status)

	if err == nil {
		err = rows.Err()
	}

	if err != nil {
		log.Errorf("Error getting all the open games from the database: %v\n", err)
		return nil, err
	}

	output := make([]model.Game, 0)

	for rows.Next() {
		var g model.Game
		var p1 model.User
		var p2 model.User
		var playerTurnId int64
		err = rows.Scan(
			&g.Key,
			&p1.Email,
			&p1.Id,
			&p1.Name,
			&p2.Email,
			&p2.Id,
			&p2.Name,
			&playerTurnId,
			&g.CreatedAt,
			&g.StartedAt,
			&g.FinishedAt,
			&g.Status,
			&g.Public,
		)

		if err != nil {
			log.Errorf("Error scanning the game row: %v\n", err)
			_ = rows.Close()
			return nil, err
		}

		g.Player1 = p1
		g.Player2 = p2
		if playerTurnId == p1.Id {
			g.PlayerTurn = 1
		} else {
			g.PlayerTurn = 2
		}
		output = append(output, g)
	}

	return output, nil
}
