package account

// Account 用户核心实体（对应数据库 accounts 表）
// GORM 会根据此结构体自动创建表
type Account struct {
	ID           uint   `gorm:"primaryKey" json:"id"`                          // 用户ID，主键
	Username     string `gorm:"unique" json:"username"`                        // 用户名，唯一不可重复
	Password     string `json:"-"`                                             // 密码，不返回给前端
	Token        string `json:"-"`                                             // 登录Token，不返回给前端
	RefreshToken string `json:"-"`                                             // 刷新Token，不返回给前端
	AvatarURL    string `gorm:"type:varchar(512)" json:"avatar_url,omitempty"` // 头像URL
	Bio          string `gorm:"type:varchar(255)" json:"bio,omitempty"`        // 个人简介
}
type Account1 struct {
	ID           uint   `gorm:"primaryKey" json:"id"`                          // 用户ID，主键
	Username     string `gorm:"unique" json:"username"`                        // 用户名，唯一不可重复
	Password     string `json:"-"`                                             // 密码，不返回给前端
	Token        string `json:"-"`                                             // 登录Token，不返回给前端
	RefreshToken string `json:"-"`                                             // 刷新Token，不返回给前端
	AvatarURL    string `gorm:"type:varchar(512)" json:"avatar_url,omitempty"` // 头像URL
	Bio          string `gorm:"type:varchar(255)" json:"bio,omitempty"`        // 个人简介
}

// CreateAccountRequest 注册请求参数
type CreateAccountRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// RenameRequest 修改用户名请求
type RenameRequest struct {
	NewUsername string `json:"new_username"` // 新用户名
}

// FindByIDRequest 根据用户ID查询请求
type FindByIDRequest struct {
	ID uint `json:"id"` // 用户ID
}

// FindByIDResponse 根据ID查询返回的用户信息（不包含敏感字段）
type FindByIDResponse struct {
	ID        uint   `json:"id"`         // 用户ID
	Username  string `json:"username"`   // 用户名
	AvatarURL string `json:"avatar_url"` // 头像
	Bio       string `json:"bio"`        // 简介
}

// FindByUsernameRequest 根据用户名查询请求
type FindByUsernameRequest struct {
	Username string `json:"username"` // 用户名
}

// FindByUsernameResponse 根据用户名查询返回的信息
type FindByUsernameResponse struct {
	ID       uint   `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	Username    string `json:"username"`     // 用户名
	OldPassword string `json:"old_password"` // 旧密码
	NewPassword string `json:"new_password"` // 新密码
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// LoginResponse 登录成功返回
type LoginResponse struct {
	Token        string `json:"token"`         // 登录凭证
	RefreshToken string `json:"refresh_token"` // 刷新凭证
	AccountID    uint   `json:"account_id"`    // 用户ID
	Username     string `json:"username"`      // 用户名
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	AvatarURL string `json:"avatar_url"` // 头像地址
	Bio       string `json:"bio"`        // 个人简介
}

// RefreshRequest 刷新Token请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"` // 刷新凭证
}

// GetProfileRequest 获取个人主页请求
type GetProfileRequest struct {
	AccountID uint `json:"account_id"` // 用户ID
}

// GetProfileResponse 个人主页完整信息（包含统计数据）
type GetProfileResponse struct {
	Account       FindByIDResponse `json:"account"`        // 用户基础信息
	VideoCount    int64            `json:"video_count"`    // 发布视频数
	TotalLikes    int64            `json:"total_likes"`    // 总获赞数
	FollowerCount int64            `json:"follower_count"` // 粉丝数
	VloggerCount  int64            `json:"vlogger_count"`  // 关注数
}
