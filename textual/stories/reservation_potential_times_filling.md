## checkpoint_to_fill_temp_time_inform_number
> checkpoint_to_fill_temp_times
* inform{"number": 2}
    - action_set_potential_hour_slot
> checkpoint_with_potential_hour_to_fill_temp_times

## checkpoint_with_potential_hour_to_fill_temp_times
> checkpoint_with_potential_hour_to_fill_temp_times
    - action_set_temp_times_from_potential_hour
    - action_clear_potential_hour

> checkpoint_with_potential_hour
    - action_set_potential_size_slot
    - slot{"potential_size": 1}
> checkpoint_with_potential_size_to_validate_and_fill_size

## checkpoint_to_fill_potential_size_inform_non_number
> checkpoint_to_fill_potential_size
* inform{"name": "dylan"}
    - utter_unhappy_doing_no_number_recognized
> checkpoint_to_fill_potential_size

## checkpoint_to_fill_potential_size_inform_none
> checkpoint_to_fill_potential_size
* inform
    - utter_unhappy_doing_no_number_recognized
> checkpoint_to_fill_potential_size

