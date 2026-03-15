package update

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/sigstore/sigstore-go/pkg/bundle"
	"github.com/sigstore/sigstore-go/pkg/root"
	"github.com/sigstore/sigstore-go/pkg/tuf"
	"github.com/sigstore/sigstore-go/pkg/verify"
)

const (
	sigstoreBundleName       = "checksums.txt.sigstore.json"
	sigstoreOIDCIssuer       = "https://token.actions.githubusercontent.com"
	sigstoreWorkflowIdentity = `^https://github.com/chill-institute/cli/\.github/workflows/release\.yml@refs/tags/v.+$`
)

func FindSigstoreBundleAsset(release Release) (ReleaseAsset, error) {
	for _, asset := range release.Assets {
		if asset.Name == sigstoreBundleName {
			return asset, nil
		}
	}
	return ReleaseAsset{}, fmt.Errorf("release asset %q not found", sigstoreBundleName)
}

func VerifySignedChecksumsBundle(_ context.Context, checksums []byte, sigstoreBundle []byte) error {
	tufClient, err := newTUFClient()
	if err != nil {
		return fmt.Errorf("create sigstore tuf client: %w", err)
	}

	trustedMaterial, err := root.GetTrustedRoot(tufClient)
	if err != nil {
		return fmt.Errorf("load sigstore trusted root: %w", err)
	}

	verifier, err := verify.NewVerifier(
		root.TrustedMaterialCollection{trustedMaterial},
		verify.WithSignedCertificateTimestamps(1),
		verify.WithTransparencyLog(1),
		verify.WithObserverTimestamps(1),
	)
	if err != nil {
		return fmt.Errorf("create sigstore verifier: %w", err)
	}

	certificateIdentity, err := verify.NewShortCertificateIdentity(
		sigstoreOIDCIssuer,
		"",
		"",
		sigstoreWorkflowIdentity,
	)
	if err != nil {
		return fmt.Errorf("create certificate identity policy: %w", err)
	}

	loadedBundle := new(bundle.Bundle)
	if err := loadedBundle.UnmarshalJSON(sigstoreBundle); err != nil {
		return fmt.Errorf("decode sigstore bundle: %w", err)
	}

	digest := sha256.Sum256(checksums)
	if _, err := verifier.Verify(
		loadedBundle,
		verify.NewPolicy(
			verify.WithArtifactDigest("sha256", digest[:]),
			verify.WithCertificateIdentity(certificateIdentity),
		),
	); err != nil {
		return fmt.Errorf("verify signed checksums bundle: %w", err)
	}

	return nil
}

func newTUFClient() (*tuf.Client, error) {
	client, err := tuf.DefaultClient()
	if err != nil {
		return nil, err
	}
	return client, nil
}
