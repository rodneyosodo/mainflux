package api

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/clients/clients"
	"github.com/mainflux/mainflux/clients/groups"
	"github.com/mainflux/mainflux/clients/postgres"
	"github.com/mainflux/mainflux/internal/apiutil"
	"github.com/mainflux/mainflux/pkg/errors"
)

const (
	StatusKey        = "status"
	OffsetKey        = "offset"
	LimitKey         = "limit"
	MetadataKey      = "metadata"
	ParentKey        = "parent_id"
	OwnerKey         = "owner_id"
	IdentifierKey    = "identifier"
	TagKey           = "tag"
	NameKey          = "name"
	TotalKey         = "total"
	SubjectKey       = "subject"
	ObjectKey        = "object"
	PolicyKey        = "policy"
	LevelKey         = "level"
	TreeKey          = "tree"
	DirKey           = "dir"
	VisibilityKey    = "visibility"
	SharedByKey      = "shared_by"
	DefTotal         = uint64(100)
	DefOffset        = 0
	DefLimit         = 10
	DefLevel         = 0
	DefStatus        = "enabled"
	DefClientStatus  = clients.Enabled
	DefGroupStatus   = groups.Enabled
	SharedVisibility = "shared"
	MyVisibility     = "mine"
	AllVisibility    = "all"
	// ContentType represents JSON content type.
	ContentType = "application/json"

	// MaxNameSize limits name size to prevent making them too complex.
	MaxNameSize = math.MaxUint8
)

// EncodeResponse encodes successful response.
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if ar, ok := response.(mainflux.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}
		w.Header().Set("Content-Type", ContentType)
		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

// EncodeError encodes an error response.
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", ContentType)
	switch {
	case errors.Contains(err, errors.ErrMalformedEntity),
		err == apiutil.ErrMissingID,
		err == apiutil.ErrEmptyList,
		err == apiutil.ErrMissingMemberType,
		errors.Contains(err, apiutil.ErrInvalidSecret),
		err == apiutil.ErrNameSize:
		w.WriteHeader(http.StatusBadRequest)
	case errors.Contains(err, errors.ErrAuthentication):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Contains(err, errors.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Contains(err, errors.ErrConflict):
		w.WriteHeader(http.StatusConflict)
	case errors.Contains(err, errors.ErrAuthorization):
		w.WriteHeader(http.StatusForbidden)
	case errors.Contains(err, postgres.ErrMemberAlreadyAssigned):
		w.WriteHeader(http.StatusConflict)
	case errors.Contains(err, errors.ErrUnsupportedContentType):
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case errors.Contains(err, errors.ErrCreateEntity),
		errors.Contains(err, errors.ErrUpdateEntity),
		errors.Contains(err, errors.ErrViewEntity),
		errors.Contains(err, errors.ErrRemoveEntity):
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if errorVal, ok := err.(errors.Error); ok {
		if err := json.NewEncoder(w).Encode(apiutil.ErrorRes{Err: errorVal.Msg()}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
