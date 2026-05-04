package handler

import (
	"agricultural-product-mall/internal/model"
	"agricultural-product-mall/pkg/database"
	"agricultural-product-mall/pkg/jwt"
	"agricultural-product-mall/pkg/response"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Type  string `json:"type" binding:"required"`
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(regex).MatchString(email)
}

func generateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(r.Intn(10))
	}
	return code
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SendVerificationCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if !isValidEmail(req.Email) {
		response.BadRequest(c, "邮箱格式不正确")
		return
	}

	var existingUser model.User
	database.DB.Where("email = ?", req.Email).First(&existingUser)
	if req.Type == "register" && existingUser.ID > 0 {
		response.BadRequest(c, "该邮箱已被注册")
		return
	}

	code := generateCode()

	verificationCode := model.VerificationCode{
		Email:     req.Email,
		Code:      code,
		Type:      req.Type,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&verificationCode).Error; err != nil {
		response.InternalServerError(c, "发送验证码失败")
		return
	}

	response.SuccessWithMessage(c, "验证码已发送", gin.H{
		"code":   code,
		"email":  req.Email,
		"expire": "10分钟",
	})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if len(req.Password) < 6 {
		response.BadRequest(c, "密码长度不能少于6位")
		return
	}

	var existingUser model.User
	database.DB.Where("username = ?", req.Username).First(&existingUser)
	if existingUser.ID > 0 {
		response.BadRequest(c, "用户名已存在")
		return
	}

	database.DB.Where("email = ?", req.Email).First(&existingUser)
	if existingUser.ID > 0 {
		response.BadRequest(c, "邮箱已被注册")
		return
	}

	var verificationCode model.VerificationCode
	database.DB.Where("email = ? AND code = ? AND type = ? AND used = ?",
		req.Email, req.Code, "register", 0).
		Where("expires_at > ?", time.Now()).
		First(&verificationCode)

	if verificationCode.ID == 0 {
		response.BadRequest(c, "验证码无效或已过期")
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		response.InternalServerError(c, "密码加密失败")
		return
	}

	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Role:     "user",
		Status:   1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		response.InternalServerError(c, "注册失败: "+err.Error())
		return
	}

	verificationCode.Used = 1
	database.DB.Save(&verificationCode)

	now := time.Now()
	user.EmailVerifiedAt = &now
	user.LastLoginAt = &now
	user.LastLoginIP = c.ClientIP()
	database.DB.Save(&user)

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.InternalServerError(c, "生成Token失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user model.User
	database.DB.Where("username = ?", req.Username).First(&user)
	if user.ID == 0 {
		database.DB.Where("email = ?", req.Username).First(&user)
	}

	if user.ID == 0 {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		response.Forbidden(c, "账号已被禁用")
		return
	}

	if !CheckPassword(req.Password, user.Password) {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = c.ClientIP()
	database.DB.Save(&user)

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.InternalServerError(c, "生成Token失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
			"avatar":   user.Avatar,
			"nickname": user.Nickname,
			"role":     user.Role,
		},
	})
}

func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"phone":    user.Phone,
		"avatar":   user.Avatar,
		"nickname": user.Nickname,
		"gender":   user.Gender,
		"role":     user.Role,
		"status":   user.Status,
	})
}

type UpdateUserInfoRequest struct {
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
}

func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Gender >= 0 {
		user.Gender = req.Gender
	}

	if err := database.DB.Save(&user).Error; err != nil {
		response.InternalServerError(c, "更新失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"phone":    user.Phone,
		"avatar":   user.Avatar,
		"nickname": user.Nickname,
		"gender":   user.Gender,
	})
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func UpdatePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if len(req.NewPassword) < 6 {
		response.BadRequest(c, "新密码长度不能少于6位")
		return
	}

	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	if !CheckPassword(req.OldPassword, user.Password) {
		response.BadRequest(c, "原密码错误")
		return
	}

	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		response.InternalServerError(c, "密码加密失败")
		return
	}

	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		response.InternalServerError(c, "更新密码失败")
		return
	}

	response.SuccessWithMessage(c, "密码更新成功", nil)
}
