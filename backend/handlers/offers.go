package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/itimofeev/task2trip/backend"
	"github.com/itimofeev/task2trip/rest/restapi/operations/offers"
	"github.com/itimofeev/task2trip/util"
)

var OffersListTaskOffersHandler = offers.ListTaskOffersHandlerFunc(func(params offers.ListTaskOffersParams, principal interface{}) middleware.Responder {
	user := principal.(*backend.User)
	offers_, err := store.ListOffers(user, params.TaskID)
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}
	return offers.NewListOffersOK().WithPayload(convertOffers(offers_))
})

var OffersCreateOfferHandler = offers.CreateOfferHandlerFunc(func(params offers.CreateOfferParams, principal interface{}) middleware.Responder {
	user := principal.(*backend.User)
	offer, err := store.CreateOffer(user, params.TaskID, params.Offer)
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}
	return offers.NewCreateOfferOK().WithPayload(convertOffer(offer))
})
