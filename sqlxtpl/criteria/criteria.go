package criteria

import (
	"fmt"
	"strings"
)

// sprintf is a helper method for fmt.Sprintf
func sprintf(s string, args ...interface{}) string {
	return fmt.Sprintf(s, args...)
}

// makeNamedVars is a helper method to make NamedBindVars map
func makeNamedVars() NamedBindVars {
	return make(NamedBindVars)
}

// NamedBindVars represents the query named vars
type NamedBindVars map[string]interface{}

// BindVars represents the query vars
type BindVars []interface{}

// Criteria represents the criteria used on sql where
type Criteria struct {
	query []string

	namedBindVars NamedBindVars
}

// New creates a Criteria instance
func New() Criteria {
	return Criteria{[]string{}, make(NamedBindVars)}
}

// Handler represents a struct resolution like an operator
type Handler func() (string, NamedBindVars)

// BothLike creates the like operator with value "%...%"
func (c Criteria) BothLike(f string, v interface{}) Handler {
	return c.like(f, v, "%%%s%%")
}

// Like creates the ANSI like operator to search for a specified pattern in a column
func (c Criteria) like(field string, v interface{}, format string) Handler {
	return func() (string, NamedBindVars) {
		query := "%s LIKE :%s"
		namedVars := makeNamedVars()

		namedVars[field] = fmt.Sprintf(format, v)

		return sprintf(query, field, field), namedVars
	}
}

// Add adds a criteria handler
func (c *Criteria) Add(h Handler) {
	s, namedVars := h()

	c.query = append(c.query, s)

	for key, v := range namedVars {
		c.namedBindVars[key] = v
	}
}

// BindVars gets the criteria bind vars used as query args
func (c Criteria) BindVars() map[string]interface{} {
	return c.namedBindVars
}

// NamedBindVars gets the criteria named bind vars used as query args
func (c Criteria) NamedBindVars() map[string]interface{} {
	return c.namedBindVars
}

// MergeSQL merges the criteria query with s (another query)
func (c Criteria) MergeSQL(s string) string {
	return s + " " + c.ToSQL()
}

// IsEmpty returns if the criteria is empty
func (c Criteria) IsEmpty() bool {
	return len(c.query) == 0
}

// ToSQL exports the sql where
func (c Criteria) ToSQL() string {
	if len(c.query) == 0 {
		return ""
	}

	return "WHERE " + strings.Join(c.query, " AND ")
}

// PostgresCriteria represents the Postgres structure
type PostgresCriteria struct {
	Criteria
}

// NewPostgresCriteria creates a PostgresCriteria instance
func NewPostgresCriteria() PostgresCriteria {
	return PostgresCriteria{New()}
}

// BothILike creates the like operator with value "%...%"
func (c PostgresCriteria) ILikeBoth(f string, v interface{}) Handler {
	return c.like(f, v, "%%%s%%")
}

// Like creates the ANSI like operator to search for a specified pattern in a column
func (c PostgresCriteria) like(field string, v interface{}, format string) Handler {
	return func() (string, NamedBindVars) {
		query := "%s ILIKE :%s"
		namedVars := makeNamedVars()

		namedVars[field] = fmt.Sprintf(format, v)

		return sprintf(query, field, field), namedVars
	}
}
