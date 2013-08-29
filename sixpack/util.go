package sixpack

import (
	"github.com/nu7hatch/gouuid"
)

func GenerateClientID() ([]byte, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return []byte(id.String()), err
}
