package galleryapi

import (
	"net/http"

	"github.com/Aeroxee/gallery-api/auth"
	"github.com/Aeroxee/gallery-api/models"
)

func getUserContext(r *http.Request) (models.User, error) {
	claims := r.Context().Value(&auth.UserAuth{}).(auth.Claims)
	user, err := models.GetUserByID(claims.Credential.UserID)
	return user, err
}
