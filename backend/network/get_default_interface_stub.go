//go:build !windows && !linux && !darwin

package network

import "errors"

func getPublicIPInfo() (*PublicIPInfo, error) {
	return nil, errors.New("default interface lookup not supported on this platform")
}
