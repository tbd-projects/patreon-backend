package rabbit

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Session struct {
	name            string
	typeEchange     string
	logger          *logrus.Entry
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
)

// New creates a new consumer state instance, and automatically
// attempts to connect to the server.
func New(logger *logrus.Entry, name string, addr string) *Session {
	session := Session{
		logger: logger,
		name:   name,
		done:   make(chan bool),
	}
	session.handleReconnect(addr)
	return &session
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (session *Session) handleReconnect(addr string) {
	session.isReady = false
	session.logger.Infoln("Attempting to connect")

	_, err := session.connect(addr)

	if err != nil {
		session.logger.Warningf("Failed to connect. with err %s", err)
		return
	}

	err = session.init(session.connection)
	if err != nil {
		session.logger.Warningf("Failed to crate channel. with err %s", err)
		return
	}
	session.isReady = true
}

func (session *Session) CheckConnection() bool {
	return session.isReady
}

// connect will create a new AMQP connection
func (session *Session) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-session.done:
			return
		case err := <- session.notifyConnClose:
			session.logger.Errorf("Failed in connection. with err %s", err)
			break
		case err := <- session.notifyChanClose:
			session.logger.Errorf("Failed in channel. with err %s", err)
			break
		}
	} ()

	session.changeConnection(conn)
	session.logger.Infoln("Connected!")
	return conn, nil
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
		"topic",
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
	session.logger.Infoln("Setup!")

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
