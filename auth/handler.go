package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mmiftahrzki/go-rest-api/response"
)

func CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	var payload AuthCreateModel
	response := response.New()

	json_decoder := json.NewDecoder(r.Body)
	err := json_decoder.Decode(&payload)
	if err != nil {
		log.Println(err)

		response.Message = http.StatusText(http.StatusBadRequest)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.ToJson())

		return
	}

	token, err := generateToken(payload)
	if err != nil {
		log.Println(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response.ToJson())

		return
	}

	response.Data["token"] = token
	response.Message = "berhasil generate token"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response.ToJson())
}

func generateToken(payload AuthCreateModel) (string, error) {
	registerd_claims := jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute))}
	claims := AuthClaimModel{
		Email:            payload.Email,
		RegisteredClaims: registerd_claims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_string, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return signed_string, err
	}

	return signed_string, nil
}
