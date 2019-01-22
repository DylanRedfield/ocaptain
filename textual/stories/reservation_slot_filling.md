## checkpoint_reservation_slot_filling_base
> checkpoint_reservation_slot_filling
    - action_brancher_reservation_slot_filling_base
> checkpoint_reservation_slot_filling_base_return

## checkpoint_reservation_slot_filling_base_return_1
> checkpoint_reservation_slot_filling_base_return
    - followup{"name": "utter_ask_for_number_on_reservation_size"}
    - utter_ask_for_number_on_reservation_size
> checkpoint_to_fill_potential_size

## checkpoint_reservation_slot_filling_base_return_2
> checkpoint_reservation_slot_filling_base_return
    - followup{"name": "action_checkpoint_with_potential_size_to_validate_and_fill_size"}
    - action_checkpoint_with_potential_size_to_validate_and_fill_size
> checkpoint_with_potential_size_to_validate_and_fill_size

## checkpoint_reservation_slot_filling_base_return_3
> checkpoint_reservation_slot_filling_base_return
    - followup{"name": "action_checkpoint_with_size_and_single_potential_times_to_fill_scheduled_time"}
    - action_checkpoint_with_size_and_single_potential_times_to_fill_scheduled_time
> checkpoint_with_size_and_single_potential_times_to_fill_scheduled_time

## checkpoint_reservation_slot_filling_base_return_4
> checkpoint_reservation_slot_filling_base_return
    - followup{"name": "utter_ask_for_time_for_potential_reservation"}
    - utter_ask_for_time_for_potential_reservation
> checkpoint_to_fill_temp_times

## checkpoint_reservation_slot_filling_base_return_5
> checkpoint_reservation_slot_filling_base_return
    - followup{"name": "action_checkpoint_with_temp_times_to_fill_potential_times"}
    - action_checkpoint_with_temp_times_to_fill_potential_times
> checkpoint_with_temp_times_to_fill_potential_times

## checkpoint_reservation_slot_filling_base_return_6
> checkpoint_reservation_slot_filling_base_return
    - followup{"name": "utter_ask_for_name"}
    - utter_ask_for_name
> checkpoint_to_fill_name

## checkpoint_reservation_slot_filling_base_return_7
> checkpoint_reservation_slot_filling_base_return
    - followup{"name": "action_brancher_to_save_new_reservation"}
    - action_brancher_to_save_new_reservation

