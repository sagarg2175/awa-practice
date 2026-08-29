package repo

import (
	"fmt"
	"main/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
)

type AwsProjRepo struct {
	db *Dbstruct
}

func NewAwsProjRepo(db *Dbstruct) *AwsProjRepo {
	return &AwsProjRepo{
		db,
	}
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (sr *AwsProjRepo) CreateAwsProj(c *gin.Context, req *domain.Master) error {
	// Build SQL - use sq.Expr for NOW()
	query := psql.Insert("scratch").
		Columns("name", "lastupdated", "domain").
		Values(req.Name, sq.Expr("NOW()"), req.Domain).
		Suffix("RETURNING name, password") // <-- must match Scan below

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	fmt.Println("SQL:", sql)
	fmt.Println("Args:", args)

	if err := sr.db.QueryRow(c, sql, args...).Scan(&req.Name, &req.Domain); err != nil {
		fmt.Println("DB Error:", err) 
		return err
	}
	return nil
}
