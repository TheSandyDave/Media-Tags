package services

import (
	"context"
	"testing"

	apierrors "github.com/TheSandyDave/Media-Tags/api_errors"
	"github.com/TheSandyDave/Media-Tags/domain"
	"github.com/TheSandyDave/Media-Tags/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestIDbModel struct {
	domain.BaseObject
	Name string
}

func TestBaseService_Create_addsModelToDatabase(t *testing.T) {
	t.Parallel()
	// Arrange
	database := utils.NewInMemoryDatabase(t, TestIDbModel{})
	ExpectedName := "expectedOBJ"
	ObjToInsert := TestIDbModel{
		Name: ExpectedName,
	}
	service := baseService[TestIDbModel]{
		Database: database,
	}

	// Act
	err := service.Create(context.Background(), &ObjToInsert)

	// Assert
	assert.NoError(t, err)
	result := TestIDbModel{}
	database.First(&result)
	assert.Equal(t, ExpectedName, result.Name)
}

func TestBaseService_Delete_removesModelFromDatabase(t *testing.T) {
	t.Parallel()
	// Arrange
	database := utils.NewInMemoryDatabase(t, TestIDbModel{})
	objToDelete := TestIDbModel{
		BaseObject: domain.BaseObject{
			ID: uuid.New(),
		},
		Name: "temp",
	}
	service := baseService[TestIDbModel]{
		Database: database,
	}
	require.NoError(t, database.Create(&objToDelete).Error)

	// Act
	err := service.Delete(context.Background(), objToDelete.ID)

	// Assert
	assert.NoError(t, err)
	result := []*TestIDbModel{}
	database.Find(&result)
	assert.Equal(t, 0, len(result))
}

func TestBaseService_GetWithID_retrievesFromDatabase(t *testing.T) {
	t.Parallel()
	// Arrange
	database := utils.NewInMemoryDatabase(t, TestIDbModel{})
	expectedName := "ExpectedObj"
	objToRetrieve := TestIDbModel{
		BaseObject: domain.BaseObject{
			ID: uuid.New(),
		},
		Name: expectedName,
	}
	service := baseService[TestIDbModel]{
		Database: database,
	}
	require.NoError(t, database.Create(&objToRetrieve).Error)

	// Act
	res, err := service.GetWithID(context.Background(), objToRetrieve.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedName, res.Name)
}

func TestBaseService_GetWithID_failsIfRecordDoesNotExist(t *testing.T) {
	t.Parallel()
	// Arrange
	database := utils.NewInMemoryDatabase(t, TestIDbModel{})
	service := baseService[TestIDbModel]{
		Database: database,
	}

	// Act
	_, err := service.GetWithID(context.Background(), uuid.New())

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &apierrors.RecordNotFoundError{}, err)
}

func TestBaseService_GetWithIDs_retrievesFromDatabase(t *testing.T) {
	t.Parallel()
	// Arrange
	database := utils.NewInMemoryDatabase(t, TestIDbModel{})
	expectedName := "ExpectedObj"
	objToRetrieve := TestIDbModel{
		BaseObject: domain.BaseObject{
			ID: uuid.New(),
		},
		Name: expectedName,
	}
	service := baseService[TestIDbModel]{
		Database: database,
	}
	require.NoError(t, database.Create(&objToRetrieve).Error)

	// Act
	res, err := service.GetWithIDs(context.Background(), []uuid.UUID{objToRetrieve.ID})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, expectedName, res[0].Name)
}

func TestBaseService_GetWithIDs_failsIfSomeRecordsDoNotExist(t *testing.T) {
	t.Parallel()
	// Arrange
	database := utils.NewInMemoryDatabase(t, TestIDbModel{})
	service := baseService[TestIDbModel]{
		Database: database,
	}
	objToRetrieve := TestIDbModel{
		BaseObject: domain.BaseObject{
			ID: uuid.New(),
		},
		Name: "temp",
	}
	require.NoError(t, database.Create(&objToRetrieve).Error)
	failedToFindID := uuid.New()
	// Act
	_, err := service.GetWithIDs(context.Background(), []uuid.UUID{failedToFindID, objToRetrieve.ID})

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &apierrors.RecordsNotFoundWithIDs{}, err)
	assert.ErrorContains(t, err, failedToFindID.String())
	assert.NotContains(t, err.Error(), objToRetrieve.ID.String())
}

func TestBaseService_Get_RetrievesAllObjects(t *testing.T) {
	t.Parallel()
	// Arrange
	database := utils.NewInMemoryDatabase(t, TestIDbModel{})
	service := baseService[TestIDbModel]{
		Database: database,
	}

	objsToRetrieve := []*TestIDbModel{
		{Name: "temp1"},
		{Name: "temp2"},
	}
	require.NoError(t, database.Create(&objsToRetrieve).Error)
	// Act
	res, err := service.Get(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}

func TestBaseService_Get_AppliesFiltersCorrectly(t *testing.T) {
	t.Parallel()
	// Arrange
	database := utils.NewInMemoryDatabase(t, TestIDbModel{})
	service := baseService[TestIDbModel]{
		Database: database,
	}
	nameToretrieve := "expectedRes"
	objsToRetrieve := []*TestIDbModel{
		{Name: nameToretrieve},
		{Name: "temp2"},
	}
	require.NoError(t, database.Create(&objsToRetrieve).Error)
	var filter Option[TestIDbModel] = func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", nameToretrieve)
	}
	// Act
	res, err := service.Get(context.Background(), filter)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, nameToretrieve, res[0].Name)
}
