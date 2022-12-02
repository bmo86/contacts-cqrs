package events

import (
	"bytes"
	"contacts_cqrs/models"
	"context"
	"encoding/gob"

	"github.com/nats-io/nats.go"
)

type NatsRepository struct {
	conn               *nats.Conn
	contactCreatedSub  *nats.Subscription
	contactCreatedChan chan CreatedContactMessage
}

func NewNats(url string) (*NatsRepository, error) {
	cn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsRepository{conn: cn}, nil
}

func (n *NatsRepository) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
	if n.contactCreatedSub != nil {
		n.contactCreatedSub.Unsubscribe()
	}
	close(n.contactCreatedChan)
}

// encode data send
func (n *NatsRepository) encodeMsg(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// event publish
func (n *NatsRepository) PublishCreatedContact(ctx context.Context, ct *models.Contact) error {
	msg := CreatedContactMessage{
		Name:       ct.Name,
		Lastname:   ct.Lastname,
		Image:      ct.Image,
		Email:      ct.Email,
		Phone:      ct.Phone,
		Status:     ct.Status,
		UpdateData: ct.UpdateData,
	}

	data, err := n.encodeMsg(msg)
	if err != nil {
		return err
	}
	return n.conn.Publish(msg.Type(), data)
}

// decode data send
func (n *NatsRepository) decodeMSg(data []byte, m interface{}) (err error) {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

// created contact
func (n *NatsRepository) OnCreatedContact(f func(CreatedContactMessage)) (err error) {
	msg := CreatedContactMessage{}
	n.contactCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.decodeMSg(m.Data, &msg)
		f(msg)
	})
	return
}

func (n *NatsRepository) SubcribeCreatedContact(ctx context.Context) (<-chan CreatedContactMessage, error) {
	msg := CreatedContactMessage{}
	n.contactCreatedChan = make(chan CreatedContactMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	n.contactCreatedSub, err = n.conn.ChanSubscribe(msg.Type(), ch)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case ms := <-ch:
				n.decodeMSg(ms.Data, &msg)
				n.contactCreatedChan <- msg
			}
		}
	}()

	/* return connection chan, chan with data, nil error */
	return (<-chan CreatedContactMessage)(n.contactCreatedChan), nil
}
