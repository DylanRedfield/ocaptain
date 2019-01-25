## action_brancher_with_potential_times_and_alternative_times_to_fill_scheduled_times_0
    - utter_ask_for_time_for_potential_reservation
    - action_blank_checkpoint_to_fill_temp_times

## action_brancher_with_potential_times_and_alternative_times_to_fill_scheduled_times_1
    - action_brancher_with_size_and_single_potential_times_query_reservation_platform

## action_brancher_with_potential_times_and_alternative_times_to_fill_scheduled_times_2
    - action_utter_ask_for_polar_on_if_single_close_alternative_time_acceptable
> checkpoint_action_utter_ask_for_polar_on_if_single_close_alternative_time_acceptable

## checkpoint_action_utter_ask_for_polar_on_if_single_close_alternative_time_acceptable_0
> checkpoint_action_utter_ask_for_polar_on_if_single_close_alternative_time_acceptable
* affirm
    - action_set_scheduled_time_slot_from_close_alternative
    - slot{"scheduled_time": "xxx"}
    - action_clear_potential_times_slot
    - slot{"potential_times": null}
    - action_brancher_reservation_slot_filling_base

## checkpoint_action_utter_ask_for_polar_on_if_single_close_alternative_time_acceptable_1
> checkpoint_action_utter_ask_for_polar_on_if_single_close_alternative_time_acceptable
* deny
    - utter_ask_for_polar_or_time_on_alternative
    - action_clear_alternative_times_slot
    - slot{"alternative_times": null}
    - action_clear_potential_times_slot
    - slot{"potential_times": null}
    - action_blank_checkpoint_with_size_and_no_alternatives

## action_brancher_with_potential_times_and_alternative_times_to_fill_scheduled_times_3
    - action_utter_ask_for_polar_or_time_or_number_on_several_alternative_times
    - action_blank_checkpoint_with_several_alternative_times_to_select

## action_brancher_with_potential_times_and_alternative_times_to_fill_scheduled_times_4
    - action_utter_ask_for_polar_on_if_single_alternative_time_acceptable
> checkpoiint_action_utter_ask_for_polar_on_if_single_alternative_times_acceptable

## checkpoiint_action_utter_ask_for_polar_on_if_single_alternative_times_acceptable_0
> checkpoiint_action_utter_ask_for_polar_on_if_single_alternative_times_acceptable
* affirm
    - action_set_scheduled_time_from_single_alternative
    - slot{"scheduled_time": "xxx"}
    - action_clear_potential_times_slot
    - slot{"potential_times": null}
    - action_clear_alternative_times_slot
    - slot{"alternative_times": null}
    - action_brancher_reservation_slot_filling_base

## checkpoiint_action_utter_ask_for_polar_on_if_single_alternative_times_acceptable_1
> checkpoiint_action_utter_ask_for_polar_on_if_single_alternative_times_acceptable
* deny
    - utter_ask_for_polar_or_time_on_alternative
    - action_clear_alternative_times_slot
    - slot{"alternative_times": null}
    - action_clear_potential_times_slot
    - slot{"potential_times": null}
    - action_blank_checkpoint_with_size_and_no_alternatives

## action_blank_checkpoint_with_several_alternative_times_to_select_0
    - action_blank_checkpoint_with_several_alternative_times_to_select
* inform{"time": "xxx"}
    - action_set_temp_times_slot
    - slot{"temp_times": ["xxx"]}
    - action_brancher_validate_temp_time_to_select_alternative_time_to_set_scheduled_time_slot

## action_blank_checkpoint_with_several_alternative_times_to_select_1
    - action_blank_checkpoint_with_several_alternative_times_to_select
* inform{"ordinal": "x"} OR inform{"ordinal": "x", "time": "xxx"}
    - action_set_temp_ordinal_slot
    - slot{"temp_ordinal": "x"}
    - action_brancher_with_alternative_times_and_ordinal_validate_ordinal_to_select_alternative_time

## action_blank_checkpoint_with_several_alternative_times_to_select_2
    - action_blank_checkpoint_with_several_alternative_times_to_select
* inform{"number": 5}
    - action_set_potential_hour_slot
    - slot{"potential_hour": 5}
    - action_set_temp_times_slot_from_potential_hour
    - slot{"temp_times": ["xxx"]}
    - action_clear_potential_hour_slot
    - slot{"potential_hour": null}
    - action_brancher_validate_temp_time_to_select_alternative_time_to_set_scheduled_time_slot

## action_brancher_validate_temp_time_to_select_alternative_time_set_scheduled_time_slot_0
    - action_brancher_validate_temp_time_to_select_alternative_time_to_set_scheduled_time_slot
    - folloup{"name": "action_need_employee"}
    - action_need_employee

## action_brancher_validate_temp_time_to_select_alternative_time_set_scheduled_time_slot_1
    - action_brancher_validate_temp_time_to_select_alternative_time_to_set_scheduled_time_slot
    - folloup{"name": "action_blank_alert_scheduled_time_slot_set"}
    - slot{"scheduled_time": "xxx"}
    - action_blank_alert_scheduled_time_slot_set
    - action_brancher_reservation_slot_filling_base

## action_brancher_validate_temp_time_to_select_alternative_time_set_scheduled_time_slot_2
    - action_brancher_validate_temp_time_to_select_alternative_time_to_set_scheduled_time_slot
    - folloup{"name": "action_blank_alert_alternative_times_slot_set"}
    - slot{"alternative_times": ["xxx"]}
    - action_blank_alert_alternative_times_slot_set
    - action_utter_ask_with_alternative_times_for_time_or_number_or_ordinal_on_more_specific_alternative_time
    - action_clear_temp_times_slot
    - slot{"temp_times": null}
    - action_brancher_validate_temp_time_to_select_alternative_time_to_set_scheduled_time_slot











