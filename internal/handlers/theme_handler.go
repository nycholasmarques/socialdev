package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nycholasmarques/socialdev/internal/entity"
	"github.com/nycholasmarques/socialdev/internal/middlewares"
	"github.com/nycholasmarques/socialdev/internal/services"
	"github.com/nycholasmarques/socialdev/utils"
)

type ThemeHandler struct {
	themeService *services.ThemeService
}

type ThemeWithUser struct {
	ThemeID uuid.UUID `json:"theme_id"`
	UserID uuid.UUID `json:"user_id"`
}

func NewThemeHandler(themeService *services.ThemeService) *ThemeHandler {
	return &ThemeHandler{themeService: themeService}
}
func (h *ThemeHandler) CreateThemeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rb := r.Body
	reader, err := io.ReadAll(rb)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to read request body", nil)
		return
	}
	defer rb.Close()

	var themeCreate entity.ThemeRequest
	err = json.Unmarshal(reader, &themeCreate)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	validate = validator.New()
	err = validate.Struct(themeCreate)
	if err != nil {
		errorMessages := entity.GetValidationThemeErrors(err)
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", errorMessages)
		return
	}

	theme, err := h.themeService.CreateTheme(themeCreate.Name, themeCreate.Description, themeCreate.PreviewImage, themeCreate.IsPremium, themeCreate.Status, themeCreate.Active, themeCreate.ExpiresIn, themeCreate.ThemeConfig)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to create theme", nil)
		return
	}

	var themeResponse entity.Theme
	
	themeResponse.ThemeID = theme.ThemeID.String()
	themeResponse.Name = theme.Name.String
	themeResponse.Description = theme.Description.String
	themeResponse.IsPremium = theme.IsPremium.Bool
	themeResponse.PreviewImage = theme.PreviewImage.String
	themeResponse.ThemeConfig = theme.ThemeConfig.RawMessage
	themeResponse.Status = theme.Status.Bool
	themeResponse.Active = theme.Active.Bool
	themeResponse.ExpiresIn = theme.ExpiresIn.Time

	utils.NewMessageResponse(w, http.StatusCreated, "Theme created successfully", themeResponse)
}

func (h *ThemeHandler) ListAllThemesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	themes, err := h.themeService.ListAllThemes()
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to list themes", nil)
		return
	}

	var themesResponse []entity.Theme

	for _, theme := range themes {
		var themeResponse entity.Theme
		
		themeResponse.ThemeID = theme.ThemeID.String()
		themeResponse.Name = theme.Name.String
		themeResponse.Description = theme.Description.String
		themeResponse.IsPremium = theme.IsPremium.Bool
		themeResponse.PreviewImage = theme.PreviewImage.String
		themeResponse.ThemeConfig = theme.ThemeConfig.RawMessage
		themeResponse.Status = theme.Status.Bool
		themeResponse.Active = theme.Active.Bool
		themeResponse.ExpiresIn = theme.ExpiresIn.Time
		
		themesResponse = append(themesResponse, themeResponse)
	}

	utils.NewMessageResponse(w, http.StatusOK, "Themes listed successfully", themesResponse)
}

func (h *ThemeHandler) GetThemeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	themeID := r.PathValue("theme_id")
	if themeID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id in URL", "")
		return
	}
	themeIDUUID, err := uuid.Parse(themeID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid theme ID", nil)
		return
	}

	theme, err := h.themeService.GetThemeByID(themeIDUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, "Theme not found", nil)
		return
	}

	var themeResponse entity.Theme
	
	themeResponse.ThemeID = theme.ThemeID.String()
	themeResponse.Name = theme.Name.String
	themeResponse.Description = theme.Description.String
	themeResponse.IsPremium = theme.IsPremium.Bool
	themeResponse.PreviewImage = theme.PreviewImage.String
	themeResponse.ThemeConfig = theme.ThemeConfig.RawMessage
	themeResponse.Status = theme.Status.Bool
	themeResponse.Active = theme.Active.Bool
	themeResponse.ExpiresIn = theme.ExpiresIn.Time

	utils.NewMessageResponse(w, http.StatusOK, "Theme retrieved successfully", themeResponse)
}

func (h *ThemeHandler) UpdateThemeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	themeID := r.PathValue("theme_id")
	if themeID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id in URL", "")
		return
	}
	themeIDUUID, err := uuid.Parse(themeID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid theme ID", nil)
		return
	}

	reader, err := io.ReadAll(r.Body)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to read request body", nil)
		return
	}
	defer r.Body.Close()

	var themeUpdate entity.ThemeRequest
	err = json.Unmarshal(reader, &themeUpdate)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	validate = validator.New()
	err = validate.Struct(themeUpdate)
	if err != nil {
		errorMessages := entity.GetValidationThemeErrors(err)
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid request body", errorMessages)
		return
	}

	theme, err := h.themeService.UpdateTheme(themeIDUUID, themeUpdate.Name, themeUpdate.Description, themeUpdate.PreviewImage, themeUpdate.IsPremium, themeUpdate.Status, themeUpdate.Active, themeUpdate.ExpiresIn, themeUpdate.ThemeConfig)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to update theme", nil)
		return
	}

	var themeResponse entity.Theme
	
	themeResponse.ThemeID = theme.ThemeID.String()
	themeResponse.Name = theme.Name.String
	themeResponse.Description = theme.Description.String
	themeResponse.IsPremium = theme.IsPremium.Bool
	themeResponse.PreviewImage = theme.PreviewImage.String
	themeResponse.ThemeConfig = theme.ThemeConfig.RawMessage
	themeResponse.Status = theme.Status.Bool
	themeResponse.Active = theme.Active.Bool
	themeResponse.ExpiresIn = theme.ExpiresIn.Time

	utils.NewMessageResponse(w, http.StatusOK, "Theme updated successfully", themeResponse)
}

func (h *ThemeHandler) DeleteThemeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	themeID := r.PathValue("theme_id")
	if themeID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id in URL", "")
		return
	}

	themeIDUUID, err := uuid.Parse(themeID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid theme ID", nil)
		return
	}

	err = h.themeService.DeleteTheme(themeIDUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to delete theme", nil)
		return
	}

	utils.NewMessageResponse(w, http.StatusOK, "Theme deleted successfully", nil)
}

func (h *ThemeHandler) AddUserThemeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	themeID := r.PathValue("theme_id")
	if themeID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id in URL", "")
		return
	}

	themeIDUUID, err := uuid.Parse(themeID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid theme ID", nil)
		return
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	reader, err := io.ReadAll(r.Body)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to read request body", "")
		return
	}

	var req struct {
		Active bool `json:"active"`
	}

	err = json.Unmarshal(reader, &req)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to unmarshal request body", "")
		return
	}

	userTheme, err := h.themeService.AddUserTheme(userIDUUID, themeIDUUID, req.Active)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to add user theme", nil)
		return
	}

	response := map[string]any{
		"user_id":    userTheme.UserID,
		"theme_id":   userTheme.ThemeID,
		"active":     userTheme.Active.Bool,
		"created_at": userTheme.CreatedAt.Time,
	}

	utils.NewMessageResponse(w, http.StatusOK, "User theme added successfully", response)
}

func (h *ThemeHandler) ListUserThemesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}
	
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	userThemes, err := h.themeService.ListUserThemes(userIDUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to list user themes", nil)
		return	
	}

	var userThemesResponse []entity.ListUserThemesRow

	for _, userTheme := range userThemes {
		userThemeResponse := entity.ListUserThemesRow{
			Theme: entity.Theme{
				ThemeID:      userTheme.ThemeID.String(),
				Name:         userTheme.Name.String,
				Description:  userTheme.Description.String,
				IsPremium:    userTheme.IsPremium.Bool,
				PreviewImage: userTheme.PreviewImage.String,
				ThemeConfig:  userTheme.ThemeConfig.RawMessage,
				Status:       userTheme.Status.Bool,
				Active:       userTheme.ThemeActive.Bool,
				ExpiresIn:    userTheme.ThemeExpiresIn.Time,
			},
			UserID:    userID,
			Useractive:    userTheme.UserActive.Bool,
			UserExpiresIn: userTheme.UserExpiresIn.Time,
		}
		userThemesResponse = append(userThemesResponse, userThemeResponse)
	}

	utils.NewMessageResponse(w, http.StatusOK, "User themes listed successfully", userThemesResponse)
}

func (h *ThemeHandler) UserHasActiveThemeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	hasTheme, err := h.themeService.UserHasActiveTheme(userIDUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to check if user has active theme", nil)
		return
	}

	res := map[string]any{
		"theme_id": hasTheme.ThemeID.String(),
		"name":     hasTheme.Name.String,
		"active":   hasTheme.Active.Bool,
		"theme_config": hasTheme.ThemeConfig.RawMessage,
	}

	utils.NewMessageResponse(w, http.StatusOK, "User has active theme", res)
}

func (h *ThemeHandler) UpdateUserThemeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	themeID := r.PathValue("theme_id")
	if themeID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id in URL", "")
		return
	}

	themeIDUUID, err := uuid.Parse(themeID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid theme ID", nil)
		return
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	reader, err := io.ReadAll(r.Body)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to read request body", "")
		return
	}

	var req struct {
		Active bool `json:"active"`
	}

	err = json.Unmarshal(reader, &req)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Failed to unmarshal request body", "")
		return
	}

	userTheme, err := h.themeService.UpdateUserTheme(userIDUUID, themeIDUUID, req.Active)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to update user theme", nil)
		return
	}

	response := map[string]any{
		"user_id":    userTheme.UserID,
		"theme_id":   userTheme.ThemeID,
		"active":     userTheme.Active.Bool,
		"created_at": userTheme.CreatedAt.Time,
	}

	utils.NewMessageResponse(w, http.StatusOK, "User theme updated successfully", response)
}

func (h *ThemeHandler) RemoveUserThemeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.PathValue("user_id")
	if userID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id in URL", "")
		return
	}

	themeID := r.PathValue("theme_id")
	if themeID == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing user_id in URL", "")
		return
	}

	themeIDUUID, err := uuid.Parse(themeID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid theme ID", nil)
		return
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	err = h.themeService.DeleteUserTheme(userIDUUID, themeIDUUID)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Failed to delete user theme", nil)
		return
	}


	utils.NewMessageResponse(w, http.StatusOK, "User theme deleted successfully", nil)
}