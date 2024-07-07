package service

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"image-go/types"
)

type UserShadowData struct {
	Hash     string
	HashType string
	Salt     string
	Password string
}

func findUserGeneric(path string, username string) *string {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, username) {
			return &line
		}
	}

	return nil
}

func findUserShadowEntry(username string) *string {
	return findUserGeneric("/etc/shadow", username)
}

func findUserPasswdEntry(username string) *string {
	return findUserGeneric("/etc/passwd", username)
}

func parseUserShadowEntry(s string) (*UserShadowData, *types.ResponseError) {
	pieces := strings.Split(s, ":")

	if len(pieces) <= 2 {
		return nil, &types.ResponseError{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		}
	}

	hash := pieces[1]

	pieces = strings.Split(hash, "$")

	if len(pieces) < 4 || pieces[1] == "y" {
		return nil, &types.ResponseError{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		}
	}

	return &UserShadowData{
		Hash:     hash,
		HashType: pieces[1],
		Salt:     pieces[2],
		Password: pieces[3],
	}, nil
}

func parseUserPasswdEntry(s string) (*types.UserPasswdData, *types.ResponseError) {

	pieces := strings.Split(s, ":")

	if len(pieces) > 7 {
		return nil, &types.ResponseError{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		}
	}

	uid, err := strconv.Atoi(pieces[2])
	if err != nil {
		return nil, &types.ResponseError{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		}
	}
	gid, err := strconv.Atoi(pieces[3])
	if err != nil {
		return nil, &types.ResponseError{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		}
	}

	return &types.UserPasswdData{
		Username: pieces[0],
		Password: pieces[1],
		Uid:      uint32(uid),
		Gid:      uint32(gid),
		Gecos:    pieces[4],
		Home:     pieces[5],
		Shell:    pieces[6],
	}, nil
}

func validateUserPassword(shadow UserShadowData, password string) *types.ResponseError {
	output, err := exec.Command(
		"openssl",
		"passwd",
		fmt.Sprintf("-%s", shadow.HashType),
		"-salt",
		shadow.Salt,
		password,
	).CombinedOutput()

	if err != nil {
		return &types.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	hash := strings.TrimSpace(string(output))

	if hash != shadow.Hash {
		return &types.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	return nil
}

func createUser(username string, password string) *types.ResponseError {
	cmd := exec.Command(
		"adduser",
		"-D",
		username,
		"-s",
		"/bin/zsh",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		return &types.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	cmd = exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("echo %s:%s | chpasswd", username, password),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return &types.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	return nil
}

type UserService struct {
	Mutex sync.RWMutex
}

func NewUserService() UserService {
	return UserService{}
}

func (user *UserService) GetUserData(username string) (*types.UserPasswdData, *types.ResponseError) {
	passwdEntry := findUserPasswdEntry(username)
	if passwdEntry == nil {
		return nil, &types.ResponseError{
			Code:    http.StatusNotFound,
			Message: "Not found",
		}
	}
	return parseUserPasswdEntry(*passwdEntry)
}

func (user *UserService) Register(username string, password string) (*types.ResponseResult, *types.ResponseError) {
	shadowEntry := findUserShadowEntry(username)
	if shadowEntry != nil {
		shadow, err := parseUserShadowEntry(*shadowEntry)
		if err != nil {
			return nil, err
		}

		err = validateUserPassword(*shadow, password)
		if err != nil {
			return nil, err
		}

		return &types.ResponseResult{
			Code:    http.StatusOK,
			Message: "Ok",
		}, nil
	} else {
		user.Mutex.Lock()
		defer func() {
			user.Mutex.Unlock()
		}()
		err := createUser(username, password)
		if err != nil {
			return nil, err
		}

		return &types.ResponseResult{
			Code:    http.StatusCreated,
			Message: "Created",
		}, nil
	}
}

func (user *UserService) Login(username string, password string) (*types.ResponseResult, *types.ResponseError) {
	shadowEntry := findUserShadowEntry(username)
	if shadowEntry != nil {
		shadow, err := parseUserShadowEntry(*shadowEntry)
		if err != nil {
			return nil, err
		}

		err = validateUserPassword(*shadow, password)
		if err != nil {
			return nil, err
		}

		return &types.ResponseResult{
			Code:    http.StatusOK,
			Message: "Ok",
		}, nil
	}

	return nil, &types.ResponseError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	}
}
