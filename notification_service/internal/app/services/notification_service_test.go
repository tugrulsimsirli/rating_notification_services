package services_test

import (
	"encoding/json"
	"errors"
	"notification_service/config"
	"notification_service/internal/app/services"
	"notification_service/internal/models/dto"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRabbitMQService struct {
	mock.Mock
}

// Define global config
var testConfig = &config.Config{
	RabbitMQ: struct {
		URL       string `yaml:"url"`
		QueueName string `yaml:"queue_name"`
	}{
		URL:       "amqp://guest:guest@localhost:5672/",
		QueueName: "notification_queue",
	},
}

func (m *MockRabbitMQService) CreateChannel(queueName string) (*amqp.Channel, <-chan amqp.Delivery, error) {
	args := m.Called(queueName)
	return args.Get(0).(*amqp.Channel), args.Get(1).(<-chan amqp.Delivery), args.Error(2)
}

func (m *MockRabbitMQService) CloseChannel(ch *amqp.Channel) error {
	args := m.Called(ch)
	return args.Error(0)
}

func TestGetLatestNotifications(t *testing.T) {
	mockRabbitMQService := new(MockRabbitMQService)

	notificationService := &services.NotificationService{
		RabbitMQService: mockRabbitMQService,
		Config:          testConfig,
	}

	mockChannel := new(amqp.Channel)
	mockMsgs := make(chan amqp.Delivery, 1)

	mockRabbitMQService.On("CreateChannel", "notification_queue").Return(mockChannel, (<-chan amqp.Delivery)(mockMsgs), nil)
	mockRabbitMQService.On("CloseChannel", mock.Anything).Return(nil)

	mockNotification := dto.NotificationDto{
		Id:         uuid.New(),
		ProviderID: uuid.New(),
		Message:    "Test Notification",
	}
	mockMsgBody, _ := json.Marshal(mockNotification)

	mockMsgs <- amqp.Delivery{Body: mockMsgBody}

	go func() {
		time.Sleep(3 * time.Millisecond)
		close(mockMsgs)
	}()

	notificationDtos, err := notificationService.GetLatestNotifications()

	assert.NoError(t, err)
	assert.Len(t, notificationDtos, 1)
	assert.Equal(t, mockNotification.Message, notificationDtos[0].Message)

	mockRabbitMQService.AssertExpectations(t)
}

func TestGetLatestNotifications_CreateChannelError(t *testing.T) {
	mockRabbitMQService := new(MockRabbitMQService)

	notificationService := &services.NotificationService{
		RabbitMQService: mockRabbitMQService,
		Config:          testConfig,
	}

	mockRabbitMQService.On("CreateChannel", "notification_queue").Return((*amqp.Channel)(nil), (<-chan amqp.Delivery)(nil), errors.New("channel creation error"))

	notificationDtos, err := notificationService.GetLatestNotifications()

	assert.Error(t, err)
	assert.Nil(t, notificationDtos)
	assert.EqualError(t, err, "channel creation error")

	mockRabbitMQService.AssertExpectations(t)
}

func TestGetLatestNotifications_EmptyQueue(t *testing.T) {
	mockRabbitMQService := new(MockRabbitMQService)

	notificationService := &services.NotificationService{
		RabbitMQService: mockRabbitMQService,
		Config:          testConfig,
	}

	mockChannel := new(amqp.Channel)
	mockMsgs := make(chan amqp.Delivery)

	mockRabbitMQService.On("CreateChannel", "notification_queue").Return(mockChannel, (<-chan amqp.Delivery)(mockMsgs), nil)
	mockRabbitMQService.On("CloseChannel", mock.Anything).Return(nil)

	go func() {
		time.Sleep(3 * time.Millisecond)
		close(mockMsgs)
	}()

	notificationDtos, err := notificationService.GetLatestNotifications()

	assert.NoError(t, err)
	assert.Empty(t, notificationDtos)

	mockRabbitMQService.AssertExpectations(t)
}

func TestGetLatestNotifications_UnmarshalError(t *testing.T) {
	mockRabbitMQService := new(MockRabbitMQService)

	notificationService := &services.NotificationService{
		RabbitMQService: mockRabbitMQService,
		Config:          testConfig,
	}

	mockChannel := new(amqp.Channel)
	mockMsgs := make(chan amqp.Delivery, 1)

	mockRabbitMQService.On("CreateChannel", "notification_queue").Return(mockChannel, (<-chan amqp.Delivery)(mockMsgs), nil)
	mockRabbitMQService.On("CloseChannel", mock.Anything).Return(nil)

	badMessage := []byte("invalid json")
	mockMsgs <- amqp.Delivery{Body: badMessage}

	go func() {
		time.Sleep(3 * time.Millisecond)
		close(mockMsgs)
	}()

	notificationDtos, err := notificationService.GetLatestNotifications()

	assert.NoError(t, err)
	assert.Empty(t, notificationDtos)

	mockRabbitMQService.AssertExpectations(t)
}

func TestGetLatestNotifications_CloseChannelError(t *testing.T) {
	mockRabbitMQService := new(MockRabbitMQService)

	notificationService := &services.NotificationService{
		RabbitMQService: mockRabbitMQService,
		Config:          testConfig,
	}

	mockChannel := new(amqp.Channel)
	mockMsgs := make(chan amqp.Delivery, 1)

	mockRabbitMQService.On("CreateChannel", "notification_queue").Return(mockChannel, (<-chan amqp.Delivery)(mockMsgs), nil)
	mockRabbitMQService.On("CloseChannel", mock.Anything).Return(errors.New("close channel error"))

	mockNotification := dto.NotificationDto{
		Id:         uuid.New(),
		ProviderID: uuid.New(),
		Message:    "Test Notification",
	}
	mockMsgBody, _ := json.Marshal(mockNotification)

	mockMsgs <- amqp.Delivery{Body: mockMsgBody}

	go func() {
		time.Sleep(3 * time.Millisecond)
		close(mockMsgs)
	}()

	notificationDtos, err := notificationService.GetLatestNotifications()

	assert.NoError(t, err)
	assert.Len(t, notificationDtos, 1)
	assert.Equal(t, mockNotification.Message, notificationDtos[0].Message)

	mockRabbitMQService.AssertExpectations(t)
}
