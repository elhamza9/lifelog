package adding

import (
	"regexp"
	"strings"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// NewTag creates the new tag and calls the service repository to store it.
// It does the following checks:
//   - checks if the given tag name is valid ( length, characters )
//   - transforms name to lowercase
//   - checks repo for tag with same name ( duplicate tags are not allowed )
func (srv Service) NewTag(name string) (domain.Tag, error) {
	// Check tag name is valid
	nameTooShort := len(name) < domain.TagNameMinLength
	nameTooLong := len(name) > domain.TagNameMaxLength
	if nameTooShort || nameTooLong {
		return domain.Tag{}, domain.ErrTagNameLen
	}
	if match, _ := regexp.Match(domain.TagNameValidCharacters, []byte(name)); !match {
		return domain.Tag{}, domain.ErrTagNameInvalidCharacters
	}
	name = strings.ToLower(name)
	// Check tag name is not duplicate
	if t, err := srv.repo.FindTagByName(name); err != nil {
		return domain.Tag{}, err
	} else if len(t.Name) > 0 {
		return domain.Tag{}, domain.ErrTagNameDuplicate
	}
	// Create & Store
	t := domain.Tag{Name: name}
	return srv.repo.AddTag(t)
}
