package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/leonf08/gophermart.git/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
)

func TestAuthenticatorService_CheckPasswordHash(t *testing.T) {
	h, err := bcrypt.GenerateFromPassword([]byte("user"), bcrypt.DefaultCost)
	require.NoError(t, err)

	a := &AuthenticatorService{}
	type args struct {
		user       *models.User
		storedUser *models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "passwords match",
			args: args{
				user: &models.User{
					Login:    "user",
					Password: "user",
				},
				storedUser: &models.User{
					Login:    "user",
					Password: string(h),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "passwords don't match",
			args: args{
				user: &models.User{
					Login:    "user",
					Password: "password",
				},
				storedUser: &models.User{
					Login:    "user",
					Password: string(h),
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, a.CheckPasswordHash(tt.args.user, tt.args.storedUser), fmt.Sprintf("CheckPasswordHash(%v, %v)", tt.args.user, tt.args.storedUser))
		})
	}
}

func TestAuthenticatorService_GenerateHashFromPassword(t *testing.T) {
	a := &AuthenticatorService{}

	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "generate hash from password",
			args: args{
				user: &models.User{
					Login:    "user",
					Password: "user",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := a.GenerateHashFromPassword(tt.args.user)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateHashFromPassword(%v)", tt.args.user)) {
				return
			}
		})
	}
}

func TestAuthenticatorService_GenerateToken(t *testing.T) {
	type fields struct {
		key string
	}
	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "generate token, no error",
			fields: fields{
				key: "key",
			},
			args: args{
				user: &models.User{
					UserID: 1,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthenticatorService{
				key: tt.fields.key,
			}
			_, err := a.GenerateToken(tt.args.user)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateToken(%v)", tt.args.user)) {
				return
			}
		})
	}
}

func TestAuthenticatorService_ValidateTokenAndExtractClaims(t *testing.T) {
	claims := &models.CustomJWTClaims{
		UserID: 1,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	type fields struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		fault   string
		wantErr bool
	}{
		{
			name: "validate token and extract claims, no error",
			fields: fields{
				key: "key",
			},
			wantErr: false,
		},
		{
			name: "validate token and extract claims, error_invalid_key",
			fields: fields{
				key: "secret",
			},
			wantErr: true,
		},
		{
			name: "validate token and extract claims, error_invalid_signature",
			fields: fields{
				key: "key",
			},
			fault:   "fault",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthenticatorService{
				key: tt.fields.key,
			}
			tk, err := token.SignedString([]byte("key"))
			require.NoError(t, err)

			_, err = a.ValidateTokenAndExtractClaims(tk + tt.fault)
			assert.Equalf(t, tt.wantErr, err != nil, "ValidateTokenAndExtractClaims(%v)", tk)
		})
	}
}

func TestNewAuthenticator(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want *AuthenticatorService
	}{
		{
			name: "new authenticator",
			args: args{
				key: "key",
			},
			want: &AuthenticatorService{
				key: "key",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAuthenticator(tt.args.key)
			assert.Truef(t, reflect.DeepEqual(got, tt.want), "NewAuthenticator(%v)", tt.args.key)
		})
	}
}
