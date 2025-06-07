package usersmongostorage

import (
	"context"
	"fmt"
	"log/slog"
	"usersservice/internal/domain/models"
	"usersservice/pkg/lib/logger/sl"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersMongoStorage struct {
	log    *slog.Logger
	client *mongo.Client
}

const UsersCollectionName = "users"
const DatabaseName = "users"

func New(log *slog.Logger, connStr string) *UsersMongoStorage {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connStr))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		panic(err)
	}

	return &UsersMongoStorage{
		log:    log,
		client: client,
	}
}

func (u *UsersMongoStorage) Close() {
	if err := u.client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

// GetUsers implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.mongo.users.GetUsers"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	collection := u.client.Database(DatabaseName).Collection(UsersCollectionName)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Error("Error fetching usres", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	users := make([]models.User, 0, 10)
	if err := cursor.All(ctx, &users); err != nil {
		log.Error("Error decode users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

// GetUserById implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.mongo.users.GetUserById"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	collection := u.client.Database(DatabaseName).Collection(UsersCollectionName)

	var user models.User

	err := collection.FindOne(ctx, bson.M{"id": uid}).Decode(&user)
	if err != nil {
		log.Error("Error finding user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Insert implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	const op = "storage.mongo.users.Insert"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	collection := u.client.Database(DatabaseName).Collection(UsersCollectionName)

	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Error("Error inserting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User OID", slog.Any("object_id", insertResult.InsertedID))
	return user, nil
}

// Update implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	const op = "storage.mongo.users.Update"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	collection := u.client.Database(DatabaseName).Collection(UsersCollectionName)

	filter := bson.M{"id": uid}

	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": user})
	if err != nil {
		log.Error("Error updating user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Delete implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.mongo.users.Delete"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	collection := u.client.Database(DatabaseName).Collection(UsersCollectionName)

	filter := bson.M{"id": uid}

	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Error("Error fetching user before deleting", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Error("Error deleting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
