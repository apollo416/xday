package crops

type CropsDatastore interface {
	Get(id string) (Crop, error)
	Put(c Crop) error
	Update(c Crop) error
}

type CropsService struct {
	datastore CropsDatastore
}

func NewCropsService(datastore CropsDatastore) *CropsService {
	return &CropsService{
		datastore: datastore,
	}
}

func (s *CropsService) Add(c Crop) error {
	return s.datastore.Put(c)
}
