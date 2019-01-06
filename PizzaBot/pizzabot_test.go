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

	if resp.Events[0].Event != FOLLOWUP || resp.Events[0].Name != "utter_ask_for_number_on_reservation_size" {
		t.Errorf("With no potential_number set as slot, should have asked for the reservation size")
	}

	// Test potential_size == 0
	req.Tracker.Slots["potential_size"] = 0

	if resp.Events[0].Event != FOLLOWUP && resp.Events[0].Name != "utter_unhappy_doing_invalid_size_AND_ask_for_size_greater_than_zero" {
		t.Errorf("With potential_size of zero, first follow up should be 'utter_unhappy_doing_invalid_size_AND_ask_for_size_greater_than_zero'")
	}

	// Test potential_size > 20
	req.Tracker.Slots["potential_size"] = 21

	if resp.Events[0].Name != "utter_unhappy_doing_request_customer_call_for_large_parties" {
		t.Errorf("Should have followed up with 'utter_unhappy_doing_request_customer_call_for_large_parties'")
	}

}
