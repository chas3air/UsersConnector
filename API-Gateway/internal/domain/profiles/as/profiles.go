package asprofiles

import (
	"api-gateway/internal/domain/models"

	authv1 "github.com/chas3air/protos/gen/go/auth"
	"github.com/google/uuid"
)

func UsrToProtoUsr(user models.User) *authv1.User {
	return &authv1.User{
		Id:       user.Id.String(),
		Login:    user.Login,
		Password: user.Password,
		Role:     user.Role,
	}
}

func ProtoUsrToUsr(proto_usr *authv1.User) (models.User, error) {
	parsedUUID, err := uuid.Parse(proto_usr.GetId())
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Id:       parsedUUID,
		Login:    proto_usr.GetLogin(),
		Password: proto_usr.GetPassword(),
		Role:     proto_usr.GetRole(),
	}, nil
}
