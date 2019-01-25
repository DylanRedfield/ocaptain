## action_brancher_with_size_and_single_potential_times_query_reservation_platform_01
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform
    - followup{"name": "utter_ask_for_number_on_reservation_size"}
    - utter_ask_for_number_on_reservation_size
    - action_blank_checkpoint_to_fill_potential_size

## action_brancher_with_size_and_single_potential_times_query_reservation_platform_02
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform
    - followup{"name": "utter_ask_for_time_for_potential_reservation"}
    - utter_ask_for_time_for_potential_reservation
    - action_blank_checkpoint_to_fill_temp_times

## action_brancher_with_size_and_single_potential_times_query_reservation_platform_03
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform
    - followup{"name": "action_need_employee"}
    - action_need_employee

## action_brancher_with_size_and_single_potential_times_query_reservation_platform_04
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform
    - followup{"name": "action_blank_alert_scheduled_time_slot_set"}
    - action_blank_alert_scheduled_time_slot_set
    - slot{"scheduled_time": "xxx"}
    - action_clear_potential_times_slot
    - action_brancher_reservation_slot_filling_base

## action_brancher_with_size_and_single_potential_times_query_reservation_platform_05
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform
    - followup{"name": "action_blank_alert_alternative_times_slot_set"}
    - action_blank_alert_alternative_times_slot_set
    - slot{"alternative_times": ["xxx"]}
    - action_brancher_with_potential_times_and_alternative_times_to_fill_scheduled_time

## action_brancher_with_size_and_single_potential_times_query_reservation_platform_06
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform
    - followup{"name": "utter_doing_no_tables_available_near_that_time_AND_ask_for_polar_or_time_on_alternative"}
    - utter_doing_no_tables_available_near_that_time_AND_ask_for_polar_or_time_on_alternative
    - action_blank_checkpoint_with_size_and_no_alternatives

## action_brancher_with_size_and_single_potential_times_query_reservation_platform_07
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform
    - followup{"name": "utter_requested_time_too_soon_AND_ask_for_polar_or_time_on_alternative"}
    - utter_requested_time_too_soon_AND_ask_for_polar_or_time_on_alternative
    - action_blank_checkpoint_with_size_and_no_alternatives

## action_blank_checkpoint_with_size_and_no_alternatives_00
    - action_blank_checkpoint_with_size_and_no_alternatives
* deny
    - utter_ask_for_next_general_request
    - action_clear_size_slot
> checkpoint_general

## action_blank_checkpoint_with_size_and_no_alternatives_01
    - action_blank_checkpoint_with_size_and_no_alternatives
* affirm
    - utter_ask_for_time_for_potential_reservation
    - action_blank_checkpoint_to_fill_temp_times

## action_blank_checkpoint_with_size_and_no_alternatives_03
    - action_blank_checkpoint_with_size_and_no_alternatives
* inform{"time": "xxx"}
    - action_set_temp_times_slot
    - slot{"temp_times": "xxx"}
    - action_brancher_with_temp_times_to_determine_next_from_times_length

## action_blank_checkpoint_with_size_and_no_alternatives_03
    - action_blank_checkpoint_with_size_and_no_alternatives
* inform{"number": 8}
    - action_set_potential_hour_slot
    - slot{"potential_hour": 8}
    - action_set_temp_times_slot_from_potential_hour
    - slot{"temp_times": "xxx"}
    - action_clear_potential_hour_slot
    - slot{"potential_hour": null}
    - action_brancher_with_temp_times_to_determine_next_from_times_length

