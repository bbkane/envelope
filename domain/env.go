package domain

import (
	"context"
	"errors"
	"time"
)

// -- Env

var ErrEnvNotFound = errors.New("env not found")

type Env struct {
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
}

type EnvCreateArgs struct {
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
}

type EnvUpdateArgs struct {
	Comment    *string
	CreateTime *time.Time
	NewName    *string
	UpdateTime *time.Time
}

// -- EnvVar

var ErrEnvVarNotFound = errors.New("local var not found")

type EnvVar struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

type EnvVarCreateArgs struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

// -- EnvLocalRef

var ErrEnvRefNotFound = errors.New("local ref not found")

type EnvRef struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	RefEnvName string
	RevVarName string
}

type EnvRefCreateArgs struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	RefEnvName string
	RefVarName string
}

// -- interface

type EnvService interface {
	EnvCreate(ctx context.Context, args EnvCreateArgs) (*Env, error)
	EnvDelete(ctx context.Context, name string) error
	EnvList(ctx context.Context) ([]Env, error)
	EnvUpdate(ctx context.Context, name string, args EnvUpdateArgs) error
	EnvShow(ctx context.Context, name string) (*Env, error)

	EnvVarCreate(ctx context.Context, args EnvVarCreateArgs) (*EnvVar, error)
	EnvVarDelete(ctx context.Context, envName string, name string) error
	EnvVarList(ctx context.Context, envName string) ([]EnvVar, error)
	EnvVarShow(ctx context.Context, envName string, name string) (*EnvVar, []EnvRef, error)

	EnvRefCreate(ctx context.Context, args EnvRefCreateArgs) (*EnvRef, error)
	EnvRefDelete(ctx context.Context, envName string, name string) error
	EnvRefList(ctx context.Context, envName string) ([]EnvRef, []EnvVar, error)
	EnvRefShow(ctx context.Context, envName string, name string) (*EnvRef, *EnvVar, error)

	KeyringEntryCreate(ctx context.Context, args KeyringEntryCreateArgs) (*KeyringEntry, error)
	KeyringEntryList(ctx context.Context) ([]KeyringEntry, []error, error)
}
