## checkpoint_reservation_slot_filling_base_return_1
    - action_brancher_reservation_slot_filling_base
    - followup{"name": "utter_ask_for_number_on_reservation_size"}
    - utter_ask_for_number_on_reservation_size
> checkpoint_to_fill_potential_size

## checkpoint_reservation_slot_filling_base_return_2
    - action_brancher_reservation_slot_filling_base
    - followup{"name": "action_brancher_validate_reservation_potential_size"}
    - action_brancher_validate_reservation_potential_size

## checkpoint_reservation_slot_filling_base_return_3
    - action_brancher_reservation_slot_filling_base
    - followup{"name": "action_brancher_validate_reservation_potential_size"}
    - action_brancher_validate_reservation_potential_size

## checkpoint_reservation_slot_filling_base_return_4
    - action_brancher_reservation_slot_filling_base
    - followup{"name": "utter_ask_for_time_for_potential_reservation"}
    - utter_ask_for_time_for_potential_reservation
> checkpoint_to_fill_temp_times

## checkpoint_reservation_slot_filling_base_return_5
    - action_brancher_reservation_slot_filling_base
    - followup{"name": "action_brancher_with_temp_times_to_determine_next_from_times_length"}
    - action_brancher_with_temp_times_to_determine_next_from_times_length

## checkpoint_reservation_slot_filling_base_return_6
    - action_brancher_reservation_slot_filling_base
    - followup{"name": "utter_ask_for_name_on_reservation"}
    - utter_ask_for_name_on_reservation
> checkpoint_to_fill_name

## checkpoint_reservation_slot_filling_base_return_7
    - action_brancher_reservation_slot_filling_base
    - followup{"name": "action_brancher_to_save_new_reservation"}
    - action_brancher_to_save_new_reservation

