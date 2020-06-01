package sqlite

import (
	"database/sql"
	"github.com/ogra1/fabrica/domain"
	"github.com/rs/xid"
	"log"
)

const createRepoTableSQL string = `
	CREATE TABLE IF NOT EXISTS repo (
		id               varchar(200) primary key not null,
		name             varchar(200) not null,
		location         varchar(200) UNIQUE not null,
		hash             varchar(200) default '',
		created          timestamp default current_timestamp,
		modified         timestamp default current_timestamp,
		branch           varchar(200) default 'master'
	)
`
const alterRepoTableSQL string = `
	ALTER TABLE repo ADD COLUMN branch varchar(200) default 'master'
`
const addRepoSQL = `
	INSERT INTO repo (id, name, location, branch) VALUES ($1, $2, $3, $4)
`
const listRepoSQL = `
	SELECT id, name, location, hash, created, modified, branch
	FROM repo
	ORDER BY name, location
`
const listRepoWatchSQL = `
	SELECT id, name, location, hash, created, modified, branch
	FROM repo
	ORDER BY modified
`
const updateRepoHashSQL = `
	UPDATE repo SET hash=$1, modified=current_timestamp WHERE id=$2
`
const getRepoSQL = `
	SELECT id, name, location, hash, created, modified, branch
	FROM repo
	WHERE id=$1
`
const deleteRepoSQL = `
	DELETE FROM repo WHERE id=$1
`

// RepoCreate creates a new repository to watch
func (db *DB) RepoCreate(name, repo, branch string) (string, error) {
	id := xid.New()
	_, err := db.Exec(addRepoSQL, id.String(), name, repo, branch)
	return id.String(), err
}

// RepoList get the list of repos
func (db *DB) RepoList(watch bool) ([]domain.Repo, error) {
	// Order the list depending on use
	sql := listRepoSQL
	if watch {
		sql = listRepoWatchSQL
	}

	records := []domain.Repo{}
	rows, err := db.Query(sql)
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		r := domain.Repo{}
		err := rows.Scan(&r.ID, &r.Name, &r.Repo, &r.LastCommit, &r.Created, &r.Modified, &r.Branch)
		if err != nil {
			return records, err
		}
		records = append(records, r)
	}

	return records, nil
}

// RepoUpdateHash updates a repo's last commit hash
func (db *DB) RepoUpdateHash(id, hash string) error {
	_, err := db.Exec(updateRepoHashSQL, hash, id)
	return err
}

// RepoGet fetches a repo from its ID
func (db *DB) RepoGet(id string) (domain.Repo, error) {
	r := domain.Repo{}
	err := db.QueryRow(getRepoSQL, id).Scan(&r.ID, &r.Name, &r.Repo, &r.LastCommit, &r.Created, &r.Modified, &r.Branch)
	switch {
	case err == sql.ErrNoRows:
		return r, err
	case err != nil:
		log.Printf("Error retrieving database repo: %v\n", err)
		return r, err
	}
	return r, nil
}

// RepoDelete removes a repo from its ID
func (db *DB) RepoDelete(id string) error {
	_, err := db.Exec(deleteRepoSQL, id)
	return err
}
