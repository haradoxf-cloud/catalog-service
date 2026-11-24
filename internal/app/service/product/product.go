package sproduct

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/haradoxf-cloud/catalog-service/internal/app/entity"
	"github.com/haradoxf-cloud/catalog-service/internal/app/repository"
	"github.com/haradoxf-cloud/catalog-service/internal/app/service"
)

type srv struct {
	repoProduct  repository.Product
	repoCategory repository.Category
}

func NewService(repoProduct repository.Product, repoCategory repository.Category) service.Product {
	return &srv{
		repoProduct:  repoProduct,
		repoCategory: repoCategory,
	}
}

func (s *srv) Create(ctx context.Context, req entity.RequestProductCreate) (entity.Product, error) {
	existing, err := s.repoProduct.List(ctx, &req.Name, nil, nil, nil)
	if err != nil {
		return entity.Product{}, err
	}
	if len(existing) > 0 {
		return entity.Product{}, entity.ErrAlreadyExists
	}
	categories, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{req.CategoryGUID})
	if err != nil {
		return entity.Product{}, err
	}
	if len(categories) == 0 {
		return entity.Product{}, entity.ErrNotFound
	}
	now := time.Now()
	product := entity.Product{
		GUID:         uuid.Must(uuid.NewV4()),
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryGUID: req.CategoryGUID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repoProduct.Create(ctx, product); err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (s *srv) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error) {
	return s.repoProduct.GetByGUIDs(ctx, guids)
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	existing, err := s.repoProduct.GetByGUIDs(ctx, []uuid.UUID{guid})
	if err != nil {
		return entity.Product{}, err
	}
	if len(existing) == 0 {
		return entity.Product{}, entity.ErrNotFound
	}
	product := existing[0]
	if req.Name != "" {
		found, err := s.repoProduct.List(ctx, &req.Name, nil, nil, nil)
		if err != nil {
			return entity.Product{}, err
		}
		for _, d := range found {
			if d.GUID != guid {
				return entity.Product{}, entity.ErrAlreadyExists
			}
		}
	}
	if !req.CategoryGUID.IsNil() {
		categories, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{req.CategoryGUID})
		if err != nil {
			return entity.Product{}, err
		}
		if len(categories) == 0 {
			return entity.Product{}, entity.ErrNotFound
		}
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if !req.CategoryGUID.IsNil() {
		product.CategoryGUID = req.CategoryGUID
	}
	product.UpdatedAt = time.Now()
	if err := s.repoProduct.Update(ctx, product); err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	existing, err := s.repoProduct.GetByGUIDs(ctx, []uuid.UUID{guid})
	if err != nil {
		return err
	}
	if len(existing) == 0 {
		return entity.ErrNotFound
	}
	if err := s.repoProduct.Delete(ctx, guid); err != nil {
		return err
	}
	return nil
}

func (s *srv) List(ctx context.Context, req entity.RequestProductList) ([]entity.Product, error) {
	return s.repoProduct.List(ctx, nil, req.CategoryGUID, req.MinPrice, req.MaxPrice)
}
