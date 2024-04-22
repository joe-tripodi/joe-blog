package database

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
)

type Post struct {
	ID      int64
	Title   string
	Url     string
	Content string
	Tag     string
  Html    template.HTML
}

func GetPost(db *sql.DB, url string) (*Post, error) {
  log.SetPrefix("database: GetPost: ")
  log.SetFlags(0)

  row := db.QueryRow("SELECT ID, Title, Url, Content, Tag FROM post WHERE url = ?", url)
  
  post := &Post{}
  err := row.Scan(&post.ID, &post.Title, &post.Url, &post.Content, &post.Tag)
  if err != nil {
    log.Print(err)
    return nil, fmt.Errorf("post: %q: %v", url, err)
  }

  err = row.Err()
  if err != nil {
    return nil, fmt.Errorf("post: %q: %v", url, err)
  }
  post.Html = template.HTML(post.Content)

	return post, nil
}

func GetAllPosts(db *sql.DB) ([]*Post, error){
  allPosts := []*Post{}
  rows, err := db.Query("SELECT * FROM post")
  if err != nil {
    return nil, fmt.Errorf("getallposts: %v", err) 
  }
  defer rows.Close()

  for rows.Next() {
    post := &Post{}
    err := rows.Scan(&post.ID, &post.Title, &post.Url, &post.Content, &post.Tag)
    if err != nil {
      return nil, fmt.Errorf("getallpost scan row error: %v", err)
    }
    allPosts = append(allPosts, post)
  }

  return allPosts, nil
}
