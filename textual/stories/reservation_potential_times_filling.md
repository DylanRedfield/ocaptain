## checkpoint_to_fill_temp_times_inform_time
    - action_blank_checkpoint_to_fill_temp_times
* inform{"time": "xxx"}
    - action_set_temp_times_slot
    - slot{"temp_times": "xxx"}
    - action_brancher_with_temp_times_to_determine_next_from_times_length

## checkpoint_to_fill_temp_times_inform_number
    - action_blank_checkpoint_to_fill_temp_times
* inform{"number": 8}
    - action_set_potential_hour_slot
    - slot{"potential_hour": 8}
    - action_set_temp_times_slot_from_potential_hour
    - slot{"temp_times": "xxx"}
    - action_clear_potential_hour_slot
    - slot{"potential_hour": null}
    - action_brancher_with_temp_times_to_determine_next_from_times_length

## action_brancher_with_temp_times_to_determine_next_from_times_length_0
    - action_brancher_with_temp_times_to_determine_next_from_times_length
    - followup{"name": "utter_ask_for_time_for_potential_reservation"}
    - utter_ask_for_time_for_potential_reservation
    - action_blank_checkpoint_to_fill_temp_times

## action_brancher_with_temp_times_to_determine_next_from_times_length_1
    - action_brancher_with_temp_times_to_determine_next_from_times_length
    - followup{"name": "action_brancher_with_temp_times_validate_single_temp_times"}
    - action_brancher_with_temp_times_validate_single_temp_times

## action_brancher_with_temp_times_validate_single_temp_times_0
    - action_brancher_with_temp_times_to_determine_next_from_times_length
    - followup{"name": "utter_ask_for_time_for_potential_reservation"}
    - utter_ask_for_time_for_potential_reservation
    - action_blank_checkpoint_to_fill_temp_times

## action_brancher_with_temp_times_validate_single_temp_times_1
    - action_brancher_with_temp_times_to_determine_next_from_times_length
    - followup{"name": "action_need_employee"}
    - action_need_employee

## action_brancher_with_temp_times_validate_single_temp_times_2
    - action_brancher_with_temp_times_to_determine_next_from_times_length
    - followup{"name": "utter_with_temp_time_ask_for_number_or_time_on_need_hour_grain_from_day"}
    - utter_with_temp_time_ask_for_number_or_time_on_need_hour_grain_from_day
> checkpoint_with_temp_times_from_reservation_day_grain

## action_brancher_with_temp_times_validate_single_temp_times_3
    - action_brancher_with_temp_times_to_determine_next_from_times_length
    - followup{"name": "utter_ask_for_polar_on_is_pm"}
    - utter_ask_for_polar_on_is_pm
> checkpoint_action_brancher_with_temp_times_validate_single_temp_times_period_ask_is_pm

# checkpoint_action_brancher_with_temp_times_validate_single_temp_times_period_ask_is_pm_0
> checkpoint_action_brancher_with_temp_times_validate_single_temp_times_period_ask_is_pm
* affirm
    - action_modify_temp_times_slot_pm
    - action_brancher_with_temp_times_validate_single_temp_times

# checkpoint_action_brancher_with_temp_times_validate_single_temp_times_period_ask_is_pm_1
> checkpoint_action_brancher_with_temp_times_validate_single_temp_times_period_ask_is_pm
* deny
    - action_modify_temp_times_slot_am
    - action_brancher_with_temp_times_validate_single_temp_times

# checkpoint_action_brancher_with_temp_times_validate_single_temp_times_period_ask_is_pm_2
> checkpoint_action_brancher_with_temp_times_validate_single_temp_times_period_ask_is_pm
* inform{"time": "xxx"}
    - action_brancher_validate_with_temp_times_and_time_entity_to_modify_temp_time_from_day_or_period_grain

# action_brancher_validate_with_temp_times_slot_and_time_entity_to_modify_temp_time_from_day_or_period_grain_0
    - action_brancher_validate_with_temp_times_and_time_entity_to_modify_temp_time_from_day_or_period_grain
    - followup{"name": "action_need_employee"}
    - action_need_employee

# action_brancher_validate_with_temp_times_slot_and_time_entity_to_modify_temp_time_from_day_or_period_grain_1
    - action_brancher_validate_with_temp_times_and_time_entity_to_modify_temp_time_from_day_or_period_grain
    - followup{"name": "action_blank_alert_temp_times_slot_set"}
    - action_blank_alert_temp_times_slot_set
    - slot{"temp_times": ["xxx"]}
    - action_need_employee

# checkpoint_action_brancher_with_temp_times_validate_single_temp_times_period_ask_is_pm_3
> checkpoint_action_brancher_with_temp_times_validate_single_temp_times_period_ask_is_pm
* inform
    - action_need_employee

## action_brancher_with_temp_times_validate_single_temp_times_4
    - action_brancher_with_temp_times_to_determine_next_from_times_length
    - followup{"name": "utter_unhappy_time_in_past_AND_ask_for_time_on_alternative"}
    - utter_unhappy_time_in_past_AND_ask_for_time_on_alternative
    - action_clear_temp_times_slot
    - action_blank_checkpoint_to_fill_temp_times

## action_brancher_with_temp_times_validate_single_temp_times_5
    - action_brancher_with_temp_times_to_determine_next_from_times_length
    - followup{"name": "utter_unhappy_time_too_far_in_future_AND_ask_for_time_for_alternative"}
    - utter_unhappy_time_too_far_in_future_AND_ask_for_time_for_alternative_potential_times
    - action_clear_temp_times_slot
    - action_blank_checkpoint_to_fill_temp_times

## action_brancher_with_temp_times_validate_single_temp_times_6
    - action_brancher_with_temp_times_to_determine_next_from_times_length
    - followup{"name": "action_blank_alert_potential_times_slot_set"}
    - action_blank_alert_potential_times_slot_set
    - slot{"potential_times": ["xxx"]}
    - action_clear_temp_times_slot
    - action_brancher_reservation_slot_filling_base

