package models

type ContactRepo interface {
	Read(name string) (*[]Contact, error)
	ReadAll() (*[]Contact, error)
	GetContact(id int) (*Contact, error)
	New(contact *Contact) error
	Update(contact *Contact) error
	Delete(contact *Contact) error
}

type ContactService struct {
	repo ContactRepo
}

func NewContactService(cr ContactRepo) *ContactService {
	return &ContactService{
		repo: cr,
	}
}
