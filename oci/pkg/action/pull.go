package action

import (
	"errors"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/tektoncd/experimental/oci/pkg/oci"
)

// Pull will perform the `pull` action by retrieving a specific named Tekton resource from the specified OCI image.
func Pull(r string, kind string, n string) error {
	// Validate the parameters.
	if r == "" || kind == "" || n == "" {
		return errors.New("must specify an image reference, kind, and resource name")
	}

	ref, err := name.ParseReference(r)
	if err != nil {
		return err
	}

	// TODO: When this is moved into the Tekton controller, authorize this
	// pull as a Service Account in the cluster, and don't rely on the
	// contents of ~/.docker/config.json (which won't exist).
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}

	contents, err := oci.Extract(img, kind, n)
	if err != nil {
		return err
	}

	fmt.Print(string(contents))
	return nil
}
