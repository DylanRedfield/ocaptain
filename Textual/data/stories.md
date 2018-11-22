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

