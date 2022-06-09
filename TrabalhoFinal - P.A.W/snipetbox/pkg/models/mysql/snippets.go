package mysql

import (
	"database/sql"

	"github.com/rmcs87/cc5m/pkg/models"
)
type SnippetModel struct{
  DB *sql.DB
}
func(m *SnippetModel)Insert(Nome, Contato, Saida string) (int, error){
  stmt := `INSERT INTO snippets (Nome, Contato, Entrada, Saida) 
            VALUES(?,?,UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

  result, err := m.DB.Exec(stmt, Nome, Contato, Saida)
  if err != nil{
    return 0, err
  }
  
  id, err := result.LastInsertId()
  if err != nil{
    return 0, err
  }
  return int(id),nil
}

func(m *SnippetModel)Update(id int, Nome, Contato string) (int, error){
  stmt := `UPDATE snippets SET Nome=?, Contato=? WHERE id = ?`

  _, err := m.DB.Exec(stmt, Nome, Contato, id)
  if err != nil{
    return 0, err
  }
  
  return int(id),nil
}

func(m *SnippetModel)Delete(id int) (int, error){
  stmt := `DELETE FROM snippets WHERE id = ?`

  _, err := m.DB.Exec(stmt, id)
  if err != nil{
    return id, err
  }
  return id,nil
}

func(m *SnippetModel) Get(id int)(*models.Snippet, error){
  stmt := `SELECT id, Nome, Contato, Entrada, Saida FROM snippets
           WHERE Saida > UTC_TIMESTAMP() AND id = ?`
  row := m.DB.QueryRow(stmt, id)

  s := &models.Snippet{}

  err := row.Scan(&s.ID, &s.Nome, &s.Contato, &s.Entrada, &s.Saida)
  if err == sql.ErrNoRows{
    return nil, models.ErrNoRecord
  }else if err != nil{
    return nil, err
  }
  
  return s, nil
}
func(m * SnippetModel) Latest()([]*models.Snippet, error){
  stmt := `SELECT id, Nome, Contato, Entrada, Saida FROM snippets
           WHERE Saida > UTC_TIMESTAMP() ORDER BY Entrada DESC LIMIT 10`

  rows, err := m.DB.Query(stmt)
  if err != nil{
    return nil, err
  }
  defer rows.Close()

  snippets := []*models.Snippet{}
  for rows.Next(){
    s := &models.Snippet{}
    err = rows.Scan(&s.ID, &s.Nome, &s.Contato, &s.Entrada, &s.Saida)
    if err != nil{
      return nil, err
    }
    snippets = append(snippets, s)
  }
  err = rows.Err()
  if err != nil {
    return  nil, err
  }
  return snippets, nil
}
