package urlstorage

type Url struct {
	Short string
	Long  string
	ID    string
}

var mockStorage = []Url{}

func NewUrl(long string, short string) Url {
	return Url{
		Short: short,
		Long:  long,
		ID:    short[len(short)-8:],
	}
}

//AddToStorage: adds url to mock database
func (u Url) AddToStorage() {
	mockStorage = append(mockStorage, u)
}

//GetLongByID: returns long version from id
func GetByID(id string) *Url {
	for _, u := range mockStorage {
		if u.ID == id {
			return &u
		}
	}
	return nil
}

//GetByLong: retruns short version by long version
func GetByLong(long string) *Url {
	for _, u := range mockStorage {
		if u.Long == long {
			return &u
		}
	}
	return nil
}
