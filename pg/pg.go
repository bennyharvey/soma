package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/gobuffalo/packr"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/bennyharvey/soma/entity"
)

type Storage struct {
	db  *sqlx.DB
	uri string

	personFaces   []entity.PersonFace
	personFacesMx sync.RWMutex

	log *logrus.Entry
}

func NewStorage(uri string) (
	*Storage, error) {
	db, err := sqlx.Open("postgres", uri)
	if err != nil {
		return nil, errors.New("open DB: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)

	err = db.PingContext(ctx)
	cancel()
	if err != nil {
		return nil, errors.New("ping DB: " + err.Error())
	}

	return &Storage{
		db:  db,
		uri: uri,
		log: logrus.WithField("subsystem", "postgres_storage"),
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

//go:generate packr

const migrationsPath = "./migrations"

func (s *Storage) Migrate() error {
	packrSource := &migration.PackrMigrationSource{
		Box: packr.NewBox(migrationsPath),
	}

	d, err := postgres.New(s.uri)
	if err != nil {
		return errors.New("create migration driver: " + err.Error())
	}

	_, err = migration.Migrate(d, packrSource, migration.Up, 0)
	if err != nil {
		return errors.New("migrate: " + err.Error())
	}

	return nil
}

func (s *Storage) LoadPersonFaces() error {
	s.personFacesMx.Lock()
	defer s.personFacesMx.Unlock()

	return s.db.Select(&s.personFaces, `select * from person_face`)
}

func (s *Storage) FindClosestPersonFace(fd entity.FaceDescriptor) (entity.PersonFace, float64, bool) {
	s.personFacesMx.RLock()
	defer s.personFacesMx.RUnlock()

	if len(s.personFaces) == 0 {
		return entity.PersonFace{}, 0, false
	}

	closestPF := s.personFaces[0]
	closestDistance := entity.FaceDescriptorDistance(closestPF.Descriptor, fd)

	for _, pf := range s.personFaces[1:] {
		distance := entity.FaceDescriptorDistance(pf.Descriptor, fd)
		if distance < closestDistance {
			closestPF = pf
			closestDistance = distance
		}
	}

	return closestPF, closestDistance, true
}

func (s *Storage) AddEvent(e entity.Event) error {
	_, err := s.db.Exec(`insert into event (time, type, passage_id , data) values ($1, $2, $3, $4)`, e.Time, e.Type, e.PassageID, e.Data)
	return err
}

func (s *Storage) User(login string) (u entity.User, err error) {
	err = s.db.QueryRowx(`
		SELECT * FROM skuder_user WHERE login = $1
	`, login).StructScan(&u)
	if err == sql.ErrNoRows {
		err = entity.ErrUserNotFound
	}
	return
}

func (s *Storage) AddUser(u entity.User) (err error) {
	_, err = s.db.Exec(`
		INSERT INTO skuder_user (login, password_hash, role) VALUES ($1, $2, $3) 
	`, u.Login, u.PasswordHash, u.Role)
	return
}

func (s *Storage) SetUser(u entity.User) (err error) {
	var setPasswordHash string
	args := []interface{}{u.Login, u.Role}
	if len(u.PasswordHash) > 0 {
		setPasswordHash = ", password_hash = $3"
		args = append(args, u.PasswordHash)
	}
	_, err = s.db.Exec(`
		UPDATE skuder_user SET role = $2`+setPasswordHash+` 
		WHERE login = $1
	`, args...)
	return
}

func (s *Storage) RemoveUser(login string) (err error) {
	_, err = s.db.Exec(`DELETE FROM skuder_user WHERE login = $1`, login)
	return
}

func (s *Storage) Users() (us []entity.User, err error) {
	err = s.db.Select(&us, `SELECT * FROM skuder_user`)
	return
}

func (s *Storage) Person(personID int64) (p entity.Person, err error) {
	err = s.db.QueryRowx(`SELECT * FROM person WHERE id = $1`, personID).StructScan(&p)
	return
}

func (s *Storage) AddPerson(p entity.Person) (entity.Person, error) {
	err := s.db.QueryRow(`
		INSERT INTO person (name, position, unit)
		VALUES ($1, $2, $3)
		RETURNING id
	`, p.Name, p.Position, p.Unit).Scan(&p.ID)
	return p, err
}

func (s *Storage) SetPerson(p entity.Person) (err error) {
	_, err = s.db.Exec(`
		UPDATE person SET name = $1, position = $2, unit = $3
		WHERE id = $4
	`, p.Name, p.Position, p.Unit, p.ID)
	return
}

func (s *Storage) removePersonFaces(personID int64) {
	s.personFacesMx.Lock()
	defer s.personFacesMx.Unlock()

	var toRemove int

	for i := len(s.personFaces) - 1; i >= 0; i-- {
		if s.personFaces[i].PersonID == personID {
			if i == len(s.personFaces)-1-toRemove {
				toRemove++
				continue
			}
			s.personFaces[i], s.personFaces[len(s.personFaces)-toRemove-1] =
				s.personFaces[len(s.personFaces)-toRemove-1], s.personFaces[i]
			toRemove++
		}
	}

	s.personFaces = s.personFaces[:len(s.personFaces)-toRemove]
}

func (s *Storage) RemovePerson(personID int64) error {
	_, err := s.db.Exec(`DELETE FROM person WHERE id = $1`, personID)
	if err != nil {
		return err
	}
	s.removePersonFaces(personID)
	return nil
}

func (s *Storage) Persons() (ps []entity.Person, err error) {
	err = s.db.Select(&ps, `
		SELECT * FROM person
	`)
	return
}

func (s *Storage) AddPersonFace(pf entity.PersonFace) (entity.PersonFace, error) {
	err := s.db.QueryRow(`
		INSERT INTO person_face(person_id, descriptor, photo_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, pf.PersonID, pf.Descriptor, pf.PhotoID).Scan(&pf.ID)
	if err != nil {
		return entity.PersonFace{}, err
	}

	s.personFacesMx.Lock()
	s.personFaces = append(s.personFaces, pf)
	s.personFacesMx.Unlock()

	return pf, nil
}

func (s *Storage) RemovePersonFace(id int64) error {
	_, err := s.db.Exec(`
		DELETE FROM person_face WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	var removeI = -1

	s.personFacesMx.Lock()

	for i, pf := range s.personFaces {
		if pf.ID == id {
			removeI = i
			break
		}
	}

	if removeI >= 0 {
		s.personFaces[removeI] =
			s.personFaces[len(s.personFaces)-1]
		s.personFaces = s.personFaces[:len(s.personFaces)-1]
	}

	s.personFacesMx.Unlock()

	return nil
}

func (s *Storage) PersonFaces(personID int64) (pfs []entity.PersonFace,
	err error) {
	err = s.db.Select(&pfs, `
		SELECT * FROM person_face WHERE person_id = $1
	`, personID)
	return
}

func (s *Storage) Events(fs ...entity.EventsFilter) (
	es []entity.Event, err error) {

	var filters entity.EventsFilters

	for _, f := range fs {
		f(&filters)
	}

	var (
		args   []interface{}
		wheres []string
	)

	if filters.Set.From {
		args = append(args, filters.From)
		wheres = append(wheres, fmt.Sprintf("time >= $%d", len(args)))
	}

	if filters.Set.To {
		args = append(args, filters.To)
		wheres = append(wheres, fmt.Sprintf("time <= $%d", len(args)))
	}

	if filters.Set.Type {
		args = append(args, filters.Type)
		wheres = append(wheres, fmt.Sprintf("type = $%d", len(args)))
	}

	if filters.Set.PassageID {
		args = append(args, filters.PassageID)
		wheres = append(wheres, fmt.Sprintf("data->>'passage_id' = $%d", len(args)))
	}

	if filters.Set.PersonID {
		args = append(args, filters.PersonID)
		wheres = append(wheres, fmt.Sprintf("(data->>'person_id')::BIGINT = $%d", len(args)))
	}

	if filters.Set.PersonName {
		args = append(args, filters.PersonName)
		wheres = append(wheres, fmt.Sprintf("data->>'person_name' ILIKE '%%' || $%d || '%%'", len(args)))
	}

	if filters.Set.PersonPosition {
		args = append(args, filters.PersonPosition)
		wheres = append(wheres, fmt.Sprintf("data->>'person_position' ILIKE '%%' || $%d || '%%'", len(args)))
	}

	if filters.Set.PersonUnit {
		args = append(args, filters.PersonUnit)
		wheres = append(wheres, fmt.Sprintf("data->>'person_unit' ILIKE '%%' || $%d || '%%'", len(args)))
	}

	where := strings.Join(wheres, " and ")
	if where != "" {
		where = "where " + where
	}

	var order string

	if filters.Set.OrderBy {
		switch filters.OrderBy {
		case "id", "time":
			order = filters.OrderBy
		case "passage_id", "person_name", "person_position", "person_unit":
			order = fmt.Sprintf("data->>'%s'", filters.OrderBy)
		case "person_id":
			order = "(data->>'person_id')::bigint"
		default:
			return nil, entity.InvalidParamErr{Param: "order_by"}
		}
	}

	if order != "" {
		order = "ORDER BY " + order
		if filters.Set.OrderDirection {
			switch filters.OrderDirection {
			case "asc", "desc":
				order += " " + filters.OrderDirection
			default:
				return nil, entity.InvalidParamErr{Param: "order_direction"}
			}
		}
	}

	var limit string

	if filters.Set.Limit {
		if filters.Limit <= 0 {
			return nil, entity.InvalidParamErr{Param: "limit"}
		}
		limit = fmt.Sprintf("limit %d", filters.Limit)
	}

	var offset string

	if filters.Set.OffSet {
		if filters.Offset < 0 {
			return nil, entity.InvalidParamErr{Param: "offset"}
		}
		offset = fmt.Sprintf("offset %d", filters.Offset)
	}

	err = s.db.Select(&es, `select * from event `+where+` `+order+` `+limit+` `+offset, args...)

	return
}
