package domain

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type PermissionRepository interface {
	GetAllForUser(userID int64) (Permissions, error)
	AddForUser(userID int64, codes ...string) error
}
