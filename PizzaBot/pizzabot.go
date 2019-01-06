package main

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"math"
	"strconv"
	"time"
)

func (bot *Bot) HandleAction(req *RasaRequest) (*RasaResponse, error) {
	log.Println(req)
	resp := NewRasaResponse()

	// TODO remove this it is just for train online testing
	bot.checkOrSetInputSlots(req, resp)

	action := req.NextAction
	log.Println(action)
	/*	switch action {
		case ACTION_CHECK_IS_OPEN:
			bot.ActionCheckIsOpen(req, resp)
		case ACTION_CHECK_IS_OPEN_ON_DAY:
			bot.ActionCheckIsOpenOnDay(req, resp)
		case ACTION_CHECK_TIME_CLOSE:
			bot.ActionCheckTimeClose(req, resp)
		case ACTION_CHECK_TIME_CLOSE_ON_DAY:
			bot.ActionCheckTimeCloseOnDay(req, resp)
		case ACTION_CHECK_RESERVATION_DATETIME:
			bot.ActionCheckReservationDatetime(req, resp)
		case ACTION_SET_SCHEDULED_TIME_SLOT:
			bot.ActionSetScheduledTimeSlot(req, resp)
		case ACTION_SET_SIZE_SLOT:
			bot.ActionSetSizeSlot(req, resp)
		case ACTION_ASK_IF_SIMILAR_TIMES_WORK:
			bot.ActionAskIfSimilarTimesWork(req, resp)
		case ACTION_UTTER_ASK_IS_OTHER_RESERVATION_TIME_OKAY:
			bot.ActionUtterAskIsOtherReservationTimeOkay(req, resp)
		case ACTION_POST_RESERVATION_SAVED:
			bot.ActionUtterPostReservationSaved(req, resp)
		case ACTION_SAVE_RESERVATION:
			bot.ActionSaveReservation(req, resp)
		case ACTION_AFFIRM_SIMILAR_TIME:
			bot.ActionAffirmSimilarTime(req, resp)
		case ACTION_AFFIRM_SIMILAR_TIME_ORDINAL:
			bot.ActionAffirmSimilarTimeOrdinal(req, resp)
		}*/

	return resp, nil
}

func (bot *Bot) ActionUtterAskForPolarOrOrdinalOrTimeOnWhichIfAnyAlternativePotentialTimesForReservationAcceptable(req *RasaRequest, resp *RasaResponse) {
	// TODO should verify potential times
	times := req.Tracker.Slots[POTENTIAL_TIMES].([]interface{})
	reply := "Nothing is available then but do any of the following times work: "

	for i, v := range times {
		datetime, _ := time.Parse(time.RFC3339, v.(string))

		if i != 0 {
			reply += ", "
		}

		reply += datetime.Format("3:04 PM")
	}
	reply += "?"
	resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot *Bot) ActionAffirmSimilarTime(req *RasaRequest, resp *RasaResponse) {
	event := Event{Event: FOLLOWUP, Name: ACTION_AFFIRM_SIMILAR_TIME_ORDINAL}
	resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionAffirmSimilarTimeOrdinal(req *RasaRequest, resp *RasaResponse) {
	// Save the correct potential time into the slot and then followup with action_save_reservation
	potential_times := req.Tracker.Slots[POTENTIAL_TIMES].([]interface{})
	name := req.Tracker.Slots[NAME]

	// Need to get the ordinal from entities
	entities := req.Tracker.LatestMessage.Entities

	ordinal := -1.0
	for _, v := range entities {
		if v.Entity == "ordinal" {
			ordinal = v.Value.(float64)
		}
	}

	if int(ordinal) > len(potential_times) {
		// TODO utter_error
	} else {
		time := potential_times[int(ordinal)-1]
		log.Println(time)

		// TODO Remove these follow ups
		event := Event{}
		switch name.(type) {
		case string:
			event = Event{Event: FOLLOWUP, Name: ACTION_SAVE_RESERVATION}
		default:
			event = Event{Event: FOLLOWUP, Name: UTTER_ASK_NAME}
		}
		resp.Events = append(resp.Events, event)

		event = Event{Event: SLOT, Name: SCHEDULED_TIME, Value: time}
		resp.Events = append(resp.Events, event)

	}
}

func (bot *Bot) ActionUtterAskIsOtherReservationTimeOkay(req *RasaRequest, resp *RasaResponse) {
	times := req.Tracker.Slots[POTENTIAL_TIMES].([]interface{})

	potentialTime := ""
	for _, v := range times {
		datetime, err := time.Parse(time.RFC3339, v.(string))

		if err != nil {
			log.Println(err)
		}

		potentialTime = datetime.Format("3:04 PM")
	}

	reply := fmt.Sprintf("Is %s close enough?", potentialTime)
	resp.Responses = append(resp.Responses, Response{Text: reply})

}

func (bot *Bot) ActionCheckReservationDatetime(req *RasaRequest, resp *RasaResponse) {
	// TODO add support for time intervals
	/* Will have a datetime, businessId, partySize, etc saved in slots */
	//recipientId := req.Tracker.Slots["recipient_id"]
	businessId := req.Tracker.Slots[BUSINESS_ID].(string)
	searchTimeStr := req.Tracker.Slots[SCHEDULED_TIME].(string)
	partySize := req.Tracker.Slots[SIZE].(string)

	searchTime, err := time.Parse(time.RFC3339, searchTimeStr)

	if err != nil {
		log.Println(err)
	}

	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		log.Println(err)
	}

	reservationResult, err := Query(business.ReservationPlatformId, searchTime, partySize)

	if err != nil {
		log.Println(err)
	}

	bot.handleReservationDatetimeQueryResult(reservationResult, req, resp)

}

func (bot *Bot) handleReservationDatetimeQueryResult(reservationResult OpenTableResult, req *RasaRequest, resp *RasaResponse) {
	searchTimeStr := req.Tracker.Slots[SCHEDULED_TIME].(string)

	searchTime, _ := time.Parse(time.RFC3339, searchTimeStr)

	// TODO do better here
	name := ""
	if req.Tracker.Slots[NAME] != nil {
		name = req.Tracker.Slots[NAME].(string)
	}

	if reservationResult.Message == "" {
		// Reservations found within 2.5 hours of request

		found := false
		// Check if one equals exactly and if so make the reservation
		for _, v := range reservationResult.Results {
			if v.Equal(searchTime) {
				found = true
			}
		}

		if found {
			// Exact match, so as long as we have a name we add the reservation to the db

			if name == "" {
				// Action ask the name
				nextAction := Event{Event: "followup", Name: UTTER_ASK_NAME}
				resp.Events = append(resp.Events, nextAction)
				return
			} else {
				// Force Action save_reservation
				nextAction := Event{Event: "followup", Name: ACTION_SAVE_RESERVATION}
				resp.Events = append(resp.Events, nextAction)
				return
			}
		}

		// Check if any are within 15 minutes and if so ask if that is fine
		lessThan15 := false
		selectedTime := reservationResult.Results[0]
		for _, v := range reservationResult.Results {
			if math.Abs(v.Sub(searchTime).Minutes()) <= 15 {
				lessThan15 = true
				selectedTime = v
			}
		}

		if lessThan15 {
			// Action is this one good?

			nextAction := Event{Event: SLOT, Name: POTENTIAL_TIMES, Value: []string{selectedTime.Format(time.RFC3339)}}
			resp.Events = append(resp.Events, nextAction)

			nextAction = Event{Event: FOLLOWUP, Name: ACTION_UTTER_ASK_IS_OTHER_RESERVATION_TIME_OKAY}
			resp.Events = append(resp.Events, nextAction)
			return
		}

		nextAction := Event{Event: SLOT, Name: POTENTIAL_TIMES, Value: reservationResult.Results}
		resp.Events = append(resp.Events, nextAction)

		// action we didn't find any at that time, but do any of these times work for you?
		nextAction = Event{Event: FOLLOWUP, Name: ACTION_ASK_IF_SIMILAR_TIMES_WORK}
		resp.Events = append(resp.Events, nextAction)
		return

	} else if reservationResult.Message == NO_AVAILABLE {
		nextAction := Event{Event: FOLLOWUP, Name: UTTER_NO_RESERVATIONS_AVAILABLE}
		resp.Events = append(resp.Events, nextAction)
		return

	} else if reservationResult.Message == IN_ADVANCE {
		nextAction := Event{Event: FOLLOWUP, Name: UTTER_REQUEST_TIME_TOO_EARLY}
		resp.Events = append(resp.Events, nextAction)
		return
	}

}

func (bot *Bot) ActionSaveReservation(req *RasaRequest, resp *RasaResponse) {

	name := req.Tracker.Slots[NAME].(string)
	size := req.Tracker.Slots[SIZE].(string)
	scheduledTime := req.Tracker.Slots[SCHEDULED_TIME].(string)
	businessId := req.Tracker.Slots[BUSINESS_ID].(string)
	recipientId := req.Tracker.Slots[RECIPIENT_ID].(string)
	contact := req.Tracker.Slots[RECIPIENT_CONTACT].(string)

	datetime, _ := time.Parse(time.RFC3339, scheduledTime)
	timeAsFloat := datetime.UnixNano() / 1000000

	numPeople, err := strconv.ParseInt(size, 0, 32)
	if err != nil {
		log.Println(err)
	}

	reservation := Reservation{
		Name:          name,
		NumPeople:     int(numPeople),
		RecipientId:   recipientId,
		ScheduledTime: timeAsFloat,
		IsVisible:     true,
		Contact:       contact,
	}

	reservationsRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Reservations)
	_, _, err = reservationsRef.Add(bot.Ctx, reservation)

	if err != nil {
		log.Println(err)
	}

}

func (bot *Bot) ActionUtterPostReservationSaved(req *RasaRequest, resp *RasaResponse) {
	//size := req.Tracker.Slots[SIZE].(string)
	scheduledTime := req.Tracker.Slots[SCHEDULED_TIME].(string)

	datetime, _ := time.Parse(time.RFC3339, scheduledTime)

	reply := fmt.Sprintf("Great, you're all set. We'll see you at %s", datetime.Format("3:04 PM"))
	resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot *Bot) checkOrSetInputSlots(req *RasaRequest, resp *RasaResponse) {
	businessId := ""
	if req.Tracker.Slots[BUSINESS_ID] == nil {
		businessId = "MewuHeThW4QJGDxD9tTr"
		nextAction := Event{Event: SLOT, Name: BUSINESS_ID, Value: businessId}
		resp.Events = append(resp.Events, nextAction)

	}
	recipientId := ""
	if req.Tracker.Slots[RECIPIENT_ID] == nil {
		recipientId = "BGeFfREAGGSRqRWrmLNx"
		nextAction := Event{Event: SLOT, Name: RECIPIENT_ID, Value: recipientId}
		resp.Events = append(resp.Events, nextAction)

	}
	recipientContact := ""
	if req.Tracker.Slots[RECIPIENT_CONTACT] == nil {
		recipientContact = "+19084771280"
		nextAction := Event{Event: "slot", Name: RECIPIENT_CONTACT, Value: recipientContact}
		resp.Events = append(resp.Events, nextAction)

	}

}

func (bot *Bot) ActionSetScheduledTimeSlot(req *RasaRequest, resp *RasaResponse) {
	scheduledTime := ""
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == TIME {
			scheduledTime = v.Value.(string)
		}
	}
	nextAction := Event{Event: SLOT, Name: SCHEDULED_TIME, Value: scheduledTime}
	resp.Events = append(resp.Events, nextAction)
}

func (bot *Bot) ActionSetPotentialSizeSlot(req *RasaRequest, resp *RasaResponse) {
	size := 0
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == NUMBER {
			switch v.Value.(type) {
			case float64:
				size = int(v.Value.(float64))
			}
		}
	}

	// TODO undo this hacky string shit
	nextAction := Event{Event: SLOT, Name: "potential_size", Value: fmt.Sprintf("%d", size)}
	resp.Events = append(resp.Events, nextAction)

}

func (bot *Bot) ActionClearPotentialSizeSlot(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "potential_size"}
	resp.Events = append(resp.Events, nextAction)
}

func (bot *Bot) ActionBrancherValidateReservationPotentialSize(req *RasaRequest, resp *RasaResponse) {
	/*potential_size := req.Tracker.Slots["potential_size"]

	if potential_size.(type) != int*/
	// TODO
}

func (bot *Bot) ActionSetSizeSlot(req *RasaRequest, resp *RasaResponse) {
	size := 0
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == NUMBER {
			switch v.Value.(type) {
			case float64:
				size = int(v.Value.(float64))
			}
		}
	}

	// TODO undo this hacky string shit
	nextAction := Event{Event: SLOT, Name: SIZE, Value: fmt.Sprintf("%d", size)}
	resp.Events = append(resp.Events, nextAction)
}

func (bot *Bot) ActionUpdateOrder(req *RasaRequest, resp *RasaResponse) {
	businessId := req.Tracker.Slots[BUSINESS_ID].(string)
	recipientId := req.Tracker.Slots[RECIPIENT_ID].(string)

	orderQuery := bot.Client.Collection(Businesses).Doc(businessId).Collection(Orders).Where("recipientId", "==", recipientId)
	orderQuery = orderQuery.Where("visible", "==", true).OrderBy("lastModificationTime", firestore.Desc)

	order := &Order{}
	iter := orderQuery.Documents(bot.Ctx)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Println(err)
			// TODO handle error
			break
		}

		err = doc.DataTo(order)
		if err != nil {
			log.Println(err)
		}

		order.Id = doc.Ref.ID
	}

	orderType := req.Tracker.Slots["type"].(string)

	if req.Tracker.Slots["address"].(string) != "" && orderType == "" {
		orderType = "DELIVERY"
		resp.Events = append(resp.Events, Event{"slot", "type", "DELIVERY"})
	}

	if order.Id == "" {
		order := Order{
			RecipientId:          req.SenderId,
			RecipientContact:     req.Tracker.Slots["recipient_contact"].(string),
			StartTime:            currentTime(),
			LastModificationTime: currentTime(),
			IsVisible:            true,
		}

		slots := req.Tracker.Slots
		if slots["address"].(string) != "" {
			order.Address = slots["address"].(string)
		}

		if slots["name"].(string) != "" {
			order.Name = slots["name"].(string)
		}

		order.Type = orderType

		if slots["contents"].(string) != "" {
			order.Content = slots["contents"].(string)
		}
		bot.saveOrder(req, &order)

	} else {
		orderRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Orders).Doc(order.Id)

		// Now we have our order and can update it
		orderRef.Update(bot.Ctx, []firestore.Update{
			{Path: "address", Value: req.Tracker.Slots["address"].(string)},
			{Path: "name", Value: req.Tracker.Slots["name"].(string)},
			{Path: "type", Value: orderType},
			{Path: "content", Value: req.Tracker.Slots["contents"].(string)},
			{Path: "lastModificationTime", Value: currentTime()},
		})

	}

}

/*func (bot *Bot) ActionAskNext(req *RasaRequest, resp *RasaResponse) {
	// Figure out which question to ask and return it as a follow up action

	slots := req.Tracker.Slots
	emptySlots := map[string]interface{}

	for k, v := range slots {
		if v == "" {
			emptySlots[k] = v
		}
	}

	// Now we have our empty slots and we can probably just choose any
}*/

func (bot *Bot) ActionCheckTimeClose(req *RasaRequest, resp *RasaResponse) {
	// Need to check the database to see if business is closed or not
	// Then modifies the RasaResponse with the correct RasaResponse

	businessId := req.Tracker.Slots["business_id"].(string)
	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		log.Println(err)
	}

	// TODO make the response dynamic
	reply := fmt.Sprintf("We close at %d", business.TimeClose(""))
	resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot *Bot) ActionCheckTimeCloseOnDay(req *RasaRequest, resp *RasaResponse) {
	// Need to check the database to see if business is closed or not
	// Then modifies the RasaResponse with the correct RasaResponse
	businessId := req.Tracker.Slots["business_id"].(string)
	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		log.Println(err)
	}

	var entity Entity
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == "time" {
			entity = v
			break
		}
	}

	t, err := time.Parse(time.RFC3339, entity.Value.(string))

	if err != nil {
		log.Println(err)
	}

	// TODO make the response dynamic
	reply := fmt.Sprintf("On %d/%d we close at %d", t.Month(), t.Day(), business.TimeClose(""))
	resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot *Bot) ActionBrancherDetermineResponseToCheckIsCurrentlyOpen(req *RasaRequest, resp *RasaResponse) {
	// TODO check if they also input a time

	businessId := req.Tracker.Slots["business_id"].(string)
	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		event := Event{Event: FOLLOWUP, Name: "action_need_employee"}
		resp.Events = append(resp.Events, event)
	}

	if business.IsOpen() {
		event := Event{Event: FOLLOWUP, Name: "action_utter_doing_affirm_currently_open"}
		resp.Events = append(resp.Events, event)
	} else {
		event := Event{Event: FOLLOWUP, Name: "action_utter_doing_deny_currently_open"}
		resp.Events = append(resp.Events, event)

	}

}

// No longer in use
func (bot *Bot) ActionCheckIsOpen(req *RasaRequest, resp *RasaResponse) {
	businessId := req.Tracker.Slots["business_id"].(string)
	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		log.Println(err)
	}

	// TODO make dynamic
	reply := ""
	if business.IsOpen() {
		reply = "Yes, we're open"
	} else {
		reply = "Sorry, no we're not"
	}
	resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot *Bot) ActionCheckIsOpenOnDay(req *RasaRequest, resp *RasaResponse) {
	businessId := req.Tracker.Slots["business_id"].(string)
	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		log.Println(err)
	}

	var entity Entity
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == "time" {
			entity = v
			break
		}
	}

	t, err := time.Parse(time.RFC3339, entity.Value.(string))

	if err != nil {
		log.Println(err)
	}

	// TODO make dynamic
	reply := ""
	if business.IsOpen() {
		reply = "Yes, we're open"
	} else {
		reply = "Sorry, no we're not"
	}

	reply = fmt.Sprintf("%s on %d/%d", reply, t.Month(), t.Day())
	resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot Bot) saveOrder(req *RasaRequest, order *Order) {
	businessId := req.Tracker.Slots["business_id"].(string)
	recipientId := req.Tracker.Slots["recipient_id"].(string)

	ordersRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Orders)

	docRef, _, err := ordersRef.Add(bot.Ctx, order)

	if err != nil {
		log.Println(err)
	}

	order.Id = docRef.ID

	recipientRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Recipients).Doc(recipientId)
	recipientRef.Update(ctx, []firestore.Update{
		{Path: RecentOrderId, Value: order.Id},
	})

}

func currentTime() int64 {
	return time.Now().UnixNano() / 1000000
}

func (bot *Bot) saveMessage(business *Business, recipeint *Recipient, message *Message) error {
	messagesRef := bot.Client.Collection(Businesses).Doc(business.Id).Collection(Messages)
	docRef, _, err := messagesRef.Add(bot.Ctx, message)

	if err != nil {
		return err
	}

	message.Id = docRef.ID
	return nil
}

func (bot *Bot) getBusinessFromId(businessId string) (Business, error) {
	var business Business
	dataSnap, err := bot.Client.Collection(Businesses).Doc(businessId).Get(bot.Ctx)

	if err != nil {
		return business, err
	}

	err = dataSnap.DataTo(&business)

	if err != nil {
		return business, err
	}

	return business, nil

}
