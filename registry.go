package distribution

import (
	"context"

	"github.com/docker/distribution/reference"
)

// Scope defines the set of items that match a namespace.
type Scope interface {
	// Contains returns true if the name belongs to the namespace.
	Contains(name string) bool
}

type fullScope struct{}

func (f fullScope) Contains(string) bool {
	return true
}

// GlobalScope represents the full namespace scope which contains
// all other scopes.
var GlobalScope = Scope(fullScope{})

// Namespace represents a collection of repositories, addressable by name.
// Generally, a namespace is backed by a set of one or more services,
// providing facilities such as registry access, trust, and indexing.
type Namespace interface {
	// Scope describes the names that can be used with this Namespace. The
	// global namespace will have a scope that matches all names. The scope
	// effectively provides an identity for the namespace.
	Scope() Scope

	// Repository should return a reference to the named repository. The
	// registry may or may not have the repository but should always return a
	// reference.

	// 返回一个操作name(比如nicktming/test仓库)的Repository
	Repository(ctx context.Context, name reference.Named) (Repository, error)

	// Repositories fills 'repos' with a lexicographically sorted catalog of repositories
	// up to the size of 'repos' and returns the value 'n' for the number of entries
	// which were filled.  'last' contains an offset in the catalog, and 'err' will be
	// set to io.EOF if there are no more entries to obtain.
	Repositories(ctx context.Context, repos []string, last string) (n int, err error)

	// Blobs returns a blob enumerator to access all blobs
	// 操作blobs
	Blobs() BlobEnumerator

	// BlobStatter returns a BlobStatter to control
	BlobStatter() BlobStatter
}

// RepositoryEnumerator describes an operation to enumerate repositories
type RepositoryEnumerator interface {
	Enumerate(ctx context.Context, ingester func(string) error) error
}

// RepositoryRemover removes given repository
type RepositoryRemover interface {
	Remove(ctx context.Context, name reference.Named) error
}

// ManifestServiceOption is a function argument for Manifest Service methods
type ManifestServiceOption interface {
	Apply(ManifestService) error
}

// WithTag allows a tag to be passed into Put
func WithTag(tag string) ManifestServiceOption {
	return WithTagOption{tag}
}

// WithTagOption holds a tag
type WithTagOption struct{ Tag string }

// Apply conforms to the ManifestServiceOption interface
func (o WithTagOption) Apply(m ManifestService) error {
	// no implementation
	return nil
}

// WithManifestMediaTypes lists the media types the client wishes
// the server to provide.
func WithManifestMediaTypes(mediaTypes []string) ManifestServiceOption {
	return WithManifestMediaTypesOption{mediaTypes}
}

// WithManifestMediaTypesOption holds a list of accepted media types
type WithManifestMediaTypesOption struct{ MediaTypes []string }

// Apply conforms to the ManifestServiceOption interface
func (o WithManifestMediaTypesOption) Apply(m ManifestService) error {
	// no implementation
	return nil
}

// Repository is a named collection of manifests and layers.
// 仓库包括了manifest, blobs 和 tags的操作
// 提供了可以获得操作上述所有类别的service
type Repository interface {
	// Named returns the name of the repository.
	//仓库名字
	Named() reference.Named

	// Manifests returns a reference to this repository's manifest service.
	// with the supplied options applied.

	// 获得操作manifest的service
	Manifests(ctx context.Context, options ...ManifestServiceOption) (ManifestService, error)

	// Blobs returns a reference to this repository's blob service.
	// 获得操作blobs的service
	Blobs(ctx context.Context) BlobStore

	// TODO(stevvooe): The above BlobStore return can probably be relaxed to
	// be a BlobService for use with clients. This will allow such
	// implementations to avoid implementing ServeBlob.

	// Tags returns a reference to this repositories tag service
	// 获得操作tags的service
	Tags(ctx context.Context) TagService
}

// TODO(stevvooe): Must add close methods to all these. May want to change the
// way instances are created to better reflect internal dependency
// relationships.
