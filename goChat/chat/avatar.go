package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNotAcatorはAvatarインスタンスがアバターのURLを返す事ができない時に返すエラー

var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// Avatar Avatarはユーザーのプロフィール画像を返す型
type Avatar interface {
	// GetAvatarURL GetAvatarは指定されたクライアントのアバターのURLを返す。
	// 問題が発生した場合、エラーを返す。URLを取得出来なかった場合、ErrNoAvatarURLを返す。
	GetAvatarURL(ChatUser) (string, error)
}
type AuthAvatar struct{}

type TryAvatars []Avatar

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {

	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
