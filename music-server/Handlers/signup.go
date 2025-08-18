package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	initializers "github.com/franzego/music-server/Initializers"
	models "github.com/franzego/music-server/Models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// The signup function
// It receives a body from the client which it decodes from json to structs(model.user)
// It takes care of the error if it shows up
// it then hashes the password with bcrypt
// it then creates the user in the avian postgres and handles subsequent errors
// finally responds with httpstatusok
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	//get email and password
	/*type Form struct {
		Email    string `json:"email"`
		Password string
	}*/
	var body models.User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create the user
	user := models.User{Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	//respond

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.SignupResponse{
		Email: body.Email,
		Msg:   "Signup Successful",
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//var body models.User

	var body models.User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// compare the user email to make sure it doesnt exist
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		http.Error(w, "Invalid Email or Password", http.StatusBadGateway)
		return
	}
	//compare the hash
	err := bcrypt.CompareHashAndPassword([]byte(body.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid Email or Password", http.StatusBadRequest)
		return
	}
	//jwt is generated here
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.ID,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret in our .env file
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	//cookies are generated here
	cookie := http.Cookie{
		Name:     "user_session",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("cookie set!"))
}
