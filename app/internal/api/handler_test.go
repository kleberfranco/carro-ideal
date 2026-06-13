package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"carro-ideal/app/models"
	"carro-ideal/app/repository"
	"carro-ideal/app/service"
)

// --- shared fakes ---

type apiSessionRepository struct {
	tokenHash string
	userID    int64
	expiresAt time.Time
}

func (r *apiSessionRepository) Create(_ context.Context, hash string, id int64, exp time.Time) error {
	r.tokenHash = hash
	r.userID = id
	r.expiresAt = exp
	return nil
}

func (r *apiSessionRepository) GetUserID(_ context.Context, hash string, now time.Time) (int64, error) {
	if hash != r.tokenHash || !now.Before(r.expiresAt) {
		return 0, repository.ErrSessionNotFound
	}
	return r.userID, nil
}

func (r *apiSessionRepository) Delete(_ context.Context, _ string) error { return nil }

func (r *apiSessionRepository) DeleteExpired(_ context.Context, _ time.Time) error { return nil }

type apiUserRepository struct {
	user   *models.User
	exists bool
	err    error
}

func (r *apiUserRepository) ExistsByEmail(_ context.Context, _ string) (bool, error) {
	return r.exists, r.err
}
func (r *apiUserRepository) Create(_ context.Context, u *models.User) error {
	u.ID = 1
	r.user = u
	return r.err
}
func (r *apiUserRepository) GetByEmail(_ context.Context, _ string) (*models.User, error) {
	if r.user == nil {
		return nil, errors.New("not found")
	}
	return r.user, nil
}
func (r *apiUserRepository) GetByID(_ context.Context, _ int64) (*models.User, error) {
	if r.user == nil {
		return nil, errors.New("not found")
	}
	return r.user, nil
}
func (r *apiUserRepository) Update(_ context.Context, u *models.User) error { return nil }
func (r *apiUserRepository) Deactivate(_ context.Context, _ int64) error    { return nil }

type apiQuestionRepository struct {
	questions []models.Question
}

func (r *apiQuestionRepository) GetActive(_ context.Context) ([]models.Question, error) {
	return r.questions, nil
}

type apiVehicleRepository struct {
	vehicles []models.Vehicle
}

func (r *apiVehicleRepository) GetActive(_ context.Context, _ int64) ([]models.Vehicle, error) {
	return r.vehicles, nil
}

func (r *apiVehicleRepository) GetByID(_ context.Context, id int64) (*models.Vehicle, error) {
	for _, v := range r.vehicles {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, repository.ErrVehicleNotFound
}

type apiRecommendationRepository struct {
	stored *models.Recommendation
}

func (r *apiRecommendationRepository) Create(_ context.Context, rec *models.Recommendation, _ []models.SubmittedAnswer) error {
	rec.ID = 1
	r.stored = rec
	return nil
}

func (r *apiRecommendationRepository) GetByUser(_ context.Context, _ int64, _, _ int) ([]models.Recommendation, int, error) {
	if r.stored == nil {
		return nil, 0, nil
	}
	return []models.Recommendation{*r.stored}, 1, nil
}

func (r *apiRecommendationRepository) GetByID(_ context.Context, id, _ int64) (*models.Recommendation, error) {
	if r.stored != nil && r.stored.ID == id {
		return r.stored, nil
	}
	return nil, repository.ErrRecommendationNotFound
}

// --- builder helpers ---

func buildHandler(userRepo *apiUserRepository, sessRepo *apiSessionRepository, questionRepo *apiQuestionRepository, vehicleRepo *apiVehicleRepository, recRepo *apiRecommendationRepository) *Handler {
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(sessRepo)
	qSvc := service.NewQuestionnaireService(questionRepo)
	vSvc := service.NewVehicleService(vehicleRepo)
	recSvc := service.NewRecommendationService(qSvc, vSvc, recRepo)
	return NewHandler(userSvc, authSvc, qSvc, recSvc, vSvc, false)
}

func postJSON(t *testing.T, handler http.Handler, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	buf, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func getWithSession(t *testing.T, handler http.Handler, path, token string, userID int64) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.AddCookie(&http.Cookie{Name: "carro_session", Value: token, Expires: time.Now().Add(time.Hour)})
	ctx := context.WithValue(req.Context(), userIDKey, userID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req.WithContext(ctx))
	return rec
}

// --- T142: API contract tests ---

func TestRegisterHandlerValid(t *testing.T) {
	userRepo := &apiUserRepository{}
	sessRepo := &apiSessionRepository{}
	h := buildHandler(userRepo, sessRepo, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})

	rec := postJSON(t, http.HandlerFunc(h.Register), "/api/auth/register", map[string]string{
		"name": "Maria Silva", "email": "maria@example.com",
		"password": "senha123", "confirm_password": "senha123",
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("Register() status = %d, want %d; body = %s", rec.Code, http.StatusCreated, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), "user") {
		t.Fatalf("Register() body should contain user: %s", rec.Body)
	}
}

func TestRegisterHandlerInvalidEmail(t *testing.T) {
	h := buildHandler(&apiUserRepository{}, &apiSessionRepository{}, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})
	rec := postJSON(t, http.HandlerFunc(h.Register), "/api/auth/register", map[string]string{
		"name": "Maria", "email": "not-an-email",
		"password": "senha123", "confirm_password": "senha123",
	})
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Register() status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

func TestRegisterHandlerDuplicateEmail(t *testing.T) {
	h := buildHandler(&apiUserRepository{exists: true}, &apiSessionRepository{}, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})
	rec := postJSON(t, http.HandlerFunc(h.Register), "/api/auth/register", map[string]string{
		"name": "Maria", "email": "maria@example.com",
		"password": "senha123", "confirm_password": "senha123",
	})
	if rec.Code != http.StatusConflict {
		t.Fatalf("Register() status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestRegisterHandlerWeakPassword(t *testing.T) {
	h := buildHandler(&apiUserRepository{}, &apiSessionRepository{}, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})
	rec := postJSON(t, http.HandlerFunc(h.Register), "/api/auth/register", map[string]string{
		"name": "Maria", "email": "maria@example.com",
		"password": "123", "confirm_password": "123",
	})
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Register() status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

func TestLoginHandlerValid(t *testing.T) {
	userRepo := &apiUserRepository{}
	sessRepo := &apiSessionRepository{}
	h := buildHandler(userRepo, sessRepo, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})

	_ = postJSON(t, http.HandlerFunc(h.Register), "/api/auth/register", map[string]string{
		"name": "Maria", "email": "maria@example.com",
		"password": "senha123", "confirm_password": "senha123",
	})

	rec := postJSON(t, http.HandlerFunc(h.Login), "/api/auth/login", map[string]string{
		"email": "maria@example.com", "password": "senha123",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("Login() status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body)
	}
}

func TestLoginHandlerWrongPassword(t *testing.T) {
	userRepo := &apiUserRepository{}
	sessRepo := &apiSessionRepository{}
	h := buildHandler(userRepo, sessRepo, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})

	_ = postJSON(t, http.HandlerFunc(h.Register), "/api/auth/register", map[string]string{
		"name": "Maria", "email": "maria@example.com",
		"password": "senha123", "confirm_password": "senha123",
	})

	rec := postJSON(t, http.HandlerFunc(h.Login), "/api/auth/login", map[string]string{
		"email": "maria@example.com", "password": "wrongpass",
	})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Login() status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestLoginHandlerEmptyCredentials(t *testing.T) {
	h := buildHandler(&apiUserRepository{}, &apiSessionRepository{}, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})
	rec := postJSON(t, http.HandlerFunc(h.Login), "/api/auth/login", map[string]string{})
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Login() status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

func TestLogoutHandler(t *testing.T) {
	sessRepo := &apiSessionRepository{}
	h := buildHandler(&apiUserRepository{}, sessRepo, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})
	rec := postJSON(t, http.HandlerFunc(h.Logout), "/api/auth/logout", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("Logout() status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestMeHandlerUnauthenticated(t *testing.T) {
	h := buildHandler(&apiUserRepository{}, &apiSessionRepository{}, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	h.Me(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Me() status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestMeHandlerAuthenticated(t *testing.T) {
	userRepo := &apiUserRepository{}
	sessRepo := &apiSessionRepository{}
	h := buildHandler(userRepo, sessRepo, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})

	regRec := postJSON(t, http.HandlerFunc(h.Register), "/api/auth/register", map[string]string{
		"name": "Maria", "email": "maria@example.com",
		"password": "senha123", "confirm_password": "senha123",
	})
	if regRec.Code != http.StatusCreated {
		t.Fatalf("Register() failed: %s", regRec.Body)
	}

	loginRec := postJSON(t, http.HandlerFunc(h.Login), "/api/auth/login", map[string]string{
		"email": "maria@example.com", "password": "senha123",
	})
	token := ""
	for _, c := range loginRec.Result().Cookies() {
		if c.Name == "carro_session" {
			token = c.Value
		}
	}
	if token == "" {
		t.Fatal("Login() did not set session cookie")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: "carro_session", Value: token, Expires: time.Now().Add(time.Hour)})
	meRec := httptest.NewRecorder()
	h.Me(meRec, req)
	if meRec.Code != http.StatusOK {
		t.Fatalf("Me() status = %d, want %d; body = %s", meRec.Code, http.StatusOK, meRec.Body)
	}
}

func TestQuestionsHandler(t *testing.T) {
	questions := []models.Question{
		{ID: 1, Text: "Qual o orçamento?", Type: "SINGLE_CHOICE", Weight: 1, Options: []models.AnswerOption{
			{ID: 10, Text: "Até R$60k", ScoreProfile: map[string]float64{"budget_low": 1}},
		}},
	}
	questionRepo := &apiQuestionRepository{questions: questions}
	h := buildHandler(&apiUserRepository{}, &apiSessionRepository{}, questionRepo, &apiVehicleRepository{}, &apiRecommendationRepository{})

	req := httptest.NewRequest(http.MethodGet, "/api/questions", nil)
	rec := httptest.NewRecorder()
	h.Questions(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Questions() status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "orçamento") {
		t.Fatalf("Questions() body should contain question: %s", rec.Body)
	}
}

func TestVehiclesHandler(t *testing.T) {
	vehicles := []models.Vehicle{{ID: 1, Brand: "Toyota", MatchProfile: map[string]float64{"urban": 1}}}
	h := buildHandler(&apiUserRepository{}, &apiSessionRepository{}, &apiQuestionRepository{}, &apiVehicleRepository{vehicles: vehicles}, &apiRecommendationRepository{})

	req := httptest.NewRequest(http.MethodGet, "/api/vehicles", nil)
	rec := httptest.NewRecorder()
	h.Vehicles(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Vehicles() status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "Toyota") {
		t.Fatalf("Vehicles() body should contain vehicle: %s", rec.Body)
	}
}

func TestGenerateRecommendationsHandler(t *testing.T) {
	questions := []models.Question{
		{ID: 1, Weight: 1, Options: []models.AnswerOption{
			{ID: 10, QuestionID: 1, ScoreProfile: map[string]float64{"urban": 1}},
		}},
	}
	vehicles := []models.Vehicle{{ID: 1, Brand: "Toyota", MatchProfile: map[string]float64{"urban": 1}}}
	h := buildHandler(
		&apiUserRepository{},
		&apiSessionRepository{},
		&apiQuestionRepository{questions: questions},
		&apiVehicleRepository{vehicles: vehicles},
		&apiRecommendationRepository{},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/recommendations/generate", bytes.NewBufferString(`{"answers":[{"question_id":1,"answer_option_id":10}]}`))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), userIDKey, int64(1))
	rec := httptest.NewRecorder()
	h.GenerateRecommendations(rec, req.WithContext(ctx))

	if rec.Code != http.StatusCreated {
		t.Fatalf("GenerateRecommendations() status = %d, want %d; body = %s", rec.Code, http.StatusCreated, rec.Body)
	}
}

func TestGenerateRecommendationsHandlerUnauthenticated(t *testing.T) {
	h := buildHandler(&apiUserRepository{}, &apiSessionRepository{}, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})
	req := httptest.NewRequest(http.MethodPost, "/api/recommendations/generate", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()
	h.GenerateRecommendations(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("GenerateRecommendations() status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRecommendationHistoryHandler(t *testing.T) {
	h := buildHandler(&apiUserRepository{}, &apiSessionRepository{}, &apiQuestionRepository{}, &apiVehicleRepository{}, &apiRecommendationRepository{})
	rec := getWithSession(t, http.HandlerFunc(h.RecommendationHistory), "/api/recommendations", "", 1)
	if rec.Code != http.StatusOK {
		t.Fatalf("RecommendationHistory() status = %d, want %d", rec.Code, http.StatusOK)
	}
}

// --- T143: E2E user flow test ---

func TestE2EUserRecommendationFlow(t *testing.T) {
	userRepo := &apiUserRepository{}
	sessRepo := &apiSessionRepository{}
	questions := []models.Question{
		{ID: 1, Weight: 1, Options: []models.AnswerOption{
			{ID: 10, QuestionID: 1, ScoreProfile: map[string]float64{"urban": 1}},
		}},
	}
	vehicles := []models.Vehicle{
		{ID: 1, Brand: "Toyota", Model: "Corolla", MatchProfile: map[string]float64{"urban": 1}},
		{ID: 2, Brand: "Jeep", Model: "Compass", MatchProfile: map[string]float64{"offroad": 1}},
	}
	recRepo := &apiRecommendationRepository{}
	h := buildHandler(userRepo, sessRepo, &apiQuestionRepository{questions: questions}, &apiVehicleRepository{vehicles: vehicles}, recRepo)

	// Step 1: Register
	regRec := postJSON(t, http.HandlerFunc(h.Register), "/api/auth/register", map[string]string{
		"name": "Maria Silva", "email": "maria@example.com",
		"password": "senha12345", "confirm_password": "senha12345",
	})
	if regRec.Code != http.StatusCreated {
		t.Fatalf("E2E Register() failed: %s", regRec.Body)
	}

	// Step 2: Get session token from cookie
	token := ""
	for _, c := range regRec.Result().Cookies() {
		if c.Name == "carro_session" {
			token = c.Value
		}
	}
	if token == "" {
		t.Fatal("E2E: no session cookie after registration")
	}

	// Step 3: Verify session via /me
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: "carro_session", Value: token, Expires: time.Now().Add(time.Hour)})
	meRec := httptest.NewRecorder()
	h.Me(meRec, req)
	if meRec.Code != http.StatusOK {
		t.Fatalf("E2E /me status = %d; body = %s", meRec.Code, meRec.Body)
	}

	// Step 4: Get questions
	qReq := httptest.NewRequest(http.MethodGet, "/api/questions", nil)
	qRec := httptest.NewRecorder()
	h.Questions(qRec, qReq)
	if qRec.Code != http.StatusOK {
		t.Fatalf("E2E Questions() status = %d", qRec.Code)
	}

	// Step 5: Generate recommendations
	genReq := httptest.NewRequest(http.MethodPost, "/api/recommendations/generate",
		bytes.NewBufferString(`{"answers":[{"question_id":1,"answer_option_id":10}]}`))
	genReq.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(genReq.Context(), userIDKey, int64(1))
	genRec := httptest.NewRecorder()
	h.GenerateRecommendations(genRec, genReq.WithContext(ctx))
	if genRec.Code != http.StatusCreated {
		t.Fatalf("E2E GenerateRecommendations() status = %d; body = %s", genRec.Code, genRec.Body)
	}
	if !strings.Contains(genRec.Body.String(), "Toyota") {
		t.Fatalf("E2E: recommendations should include Toyota Corolla (best urban match); got: %s", genRec.Body)
	}

	// Step 6: View recommendation history
	histReq := httptest.NewRequest(http.MethodGet, "/api/recommendations", nil)
	histCtx := context.WithValue(histReq.Context(), userIDKey, int64(1))
	histRec := httptest.NewRecorder()
	h.RecommendationHistory(histRec, histReq.WithContext(histCtx))
	if histRec.Code != http.StatusOK {
		t.Fatalf("E2E RecommendationHistory() status = %d", histRec.Code)
	}

	// Step 7: Logout
	logoutRec := postJSON(t, http.HandlerFunc(h.Logout), "/api/auth/logout", nil)
	if logoutRec.Code != http.StatusOK {
		t.Fatalf("E2E Logout() status = %d", logoutRec.Code)
	}
}

// --- T144: Admin workflow test ---

func TestE2EAdminCannotAccessWithUserRole(t *testing.T) {
	userRepo := &apiUserRepository{
		user: &models.User{ID: 1, Role: "user", Active: true},
	}
	sessRepo := &apiSessionRepository{}
	authSvc := service.NewAuthService(sessRepo)
	userSvc := service.NewUserService(userRepo)

	token, expiresAt, _ := authSvc.CreateSession(context.Background(), 1)
	sessRepo.expiresAt = expiresAt

	adminMiddleware := RequireAdmin(userSvc, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/dashboard", nil)
	req.AddCookie(&http.Cookie{Name: "carro_session", Value: token, Expires: expiresAt})
	ctx := context.WithValue(req.Context(), userIDKey, int64(1))
	rec := httptest.NewRecorder()
	adminMiddleware.ServeHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("Admin route should return 403 for user role, got %d", rec.Code)
	}
}

func TestE2EAdminCanAccessWithAdminRole(t *testing.T) {
	userRepo := &apiUserRepository{
		user: &models.User{ID: 1, Role: "admin", Active: true},
	}
	sessRepo := &apiSessionRepository{}
	authSvc := service.NewAuthService(sessRepo)
	userSvc := service.NewUserService(userRepo)

	token, expiresAt, _ := authSvc.CreateSession(context.Background(), 1)
	sessRepo.expiresAt = expiresAt

	adminMiddleware := RequireAdmin(userSvc, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/dashboard", nil)
	req.AddCookie(&http.Cookie{Name: "carro_session", Value: token, Expires: expiresAt})
	ctx := context.WithValue(req.Context(), userIDKey, int64(1))
	rec := httptest.NewRecorder()
	adminMiddleware.ServeHTTP(rec, req.WithContext(ctx))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("Admin route should return 204 for admin role, got %d", rec.Code)
	}
}
