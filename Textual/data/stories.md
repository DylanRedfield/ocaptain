## greet
* greet
   - utter_greet

## Generated Story -5057898098459000086
* greet
    - utter_greet
* goodbye
    - utter_goodbye
* make_reservation
    - utter_ask_time
* inform{"time": "2018-11-22T18:00:00.000-08:00"}
    - action_set_scheduled_time_slot
    - slot{"scheduled_time": "2018-11-22T18:00:00.000-08:00"}
    - utter_ask_size
* inform{"number": 4}
    - action_set_size_slot
    - action_check_reservation_datetime

## test
* goodbye
    - utter_goodbye
## Generated Story 8573145562614629515
* make_reservation
    - utter_ask_time

## Generated Story -8705239742651640220
* make_reservation
    - utter_ask_time
* inform{"time": "2018-11-22T18:00:00.000-08:00"}
    - action_set_scheduled_time_slot
    - slot{"scheduled_time": "2018-11-22T18:00:00.000-08:00"}
    - utter_ask_size
* inform{"number": 4}
    - action_set_size_slot
    - slot{"size": 4}
    - action_check_reservation_datetime
## Generated Story -292295632353312686
* make_reservation{"number": 4, "time": "2018-11-23T18:30:00.000+00:00"}
    - action_set_scheduled_time_slot
    - slot{"scheduled_time": "2018-11-23T18:30:00.000+00:00"}
    - action_set_size_slot
    - slot{"size": "4.000000"}
    - action_check_reservation_datetime

## Generated Story 2483609456997642852
* make_reservation{"number": 4}
    - action_set_size_slot
    - slot{"size": "4.000000"}
    - utter_ask_time
* inform{"time": "2018-11-22T18:00:00.000+00:00"}
    - action_set_scheduled_time_slot
    - slot{"scheduled_time": "2018-11-22T18:00:00.000+00:00"}
    - action_check_reservation_datetime

## Generated Story -4263754831249927402
* inform{"time": "2018-11-24T18:00:00.000+00:00", "number": 4}
    - action_set_scheduled_time_slot
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "hxk7QAAgWVi47Qt05s7o"}
    - slot{"recipient_contact": "+19084771280"}
    - slot{"scheduled_time": "2018-11-24T18:00:00.000+00:00"}
    - action_set_size_slot
    - slot{"size": "4.000000"}
    - action_check_reservation_datetime
    - slot{"potential_times": ["2018-11-24T17:45:00Z"]}
    - followup{"name": "action_utter_ask_is_other_reservation_time_okay"}
    - action_utter_ask_is_other_reservation_time_okay
* affirm
    - utter_ask_name
* inform{"name": "dylan"}
    - slot{"name": "dylan"}
    - action_save_reservation
    - action_post_reservation_saved

## Generated Story -7876375454684602528
* make_reservation{"time": "2018-11-24T18:00:00.000+00:00", "number": 4}
    - action_set_scheduled_time_slot
    - slot{"business_id": "MewuHeThW4QJGDxD9tTr"}
    - slot{"recipient_id": "hxk7QAAgWVi47Qt05s7o"}
    - slot{"recipient_contact": "+19084771280"}
    - slot{"scheduled_time": "2018-11-24T18:00:00.000+00:00"}
    - action_set_size_slot
    - slot{"size": "4.000000"}
    - action_check_reservation_datetime
    - slot{"potential_times": ["2018-11-24T17:45:00Z"]}
    - followup{"name": "action_utter_ask_is_other_reservation_time_okay"}
    - action_utter_ask_is_other_reservation_time_okay
* affirm
    - utter_ask_name
    - action_save_reservation
    - action_post_reservation_saved

