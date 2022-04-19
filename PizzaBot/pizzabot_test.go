package main

import (
	"testing"
)

func TestActionBrancherValidateReservationsPotentialSize(t *testing.T) {
	slots := make(map[string]interface{})
	tracker := Tracker{Slots: slots}
	req := &RasaRequest{Tracker: tracker}
	resp := NewRasaResponse()

	// Test happy
	req.Tracker.Slots["potential_size"] = 1

	bot.ActionBrancherValidateReservationPotentialSize(req, resp)

	if resp.Events[0].Event != SLOT {
		t.Errorf("Happy path should have set the size slot")
	}

	// Test not set at all
	delete(req.Tracker.Slots, "potential_size")

	resp = NewRasaResponse()
	bot.ActionBrancherValidateReservationPotentialSize(req, resp)
	if resp.Events[0].Event != FOLLOWUP || resp.Events[0].Name != "utter_ask_for_number_on_reservation_size" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

	// Test potential_size == 0
	req.Tracker.Slots["potential_size"] = 0

	resp = NewRasaResponse()
	bot.ActionBrancherValidateReservationPotentialSize(req, resp)

	if resp.Events[0].Name != "utter_unhappy_doing_invalid_size_AND_ask_for_size_greater_than_zero" {

		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

	// Test potential_size > 20
	req.Tracker.Slots["potential_size"] = 21

	resp = NewRasaResponse()
	bot.ActionBrancherValidateReservationPotentialSize(req, resp)
	if resp.Events[0].Name != "utter_unhappy_doing_request_customer_call_for_large_parties" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

	// Test potential_size > 20
	req.Tracker.Slots["potential_size"] = .5

	resp = NewRasaResponse()
	bot.ActionBrancherValidateReservationPotentialSize(req, resp)
	if resp.Events[0].Name != "utter_unhappy_doing_invalid_size_AND_ask_for_size_greater_than_zero" {

		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

}

func TestActionBrancherWithTempTimesToDetermineNextFromMTimesLength(t *testing.T) {
	slots := make(map[string]interface{})
	tracker := Tracker{Slots: slots}
	req := &RasaRequest{Tracker: tracker}
	resp := NewRasaResponse()

	bot.ActionBrancherWithTempTimesToDetermineNextFromTimesLength(req, resp)

	if resp.Events[0].Name != "utter_ask_for_time_for_potential_reservation" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

	resp = NewRasaResponse()
	req.Tracker.Slots["temp_times"] = []interface{}{}
	req.Tracker.Slots["temp_times"] = append(req.Tracker.Slots["temp_times"].([]interface{}), "xxx")

	bot.ActionBrancherWithTempTimesToDetermineNextFromTimesLength(req, resp)
	if resp.Events[0].Name != "action_brancher_with_temp_times_validate_single_temp_times" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

}

func TestActionBrancherWithTempTimesValidateSingleTempTimes(t *testing.T) {
	slots := make(map[string]interface{})
	tracker := Tracker{Slots: slots}
	req := &RasaRequest{Tracker: tracker}
	resp := NewRasaResponse()

	tempTime := map[string]interface{}{}
	tempTime["value"] = "2019-04-01T00:00:00.000+00:00"
	tempTime["grain"] = "month"
	tempTime["type"] = "value"

	req.Tracker.Slots["temp_times"] = []map[string]interface{}{}

	req.Tracker.Slots["temp_times"] = append(req.Tracker.Slots["temp_times"].([]map[string]interface{}), tempTime)

	bot.ActionBrancherWithTempTimesValidateSingleTempTimes(req, resp)
	if resp.Events[0].Name != "action_need_employee" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

	resp = NewRasaResponse()

	tempTime["grain"] = "day"

	bot.ActionBrancherWithTempTimesValidateSingleTempTimes(req, resp)
	if resp.Events[0].Name != "utter_with_temp_time_ask_for_number_or_time_on_need_hour_grain_from_day" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

	resp = NewRasaResponse()

	tempTime["grain"] = "period"

	bot.ActionBrancherWithTempTimesValidateSingleTempTimes(req, resp)
	if resp.Events[0].Name != "utter_ask_for_polar_on_is_pm" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

	resp = NewRasaResponse()
	// In past
	tempTime["grain"] = "hour"
	tempTime["value"] = "2019-02-01T00:00:00.000+00:00"

	bot.ActionBrancherWithTempTimesValidateSingleTempTimes(req, resp)
	if resp.Events[0].Name != "utter_unhappy_time_in_past_AND_ask_for_time_on_alternative" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

	tempTime["grain"] = "minute"
	tempTime["value"] = "2020-01-01T00:00:00.000+00:00"
	resp = NewRasaResponse()

	bot.ActionBrancherWithTempTimesValidateSingleTempTimes(req, resp)
	if resp.Events[0].Name != "utter_unhappy_time_too_far_in_future_AND_ask_for_time_on_alternative" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

	tempTime["grain"] = "minute"
	tempTime["value"] = "2019-04-27T00:00:00.000+00:00"
	resp = NewRasaResponse()

	bot.ActionBrancherWithTempTimesValidateSingleTempTimes(req, resp)
	if resp.Events[0].Name != "action_blank_alert_potential_times_slot_set" {
		t.Errorf("followed with: %s", resp.Events[0].Name)
	}

}
func TestActionBrancherValidatePotentialHourSlot(t *testing.T) {
}
func TestActionBrancherValidateWithTempTimesAndSinglePotentialTimesQueryReservationPlatform(t *testing.T) {
}
func TestActionBrancherWithPotentialTimesAndAlternativeTimesToFillScheduledTime(t *testing.T) {
}
func TestActionBrancherValidateTempTimeToSelectAlternativeTimeToSetScheduledTimeSlot(t *testing.T) {
}
func TestActionBrancherWithAlternativeTimesAndOrdinalValidateOrdinalToSelectAlternativeTime(t *testing.T) {
}
func TestActionBrancherToSaveNewReservation(t *testing.T) {
}
func TestActionBrancherReservationSlotFillingBase(t *testing.T) {
}
