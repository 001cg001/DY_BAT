package db_mysql

import (
	"DY_BAT/kitex_gen/user"
	"DY_BAT/pkg/tools"
	"errors"
	"sync"
)

var (
	userService     UserService
	userServiceOnce sync.Once
)

type UserService interface {
	UserLogin(username string, password string) (*User, error)
	UserRegister(username string, password string, salt string) error
	FindLastUserId() int64
	GetUserById(userId int64, MyUserid int64) (userResp *user.User, err error)
	UpdateUserFollow(userId int64, followDiff int64) error
	UpdateUserFollower(userId int64, followDiff int64) error
}

type UserServiceImpl struct {
	userDao UserDao
}

func GetUserService() UserService {
	userServiceOnce.Do(func() {
		userService = &UserServiceImpl{
			userDao: GetUserDao(),
		}
	})
	return userService
}

func (u *UserServiceImpl) UserLogin(username string, password string) (*User, error) {
	userGet, err := u.userDao.FindByName(username)
	if err != nil {
		return nil, errors.New("userGet does not exist")
	}
	//MD5加密验证
	password = tools.Md5Util(password, userGet.Salt)
	if password != userGet.Password {
		return nil, errors.New("password error")
	}
	return userGet, nil
}

func (u *UserServiceImpl) UserRegister(username string, password string, salt string) error {
	//判断用户是否已经注册
	_, err := u.userDao.FindByName(username)
	if err == nil {
		return errors.New("userAdd does not exist")
	}
	//添加用户
	userAdd := User{
		Username:      username,
		Password:      password,
		Salt:          salt,
		FollowCount:   0,
		FollowerCount: 0,
	}
	e := u.userDao.AddUser(&userAdd)
	if e != nil {
		return errors.New("userAdd register failed")
	}
	return nil
}

// 返回当前最大的用户ID

func (u *UserServiceImpl) FindLastUserId() int64 {
	return u.userDao.LastId()
}

// 查找用户信息

func (u *UserServiceImpl) GetUserById(userId int64, MyUserid int64) (*user.User, error) {
	userGet, err := u.userDao.FindById(userId)
	var userResp = user.NewUser()
	if err == nil {
		userResp = user.NewUser()
		userResp.Id = userGet.UserId
		userResp.Name = userGet.Username
		userResp.FollowCount = &userGet.FollowCount
		userResp.FollowerCount = &userGet.FollowerCount
		//if isFollow, err := rpc.ExistRelation(MyUserid, userId); err == nil {
		//	userResp.IsFollow = isFollow
		//} else {
		//	fmt.Println(err.Error())
		//	userResp.IsFollow = false
		//}
	}
	return userResp, err
}

func (u *UserServiceImpl) UpdateUserFollow(userId int64, followDiff int64) error {
	return u.userDao.UpdateFollow(userId, followDiff)
}
func (u *UserServiceImpl) UpdateUserFollower(userId int64, followerDiff int64) error {
	return u.userDao.UpdateFollower(userId, followerDiff)
}
