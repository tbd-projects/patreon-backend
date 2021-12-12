package rabbit

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Session struct {
	name            string
	typeEchange     string
	logger          *logrus.Logger
	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	isReady         bool
}

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second

	// When setting up the channel after a channel exception
	reInitDelay = 2 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second
)

// New creates a new consumer state instance, and automatically
// attempts to connect to the server.
func New(logger *logrus.Logger, name string, addr string) *Session {
	session := Session{
		logger: logger,
		name:   name,
		done:   make(chan bool),
	}
	go session.handleReconnect(addr)
	return &session
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (session *Session) handleReconnect(addr string) {
	for {
		session.isReady = false
		log.Println("Attempting to connect")

		conn, err := session.connect(addr)

		if err != nil {
			log.Println("Failed to connect. Retrying...")

			select {
			case <-session.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := session.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (session *Session) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	session.changeConnection(conn)
	log.Println("Connected!")
	return conn, nil
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (session *Session) handleReInit(conn *amqp.Connection) bool {
	for {
		session.isReady = false

		err := session.init(conn)

		if err != nil {
			log.Println("Failed to initialize channel. Retrying...")

			select {
			case <-session.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-session.done:
			return true
		case <-session.notifyConnClose:
			log.Println("Connection closed. Reconnecting...")
			return false
		case <-session.notifyChanClose:
			log.Println("Channel closed. Re-running init...")
		}
	}
}

// init will initialize channel & declare queue
func (session *Session) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	err = ch.Confirm(false)

	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(
		session.name,
		session.typeEchange,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	session.changeChannel(ch)
	session.isReady = true
	log.Println("Setup!")

	return nil
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (session *Session) changeConnection(connection *amqp.Connection) {
	session.connection = connection
	session.notifyConnClose = make(chan *amqp.Error)
	session.connection.NotifyClose(session.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (session *Session) changeChannel(channel *amqp.Channel) {
	session.channel = channel
	session.notifyChanClose = make(chan *amqp.Error)
	session.channel.NotifyClose(session.notifyChanClose)
}

func (session *Session) GetName() string {
	return session.name
}

func (session *Session) GetChannel() (*amqp.Channel, error) {
	if !session.isReady {
		return nil, ErrGetUninnitChanel
	}
	return session.channel, nil
}

// Close will cleanly shutdown the channel and connection.
func (session *Session) Close() error {
	if !session.isReady {
		return ErrAlreadyClosed
	}
	err := session.channel.Close()
	if err != nil {
		return err
	}
	err = session.connection.Close()
	if err != nil {
		return err
	}
	close(session.done)
	session.isReady = false
	return nil
}
