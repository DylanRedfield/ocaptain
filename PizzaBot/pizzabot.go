package main

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"
)

func (bot *Bot) HandleAction(req *RasaRequest) (*RasaResponse, error) {
	log.Println(req.NextAction)
	resp := NewRasaResponse()

	// TODO remove this it is just for train online testing
	bot.checkOrSetInputSlots(req, resp)

	action := req.NextAction

	log.Println(action)
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
	case "action_clear_name_slot":
		bot.ActionClearNameSlot(req, resp)
	case "action_set_potential_time_slot":
		bot.ActionSetPotentialTimeSlot(req, resp)
	case "action_clear_potential_times_slot":
		bot.ActionClearPotentialTimesSlot(req, resp)
	case "action_clear_temp_ordinal_slot":
		bot.ActionClearTempOrdinalSlot(req, resp)
	case "action_clear_temp_times_slot":
		bot.ActionClearTempTimeSlot(req, resp)
	case "action_clear_scheduled_time_slot":
		bot.ActionClearScheduledTimeSlot(req, resp)
	case "action_test_bed":
		bot.ActionTestBed(req, resp)
	case "action_set_temp_times_slot":
		bot.ActionSetTempTimesSlot(req, resp)
	case "action_set_temp_times_slot_from_potential_hour":
		bot.ActionSetTempTimesSlotFromPotentialHour(req, resp)
	case "action_set_potential_hour_slot":
		bot.ActionSetPotentialHourSlot(req, resp)
	case "action_brancher_with_temp_times_to_determine_next_from_times_length":
		bot.ActionBrancherWithTempTimesToDetermineNextFromTimesLength(req, resp)
	case "action_modify_temp_times_slot_pm":
		bot.ActionModifyTempTimesSlotPm(req, resp)
	case "action_modify_temp_times_slot_am":
		bot.ActionModifyTempTimesSlotAm(req, resp)
	case "action_brancher_with_temp_times_validate_single_temp_times":
		bot.ActionBrancherWithTempTimesValidateSingleTempTimes(req, resp)
	case "action_brancher_validate_potential_hour_slot":
		bot.ActionBrancherValidatePotentialHourSlot(req, resp)
	case "action_brancher_validate_with_temp_times_and_time_entity_to_modify_temp_times_from_day_or_period_grain":
		bot.ActionBrancherValidateWithTempTimesAndTimeEntityToModifyTempTimeFromDayOrPeriodGrain(req, resp)
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
	case "action_utter_ask_for_polar_on_if_single_close_alternative_time_acceptable":
		bot.ActionUtterAskForPolarOnIfSingleCloseAlternativeTimeAcceptable(req, resp)
	case "action_utter_ask_for_polar_or_time_or_number_on_serveral_alternative_times":
		bot.ActionUtterAskForPolarOrTimeOrNumberOnSeveralAlternativeTimes(req, resp)
	case "action_utter_ask_with_alternative_times_for_time_or_number_or_ordinal_on_more_specific_alternative_time":
		bot.ActionUtterAskWithAlternativeTimesForTimeOrNumberOrOrdinalOnMoreSpecificAlternativeTime(req, resp)
	case "action_need_employee_because_error":
		bot.ActionNeedEmployeeBecauseError(req, resp)
	case "action_brancher_with_alternative_times_and_ordinal_validate_ordinal_to_select_alternative_time":
		bot.ActionBrancherWithAlternativeTimesAndOrdinalValidateOrdinalToSelectAlternativeTime(req, resp)
	case "action_brancher_validate_temp_time_to_select_alternative_time_to_set_scheduled_time_slot":
		bot.ActionBrancherValidateTempTimeToSelectAlternativeTimeToSetScheduledTimeSlot(req, resp)
	case "action_clear_alternative_times_slot":
		bot.ActionClearAlternativeTimes(req, resp)
	case "action_set_temp_ordinal_slot":
		bot.ActionSetTempOrdinalSlot(req, resp)
	case "action_utter_post_reservation_save_AND_ask_for_next_general_request":
		bot.ActionUtterPostReservationAndAskForNextGeneralRequest(req, resp)
	case "action_utter_answer_time":
		bot.ActionUtterAnswerTime(req, resp)
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
	log.Println("base")
	slots := req.Tracker.Slots

	size := slots["size"]
	potentialSize := slots["potential_size"]

	scheduledTime := slots["scheduled_time"]
	tempTimes := slots["temp_times"]
	potentialTimes := slots["potential_times"]

	name := slots["name"]

	event := Event{}
	//  zero := reflect.Zero(reflect.TypeOf(size))

	event.Event = FOLLOWUP
	if reflect.TypeOf(size) == nil {
		log.Println("size_base")
		if reflect.TypeOf(potentialSize) == nil {
			event.Name = "utter_ask_for_number_on_reservation_size"
		} else {
			event.Name = "action_brancher_validate_reservation_potential_size"
		}
	} else if reflect.TypeOf(scheduledTime) == nil {

		if reflect.TypeOf(potentialTimes) == nil {

			if reflect.TypeOf(tempTimes) == nil {
				event.Name = "utter_ask_for_time_for_potential_reservation"
			} else {
				event.Name = "action_brancher_with_temp_times_to_determine_next_from_times_length"
			}
		} else {
			event.Name = "action_brancher_with_size_and_single_potential_times_query_reservation_platform"
		}
	} else if reflect.TypeOf(name) == nil {
		log.Println("name_base")
		event.Name = "utter_ask_for_name_on_reservation"
	} else {
		log.Println("save_base")
		event.Name = "action_brancher_to_save_new_reservation"
	}

	resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionBrancherWithAlternativeTimesAndOrdinalValidateOrdinalToSelectAlternativeTime(req *RasaRequest, resp *RasaResponse) {
	slots := req.Tracker.Slots

	rawOrdinal := slots["temp_ordinal"]
	rawAlternativeTimes := slots["alternative_times"]

	event := Event{Event: FOLLOWUP}

	if reflect.TypeOf(rawAlternativeTimes) == nil {
		event.Name = "action_need_employee"
	} else if reflect.TypeOf(rawOrdinal) == nil {
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
			event.Value = alternativeTimesStr[ordinalValue-1]
			resp.Events = append(resp.Events, event)

			event = Event{Event: FOLLOWUP, Name: "action_blank_alert_scheduled_time_slot_set"}

		}
	}

	resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionBrancherValidateTempTimeToSelectAlternativeTimeToSetScheduledTimeSlot(req *RasaRequest, resp *RasaResponse) {
	slots := req.Tracker.Slots

	rawTempTimes := slots["temp_times"]
	rawAltTimes := slots["alternative_times"]

	event := Event{Event: FOLLOWUP}

	if reflect.TypeOf(rawTempTimes) == nil {
		event.Name = "action_need_employee"
	} else if reflect.TypeOf(rawAltTimes) == nil {
		event.Name = "action_need_employee"
	} else {
		tempTime := rawTempTimes.([]interface{})[0]

		var rasaTime RasaTime
		err := mapstructure.Decode(tempTime.(map[string]interface{}), &rasaTime)

		if err != nil {
			log.Println(err)
			event = Event{Event: FOLLOWUP, Name: "action_need_employee_because_error"}
			resp.Events = append(resp.Events, event)
			return
		}

		tempTimeObj, err := time.Parse(time.RFC3339, rasaTime.Value)
		if err != nil {
			log.Println(err)
			event = Event{Event: FOLLOWUP, Name: "action_need_employee_because_error"}
			resp.Events = append(resp.Events, event)
			return
		}

		tempHour := tempTimeObj.Hour()
		tempMinutes := tempTimeObj.Minute()

		// Assume PM
		if rasaTime.Grain == "period" {
			if tempHour >= 0 && tempHour <= 11 {
				tempHour += 12
			}
		}

		// convert alt time strings into time objects
		altTimesStr := rawAltTimes.([]string)
		altTimes := []time.Time{}

		for _, v := range altTimesStr {
			altTime, err := time.Parse(time.RFC3339, v)
			if err != nil {
				log.Println(err)
				event = Event{Event: FOLLOWUP, Name: "action_need_employee_because_error"}
				resp.Events = append(resp.Events, event)
				return
			}

			altTimes = append(altTimes, altTime)

		}

		minuteMatch := false
		hourMatches := 0
		index := -1
		// Try to match
		for i, v := range altTimes {
			if tempHour == v.Hour() {
				index = i
				if tempMinutes == v.Minute() {
					minuteMatch = true
					break
				} else {
					hourMatches++
				}
			}
		}

		if minuteMatch {
			event.Event = SLOT
			event.Name = "scheduled_time"
			event.Value = altTimesStr[index]
			resp.Events = append(resp.Events, event)
			event = Event{Event: FOLLOWUP, Name: "action_blank_alert_scheduled_time_set"}
		} else {
			if hourMatches == 0 {
				event.Name = "action_blank_alert_employee"
			} else if hourMatches == 1 {
				event.Event = SLOT
				event.Name = "scheduled_time"
				event.Value = altTimesStr[index]
				resp.Events = append(resp.Events, event)
				event = Event{Event: FOLLOWUP, Name: "action_blank_alert_scheduled_time_set"}
			} else {
				newAltTimesStr := []string{}

				for i, v := range altTimes {
					if tempHour == v.Hour() {
						newAltTimesStr = append(newAltTimesStr, altTimesStr[i])
					}
				}

				event.Event = SLOT
				event.Name = "alternative_times"
				event.Value = altTimesStr
				resp.Events = append(resp.Events, event)
				event = Event{Event: FOLLOWUP, Name: "action_blank_alert_alternative_times_set"}
			}
		}

	}

	resp.Events = append(resp.Events, event)
}

func (bot *Bot) ActionBrancherWithPotentialTimesAndAlterativeTimesToFillScheduledTime(req *RasaRequest, resp *RasaResponse) {
	slots := req.Tracker.Slots

	rawPotentialTimes := slots["potential_times"]
	rawAlternativeTimes := slots["alternative_times"]

	event := Event{Event: FOLLOWUP}

	if reflect.TypeOf(rawPotentialTimes) == nil {
		event.Name = "utter_ask_for_time_on_potential_reservation"
	} else if reflect.TypeOf(rawAlternativeTimes) == nil {
		event.Name = "action_brancher_with_size_and_single_potential_times_query_reservation_platform"
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

		// TODO IDK WHY I COMMENT THIS OUT IT JUST WASN"T COMPILING
		/*if name == "" {
			// Action ask the name
			nextAction := Event{Event: "followup", Name: UTTER_ASK_NAME}
			resp.Events = append(resp.Events, nextAction)
			log.Println("Test")
			return
		} else {
			// Force Action save_reservation
			nextAction := Event{Event: "followup", Name: ACTION_SAVE_RESERVATION}
			resp.Events = append(resp.Events, nextAction)
		}
		altTimesArr := rawAlternativeTimes.([]interface{})
		for _, v := range altTimesArr {
			altTime, err := time.Parse(time.RFC3339, v.(string))

			if err != nil {
				event.Name = "action_need_employee_because_error"

				resp.Events = append(resp.Events, event)
				return
			}

			altTimes = append(altTimes, altTime)
		}*/

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

	resp.Events = append(resp.Events, event)
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
	if reflect.TypeOf(rawSize) == nil || reflect.TypeOf(rawName) == nil || reflect.TypeOf(rawScheduledTime) == nil {
		event.Event = FOLLOWUP
		event.Name = "action_brancher_reservation_slot_filling"
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

	if reflect.TypeOf(rawSize) == nil {
		event.Name = "utter_ask_for_number_on_reservation_size"
	} else if reflect.TypeOf(rawPotentialTimes) == nil {
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
				for _, result := range reservationResult.Results {
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
					event.Name = "action_blank_alert_scheduled_time_slot_set"
				} else {
					event.Name = "action_blank_alert_alternative_times_slot_set"
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

	newTime := time.Date(now.Year(), now.Month(), now.Day(), int(potentialHour), 0, 0, 0, time.Local)
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

func (bot *Bot) ActionUtterPostReservationAndAskForNextGeneralRequest(req *RasaRequest, resp *RasaResponse) {
	//size := req.Tracker.Slots[SIZE].(string)
	scheduledTime := req.Tracker.Slots[SCHEDULED_TIME].(string)

	datetime, _ := time.Parse(time.RFC3339, scheduledTime)

	reply := fmt.Sprintf("Great, you're all set. We'll see you at %s. Is there anything else we can do for you?", datetime.Format("3:04 PM"))
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
	for _, v := range entities {
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

	for _, v := range entities {
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
	rawTempTimes := req.Tracker.Slots["temp_times"]

	if reflect.TypeOf(rawTempTimes) == nil {
		event := Event{Event: FOLLOWUP, Name: "utter_ask_for_time_for_potential_reservation"}
		resp.Events = append(resp.Events, event)
	} else {
		tempTimes := rawTempTimes.([]interface{})
		if len(tempTimes) == 0 {
			event := Event{Event: FOLLOWUP, Name: "utter_ask_for_time_for_potential_reservation"}
			resp.Events = append(resp.Events, event)
		} else if len(tempTimes) == 1 {
			event := Event{Event: FOLLOWUP, Name: "action_brancher_with_temp_times_validate_single_temp_times"}
			resp.Events = append(resp.Events, event)
		} else {
			event := Event{Event: FOLLOWUP, Name: "action_brancher_with_temp_times_validate_single_temp_times"}
			resp.Events = append(resp.Events, event)
		}
	}
}

func (bot *Bot) ActionModifyTempTimesSlotAm(req *RasaRequest, resp *RasaResponse) {
	rawTempTime := req.Tracker.Slots["temp_times"].([]interface{})[0]

	time_map := rawTempTime.(map[string]interface{})
	var rasaTime RasaTime
	err := mapstructure.Decode(time_map, &rasaTime)

	if err != nil {
		log.Println(err)
		return
	}

	timeObj, err := time.Parse(time.RFC3339, rasaTime.Value)

	if err != nil {
		log.Println(err)
		return
	}

	hour := timeObj.Hour()

	if hour >= 12 {
		hour = hour - 12
	}

	newTime := time.Date(timeObj.Year(), timeObj.Month(), timeObj.Day(), hour, timeObj.Minute(), 0, 0, time.Local)

	rasaTime.Value = newTime.Format(time.RFC3339)

	event := Event{Event: SLOT, Name: "temp_times", Value: []RasaTime{rasaTime}}
	resp.Events = append(resp.Events, event)

}

func (bot *Bot) ActionModifyTempTimesSlotPm(req *RasaRequest, resp *RasaResponse) {
	rawTempTime := req.Tracker.Slots["temp_times"].([]interface{})[0]

	time_map := rawTempTime.(map[string]interface{})
	var rasaTime RasaTime
	err := mapstructure.Decode(time_map, &rasaTime)

	if err != nil {
		log.Println(err)
		return
	}

	timeObj, err := time.Parse(time.RFC3339, rasaTime.Value)

	if err != nil {
		log.Println(err)
		return
	}

	hour := timeObj.Hour()

	if hour >= 0 && hour <= 11 {
		hour = hour + 12
	}

	newTime := time.Date(timeObj.Year(), timeObj.Month(), timeObj.Day(), hour, timeObj.Minute(), 0, 0, time.Local)

	rasaTime.Value = newTime.Format(time.RFC3339)

	event := Event{Event: SLOT, Name: "temp_times", Value: []RasaTime{rasaTime}}
	resp.Events = append(resp.Events, event)

}

func (bot *Bot) ActionBrancherWithTempTimesValidateSingleTempTimes(req *RasaRequest, resp *RasaResponse) {
	rawTempTimes := req.Tracker.Slots["temp_times"]

	if reflect.TypeOf(rawTempTimes) == nil {
		event := Event{Event: FOLLOWUP, Name: "utter_ask_for_time_for_potential_reservation"}
		resp.Events = append(resp.Events, event)
	} else {
		time_map := rawTempTimes.([]interface{})[0]

		var rasaTime RasaTime
		err := mapstructure.Decode(time_map, &rasaTime)

		if err != nil {
			log.Println(err)
			// TODO
			return
		}

		if rasaTime.Grain == "week" || rasaTime.Grain == "month" || rasaTime.Grain == "year" {
			event := Event{Event: FOLLOWUP, Name: "action_need_employee"}
			resp.Events = append(resp.Events, event)
		} else if rasaTime.Grain == "day" {
			event := Event{Event: FOLLOWUP, Name: "utter_with_temp_time_ask_for_number_or_time_on_need_hour_grain_from_day"}
			resp.Events = append(resp.Events, event)
		} else if rasaTime.Grain == "period" {
			event := Event{Event: FOLLOWUP, Name: "utter_ask_for_polar_on_is_pm"}
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
				// So set it as the first item in temp_temps
				event := Event{Event: SLOT, Name: "potential_times", Value: []string{rasaTime.Value}}
				resp.Events = append(resp.Events, event)

				event = Event{Event: FOLLOWUP, Name: "action_blank_alert_potential_times_slot_set"}
				resp.Events = append(resp.Events, event)
			}

		}
	}
	/*switch v := req.Tracker.Slots["temp_times"].(type) {
	case []interface{}:
		switch time_map := v[0].(type) {
		case map[string]interface{}:

			}

		}
	default:
		event := Event{Event: FOLLOWUP, Name: "utter_ask_for_time_for_potential_reservation"}
		resp.Events = append(resp.Events, event)
	}*/
}

func (bot *Bot) ActionBrancherValidatePotentialHourSlot(req *RasaRequest, resp *RasaResponse) {
	temp_times := req.Tracker.Slots["temp_times"]
	potential_hour := req.Tracker.Slots["potential_hour"]

	event := Event{}
	if reflect.TypeOf(temp_times) == nil {
		event = Event{Event: FOLLOWUP, Name: "utter_ask_for_time_for_potential_reservation"}
	} else if reflect.TypeOf(potential_hour) == nil {
		event = Event{Event: FOLLOWUP, Name: "action_brancher_with_temp_times_validate_single_temp_times"}
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

func (bot *Bot) ActionBrancherValidateWithTempTimesAndTimeEntityToModifyTempTimeFromDayOrPeriodGrain(req *RasaRequest, resp *RasaResponse) {
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
	if reflect.TypeOf(temp_times) == nil {
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

					newTime := time.Date(tempTime.Year(), tempTime.Month(), tempTime.Day(), entityTimeHours, tempTime.Minute(), 0, 0, time.Local)

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

func (bot *Bot) ActionUtterAskForPolarOrTimeOrNumberOnSeveralAlternativeTimes(req *RasaRequest, resp *RasaResponse) {
	times := req.Tracker.Slots["alternative_times"].([]interface{})
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

func (bot *Bot) ActionUtterAskWithAlternativeTimesForTimeOrNumberOrOrdinalOnMoreSpecificAlternativeTime(req *RasaRequest, resp *RasaResponse) {
	times := req.Tracker.Slots["alternative_times"].([]interface{})
	reply := "Sorry which of these two: "

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

func (bot *Bot) ActionUtterAskForPolarOnIfSingleAlternativeTimeAcceptable(req *RasaRequest, resp *RasaResponse) {
	rawAltTime := req.Tracker.Slots["alternative_times"].([]interface{})[0].(string)

	datetime, err := time.Parse(time.RFC3339, rawAltTime)

	if err != nil {
		log.Println(err)
	}

	altTime := datetime.Format("3:04 PM")

	reply := fmt.Sprintf("The only close time is %s. Does that work?", altTime)
	resp.Responses = append(resp.Responses, Response{Text: reply})

}
func (bot *Bot) ActionUtterAskForPolarOnIfSingleCloseAlternativeTimeAcceptable(req *RasaRequest, resp *RasaResponse) {
	rawAltTime := req.Tracker.Slots["alternative_times"].([]interface{})[0].(string)

	datetime, err := time.Parse(time.RFC3339, rawAltTime)

	if err != nil {
		log.Println(err)
	}

	altTime := datetime.Format("3:04 PM")

	reply := fmt.Sprintf("Is %s close enough?", altTime)
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
				nextAction := Event{Event: SLOT, Name: "potential_times", Value: s}
				resp.Events = append(resp.Events, nextAction)
			}
		}
	}
}

func (bot *Bot) ActionSetTempOrdinalSlot(req *RasaRequest, resp *RasaResponse) {
	ordinal := ""
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == "ordinal" {
			ordinal = v.Value.(string)
		}
	}

	nextAction := Event{Event: SLOT, Name: "temp_ordinal", Value: ordinal}
	resp.Events = append(resp.Events, nextAction)

}

func (bot *Bot) ActionClearPotentialTimesSlot(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "potential_times"}
	resp.Events = append(resp.Events, nextAction)
}
func (bot *Bot) ActionClearNameSlot(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "name"}
	resp.Events = append(resp.Events, nextAction)
}
func (bot *Bot) ActionClearPotentialSizeSlot(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "potential_size"}
	resp.Events = append(resp.Events, nextAction)
}

func (bot *Bot) ActionClearTempTimeSlot(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "temp_times"}
	resp.Events = append(resp.Events, nextAction)
}
func (bot *Bot) ActionClearTempOrdinalSlot(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "temp_ordinal"}
	resp.Events = append(resp.Events, nextAction)
}

func (bot *Bot) ActionClearAlternativeTimes(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "alternative_times"}
	resp.Events = append(resp.Events, nextAction)
}

func (bot *Bot) ActionClearScheduledTimeSlot(req *RasaRequest, resp *RasaResponse) {
	nextAction := Event{Event: SLOT, Name: "scheduled_time"}
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
		potential_size = int(v)
		found = true
	default:
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

func (bot *Bot) ActionUtterAnswerTime(req *RasaRequest, resp *RasaResponse) {
	// “ We are currently [open | closed] and will be [closed | open](requested) [from xx-xx](if requested is open that day)”
	// Thus I need to see if they are currently open, check if there is an entiy (and if not se the day to today) and
	// check if they are open on that day then print in the above format
	businessId := req.Tracker.Slots["business_id"].(string)
	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		event := Event{Event: FOLLOWUP, Name: "action_need_employee"}
		resp.Events = append(resp.Events, event)
		return
	}

	// First logic to determine if the customer is asking about today by checking for entities and if that entity is today
	var entity Entity
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == "time" {
			entity = v
			break
		}
	}

	requestTime := time.Now()
	requestedToday := true
	if entity.Entity == "time" {
		// time entity was input
		entityTime, err := time.Parse(time.RFC3339, entity.Value.(string))

		if err != nil {
			event := Event{Event: FOLLOWUP, Name: "action_need_employee"}
			resp.Events = append(resp.Events, event)
			return
		}

		if entityTime.Year() != time.Now().Year() || entityTime.Month() != time.Now().Month() || entityTime.Day() == time.Now().Day() {
			// Time is not today
			requestTime = entityTime
			requestedToday = false
		}
	}

	isOpen := business.IsOpenOnDay(requestTime)
	if requestedToday {

		if isOpen {
			// get the open close strings to print
			openString, err := business.TimeOpenOnDayString(requestTime)
			closeString, err2 := business.TimeCloseOnDayString(requestTime)

			log.Println("Open")
			if err != nil || err2 != nil {
				event := Event{Event: FOLLOWUP, Name: "action_need_employee"}
				resp.Events = append(resp.Events, event)
				return
			}

			reply := fmt.Sprintf("We are currently open and will be %s - %s", openString, closeString)
			resp.Responses = append(resp.Responses, Response{Text: reply})
			return
		} else {
			// Find next available open day and times

			nextOpenDay := business.GetNextOpenDayAfter(requestTime)

			openString, err := business.TimeOpenOnDayString(nextOpenDay)
			closeString, err2 := business.TimeCloseOnDayString(nextOpenDay)

			if err != nil || err2 != nil {
				event := Event{Event: FOLLOWUP, Name: "action_need_employee"}
				resp.Events = append(resp.Events, event)
				return
			}

			dayOfWeek := nextOpenDay.Weekday().String()

			// We are currently closed but will be open Thursday (5/25) from 8:00am - 8:00pm
			reply := fmt.Sprintf("We are currently closed but will be open %s (%d/%d) from %s - %s",
				dayOfWeek, nextOpenDay.Month(), nextOpenDay.Day(), openString, closeString)

			resp.Responses = append(resp.Responses, Response{Text: reply})
			return

		}

	} else {

		dayOfWeek := requestTime.Weekday().String()
		openString, err := business.TimeOpenOnDayString(requestTime)
		closeString, err2 := business.TimeCloseOnDayString(requestTime)

		if err != nil || err2 != nil {
			event := Event{Event: FOLLOWUP, Name: "action_need_employee"}
			resp.Events = append(resp.Events, event)
			return
		}

		if isOpen {

			// "On Thursday (5/25) we'll be open from 8:00am - 8:00pm"
			reply := fmt.Sprintf("On %s (%d/%d) we'll be open from %s - %s", dayOfWeek, requestTime.Month(), requestTime.Day(), openString, closeString)

			resp.Responses = append(resp.Responses, Response{Text: reply})
			return

		} else {
			// "Sorry we'll be closed on Thursday (5/25)"
			reply := fmt.Sprintf("Sorry we'll be closed on %s (%d/%d)", dayOfWeek, requestTime.Month(), requestTime.Day())

			resp.Responses = append(resp.Responses, Response{Text: reply})
			return

		}
	}

}

func (bot *Bot) saveOrder(req *RasaRequest, order *Order) {
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
