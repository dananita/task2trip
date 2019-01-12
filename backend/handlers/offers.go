package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/itimofeev/task2trip/backend"
	"github.com/itimofeev/task2trip/rest/restapi/operations/offers"
	"github.com/itimofeev/task2trip/util"
)

var OffersListTaskOffersHandler = offers.ListTaskOffersHandlerFunc(func(params offers.ListTaskOffersParams, principal interface{}) middleware.Responder {
	user := principal.(*backend.User)
	offers_, err := Store.ListOffers(user, params.TaskID)
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}
	return offers.NewListTaskOffersOK().WithPayload(convertOffers(offers_))
})

var OffersCreateOfferHandler = offers.CreateOfferHandlerFunc(func(params offers.CreateOfferParams, principal interface{}) middleware.Responder {
	user := principal.(*backend.User)
	offer, err := Store.CreateOffer(user, params.TaskID, params.Offer)
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}
	return offers.NewCreateOfferOK().WithPayload(convertOffer(offer))
})
