## greet
* greet
    - utter_greet

## thank
* thank
    - utter_thank

## goodbye
* goodbye
    - utter_goodbye

## need_reservation
* need_reservation
    - utter_need_reservation

## allow_reservations
* allow_reservations
    - utter_allow_reservations_and_ask

## order
* order
    - action_alert_business

## time_close
* time_close
    - action_check_time_close
* inform{"time": "xxx"}
    - action_check_time_close_on_day

## time_close with param
* time_close{"time":"xxxx"}
    - action_check_time_close_on_day
* inform{"time": "xxx"}
    - action_check_time_close_on_day

## time_open
* time_open
    - action_check_time_open
* inform{"time": "xxx"}
    - action_check_time_open_on_day

## time_open with param
* time_open{"time": "xxxxx"}
    - action_check_time_open_on_day
* inform{"time": "xxx"}
    - action_check_time_open_on_day

## is_open
* is_open
    - action_check_is_open
* inform{"time": "xxx"}
    - action_check_is_open_on_day

## is open with param
* is_open{"time": "2018-xxxxxx"}
    - action_check_is_open_on_day
* inform{"time": "xxx"}
    - action_check_is_open_on_day
