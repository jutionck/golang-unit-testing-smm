package repository

import (
	"database/sql"
	"errors"
	"log"
	"testing"

	"enigmacamp.com/golatihanlagi/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var dummyCustomers = []model.Customer{
	{
		Id:      "C001",
		Nama:    "Dummy Name 1",
		Address: "Dummy Address 1",
	},
	{
		Id:      "C002",
		Nama:    "Dummy Name 2",
		Address: "Dummy Address 2",
	},
}

type CustomerRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock // github.com/DATA-DOG/go-sqlmock
}

func (suite *CustomerRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("An error when opening a stub database connection", err)
	}
	suite.mockDb = mockDb
	suite.mockSql = mockSql
}

func (suite *CustomerRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func (suite *CustomerRepositoryTestSuite) TestCustomerRetrieveAll_Success() {
	// siapkan column (sama seperti di field table customer)
	rows := sqlmock.NewRows([]string{"id", "nama", "address"})
	for _, v := range dummyCustomers {
		rows.AddRow(v.Id, v.Nama, v.Address)
	}

	// buat query mock nya (menggunakan regex -> (.+)
	suite.mockSql.ExpectQuery("select \\* from customer").WillReturnRows(rows)

	// panggil repository aslinya
	repo := NewCustomerDbRepository(suite.mockDb)

	// panggil method yang mau dtest
	actual, err := repo.RetrieveAll()

	// buat test assertion
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, len(actual))
	assert.Equal(suite.T(), "C001", actual[0].Id) // cek dengan id=C001
}

func (suite *CustomerRepositoryTestSuite) TestCustomerRetrieveAll_Failed() {
	// siapkan column (sama seperti di field table customer)
	rows := sqlmock.NewRows([]string{"ids", "namaaaaa", "addresssss"})
	for _, v := range dummyCustomers {
		rows.AddRow(v.Id, v.Nama, v.Address)
	}

	// buat query mock nya (menggunakan regex -> (.+)
	suite.mockSql.ExpectQuery("select \\* from customer").WillReturnError(errors.New("failed"))

	// panggil repository aslinya
	repo := NewCustomerDbRepository(suite.mockDb)

	// panggil method yang mau dtest
	actual, err := repo.RetrieveAll()

	// buat test assertion
	assert.Nil(suite.T(), actual)
	assert.Error(suite.T(), err)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerCreate_Success() {
	dummyCustomer := dummyCustomers[0]
	suite.mockSql.ExpectExec("insert into customer values").WithArgs(dummyCustomers[0].Id, dummyCustomers[0].Nama, dummyCustomers[0].Address).WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewCustomerDbRepository(suite.mockDb)
	err := repo.Create(dummyCustomer)

	assert.Nil(suite.T(), err)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerCreate_Failed() {
	dummyCustomer := dummyCustomers[0]
	suite.mockSql.ExpectExec("insert into customer values").WillReturnError(errors.New("failed"))
	repo := NewCustomerDbRepository(suite.mockDb)
	err := repo.Create(dummyCustomer)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), errors.New("failed"), err)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerFindById_Success() {
	dummyCustomer := dummyCustomers[0]
	rows := sqlmock.NewRows([]string{"id", "nama", "address"})
	rows.AddRow(dummyCustomer.Id, dummyCustomer.Nama, dummyCustomer.Address)
	// buat query mock nya (menggunakan regex -> (.+)
	suite.mockSql.ExpectQuery("select \\* from customer where id").WillReturnRows(rows)

	// panggil repository aslinya
	repo := NewCustomerDbRepository(suite.mockDb)

	// panggil method yang mau dtest
	actual, err := repo.FindById(dummyCustomer.Id)

	// buat test assertion
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), actual)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerFindById_Failed() {
	dummyCustomer := dummyCustomers[0]
	rows := sqlmock.NewRows([]string{"ids", "namaaaa", "addresssss"})
	rows.AddRow(dummyCustomer.Id, dummyCustomer.Nama, dummyCustomer.Address)
	// buat query mock nya (menggunakan regex -> (.+)
	suite.mockSql.ExpectQuery("select \\* from customer where id").WillReturnError(errors.New("failed"))

	// panggil repository aslinya
	repo := NewCustomerDbRepository(suite.mockDb)

	// panggil method yang mau dtest
	actual, err := repo.FindById(dummyCustomer.Id)

	// buat test assertion
	func() {
		defer func() {
			if r := recover(); r == nil {
				assert.Error(suite.T(), err)
			}
		}()
		// This function should cause a panic
		repo.FindById(dummyCustomer.Id)
	}()
	assert.NotEqual(suite.T(), dummyCustomer, actual)
	assert.Error(suite.T(), err)
}

func TestCustomerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}
