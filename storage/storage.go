package storage

import "market/models"

type StorageI interface {
	Close()
	Branch() BranchRepoI
	Category() CategoryRepoI
	Product() ProductRepoI
	ComingTable() ComingTableRepoI
	ComingTableProduct() ComingTableProductRepoI
	Remaining() RemainingRepoI
}

type BranchRepoI interface {
	Create(*models.CreateBranch) (string, error)
	GetByID(*models.BranchPrimaryKey) (*models.Branch, error)
	GetList(*models.BranchGetListRequest) (*models.BranchGetListResponse, error)
	Update(*models.UpdateBranch) (string, error)
	Delete(*models.BranchPrimaryKey) error
}

type CategoryRepoI interface {
	Create(*models.CreateCategory) (string, error)
	GetByID(*models.CategoryPrimaryKey) (*models.Category, error)
	GetList(*models.CategoryGetListRequest) (*models.CategoryGetListResponse, error)
	Update(*models.UpdateCategory) (string, error)
	Delete(*models.CategoryPrimaryKey) error
}

type ProductRepoI interface {
	Create(*models.CreateProduct) (string, error)
	GetByID(*models.ProductPrimaryKey) (*models.Product, error)
	GetList(*models.ProductGetListRequest) (*models.ProductGetListResponse, error)
	Update(*models.UpdateProduct) (string, error)
	Delete(*models.ProductPrimaryKey) error

	GetByBarcode(req *models.ProductBarcodeRequest) (*models.ProductBarcodeResponse, error)
}

type ComingTableRepoI interface {
	Create(*models.CreateComingTable) (string, error)
	GetByID(*models.ComingTablePrimaryKey) (*models.ComingTable, error)
	GetList(*models.ComingTableGetListRequest) (*models.ComingTableGetListResponse, error)
	Update(*models.UpdateComingTable) (string, error)
	UpdateStatus(*models.ComingTablePrimaryKey) (string, error)
	Delete(*models.ComingTablePrimaryKey) error

	GetStatus(*models.ComingTablePrimaryKey) (string, error)
}

type ComingTableProductRepoI interface {
	Create(*models.CreateComingTableProduct) (string, error)
	GetByID(*models.ComingTableProductPrimaryKey) (*models.ComingTableProduct, error)
	GetList(*models.ComingTableProductGetListRequest) (*models.ComingTableProductGetListResponse, error)
	Update(*models.UpdateComingTableProduct) (string, error)
	Delete(*models.ComingTableProductPrimaryKey) error

	CheckExistProduct(*models.ComingTableProductBarcode) (string, error)
	UpdateIdExists(req *models.UpdateComingTableProduct) (string, error)
	GetByComingTableId(req *models.ComingTableProductPrimaryKey) (*models.ComingTableProduct, error)
}

type RemainingRepoI interface {
	Create(*models.CreateRemaining) (string, error)
	GetByID(*models.RemainingPrimaryKey) (*models.Remaining, error)
	GetList(*models.RemainingGetListRequest) (*models.RemainingGetListResponse, error)
	Update(*models.UpdateRemaining) (string, error)
	Delete(*models.RemainingPrimaryKey) error

	CheckRemaing(*models.CheckingRemaining) (string, error)
	UpdateExists(req *models.UpdateRemaining) (string, error)
}
