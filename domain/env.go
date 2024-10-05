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
	Name       *string
	UpdateTime *time.Time
}

// -- Var

var ErrVarNotFound = errors.New("local var not found")

type Var struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

type VarCreateArgs struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	Value      string
}

type VarUpdateArgs struct {
	Comment    *string
	CreateTime *time.Time
	EnvName    *string
	Name       *string
	UpdateTime *time.Time
	Value      *string
}

// -- VarRef

var ErrVarRefNotFound = errors.New("local ref not found")

type VarRef struct {
	EnvName    string
	Name       string
	Comment    string
	CreateTime time.Time
	UpdateTime time.Time
	RefEnvName string
	RevVarName string
}

type VarRefCreateArgs struct {
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

	VarCreate(ctx context.Context, args VarCreateArgs) (*Var, error)
	VarDelete(ctx context.Context, envName string, name string) error
	VarList(ctx context.Context, envName string) ([]Var, error)
	VarUpdate(ctx context.Context, envName string, name string, args VarUpdateArgs) error
	VarShow(ctx context.Context, envName string, name string) (*Var, []VarRef, error)

	VarRefCreate(ctx context.Context, args VarRefCreateArgs) (*VarRef, error)
	VarRefDelete(ctx context.Context, envName string, name string) error
	VarRefList(ctx context.Context, envName string) ([]VarRef, []Var, error)
	VarRefShow(ctx context.Context, envName string, name string) (*VarRef, *Var, error)
}
