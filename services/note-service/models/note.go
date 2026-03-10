package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category string

const (
	CategoryProject Category = "專案"
	CategoryTech    Category = "技術筆記"
	CategoryDaily   Category = "日常隨筆"
	CategoryResume  Category = "履歷"
)

func (c Category) IsValid() bool {
	switch c {
	case CategoryProject, CategoryTech, CategoryDaily, CategoryResume:
		return true
	}
	return false
}

// Note is the MongoDB document.
// 各類別使用的欄位不同，未使用的欄位以 omitempty 省略。
type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"      json:"id,omitempty"`
	Category  Category           `bson:"category"           json:"category"`
	CreatedAt time.Time          `bson:"created_at"         json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"         json:"updated_at"`

	// 技術筆記 & 日常隨筆
	Title   string `bson:"title,omitempty"   json:"title,omitempty"`
	Content string `bson:"content,omitempty" json:"content,omitempty"`

	// 技術筆記
	Tags []string `bson:"tags,omitempty" json:"tags,omitempty"`

	// 專案
	ProjectName string   `bson:"project_name,omitempty" json:"project_name,omitempty"`
	Description string   `bson:"description,omitempty"  json:"description,omitempty"`
	TechStack   []string `bson:"tech_stack,omitempty"   json:"tech_stack,omitempty"`
	GithubURL   string   `bson:"github_url,omitempty"   json:"github_url,omitempty"`
	Year        int      `bson:"year,omitempty"         json:"year,omitempty"`

	// 履歷（TechStack 與專案共用）
	Company  string `bson:"company,omitempty"  json:"company,omitempty"`
	Position string `bson:"position,omitempty" json:"position,omitempty"`
	Period   string `bson:"period,omitempty"   json:"period,omitempty"`
	Projects string `bson:"projects,omitempty" json:"projects,omitempty"` // 負責專案
}

// NoteResponse is the API response (ID as hex string).
type NoteResponse struct {
	ID        string    `json:"id"`
	Category  Category  `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title   string   `json:"title,omitempty"`
	Content string   `json:"content,omitempty"`
	Tags    []string `json:"tags,omitempty"`

	ProjectName string   `json:"project_name,omitempty"`
	Description string   `json:"description,omitempty"`
	TechStack   []string `json:"tech_stack,omitempty"`
	GithubURL   string   `json:"github_url,omitempty"`
	Year        int      `json:"year,omitempty"`

	Company  string `json:"company,omitempty"`
	Position string `json:"position,omitempty"`
	Period   string `json:"period,omitempty"`
	Projects string `json:"projects,omitempty"`
}

func (n *Note) ToResponse() NoteResponse {
	return NoteResponse{
		ID:          n.ID.Hex(),
		Category:    n.Category,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
		Title:       n.Title,
		Content:     n.Content,
		Tags:        n.Tags,
		ProjectName: n.ProjectName,
		Description: n.Description,
		TechStack:   n.TechStack,
		GithubURL:   n.GithubURL,
		Year:        n.Year,
		Company:     n.Company,
		Position:    n.Position,
		Period:      n.Period,
		Projects:    n.Projects,
	}
}

// CreateNoteRequest is the payload for POST /notes.
type CreateNoteRequest struct {
	Category Category `json:"category" binding:"required"`

	// 技術筆記 & 日常隨筆
	Title   string `json:"title"`
	Content string `json:"content"`

	// 技術筆記
	Tags []string `json:"tags"`

	// 專案
	ProjectName string   `json:"project_name"`
	Description string   `json:"description"`
	TechStack   []string `json:"tech_stack"`
	GithubURL   string   `json:"github_url"`
	Year        int      `json:"year"`

	// 履歷
	Company  string `json:"company"`
	Position string `json:"position"`
	Period   string `json:"period"`
	Projects string `json:"projects"`
}

// UpdateNoteRequest uses pointer/slice fields; nil means "do not update".
type UpdateNoteRequest struct {
	Category *Category `json:"category"`

	Title   *string `json:"title"`
	Content *string `json:"content"`
	Tags    []string `json:"tags"`

	ProjectName *string  `json:"project_name"`
	Description *string  `json:"description"`
	TechStack   []string `json:"tech_stack"`
	GithubURL   *string  `json:"github_url"`
	Year        *int     `json:"year"`

	Company  *string `json:"company"`
	Position *string `json:"position"`
	Period   *string `json:"period"`
	Projects *string `json:"projects"`
}
