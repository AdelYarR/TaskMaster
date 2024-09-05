package apiserver

import (
	"TaskMaster/configs/config"
	"TaskMaster/pkg/handler"
	"TaskMaster/pkg/repository"
	"TaskMaster/pkg/service"
	"strings"

	"log/slog"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type APIServer struct {
	config  *config.Config
	router  *mux.Router
	logger  *slog.Logger
	handler *handler.Handler
}

func New(cfg *config.Config) *APIServer {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	}))

	repo := repository.NewRepo(cfg.Store.DBUrl, logger)
	serv := service.NewService(repo, logger)
	hand := handler.NewHandler(serv, logger)

	return &APIServer{
		config:  cfg,
		router:  mux.NewRouter(),
		logger:  logger,
		handler: hand,
	}
}

func (s *APIServer) Start() error {
	s.Handle()

	if err := http.ListenAndServe(s.config.BindAddr, s.router); err != nil {
		s.logger.Error("Failed to start server")
	}

	return nil
}

func (s *APIServer) Handle() {
	s.router.HandleFunc("/signup", s.handler.SignUp())
	s.router.HandleFunc("/signin", s.handler.SignIn())

	{
		s.router.HandleFunc("/tasks", s.JwtParse(s.handler.Tasks))
	}
}

func (s *APIServer) JwtParse(next func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			s.logger.Error("Invalid token format")
		}

		tokenString := bearerToken[1]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func (token *jwt.Token) (interface{}, error) {
			return []byte("ekjgerj5et45F@wEDFRge$*riwe934urHsajd*W!@ffklgjmklVWbdsklcnkJBJKFSWGVWEF"), nil
		})
		if err != nil {
			s.logger.Error("Failed to parse")
		}

		if !token.Valid {
			s.logger.Error("Token is not valid")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Failed to parse", http.StatusUnauthorized)
			return
		}

		next(w, r, int(userID))
	}
}
