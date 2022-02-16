package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/serializers"
	"clean/app/utils/methodsutil"
	"clean/app/utils/msgutil"
	"clean/infra/errors"
	"clean/infra/logger"
	"strings"
	"time"

	"gorm.io/gorm"
)

type users struct {
	*gorm.DB
}

// NewMySqlUsersRepository will create an object that represent the User.Repository implementations
func NewMySqlUsersRepository(db *gorm.DB) repository.IUsers {
	return &users{
		DB: db,
	}
}

func (r *users) Save(user *domain.User) (*domain.User, *errors.RestErr) {
	res := r.DB.Model(&domain.User{}).Create(&user)

	if res.Error != nil {
		logger.Error("error occurred when create user", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return user, nil
}

func (r *users) GetUser(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	//var resp domain.User
	var intUser domain.IntermediateUserWithPermissions
	var userWithParams domain.UserWithPerms

	sections := `
		users.*,
		roles.name role_name
	`
	if withPermission {
		sections += ",GROUP_CONCAT(DISTINCT permissions.name) AS permissions"
	}

	query := r.DB.Model(&domain.User{}).
		Select(sections).
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Where("users.deleted_at IS NULL")

	if withPermission {
		query = query.
			Joins("JOIN role_permissions ON users.role_id = role_permissions.role_id").
			Joins("JOIN permissions ON role_permissions.permission_id = permissions.id")
	}

	query.Group("users.id")

	res := query.Where("users.id = ?", userID).Find(&intUser)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("getting user with permission"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodsutil.StructToStruct(intUser, &userWithParams.User)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	userWithParams.RoleName = intUser.RoleName

	if withPermission {
		userWithParams.Permissions = strings.Split(intUser.Permissions, ",")
	}

	return &userWithParams, nil
}

func (r *users) GetUserByID(userID uint) (*domain.User, *errors.RestErr) {
	var resp domain.User

	res := r.DB.Model(&domain.User{}).Where("id = ?", userID).First(&resp)

	if res.RowsAffected == 0 {
		logger.Error("error occurred when getting user by user id", res.Error)
		return nil, errors.NewNotFoundError(errors.ErrRecordNotFound)
	}

	if res.Error != nil {
		logger.Error("error occurred when getting user by user id", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (r *users) Update(user *domain.User) *errors.RestErr {
	res := r.DB.Model(&domain.User{}).Omit("company_id", "password", "app_key").Where("id = ? AND company_id = ?", user.ID, user.CompanyID).Updates(&user)

	if res.Error != nil {
		logger.Error("error occurred when updating user by user id", res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *users) UpdatePassword(userID uint, companyID uint, updateValues map[string]interface{}) *errors.RestErr {
	res := r.DB.Model(&domain.User{}).Where("id = ? AND company_id = ?", userID, companyID).Updates(&updateValues)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("updating user by user id"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *users) GetUserByAppKey(appKey string) (*domain.User, *errors.RestErr) {
	var resp domain.User

	res := r.DB.Model(&domain.User{}).Where("app_key = ?", appKey).First(&resp)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError("no user found")
	}

	if res.Error != nil {
		logger.Error("error occurred when getting user by app key", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (r *users) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	res := r.DB.Model(&domain.User{}).Where("email = ?", email).Find(&user)
	if res.RowsAffected == 0 {
		logger.Error("no user found by this email", res.Error)
		return nil, errors.NewError(errors.ErrRecordNotFound)
	}
	if res.Error != nil {
		logger.Error("error occurred when trying to get user by email", res.Error)
		return nil, errors.NewError(errors.ErrSomethingWentWrong)
	}

	return user, nil
}

func (r *users) GetUsersByCompanyIdAndRole(companyID, roleID uint,
	pagination *serializers.Pagination) ([]*domain.IntermediateUserResp, *errors.RestErr) {
	var resp []*domain.IntermediateUserResp

	var totalRows int64 = 0
	tableName := "users"

	stmt := GenerateFilteringCondition(r.DB, tableName, pagination, false)

	stmt = stmt.Model(&domain.Users{}).
		Select("companies.id company_id, companies.name company_name,"+
			"users.id, users.user_name, users.first_name, users.last_name, users.email, users.phone, users.profile_pic, "+
			"users.role_id, users.created_at, users.updated_at, users.last_login_at, users.first_login").
		Joins("LEFT JOIN companies ON users.company_id = companies.id").
		Where("company_id = ? AND role_id = ?", companyID, roleID).
		Find(&resp)

	if len(pagination.QueryString) > 0 {
		searchStmt := "users.user_name LIKE ? OR users.first_name LIKE ? OR users.last_name LIKE ? OR users.email LIKE ? OR users.phone LIKE ?"
		searchTerm := "%" + pagination.QueryString + "%"
		stmt.Where(searchStmt, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}
	res := stmt.Find(&resp)
	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError("no users found")
	}

	if res.Error != nil {
		logger.Error("error occurred when getting users by company_id and role id", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	pagination.Rows = resp

	stmt = GenerateFilteringCondition(r.DB, tableName, pagination, true)
	errCount := r.DB.Model(&domain.Users{}).
		Joins("LEFT JOIN companies ON users.company_id = companies.id").
		Where("company_id = ? AND role_id = ?", companyID, roleID).
		Count(&totalRows).Error

	if errCount != nil {
		logger.Error("error occurred when getting total users by company_id and role id", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	pagination.TotalRows = totalRows
	totalPages := CalculateTotalPageAndRows(pagination, totalRows)
	pagination.TotalPages = totalPages
	return resp, nil
}

func (r *users) SetLastLoginAt(user *domain.User) error {
	lastLoginAt := time.Now().UTC()

	err := r.DB.Model(&user).Update("last_login_at", lastLoginAt).Error

	if err != nil {
		logger.Error(err.Error(), err)
		return err
	}

	return nil
}

func (r *users) HasRole(userID, roleID uint) bool {
	var count int64
	count = 0

	r.DB.Model(&domain.User{}).
		Where("id = ? AND role_id = ?", userID, roleID).
		Count(&count)

	return count > 0
}

func (r *users) ResetPassword(userID int, hashedPass []byte) error {
	err := r.DB.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("password", hashedPass).
		Error

	if err != nil {
		logger.Error("error occur when reset password", err)
		return err
	}

	return nil
}

func (r *users) GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr) {
	tempUser := &domain.TempVerifyTokenResp{}
	var vtUser domain.VerifyTokenResp

	query := r.tokenUserFetchQuery()

	res := query.Where("users.id = ?", id).Find(&tempUser)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("get token user"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodsutil.StructToStruct(tempUser, &vtUser.BaseVerifyTokenResp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	vtUser.Permissions = strings.Split(tempUser.Permissions, ",")

	return &vtUser, nil
}

func (r *users) GetUserWithPermissions(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	var intUser domain.IntermediateUserWithPermissions
	var userWithParams domain.UserWithPerms

	sections := `
		users.*,
		roles.name role_name
	`
	if withPermission {
		sections += ",GROUP_CONCAT(DISTINCT permissions.name) AS permissions"
	}

	query := r.DB.Model(&domain.User{}).
		Select(sections).
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Where("users.deleted_at IS NULL")

	if withPermission {
		query = query.
			Joins("JOIN role_permissions ON users.role_id = role_permissions.role_id").
			Joins("JOIN permissions ON role_permissions.permission_id = permissions.id")
	}

	query.Group("users.id")

	res := query.Where("users.id = ?", userID).Find(&intUser)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("getting user with permission"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodsutil.StructToStruct(intUser, &userWithParams.User)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	userWithParams.RoleName = intUser.RoleName

	if withPermission {
		userWithParams.Permissions = strings.Split(intUser.Permissions, ",")
	}

	return &userWithParams, nil
}

func (r *users) tokenUserFetchQuery() *gorm.DB {
	selections := `
		users.id,
		users.first_name,
		users.last_name,
		users.email,
		users.phone,
		users.profile_pic,
		companies.business_id,
		businesses.name business_name,
		companies.id company_id,
		companies.name company_name,
		(
			CASE
				WHEN 1 IN (GROUP_CONCAT(DISTINCT users.role_id)) THEN 1 ELSE 0
			END
		) AS admin,
		(
			CASE
				WHEN 3 IN (GROUP_CONCAT(DISTINCT users.role_id)) THEN 1 ELSE 0
			END
		) AS super_admin,
		GROUP_CONCAT(DISTINCT permissions.name) AS permissions
	`

	return r.DB.Table("users").
		Select(selections).
		Joins("LEFT JOIN companies ON users.company_id = companies.id").
		Joins("LEFT JOIN businesses ON companies.business_id = businesses.id").
		Joins("JOIN roles ON users.role_id = roles.id").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("users.deleted_at IS NULL").
		Group("users.id")
}
