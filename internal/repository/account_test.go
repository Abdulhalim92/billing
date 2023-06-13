package repository

import (
	"billing/internal/model"
	"billing/logging"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	l := logging.GetLogger()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("cannot open database due error: %v", err)
	}

	r := NewAccountRepository(db, l)

	tests := []struct {
		name    string
		mock    func()
		input   model.Account
		want    int
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				// mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO accounts").
					WithArgs("123456", 100.00, "active").WillReturnRows(rows)
				// mock.ExpectCommit()
			},
			input: model.Account{
				Number:  "123456",
				Balance: 100.00,
				Status:  "active",
			},
			want: 1,
		},
		{
			name: "Empty Fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO accounts").
					WithArgs("123456", 100.00, "").WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: model.Account{
				Number:  "123456",
				Balance: 100.00,
				Status:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateAccount(&tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
