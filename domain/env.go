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

// -- EnvLocalVar

var ErrEnvLocalVarNotFound = errors.New("local var not found")

type EnvLocalVar struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

type EnvLocalVarCreateArgs struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

// -- EnvLocalRef

var ErrEnvLocalRefNotFound = errors.New("local ref not found")

type EnvLocalRef struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	RefEnvName string
	RevVarName string
}

type EnvLocalRefCreateArgs struct {
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

	EnvLocalVarCreate(ctx context.Context, args EnvLocalVarCreateArgs) (*EnvLocalVar, error)
	EnvLocalVarDelete(ctx context.Context, envName string, name string) error
	EnvLocalVarList(ctx context.Context, envName string) ([]EnvLocalVar, error)
	EnvLocalVarShow(ctx context.Context, envName string, name string) (*EnvLocalVar, error)

	EnvRefCreate(ctx context.Context, args EnvLocalRefCreateArgs) (*EnvLocalRef, error)

	KeyringEntryCreate(ctx context.Context, args KeyringEntryCreateArgs) (*KeyringEntry, error)
	KeyringEntryList(ctx context.Context) ([]KeyringEntry, []error, error)
}
