package dbsql

import (
	"context"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestColumnMapper_Columns(t *testing.T) {
	td := TestingDataType{}
	require.NotEmpty(t, td.ColumnBinders())
}

func TestMapColumn(t *testing.T) {
	db := ConnectToDatabase(t)
	defer CloseDB(t, db)

	preparedStatement, err := PrepareStatement(
		insertTestingDataTypeQuery,
	)
	require.NoError(t, err)
	require.NotNil(t, preparedStatement)

	testData := NewFakeTestingDataType(t)

	result, execErr := ExecContext(
		context.Background(),
		db,
		preparedStatement,
		testData.ParameterValues()...,
	)
	require.NoError(t, execErr)
	require.NotNil(t, result)

	t.Run("valid types", func(t *testing.T) {
		preparedStatement, err = PrepareStatement(selectTestingDataTypeQuery)
		require.NoError(t, err)
		require.NotNil(t, preparedStatement)

		row, err := QueryRowContext(
			context.TODO(),
			db,
			preparedStatement,
			BindParameterValue("uuids", pq.Array([]string{testData.UUID})),
		)
		require.NoError(t, err)
		require.NotNil(t, row)
		columns := Columns{
			"testing_datatype_id",
			"testing_datatype_uuid",
			"word",
			"paragraph",
			"metadata",
			"created_at",
		}
		mappedRow, err := MapRow(row, columns)
		require.NoError(t, err)
		require.NotNil(t, mappedRow)

		td := TestingDataType{}
		err = td.BindMappedRow(mappedRow)
		require.NoError(t, err)

		require.NotEmpty(t, td.ID)
		require.Equal(t, testData.UUID, td.UUID)
		require.Equal(t, testData.Word, td.Word)
		require.Equal(t, testData.Paragraph, td.Paragraph)
	})

	t.Run("invalid type", func(t *testing.T) {
		invalidDataType := `
			SELECT
				CAST(td.testing_datatype_id AS VARCHAR) AS test_data_type_id,
				td.testing_datatype_uuid,
				td.word,
				td.paragraph,
				td.metadata,
				td.created_at
			FROM testing_datatypes td
			WHERE (nullif(@id, NULL) IS NULL OR td.testing_datatype_id = @id)
			AND (nullif(@ids, '{}') IS NULL OR td.testing_datatype_id = ANY(@ids))
			AND (nullif(@uuid, NULL) IS NULL OR td.testing_datatype_uuid = @uuid)
			AND (nullif(@uuids, '{}') IS NULL OR td.testing_datatype_uuid = ANY(@uuids))
			ORDER BY td.created_at
		`

		preparedStatement, err = PrepareStatement(invalidDataType)
		require.NoError(t, err)
		require.NotNil(t, preparedStatement)

		row, err := QueryRowContext(
			context.TODO(),
			db,
			preparedStatement,
			BindParameterValue(
				"uuids",
				pq.Array([]string{testData.UUID}),
			),
		)
		require.NoError(t, err)
		require.NotNil(t, row)
		columns := Columns{
			"testing_datatype_id",
			"testing_datatype_uuid",
			"word",
			"paragraph",
			"metadata",
			"created_at",
		}
		mappedRow, err := MapRow(row, columns)
		require.NoError(t, err)
		require.NotNil(t, mappedRow)

		td := TestingDataType{}
		err = td.BindMappedRow(mappedRow)
		require.Error(t, err)
		require.Equal(
			t,
			err.Error(),
			"column testing_datatype_id has a type of string and does not match asserted type: int64",
		)
	})
}
