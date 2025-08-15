package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// testCreateEvent - тест на создание события в хранилище
func TestCreateEvent(t *testing.T) {
	storage := NewStorage()

	event := &Event{
		Title: "Какое-то событие",
	}

	storage.CreateEvent(event)
	assert.Equal(t, 1, event.EventId)
	assert.Contains(t, storage.storageEvents, event.EventId)
	assert.Len(t, storage.storageEvents, 1)

}

// testUpdateEvent - тест на обновление события в хранилище
func TestUpdateEvent(t *testing.T) {
	storage := NewStorage()
	currentEvent := &Event{EventId: 0, Title: "Исходное событие"}
	storage.storageEvents[0] = currentEvent

	updatedEvent := &Event{EventId: 0, Title: "Обновленное событие"}
	err := storage.UpdateEvent(updatedEvent)
	assert.NoError(t, err)
	assert.Equal(t, updatedEvent.Title, storage.storageEvents[0].Title)
}

// testDeleteEvent - тест на удаление собыытия из хранилища
func TestDeleteEvent(t *testing.T) {
	storage := NewStorage()
	removedEvent := &Event{EventId: 0, Title: "Удаляемое событие"}
	storage.storageEvents[removedEvent.EventId] = removedEvent

	result, err := storage.DeleteEvent(removedEvent.EventId)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, removedEvent, result)
	assert.NotContains(t, storage.storageEvents, removedEvent.EventId)
}
