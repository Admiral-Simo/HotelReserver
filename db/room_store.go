package db

import (
	"context"
	"fmt"

	"github.com/Admiral-Simo/HotelReserver/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	roomColl = "rooms"
)

type RoomStore interface {
	Dropper

	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*types.Room, error)
	GetRoomById(ctx context.Context, id primitive.ObjectID) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection(roomColl),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetRoomById(ctx context.Context, id primitive.ObjectID) (*types.Room, error) {
	var result *types.Room
	filter := bson.M{"_id": id}
	if err := s.coll.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	// update the hotel with this room id
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}

	if err := s.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*types.Room, error) {
	cur, err := s.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var rooms []*types.Room
	if err = cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping rooms collection")
	return s.coll.Drop(ctx)
}
