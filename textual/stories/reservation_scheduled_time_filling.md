## action_brancher_with_size_and_single_potential_times_query_reservation_platform_01
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform
    - followup{"name": "utter_ask_for_number_on_reservation_size"}
    - utter_ask_for_number_on_reservation_size
> checkpoint_to_fill_potential_size

## action_brancher_with_size_and_single_potential_times_query_reservation_platform_02
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform
    - followup{"name": "utter_ask_for_time_for_potential_reservation"}
    - utter_ask_for_time_for_potential_reservation
> checkpoint_to_fill_temp_times

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

