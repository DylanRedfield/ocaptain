package main

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"
)

func (bot *Bot) HandleAction(req *RasaRequest) (*RasaResponse, error) {
	resp := NewRasaResponse()

	// TODO remove this it is just for train online testing
	bot.checkOrSetInputSlots(req, resp)

	action := req.NextAction

	switch action {
	case "action_set_potential_size_slot":
		bot.ActionSetPotentialSizeSlot(req, resp)
	case "action_brancher_validate_reservation_potential_size":
		bot.ActionBrancherValidateReservationPotentialSize(req, resp)
	case "action_clear_potential_size_slot":
		bot.ActionClearPotentialSizeSlot(req, resp)
	case "action_need_employee":
		bot.ActionNeedEmployee(req, resp)
	case "action_set_size_slot":
		bot.ActionSetSizeSlot(req, resp)
	case "action_set_potential_time_slot":
		bot.ActionSetPotentialTimeSlot(req, resp)
	case "action_clear_potential_time_slot":
		bot.ActionClearPotentialTimeSlot(req, resp)
	case "action_test_bed":
		bot.ActionTestBed(req, resp)
  case "action_set_temp_times_slot":
    bot.ActionSetTempTimesSlot(req, resp)
  case "action_set_potential_hour_slot":
    bot.ActionSetPotentialHourSlot(req, resp)
  case "action_brancher_with_temp_times_to_determine_next_from_times_length":
    bot.ActionBrancherWithTempTimesToDetermineNextFromTimesLength(req, resp)
  case "action_brancher_with_temp_times_validate_single_temp_times":
    bot.ActionBrancherWithTempTimesValidateSingleTempTimes(req, resp)
  case "action_brancher_validate_potential_hour_slot":
    bot.ActionBrancherValidatePotentialHourSlot(req, resp)
  case "action_brancher_validate_with_temp_times_and_time_entity_to_modify_temp_times_from_day_grain":
    bot.ActionBrancherValidateWithTempTimesAndTimeEntityToModifyTempTimeFromDayGrain(req, resp)
  case "action_brancher_reservation_slot_filling_base":
    bot.ActionBrancherReservationSlotFillingBase(req, resp)
  case "action_brancher_to_save_new_reservation":
    bot.ActionBrancherToSaveNewReservation(req, resp)
  case "action_need_employee_because_error_saving":
    bot.ActionNeedEmployeeBecauseErrorSaving(req, resp)
  case "action_brancher_with_size_and_single_potential_times_query_reservation_platform":
    bot.ActionBrancherWithSizeAndSinglePotentialTimesQueryReservationPlatform(req, resp)
  case "action_brancher_with_potential_times_and_alternative_times_to_fill_scheduled_time":
      bot.ActionBrancherWithPotentialTimesAndAlterativeTimesToFillScheduledTime(req, resp)
	default:
		log.Println(action)
	}

	return resp, nil
}

func (bot *Bot) ActionTestBed(req *RasaRequest, resp *RasaResponse) {
	entities := req.Tracker.LatestMessage.Entities

	for _, entity := range entities {
		log.Println(entity.Start)
		/*switch v := entity.Value.(type) {
		  default:
		    log.Printf("%T\n", v)
		  } */
	}
}

func (bot *Bot) ActionBrancherReservationSlotFillingBase(req *RasaRequest, resp *RasaResponse) {
  slots := req.Tracker.Slots

  size := slots["size"]
  potentialSize := slots["potential_size"]

  scheduledTime := slots["scheduled_time"]
  tempTimes := slots["temp_times"]
  potentialTimes := slots["potential_times"]

  name := slots["name"]

  event := Event{}

  if reflect.ValueOf(size).IsNil() {
    event.Event = FOLLOWUP
    if reflect.ValueOf(potentialSize).IsNil() {
      event.Name = "utter_ask_for_number_on_reservation_size"
    } else {
      event.Name = "action_checkpoint_with_potential_size_to_validate_and_fill_size"
    }
  } else if reflect.ValueOf(scheduledTime).IsNil() {
    event.Event = FOLLOWUP

    if reflect.ValueOf(potentialTimes).IsNil() {

      if reflect.ValueOf(tempTimes).IsNil() {
        event.Name = "utter_ask_for_time_for_potential_reservation"
      } else {
        event.Name = "action_checkpoint_with_temp_times_to_fill_potential_times"
      }
    } else {
      event.Name = "action_checkpoint_with_size_and_single_potential_times_to_fill_scheduled_time"
    }
  } else if reflect.ValueOf(name).IsNil() {
    event.Name = "utter_ask_for_name"
  } else {
    event.Name = "action_brancher_to_save_new_reservation"
  }

  resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionBrancherWithAlternativeTimesAndOrdinalValidateOrdinalToSelectAlternativeTime(req *RasaRequest, resp *RasaResponse) {
  slots := req.Tracker.Slots

  rawOrdinal := slots["temp_ordinal"]
  rawAlternativeTimes := slots["alternative_times"]

  event := Event{Event: FOLLOWUP}

  if reflect.ValueOf(rawAlternativeTimes).IsNil() {
    event.Name = "action_need_employee"
  } else if reflect.ValueOf(rawOrdinal).IsNil() {
    event.Name = "action_need_employee"
  } else {
    ordinalStr := rawOrdinal.(string)

    ordinalValue, err := strconv.Atoi(ordinalStr)
    if err != nil {
      event.Name = "action_need_employee_because_error"
      resp.Events = append(resp.Events, event)
      return
    }

    

    alternativeTimesArr := rawAlternativeTimes.([]interface{})

    if ordinalValue <= 0 || ordinalValue >= len(alternativeTimesArr) {
      event.Name = "action_need_employee"
    } else {
      // Valid
      alternativeTimesStr := rawAlternativeTimes.([]string)
      event.Event = SLOT
      event.Name = "scheduled_time"
      event.Value = alternativeTimesStr[ordinalValue - 1]
      resp.Events = append(resp.Events, event)

      event = Event { Event : FOLLOWUP, Name: "action_blank_alert_scheduled_time_slot_set"}

    }
  }
  
  resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionBrancherWithPotentialTimesAndAlterativeTimesToFillScheduledTime(req *RasaRequest, resp *RasaResponse) {
  slots := req.Tracker.Slots

  rawPotentialTimes := slots["potential_times"]
  rawAlternativeTimes := slots["alternative_times"]

  event := Event{Event: FOLLOWUP}

  if reflect.ValueOf(rawPotentialTimes).IsNil() {
    event.Name = "utter_ask_for_time_on_potential_reservation"
  } else if reflect.ValueOf(rawAlternativeTimes).IsNil() {
    event.Name = "action_checkpoint_with_size_and_single_potential_times_to_fill_scheduled_time"
  } else {
    potentialTimes := rawPotentialTimes.([]interface{})
    rawPotentialTime := potentialTimes[0]
    potentialTimeStr := rawPotentialTime.(string)

    potentialTime, err := time.Parse(time.RFC3339, potentialTimeStr)

    if err != nil {
      event.Name = "action_need_employee_because_error"
      
      resp.Events = append(resp.Events, event)
      return
    }

    altTimes := []time.Time{}

    altTimesArr := rawAlternativeTimes.([]interface{})
    for _, v := range altTimesArr {
      altTime, err := time.Parse(time.RFC3339, v.(string))

      if err != nil {
        event.Name = "action_need_employee_because_error"
        
        resp.Events = append(resp.Events, event)
        return
      }

      altTimes = append(altTimes, altTime)
    }


		// Check if any are within 15 minutes and if so ask if that is fine
		lessThan15 := false
		for _, v := range altTimes {
			if math.Abs(v.Sub(potentialTime).Minutes()) <= 15 {
				lessThan15 = true
			}
		}

    if lessThan15 {
      event.Name = "action_utter_ask_for_polar_on_if_single_close_alternative_time_acceptable"
    } else if len(altTimes) > 1 {
      event.Name = "action_utter_ask_for_polar_or_time_or_number_or_ordinal_on_serveral_alternative_times"
    } else {
      event.Name = "action_utter_ask_for_polar_on_if_single_alternative_times_acceptable"
    }

  }
}

func (bot *Bot) ActionBrancherToSaveNewReservation(req *RasaRequest, resp *RasaResponse) {
  slots := req.Tracker.Slots

  rawSize := slots["size"]
  rawName := slots["name"]
  rawScheduledTime := slots["scheduled_time"]

	businessId := req.Tracker.Slots[BUSINESS_ID].(string)
	recipientId := req.Tracker.Slots[RECIPIENT_ID].(string)
	contact := req.Tracker.Slots[RECIPIENT_CONTACT].(string)

  event := Event{}
  if reflect.ValueOf(rawSize).IsNil() || reflect.ValueOf(rawName).IsNil() || reflect.ValueOf(rawScheduledTime).IsNil() {
    event.Event = FOLLOWUP
    event.Name = "action_checkpoint_reservation_slot_filling"
  } else {
    name := rawName.(string)
    scheduledTimeStr := rawScheduledTime.(string)

    // We want to round down anyway
    size := int(rawSize.(float64))

    scheduledTime, err := time.Parse(time.RFC3339, scheduledTimeStr)

    if err != nil {
      event.Event = FOLLOWUP
      event.Name = "action_need_employee_because_error_saving"
      
      resp.Events = append(resp.Events, event)
      return
    }

    timeAsFloat := scheduledTime.UnixNano() / 1000000

    reservation := Reservation{
      Name:          name,
      NumPeople:     size,
      RecipientId:   recipientId,
      ScheduledTime: timeAsFloat,
      IsVisible:     true,
      Contact:       contact,
    }

    reservationsRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Reservations)
    _, _, err = reservationsRef.Add(bot.Ctx, reservation)

    event.Event = FOLLOWUP
    if err != nil {
      event.Name = "action_need_employee_because_error_saving"
    } else {
      event.Name = "action_utter_post_reservation_save_AND_ask_for_next_general_request"
    }

  }

  resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionBrancherWithSizeAndSinglePotentialTimesQueryReservationPlatform(req *RasaRequest, resp *RasaResponse) {
  slots := req.Tracker.Slots

  rawSize := slots["size"]
  rawPotentialTimes := slots["potential_times"]

	businessId := req.Tracker.Slots[BUSINESS_ID].(string)

  event := Event{}
  event.Event = FOLLOWUP

  if reflect.ValueOf(rawSize).IsNil() {
    event.Name = "utter_ask_for_number_on_reservation_size"
  } else if reflect.ValueOf(rawPotentialTimes).IsNil() {
    event.Name = "utter_ask_for_time_for_potential_reservation"
  } else {
    business, err := bot.getBusinessFromId(businessId)

    if err != nil {
      event.Name = "action_need_employee"

      resp.Events = append(resp.Events, event)
    }

    // open table query takes in a string for now
    size := string(int(rawSize.(float64)))

    potentialTimes := rawPotentialTimes.([]interface{})
    rawPotentialTime := potentialTimes[0]
    potentialTimeStr := rawPotentialTime.(string)

    potentialTime, err := time.Parse(time.RFC3339, potentialTimeStr)

    if err != nil {
      event.Name = "action_need_employee_because_error"
      resp.Events = append(resp.Events, event)
      return
    }
    reservationResult, err := Query(business.ReservationPlatformId, potentialTime, size)

    if err != nil {
      event.Name = "action_need_employee_because_error"
      return
    }

    if reservationResult.Message == NO_AVAILABLE {
      event.Name = "utter_doing_no_tables_available_near_that_time_AND_ask_for_polar_or_time_on_alternative"
    } else if reservationResult.Message == IN_ADVANCE {
      event.Name = "utter_requested_time_too_soon_AND_ask_for_polar_or_time_on_alternative"
    } else {

      if len(reservationResult.Results) == 0 {
        event.Name = "action_need_employee"
      } else {

        found := false
        for _, result := range(reservationResult.Results) {
          if result.Equal(potentialTime) {
            found = true

		        matchTimeStr := result.Format(time.RFC3339)

            event.Event = SLOT
            event.Name = "scheduled_time"
            event.Value = matchTimeStr

	          resp.Events = append(resp.Events, event)
            event.Event = FOLLOWUP
            event.Value = nil
            break
          }
        }

        if found {
          event.Name = "action_alert_scheduled_time_slot_set"
        } else {
          event.Name = "action_alert_alternative_times_slot_set"
        }
      }
    }

  }

	resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionNeedEmployeeBecauseErrorSaving(req *RasaRequest, resp *RasaResponse) {
  // TODO
	event := Event{Event: "pause"}
	resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionNeedEmployeeBecauseError(req *RasaRequest, resp *RasaResponse) {
  // TODO
	event := Event{Event: "pause"}
	resp.Events = append(resp.Events, event)
}
func (bot *Bot) ActionSetTempTimesSlotFromPotentialHour(req *RasaRequest, resp *RasaResponse) {
  // Set's the temp_times[0] to the current time but with the potential time as the potential_hour + 12 if less than 12
  // and greater than 0
  potentialHour := req.Tracker.Slots["potential_hour"].(float64)

  now := time.Now()

  newTime := time.Date(now.Year(), now.Month(), now.Day(), int(potentialHour), 0, 0, 0, time.UTC)
  newTimeStr := newTime.Format(time.RFC3339)

  tempTime := RasaTime{Value: newTimeStr, Grain: "hour", Type: "value"}

  event := Event{Event: SLOT, Name: "temp_times", Value: []interface{}{tempTime}}
  resp.Events = append(resp.Events, event)


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

func (bot *Bot) ActionNeedEmployee(req *RasaRequest, resp *RasaResponse) {
	// TODO
	event := Event{Event: "pause"}
	resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionSetTempTimesSlot(req *RasaRequest, resp *RasaResponse) {
  entities := req.Tracker.LatestMessage.Entities

  entityTime := Entity{}
  for _, v := range(entities) {
    if v.Entity == "time" {
      entityTime = v
    }
  }

  event := Event{Event: SLOT, Name: "temp_times", Value: []interface{}{entityTime.AdditionalInfo}}
  resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionSetPotentialHourSlot(req *RasaRequest, resp *RasaResponse) {
  entities := req.Tracker.LatestMessage.Entities

  potentialHour := 0.0

  for _, v := range(entities) {
    if v.Entity == "number" {
      potentialHour = v.Value.(float64)
    }
  }

  potentialHour = math.Abs(potentialHour)

  if potentialHour >= 1 && potentialHour <= 11 {
    potentialHour += 12
  }


  event := Event{Event: SLOT, Name: "potential_hour", Value: potentialHour}
  resp.Events = append(resp.Events, event)
}
func (bot *Bot) ActionBrancherWithTempTimesToDetermineNextFromTimesLength(req *RasaRequest, resp *RasaResponse) {
	switch v := req.Tracker.Slots["time_times"].(type) {
	case []interface{}:
		if len(v) == 1 {
			event := Event{Event: FOLLOWUP, Name: "action_brancher_validate_single_temp_times"}
			resp.Events = append(resp.Events, event)
		} else {
			event := Event{Event: FOLLOWUP, Name: "action_brancher_validate_single_temp_times"}
			resp.Events = append(resp.Events, event)
		}
	default:
		event := Event{Event: FOLLOWUP, Name: "utter_ask_for_time_for_potential_reservation"}
		resp.Events = append(resp.Events, event)
	}
}

func (bot *Bot) ActionBrancherWithTempTimesValidateSingleTempTimes(req *RasaRequest, resp *RasaResponse) {
	switch v := req.Tracker.Slots["temp_times"].(type) {
	case []interface{}:
		switch time_map := v[0].(type) {
		case map[string]interface{}:
			var rasaTime RasaTime
			err := mapstructure.Decode(time_map, &rasaTime)

			if err != nil {
				log.Println(err)
				break
			}

			if rasaTime.Grain == "week" || rasaTime.Grain == "month" || rasaTime.Grain == "year" {
				event := Event{Event: FOLLOWUP, Name: "action_need_employee"}
				resp.Events = append(resp.Events, event)
			} else if rasaTime.Grain == "day" {
				event := Event{Event: FOLLOWUP, Name: "utter_with_temp_time_ask_for_number_or_time_on_need_hour_grain_from_day"}
				resp.Events = append(resp.Events, event)
			} else if rasaTime.Grain == "period" {
				event := Event{Event: FOLLOWUP, Name: "utter_ask_for_polar_is_pm"}
				resp.Events = append(resp.Events, event)
			} else {
				searchTimeStr := rasaTime.Value
				searchTime, err := time.Parse(time.RFC3339, searchTimeStr)

				if err != nil {
					log.Println(err)
				}

				MAX_TIME := float64(24 * 90)

				if searchTime.Before(time.Now()) {
					event := Event{Event: FOLLOWUP, Name: "utter_unhappy_time_in_past_AND_ask_for_time_on_alternative"}
					resp.Events = append(resp.Events, event)
				} else if searchTime.Sub(time.Now()).Hours() > MAX_TIME {
					event := Event{Event: FOLLOWUP, Name: "utter_unhappy_time_too_far_in_future_AND_ask_for_time_on_alternative"}
					resp.Events = append(resp.Events, event)
				} else {
					event := Event{Event: FOLLOWUP, Name: "action_blank_alert_potential_times_slot_set"}
					resp.Events = append(resp.Events, event)
					// So set it as the first item in temp_temps
					event = Event{Event: SLOT, Name: "potential_times", Value: v[0]}
					resp.Events = append(resp.Events, event)
				}

			}

		}
	default:
		event := Event{Event: FOLLOWUP, Name: "utter_ask_for_time_for_potential_reservation"}
		resp.Events = append(resp.Events, event)
	}
}

func (bot *Bot) ActionBrancherValidatePotentialHourSlot(req *RasaRequest, resp *RasaResponse) {
	temp_times := req.Tracker.Slots["temp_times"]
	potential_hour := req.Tracker.Slots["potential_hour"]

	event := Event{}
	if reflect.ValueOf(temp_times).IsNil() {
		event = Event{Event: FOLLOWUP, Name: "utter_ask_for_time_for_potential_reservation"}
	} else if reflect.ValueOf(potential_hour).IsNil() {
		event = Event{Event: FOLLOWUP, Name: "action_checkpoint_with_single_temp_times_to_fill_potential_times"}
	} else {
		switch v := potential_hour.(type) {
		case float64:
			if v >= 24 || v == 0 {
				event = Event{Event: FOLLOWUP, Name: "action_need_employee"}
			} else {
				if v >= 1 && v <= 11 {
					v += 12
				}
				event = Event{Event: SLOT, Name: "potential_hour", Value: v}
        event := Event{Event: FOLLOWUP, Name: "action_blank_alert_potential_hour_slot_set"}
        resp.Events = append(resp.Events, event)

			}
		default:
			event = Event{Event: FOLLOWUP, Name: "action_need_employee"}
		}
	}
	resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionBrancherValidateWithTempTimesAndTimeEntityToModifyTempTimeFromDayGrain(req *RasaRequest, resp *RasaResponse) {
	temp_times := req.Tracker.Slots["temp_times"]

	entities := req.Tracker.LatestMessage.Entities

	entityTime := Entity{}

	found := false
	for _, entity := range entities {
		if entity.Entity == "time" {
			entityTime = entity
			found = true
		}
	}

	event := Event{}
	if reflect.ValueOf(temp_times).IsNil() {
		event = Event{Event: FOLLOWUP, Name: "utter_ask_for_time_for_potential_reservation"}
	} else if !found {
		event = Event{Event: FOLLOWUP, Name: "action_need_employee"}
	} else {
		switch typed_times := temp_times.(type) {
		case []interface{}:
			switch time_map := typed_times[0].(type) {
			case map[string]interface{}:
				var rasaTime RasaTime
				err := mapstructure.Decode(time_map, &rasaTime)

				if err != nil {
					log.Println(err)
					event = Event{Event: FOLLOWUP, Name: "action_need_employee"}
					break
				}

				entityTimeMap := entityTime.Value.(map[string]interface{})
				entityTimeStr := entityTimeMap["value"].(string)
				entityTimeGrain := entityTimeMap["grain"].(string)

				entityTimeObj, _ := time.Parse(time.RFC3339, entityTimeStr)

				if entityTimeGrain == "week" || entityTimeGrain == "month" || entityTimeGrain == "year" {
					event = Event{Event: FOLLOWUP, Name: "action_need_employee"}
				} else {
					// Need to take the hour from the timeEntity and add it to the temp_time
					tempTimeStr := rasaTime.Value
					tempTime, _ := time.Parse(time.RFC3339, tempTimeStr)

					entityTimeHours := entityTimeObj.Hour()

					if entityTimeGrain == "period" {
						if entityTimeHours >= 1 && entityTimeHours <= 11 {
							entityTimeHours += 12
						}
					}

					newTime := time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), entityTimeHours, tempTime.Minute(), 0, 0, time.UTC)

					time_map["value"] = newTime.Format(time.RFC3339)

					event = Event{Event: SLOT, Name: "temp_times", Value: []interface{}{time_map}}
					event := Event{Event: FOLLOWUP, Name: "action_blank_alert_temp_times_slot_set"}
					resp.Events = append(resp.Events, event)

				}
			}
		}
	}
	resp.Events = append(resp.Events, event)
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
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == NUMBER {
			switch s := v.Value.(type) {
			case float64:
				nextAction := Event{Event: SLOT, Name: "potential_size", Value: s}
				resp.Events = append(resp.Events, nextAction)
			case int:
				nextAction := Event{Event: SLOT, Name: "potential_size", Value: float64(s)}
				resp.Events = append(resp.Events, nextAction)
			}
		}
	}

}

func (bot *Bot) ActionSetPotentialTimeSlot(req *RasaRequest, resp *RasaResponse) {
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == TIME {
			switch s := v.Value.(type) {
			case string:
				nextAction := Event{Event: SLOT, Name: "potential_time", Value: s}
				resp.Events = append(resp.Events, nextAction)
			}
		}
	}
}

func (bot *Bot) ActionClearPotentialTimeSlot(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "potential_time"}
	resp.Events = append(resp.Events, nextAction)
}

func (bot *Bot) ActionClearPotentialSizeSlot(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "potential_size"}
	resp.Events = append(resp.Events, nextAction)
}

func (bot *Bot) ActionBrancherValidateReservationPotentialSize(req *RasaRequest, resp *RasaResponse) {

	potential_size := 0
	found := false

	// Make sure potential_size exists
	switch v := req.Tracker.Slots["potential_size"].(type) {
	case int:
		potential_size = v
		found = true
	case float64:
		log.Println("Wanted int but got float")
		potential_size = int(v)
		found = true
	default:
		log.Printf("%T\n", v)
		event := Event{Event: FOLLOWUP, Name: "utter_ask_for_number_on_reservation_size"}
		resp.Events = append(resp.Events, event)
	}

	if found {
		if potential_size == 0 {
			event := Event{Event: FOLLOWUP, Name: "utter_unhappy_doing_invalid_size_AND_ask_for_size_greater_than_zero"}
			resp.Events = append(resp.Events, event)
		} else if potential_size > 20 {
			event := Event{Event: FOLLOWUP, Name: "utter_unhappy_doing_request_customer_call_for_large_parties"}
			resp.Events = append(resp.Events, event)
		} else {
			event := Event{Event: SLOT, Name: "size", Value: Abs(potential_size)}
			resp.Events = append(resp.Events, event)
      event = Event{Event: FOLLOWUP, Name: "action_blank_alert_size_slot_set"}
      resp.Events = append(resp.Events, event)

		}
	}

}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
