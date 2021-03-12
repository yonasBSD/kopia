package auth_test

import (
	"context"
	"testing"

	"github.com/kopia/kopia/internal/auth"
	"github.com/kopia/kopia/internal/repotesting"
	"github.com/kopia/kopia/internal/testlogging"
	"github.com/kopia/kopia/internal/user"
	"github.com/kopia/kopia/repo"
)

func TestRepositoryAuthenticator(t *testing.T) {
	a := auth.AuthenticateRepositoryUsers()
	ctx := testlogging.Context(t)

	var env repotesting.Environment
	defer env.Setup(t).Close(ctx, t)

	must(t, repo.WriteSession(ctx, env.Repository, repo.WriteSessionOptions{},
		func(w repo.RepositoryWriter) error {
			p := &user.Profile{
				Username: "user1@host1",
			}

			p.SetPassword("password1")

			return user.SetUserProfile(ctx, w, p)
		}))

	verifyRepoAuthenticator(ctx, t, a, env.Repository, "user1@host1", "password1", true)
	verifyRepoAuthenticator(ctx, t, a, env.Repository, "user1@host1", "password2", false)
	verifyRepoAuthenticator(ctx, t, a, env.Repository, "user1@host1", "password11", false)
	verifyRepoAuthenticator(ctx, t, a, env.Repository, "user1@host1a", "password1", false)
	verifyRepoAuthenticator(ctx, t, a, env.Repository, "user1@host1a", "password1a", false)
}

func verifyRepoAuthenticator(ctx context.Context, t *testing.T, a auth.Authenticator, r repo.Repository, username, password string, want bool) {
	t.Helper()

	if got := a(ctx, r, username, password); got != want {
		t.Errorf("invalid authenticator result for %v/%v: %v, want %v", username, password, got, want)
	}
}

func must(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}