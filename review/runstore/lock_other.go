//go:build !unix

package runstore

func holdingTheRecord(_ string, do func() error) error {
	return do()
}
