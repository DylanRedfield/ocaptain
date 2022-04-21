package main

// Program
const (
	PROD_STATE                  string = "prod"
	DEV_STATE_LOCAL             string = "dev_local"
	DEV_STATE_ONLINE            string = "dev"
	FACEBOOK_MESSENGER_PLATFORM        = "FACEBOOK_MESSENGER_PLATFORM"
	TWILIO_PLATFORM                    = "TWILIO_PLATFORM"
)

// Business
const (
	PhoneNumber         string = "phoneNumber"
	Approved            string = "approved"
	Password            string = "password"
	FacebookMessengerId string = "facebookMessengerId"
)

// Recipient
const (
	Contact       string = "contact"
	RecentMessage        = "recentMessage"
	RecentOrderId        = "recentOrderId"
)

// Order
const (
	Type string = "type"
	RecipientId
	Contents = "contents"
	Name     = "name"
	Address  = "address"
)

// Collections
const (
	Businesses   string = "businesses"
	Recipients          = "recipients"
	Messages            = "messages"
	Orders              = "orders"
	Reservations        = "reservations"
)

// Actions
const (
	UTTER_GREET                                     string = "utter_greet"
	UTTER_GOODBYE                                          = "utter_goodbye"
	UTTER_YOUR_WELCOME                                     = "action_utter_respond_your_welcome_AND_ask_for_next_general_request_while_clearing_slots"
	UTTER_ASK_ADDRESS                                      = "utter_ask_address"
	UTTER_ASK_NAME                                         = "utter_ask_name"
	UTTER_THANK                                            = "utter_thank"
	UTTER_ASK_ORDER_CONTENTS                               = "utter_ask_order_contents"
	UTTER_ASK_CONFIRMATION_DELIVERY                        = "utter_ask_confirmation_delivery"
	UTTER_ASK_CONFIRMATION_PICK_UP                         = "utter_ask_confirmation_pick_up"
	UTTER_ASK_TYPE                                         = "utter_ask_type"
	UTTER_AFTER_ORDER                                      = "utter_after_order"
	UTTER_ASK_IS_ALL                                       = "utter_ask_is_all"
	ACTION_LISTEN                                          = "action_listen"
	ACTION_START_ORDER                                     = "action_start_order"
	ACTION_START_ORDER_WITH_INPUTS                         = "action_start_order_with_inputs"
	ACTION_SET_TYPE                                        = "action_set_type"
	ACTION_SET_ADDRESS                                     = "action_set_address"
	ACTION_SET_CONTENT                                     = "action_set_content"
	ACTION_SET_NAME                                        = "action_set_name"
	ACTION_CHECK_IS_OPEN                                   = "action_check_is_open"
	ACTION_CHECK_IS_OPEN_ON_DAY                            = "action_check_is_open_on_day"
	ACTION_CHECK_TIME_CLOSE                                = "action_check_time_close"
	ACTION_CHECK_TIME_CLOSE_ON_DAY                         = "action_check_time_close_on_DAY"
	ACTION_CHECK_RESERVATION_DATETIME                      = "action_check_reservation_datetime"
	ACTION_ASK_IF_SIMILAR_TIMES_WORK                       = "action_ask_if_any_similar_times_work"
	ACTION_UPDATE_ORDER                                    = "action_update_order"
	ACTION_RESET_SLOTS                                     = "action_reset_slots"
	ACTION_RESTART_SLOTS                                   = "action_restart"
	ACTION_SET_SCHEDULED_TIME_SLOT                         = "action_set_scheduled_time_slot"
	ACTION_SET_SIZE_SLOT                                   = "action_set_size_slot"
	ACTION_UTTER_ASK_IS_OTHER_RESERVATION_TIME_OKAY        = "action_utter_ask_is_other_reservation_time_okay"
	UTTER_REQUEST_TIME_TOO_EARLY                           = "utter_request_time_too_early"
	UTTER_NO_RESERVATIONS_AVAILABLE                        = "utter_no_reservations_available"
	ACTION_SAVE_RESERVATION                                = "action_save_reservation"
	ACTION_POST_RESERVATION_SAVED                          = "action_post_reservation_saved"
	ACTION_AFFIRM_SIMILAR_TIME_ORDINAL                     = "action_affirm_similar_time_ordinal"
	ACTION_AFFIRM_SIMILAR_TIME                             = "action_affirm_similar_time"
)

// Entities
const (
	TIME   = "time"
	NUMBER = "number"
)

// Slots
const (
	NAME              string = "name"
	SCHEDULED_TIME           = "scheduled_time"
	SIZE                     = "size"
	BUSINESS_ID              = "business_id"
	RECIPIENT_ID             = "recipient_id"
	RECIPIENT_CONTACT        = "recipient_contact"
	POTENTIAL_TIMES          = "potential_times"
)

// Event
const (
	SLOT     string = "slot"
	FOLLOWUP        = "followup"
)
