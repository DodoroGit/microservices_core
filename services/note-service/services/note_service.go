package services

import (
	"errors"

	"note-service/models"
	"note-service/repository"
)

type NoteServiceInterface interface {
	GetAll(category string) ([]models.Note, error)
	GetByID(id string) (*models.Note, error)
	Create(req models.CreateNoteRequest) (*models.Note, error)
	Update(id string, req models.UpdateNoteRequest) (*models.Note, error)
	Delete(id string) error
}

type NoteService struct {
	repo repository.NoteRepositoryInterface
}

func NewNoteService(repo repository.NoteRepositoryInterface) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) GetAll(category string) ([]models.Note, error) {
	return s.repo.FindAll(category)
}

func (s *NoteService) GetByID(id string) (*models.Note, error) {
	note, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, errors.New("note not found")
	}
	return note, nil
}

func (s *NoteService) Create(req models.CreateNoteRequest) (*models.Note, error) {
	if !req.Category.IsValid() {
		return nil, errors.New("invalid category")
	}
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	note := &models.Note{
		Category:    req.Category,
		Title:       req.Title,
		Content:     req.Content,
		Tags:        req.Tags,
		ProjectName: req.ProjectName,
		Description: req.Description,
		TechStack:   req.TechStack,
		GithubURL:   req.GithubURL,
		Year:        req.Year,
		Company:     req.Company,
		Position:    req.Position,
		Period:      req.Period,
		Projects:    req.Projects,
	}
	return s.repo.Create(note)
}

func (s *NoteService) Update(id string, req models.UpdateNoteRequest) (*models.Note, error) {
	if req.Category != nil && !req.Category.IsValid() {
		return nil, errors.New("invalid category")
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("note not found")
	}

	updates := make(map[string]interface{})

	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.ProjectName != nil {
		updates["project_name"] = *req.ProjectName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.TechStack != nil {
		updates["tech_stack"] = req.TechStack
	}
	if req.GithubURL != nil {
		updates["github_url"] = *req.GithubURL
	}
	if req.Year != nil {
		updates["year"] = *req.Year
	}
	if req.Company != nil {
		updates["company"] = *req.Company
	}
	if req.Position != nil {
		updates["position"] = *req.Position
	}
	if req.Period != nil {
		updates["period"] = *req.Period
	}
	if req.Projects != nil {
		updates["projects"] = *req.Projects
	}

	if len(updates) == 0 {
		return existing, nil
	}

	updated, err := s.repo.Update(id, updates)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errors.New("note not found")
	}
	return updated, nil
}

func (s *NoteService) Delete(id string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("note not found")
	}
	return s.repo.Delete(id)
}

// validateCreateRequest 根據類別驗證必填欄位
func validateCreateRequest(req models.CreateNoteRequest) error {
	switch req.Category {
	case models.CategoryProject:
		if req.ProjectName == "" {
			return errors.New("invalid input: 專案名稱為必填")
		}
		if req.Description == "" {
			return errors.New("invalid input: 專案介紹為必填")
		}
		if req.Year == 0 {
			return errors.New("invalid input: 建置時間為必填")
		}
	case models.CategoryTech:
		if req.Title == "" {
			return errors.New("invalid input: 標題為必填")
		}
		if req.Content == "" {
			return errors.New("invalid input: 內容為必填")
		}
	case models.CategoryResume:
		if req.Company == "" {
			return errors.New("invalid input: 任職公司為必填")
		}
		if req.Position == "" {
			return errors.New("invalid input: 任職職位為必填")
		}
		if req.Period == "" {
			return errors.New("invalid input: 任職日期為必填")
		}
	case models.CategoryDaily:
		if req.Title == "" {
			return errors.New("invalid input: 標題為必填")
		}
		if req.Content == "" {
			return errors.New("invalid input: 內容為必填")
		}
	}
	return nil
}
