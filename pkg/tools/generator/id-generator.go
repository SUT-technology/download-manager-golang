package generator

import "github.com/google/uuid"

func IdGenerator(ids []string) string {
	var id string

	for {
		id = uuid.New().String()
		flag := true
		for _, i := range ids {
			if i == id {
				flag = false
				break
			}
		}
		if flag {
			break
		}
	}

	return id
}
