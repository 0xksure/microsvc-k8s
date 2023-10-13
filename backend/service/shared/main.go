package shared

import (
	"database/sql"
	"fmt"
	"log"
)

type Args struct{}
type Reply string
type MessageServer struct {
	// this is fuckin scary
	Db *sql.DB
}

func (t *MessageServer) GetMessage(args *Args, reply *Reply) error {

	res, err := t.Db.Exec(`
		INSERT INTO records(name) 
		VALUES('getMessage')
	`)
	if err != nil {
		log.Print("failed to insert into records. Cause: ", err)
		return nil
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Print("failed to get affected rows. Cause: ", err)
		return nil
	}
	*reply = Reply(fmt.Sprintf("hello from your server. You affected %d rows", rowsAffected))

	return nil
}
