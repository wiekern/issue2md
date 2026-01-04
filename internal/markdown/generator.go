package markdown

// Generator defines the markdown generator interface.
type Generator interface {
	// Generate generates complete markdown content from Issue/PR data.
	Generate(data interface{}) (string, error)

	// GenerateFrontmatter generates only the YAML frontmatter.
	GenerateFrontmatter(data interface{}) string

	// SanitizeFilename sanitizes a title to a valid filename.
	SanitizeFilename(title string) string

	// DefaultFilename returns the default filename for the given data.
	DefaultFilename(data interface{}) string
}

type generator struct {
	includeComments bool
	includeMeta     bool
}

// NewGenerator creates a new Generator instance.
func NewGenerator(includeComments, includeMeta bool) Generator {
	return &generator{
		includeComments: includeComments,
		includeMeta:     includeMeta,
	}
}

// Generate generates complete markdown content from Issue/PR data.
func (g *generator) Generate(data interface{}) (string, error) {
	return "", nil
}

// GenerateFrontmatter generates only the YAML frontmatter.
func (g *generator) GenerateFrontmatter(data interface{}) string {
	return ""
}

// SanitizeFilename sanitizes a title to a valid filename.
func (g *generator) SanitizeFilename(title string) string {
	return ""
}

// DefaultFilename returns the default filename for the given data.
func (g *generator) DefaultFilename(data interface{}) string {
	return ""
}
