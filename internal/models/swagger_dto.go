package models

// ----- Auth DTOs -----

type RegisterRequest struct {
    Email    string `json:"email" example:"user@example.com"`
    Password string `json:"password" example:"strongpassword"`
}

type UserResponse struct {
    ID    uint   `json:"id" example:"1"`
    Email string `json:"email" example:"user@example.com"`
}

type LoginRequest struct {
    Email    string `json:"email" example:"user@example.com"`
    Password string `json:"password" example:"strongpassword"`
}

type LoginResponse struct {
    Token string `json:"token" example:"jwt.token.here"`
}

// ----- Todo DTOs -----

type CreateTodoRequest struct {
    Title string `json:"title" example:"Buy milk"`
}

type UpdateTodoRequest struct {
    Title     string `json:"title" example:"Buy bread"`
    Completed bool   `json:"completed" example:"false"`
}
