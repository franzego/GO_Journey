package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	initializers "github.com/franzego/music-server/Initializers"
	models "github.com/franzego/music-server/Models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lpernett/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Interface for different login services
type Loginservice interface {
	Login(w http.ResponseWriter, r *http.Request) error
}

// Structs for different login services
type TwitterLogin struct {
	Twitter string `json:"twitter"`
}
type GithubLogin struct {
	Github string `json:"github"`
}
type GoogleLogin struct {
	Google string `json:"google"`
}
type FacebookLogin struct {
	Facebook string `json:"facebook"`
}

// Service struct to group all login services

type Service struct {
	TwitterLogin
	GithubLogin
	GoogleLogin
	FacebookLogin
}

// URLs from .env
var t = os.Getenv("TWITTER_AUTH_URL")
var g = os.Getenv("GITHUB_AUTH_URL")
var o = os.Getenv("GOOGLE_AUTH_URL")
var f = os.Getenv("FACEBOOK_AUTH_URL")

// Implement Login method for each service
func (tl TwitterLogin) Login(w http.ResponseWriter, r *http.Request) error {
	// Handle Twitter OAuth callback and user authentication here

	http.Redirect(w, r, t, http.StatusTemporaryRedirect)

	log.Println("Logged in with Twitter!")
	return nil
}

func (gl GithubLogin) Login(w http.ResponseWriter, r *http.Request) error {
	// Handle Github OAuth callback and user authentication here

	http.Redirect(w, r, g, http.StatusTemporaryRedirect)

	log.Println("Logged in with Github!")
	return nil
}

func (gol GoogleLogin) Login(w http.ResponseWriter, r *http.Request) error {
	// Handle Google OAuth callback and user authentication here

	http.Redirect(w, r, o, http.StatusTemporaryRedirect)

	log.Println("Logged in with Google!")
	return nil
}
func (fl FacebookLogin) Login(w http.ResponseWriter, r *http.Request) error {

	http.Redirect(w, r, f, http.StatusTemporaryRedirect)
	log.Println("Logged in with Facebook!")
	return nil
}

// Map of providers to their respective login services
var providers = map[string]Loginservice{
	"twitter":  TwitterLogin{Twitter: t},
	"github":   GithubLogin{Github: g},
	"google":   GoogleLogin{Google: o},
	"facebook": FacebookLogin{Facebook: f},
}

// oauth login handler
func OauthLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract provider from query parameters
	q := r.URL.Query()
	provider := q.Get("login")

	if service, ok := providers[provider]; ok {
		service.Login(w, r)
		return
	}
	http.Error(w, "Unsupported", http.StatusBadRequest)

}

/*
	var l = &Service{
		TwitterLogin: TwitterLogin{Twitter: t},
		GithubLogin:  GithubLogin{Github: g},
		GoogleLogin:  GoogleLogin{Google: o},
	}
*/

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	//get email and password
	/*type Form struct {
		Email    string `json:"email"`
		Password string
	}*/
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to system environment")
	}
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

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to system environment")
	}
	var body models.User // we only need email and password from the body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	////basic validation

	/*if body.Email == "" || body.Password == "" {
		http.Error(w, "Must input credentials", http.StatusUnauthorized)
		return
	}*/

	// compare the user email to make sure it doesnt exist
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		http.Error(w, "Invalid Email", http.StatusUnauthorized)
		return
	}
	//compare the hash

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}

	//jwt is generated here
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Email,
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
	//w.Write([]byte("cookie set!"))
	//fmt.Println(user.Password)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"msg":   "Login successful",
		"email": user.Email,
	})
}
