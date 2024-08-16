package mp3

import (
	"database/sql"
	"fmt"

	"github.com/NicoHernandezR/Back-end-spotychafa-go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetMp3ByID(id int) (*types.Mp3, error) {
	rows, err := s.db.Query("SELECT * FROM mp3 WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	mp3 := new(types.Mp3)

	for rows.Next() {
		mp3, err = scanRowIntoMP3(rows)
		if err != nil {
			return nil, err
		}
	}

	if mp3.ID == 0 {
		return nil, fmt.Errorf("MP3 not found")
	}

	return mp3, nil
}
func (s *Store) InsertMp3(mp3 types.Mp3) error {

	//De momento al mandar la solicitud se enviara un string, a futuro se manejara
	//Cuando se mande el archivo .mp3
	//Que significa manejar donde se guardara el .mp3
	//Y generar la url para acceder a el
	_, err := s.db.Exec("INSERT INTO mp3 (id, title, artist, mp3File) VALUES (?,?,?,?)", mp3.ID, mp3.Title, mp3.Artist, mp3.Mp3File)

	if err != nil {
		return err
	}

	return nil
}

func scanRowIntoMP3(rows *sql.Rows) (*types.Mp3, error) {
	mp3 := new(types.Mp3)

	err := rows.Scan(
		&mp3.ID,
		&mp3.Title,
		&mp3.Artist,
		&mp3.Mp3File,
	)

	if err != nil {
		return nil, err
	}

	return mp3, nil

}
