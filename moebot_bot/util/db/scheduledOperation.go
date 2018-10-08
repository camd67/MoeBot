package db

import (
	"log"
	"time"
)

type SchedulerType int

const (
	SchedulerChannelRotation SchedulerType = iota
)

type ScheduledOperation struct {
	ID                   int64
	ServerID             int
	Type                 SchedulerType
	PlannedExecutionTime time.Time
}

const (
	scheduledOperationTable = `CREATE TABLE IF NOT EXISTS scheduled_operation(
		id SERIAL NOT NULL PRIMARY KEY,
		server_id INTEGER NOT NULL REFERENCES server(id) ON DELETE CASCADE,
		type INTEGER NOT NULL,
		planned_execution_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIME,
		execution_interval INTERVAL NOT NULL
	)`

	scheduledOperationQueryNow = `SELECT id, server_id, planned_execution_time FROM scheduled_operation WHERE planned_execution_time < CURRENT_TIMESTAMP`

	scheduledOperationQueryServer = `SELECT id, server_id, planned_execution_time FROM scheduled_operation WHERE server_id = $1`

	scheduledOperationUpdate = `UPDATE scheduled_operation SET planned_execution_time = planned_execution_time + execution_interval WHERE id = $1`

	scheduledOperationDelete = `DELETE FROM scheduled_operation WHERE id = $1 AND server_id = $2`

	scheduledOperationInsert = `INSERT INTO scheduled_operation (server_id, type, execution_interval) VALUES ($1, $2, $3) RETURNING id, CURRENT_TIME + execution_interval`
)

func scheduledOperationCreateTable() {
	moeDb.Exec(scheduledOperationTable)
}

func ScheduledOperationQueryNow() ([]*ScheduledOperation, error) {
	rows, err := moeDb.Query(scheduledOperationQueryNow)
	if err != nil {
		log.Println("Error querying for scheduled operations", err)
		return nil, err
	}
	result := []*ScheduledOperation{}
	for rows.Next() {
		operation := new(ScheduledOperation)
		err = rows.Scan(&operation.ID, &operation.ServerID, &operation.PlannedExecutionTime)
		if err != nil {
			log.Println("Error querying for scheduled operations", err)
			return nil, err
		}
		result = append(result, operation)
	}
	return result, nil
}

func ScheduledOperationQueryServer(serverID int) ([]*ScheduledOperation, error) {
	rows, err := moeDb.Query(scheduledOperationQueryServer, serverID)
	if err != nil {
		log.Println("Error querying for scheduled operations", err)
		return nil, err
	}
	result := []*ScheduledOperation{}
	for rows.Next() {
		operation := new(ScheduledOperation)
		err = rows.Scan(&operation.ID, &operation.ServerID, &operation.PlannedExecutionTime)
		if err != nil {
			log.Println("Error querying for scheduled operations", err)
			return nil, err
		}
		result = append(result, operation)
	}
	return result, nil
}

func ScheduledOperationUpdateTime(operationID int64) error {
	_, err := moeDb.Exec(scheduledOperationUpdate, operationID)
	if err != nil {
		log.Println("Error updating scheduled operations", err)
		return err
	}
	return nil
}

func ScheduledOperationDelete(operationID int64, serverID int) (bool, error) {
	r, err := moeDb.Exec(scheduledOperationDelete, operationID, serverID)
	if err != nil {
		log.Println("Error deleting scheduled operations", err)
		return false, err
	}
	rowsAffected, err := r.RowsAffected()
	return rowsAffected > 0, err
}

func scheduledOperationInsertNew(serverID int, operationType SchedulerType, interval string) (*ScheduledOperation, error) {
	var insertID int64
	var nextExecution time.Time
	err := moeDb.QueryRow(scheduledOperationInsert, serverID, operationType, interval).Scan(&insertID, &nextExecution)
	if err != nil {
		log.Println("Error creating scheduled operation", err)
		return nil, err
	}
	err = ScheduledOperationUpdateTime(insertID)
	if err != nil {
		return nil, err
	}
	result := &ScheduledOperation{ID: insertID, ServerID: serverID, Type: operationType, PlannedExecutionTime: nextExecution}
	return result, nil
}
